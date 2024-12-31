package adapters

import (
	"dicio-scrapper/internal/adapters/amqp"
	"dicio-scrapper/internal/adapters/db"
	"dicio-scrapper/internal/adapters/http"
	"dicio-scrapper/internal/adapters/scrapper"

	"go.uber.org/fx"
)

var Module = fx.Options(
	http.Module,
	amqp.Module,
	scrapper.Module,
	db.Module,
)
