package adapters

import (
	"dicio-scrapper/internal/adapters/amqp"
	"dicio-scrapper/internal/adapters/db/postgres"
	"dicio-scrapper/internal/adapters/http"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		postgres.NewConnection,
		amqp.NewConnection,
	),
	http.Module,
)
