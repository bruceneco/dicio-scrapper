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
		publishers.NewWord,
		func(p *publishers.Word) wordports.Publisher {
			return p
		},
	),
	fx.Invoke(
		consumers.StartExtractWord,
	),
)
