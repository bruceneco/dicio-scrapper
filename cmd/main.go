package main

import (
	"dicio-scrapper/config"
	"dicio-scrapper/internal/adapters"
	"dicio-scrapper/internal/domain"

	"github.com/ipfans/fxlogger"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
)

func main() {
	env := config.LoadEnv()
	config.SetupLogger(env)
	fx.New(
		fx.Provide(func() *config.EnvConfig { return env }),
		fx.WithLogger(fxlogger.WithZerolog(log.Logger)),
		adapters.Module,
		domain.Module,
	).Run()
}
