package scrapper

import "go.uber.org/fx"

var Module = fx.Provide(
	NewScrapper,
	NewWord,
)
