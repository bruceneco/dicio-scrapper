package scrapper

import (
	"dicio-scrapper/internal/ports/wordports"

	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewWord,
	func(s *Word) wordports.Scrapper {
		return s
	},
	NewScrapper,
)
