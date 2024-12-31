package word

import (
	"context"
	"dicio-scrapper/internal/ports/wordports"
	"errors"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

type (
	Servicer interface {
		EnqueueExtraction(ctx context.Context, word string) error
		EnqueueMostSearched(ctx context.Context, page int) error
		Extract(word string) error
	}

	ServicerParams struct {
		fx.In
		Publisher wordports.Publisher
		Scrapper  wordports.Scrapper
		Repo      wordports.Repo
	}

	Service struct {
		publisher wordports.Publisher
		scrapper  wordports.Scrapper
		repo      wordports.Repo
	}
)

func NewService(params ServicerParams) *Service {
	s := &Service{
		publisher: params.Publisher,
		scrapper:  params.Scrapper,
		repo:      params.Repo,
	}
	return s
}

func (s *Service) EnqueueExtraction(ctx context.Context, word string) error {
	return s.publisher.ExtractWord(ctx, word)
}

func (s *Service) EnqueueMostSearched(ctx context.Context, page int) error {
	words := s.scrapper.MostSearched(page)
	log.Info().Int("words", len(words)).Int("page", page).Msg("enqueuing most searched words")

	for _, word := range words {
		err := s.publisher.ExtractWord(ctx, strings.TrimSpace(word))
		if err != nil {
			log.Error().Str("word", word).Err(err).Msg("failed to enqueue word")

			return err
		}
	}

	return nil
}

func (s *Service) Extract(word string) error {
	ctx, cancel := s.MakeCtx()
	defer cancel()

	l := log.With().Str("word", word).Logger()

	if err := s.checkDuplication(ctx, word); err != nil {
		return err
	}

	scrappedWord, err := s.scrapper.Scrape(word)
	if err != nil {
		l.Err(err).Msg("failed to scrape word")
		return err
	}

	_, err = s.repo.Insert(ctx, scrappedWord)
	if err != nil {
		l.Err(err).Msg("failed to insert word into repo")
		return err
	}

	return nil
}

func (s *Service) MakeCtx() (context.Context, func()) {
	maxTimeout := 5
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(maxTimeout)*time.Second)
	return ctx, cancel
}

func (s *Service) checkDuplication(ctx context.Context, word string) error {
	_, err := s.repo.FindByContent(ctx, word)
	if err != nil {
		if errors.Is(err, wordports.ErrWordNotFound) {
			return nil
		}

		return err
	}

	return wordports.ErrWordAlreadyExists
}
