package db

import (
	"dicio-scrapper/internal/adapters/db/mongodb"

	"go.uber.org/fx"
)

var Module = fx.Options(
	mongodb.Module,
)
