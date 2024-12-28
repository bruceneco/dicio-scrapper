package word

import (
	"context"
	"dicio-scrapper/internal/ports/wordports"

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
	}

	Service struct {
		publisher wordports.Publisher
		scrapper  wordports.Scrapper
	}
)

func NewService(params ServicerParams) Servicer {
	s := &Service{
		publisher: params.Publisher,
		scrapper:  params.Scrapper,
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
		err := s.publisher.ExtractWord(ctx, word)
		if err != nil {
			log.Error().Str("word", word).Err(err).Msg("failed to enqueue word")

			return err
		}
	}

	return nil
}

func (s *Service) Extract(word string) error {
	log.Info().Str("word", word).Msg("extracting")

	_, err := s.scrapper.Scrape(word)
	if err != nil {
		return err
	}

	return nil
}
