package consumers

import (
	"context"
	"dicio-scrapper/internal/adapters/amqp/settings"
	"dicio-scrapper/internal/domain/word"

	"github.com/rs/zerolog/log"
	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/fx"
)

type (
	ExtractWordParams struct {
		fx.In
		LC      fx.Lifecycle
		Conn    *settings.Connection
		Service word.Servicer
	}
	extractWordHandler struct {
		consumer *rabbitmq.Consumer
		service  word.Servicer
	}
)

func StartExtractWord(params ExtractWordParams) {
	const numConcurrency = 16

	consumer, err := params.Conn.MakeConsumer("extract_word", settings.ConsumerOpts{
		QueueName:   "extract_word",
		Concurrency: numConcurrency,
	}, settings.ExchangeOpts{
		Name: settings.DefaultExchangeName,
		Kind: settings.DefaultExchangeType,
	})

	if err != nil {
		log.Fatal().Err(err).Msg("cannot initialize extract word amqp consumer")
	}

	params.LC.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			consumer.CloseWithContext(ctx)
			return nil
		},
	})

	handler := &extractWordHandler{
		consumer: consumer,
		service:  params.Service,
	}

	go func() {
		err = consumer.Run(handler.handle)
		if err != nil {
			log.Fatal().Err(err).Msg("cannot run extract word amqp consumer")
		}
	}()
}

func (h *extractWordHandler) handle(d rabbitmq.Delivery) rabbitmq.Action {
	w := string(d.Body)
	if err := h.service.Extract(w); err != nil {
		if d.Redelivered {
			log.Error().Str("word", w).Err(err).Msg("cannot extract word, discarding")

			return rabbitmq.NackDiscard
		}

		return rabbitmq.NackRequeue
	}

	return rabbitmq.Ack
}
