package domain

import (
	"dicio-scrapper/internal/domain/word"

	"go.uber.org/fx"
)

var Module = fx.Provide(
	word.NewService,
)
