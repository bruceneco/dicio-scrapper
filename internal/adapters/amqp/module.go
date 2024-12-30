package amqp

import (
	"dicio-scrapper/internal/adapters/amqp/consumers"
	"dicio-scrapper/internal/adapters/amqp/publishers"
	"dicio-scrapper/internal/adapters/amqp/settings"
	"dicio-scrapper/internal/ports/wordports"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		settings.NewConnection,
		fx.Annotate(publishers.NewWord, fx.As(new(wordports.Publisher))),
	),
	fx.Invoke(
		consumers.StartExtractWord,
	),
)
