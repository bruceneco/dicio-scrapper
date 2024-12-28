package adapters

import (
	"dicio-scrapper/internal/adapters/amqp"
	"dicio-scrapper/internal/adapters/db/postgres"
	"dicio-scrapper/internal/adapters/http"
	"dicio-scrapper/internal/adapters/scrapper"

	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		postgres.NewConnection,
	),
	http.Module,
	amqp.Module,
	scrapper.Module,
)
