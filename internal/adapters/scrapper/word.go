package scrapper

import (
	"dicio-scrapper/internal/ports/wordports"
	"fmt"
	"strings"
	"sync"

	"github.com/gocolly/colly/v2"
	"github.com/kennygrant/sanitize"
	"github.com/rs/zerolog/log"
)

type Word struct {
	scrapper *Scrapper
}

func NewWord(scrapper *Scrapper) wordports.Scrapper {
	return &Word{scrapper: scrapper}
}

func (w *Word) Scrape(word string) ([]string, error) {
	word = strings.ReplaceAll(strings.TrimSpace(sanitize.Accents(word)), " ", "-")

	c := w.scrapper.Collector()

	c.OnHTML(
		"#content > div.col-xs-12.col-sm-7.col-md-8.p0.mb20"+
			" > div.card.card-main.mb10 > p > span:nth-child(1)",
		func(e *colly.HTMLElement) {
			log.Info().Str("word", word).Str("type", e.Text).Send()
		})

	c.OnError(func(_ *colly.Response, cErr error) {
		log.Info().Err(cErr).Send()
	})

	err := c.Visit(fmt.Sprintf("https://www.dicio.com.br/%s", word))
	if err != nil {
		log.Error().Err(err).Msg("failed to scrape word")
		return nil, err
	}

	return []string{word}, nil
}

func (w *Word) MostSearched(page int) []string {
	var words []string
	var mu sync.Mutex
	c := w.scrapper.Collector()

	c.OnHTML("ul.list > li > a", func(e *colly.HTMLElement) {
		mu.Lock()
		defer mu.Unlock()

		words = append(words, e.Text)
	})

	c.OnRequest(func(r *colly.Request) {
		log.Info().Msgf("Scrapping %s", r.URL.String())
	})

	err := c.Visit(fmt.Sprintf("https://www.dicio.com.br/palavras-mais-buscadas/%d/", page))
	if err != nil {
		log.Error().Err(err).Msg("failed to scrape page")
		return nil
	}

	return words
}
