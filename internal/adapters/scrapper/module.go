package scrapper

import (
	"dicio-scrapper/internal/ports/wordports"

	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewWord,
	fx.Annotate(NewWord, fx.As(new(wordports.Scrapper))),
	NewScrapper,
)
