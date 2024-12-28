package amqp

import (
	"dicio-scrapper/internal/adapters/amqp/consumers"
	"dicio-scrapper/internal/adapters/amqp/publishers"
	"dicio-scrapper/internal/adapters/amqp/settings"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		settings.NewConnection,
		publishers.NewWord,
	),
	fx.Invoke(
		consumers.StartExtractWord,
	),
)
