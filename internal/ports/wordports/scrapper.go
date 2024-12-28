package wordports

type Scrapper interface {
	Scrape(word string) ([]string, error)
	MostSearched(page int) []string
}
