package scrapper

import (
	"dicio-scrapper/internal/domain/core"
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

func NewWord(scrapper *Scrapper) *Word {
	return &Word{scrapper: scrapper}
}

func (w *Word) Scrape(searchSlug string) (core.Word, error) {
	searchSlug = strings.ReplaceAll(strings.TrimSpace(sanitize.Accents(searchSlug)), " ", "-")

	c := w.scrapper.Collector()

	var (
		mu   sync.Mutex
		word core.Word
	)

	setContent(c, &mu, &word)
	setMeanings(c, &mu, &word)
	setEtymologies(c, &mu, &word)
	setPhrases(c, &mu, &word)
	setSynonyms(c, &mu, &word)

	c.OnError(func(_ *colly.Response, cErr error) {
		log.Info().Err(cErr).Send()
	})

	err := c.Visit(fmt.Sprintf("https://www.dicio.com.br/%s", searchSlug))
	if err != nil {
		log.Error().Err(err).Msg("failed to scrape word")
		return core.Word{}, err
	}

	return word, nil
}

func setContent(c *colly.Collector, mu *sync.Mutex, word *core.Word) {
	c.OnHTML("div.title-header > h1", func(e *colly.HTMLElement) {
		mu.Lock()
		defer mu.Unlock()

		word.Content = strings.TrimSpace(strings.Trim(e.Text, "\n"))
	})
}

func setMeanings(c *colly.Collector, mu *sync.Mutex, word *core.Word) {
	c.OnHTML("p.significado > span:not(.cl):not(.etim)", func(e *colly.HTMLElement) {
		mu.Lock()
		defer mu.Unlock()

		tag := e.DOM.Find("span.tag").Text()

		definition := core.Meaning{
			Tag:     strings.TrimSuffix(strings.TrimPrefix(tag, "["), "]"),
			Content: strings.TrimSpace(strings.Replace(e.Text, tag, "", 1)),
		}

		word.Meanings = append(word.Meanings, definition)
	})
}

func setEtymologies(c *colly.Collector, mu *sync.Mutex, word *core.Word) {
	c.OnHTML("p.significado > span.etim", func(e *colly.HTMLElement) {
		mu.Lock()
		defer mu.Unlock()
		word.Etymologies = append(word.Etymologies, e.Text)
	})
}

func setPhrases(c *colly.Collector, mu *sync.Mutex, word *core.Word) {
	c.OnHTML(".frases > .frase", func(e *colly.HTMLElement) {
		mu.Lock()
		defer mu.Unlock()

		by := e.DOM.Find("em").Text()
		content := strings.ReplaceAll(e.Text, by, "")

		by = strings.TrimPrefix(by, "- ")
		content = strings.ReplaceAll(content, "\n", "")
		content = strings.Trim(content, " ")

		phrase := core.Phrase{
			By:      by,
			Content: content,
		}

		word.Phrases = append(word.Phrases, phrase)
	})
}
func setSynonyms(c *colly.Collector, mu *sync.Mutex, word *core.Word) {
	c.OnHTML("p.sinonimos > a", func(e *colly.HTMLElement) {
		mu.Lock()
		defer mu.Unlock()

		word.Synonyms = append(word.Synonyms, e.Text)
	})
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
