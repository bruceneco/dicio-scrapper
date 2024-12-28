package scrapper

import "github.com/gocolly/colly/v2"

type Scrapper struct {
	collector *colly.Collector
}

func NewScrapper() *Scrapper {
	return &Scrapper{
		collector: colly.NewCollector(colly.AllowURLRevisit()),
	}
}

func (s *Scrapper) Collector() *colly.Collector {
	return s.collector.Clone()
}
