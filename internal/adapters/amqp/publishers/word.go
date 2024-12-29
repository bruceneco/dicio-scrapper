package publishers

import (
	"context"
	"dicio-scrapper/internal/adapters/amqp/settings"

	"github.com/rs/zerolog/log"
	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/fx"
)

type Word struct {
	p *rabbitmq.Publisher
}

func NewWord(conn *settings.Connection, lc fx.Lifecycle) *Word {
	publisher, err := conn.MakePublisher(settings.ExchangeOpts{
		Name: settings.DefaultExchangeName,
		Kind: settings.DefaultExchangeType,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("cannot initialize extract word amqp publisher")
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			publisher.Close()
			return nil
		},
	})

	return &Word{p: publisher}
}

func (e *Word) ExtractWord(ctx context.Context, word string) error {
	return e.p.PublishWithContext(ctx, []byte(word), []string{"extract_word"},
		rabbitmq.WithPublishOptionsPersistentDelivery,
	)
}
