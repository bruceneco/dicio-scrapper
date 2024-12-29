package wordports

import (
	"dicio-scrapper/internal/domain/core"
)

type Scrapper interface {
	Scrape(word string) (core.Word, error)
	MostSearched(page int) []string
}
