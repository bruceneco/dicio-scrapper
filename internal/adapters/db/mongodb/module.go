package mongodb

import (
	"dicio-scrapper/internal/ports/wordports"

	"go.uber.org/fx"
)

var Module = fx.Provide(
	NewConnection,
	fx.Annotate(NewWord, fx.As(new(wordports.Repo))),
)
