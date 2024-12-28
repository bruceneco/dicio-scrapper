package settings

import (
	"context"
	"dicio-scrapper/config"

	"github.com/rs/zerolog/log"
	"github.com/wagslane/go-rabbitmq"
	"go.uber.org/fx"
)

type (
	Connection struct {
		conn *rabbitmq.Conn
	}

	ExchangeOpts struct {
		Name ExchangeName
		Kind ExchangeType
	}

	ConsumerOpts struct {
		QueueName   string
		Exclusive   bool
		NoWait      bool
		Concurrency int
	}
)

func NewConnection(lc fx.Lifecycle, cfg *config.EnvConfig) *Connection {
	conn, err := rabbitmq.NewConn(
		cfg.AMQPHost,
		rabbitmq.WithConnectionOptionsLogger(NewLoggerAdapter(&log.Logger)),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot initialize amqp connection")
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			if conn == nil {
				return nil
			}
			return conn.Close()
		},
	})

	return &Connection{conn: conn}
}

func (c *Connection) MakePublisher(opts ExchangeOpts) (*rabbitmq.Publisher, error) {
	p, err := rabbitmq.NewPublisher(
		c.conn,
		rabbitmq.WithPublisherOptionsLogger(NewLoggerAdapter(&log.Logger)),
		rabbitmq.WithPublisherOptionsExchangeName(opts.Name.String()),
		rabbitmq.WithPublisherOptionsExchangeDurable,
		rabbitmq.WithPublisherOptionsExchangeDeclare,
		rabbitmq.WithPublisherOptionsExchangeKind(opts.Kind.String()),
	)
	if err != nil {
		return nil, err
	}

	p.NotifyReturn(func(r rabbitmq.Return) {
		log.Info().Interface("content", r).Msg("message returned from exchange")
	})

	return p, nil
}

func (c *Connection) MakeConsumer(name string, cOpts ConsumerOpts, exOpts ExchangeOpts) (*rabbitmq.Consumer, error) {
	consumer, err := rabbitmq.NewConsumer(
		c.conn,
		name,
		rabbitmq.WithConsumerOptionsQOSPrefetch(cOpts.Concurrency*PrefetchCount),
		rabbitmq.WithConsumerOptionsLogger(NewLoggerAdapter(&log.Logger)),
		rabbitmq.WithConsumerOptionsExchangeName(exOpts.Name.String()),
		rabbitmq.WithConsumerOptionsExchangeDurable,
		rabbitmq.WithConsumerOptionsExchangeDeclare,

		rabbitmq.WithConsumerOptionsExchangeKind(exOpts.Kind.String()),
		rabbitmq.WithConsumerOptionsConcurrency(cOpts.Concurrency),
		rabbitmq.WithConsumerOptionsQueueQuorum,
		func(options *rabbitmq.ConsumerOptions) {
			options.QueueOptions = rabbitmq.QueueOptions{
				Name:       cOpts.QueueName,
				Durable:    true,
				AutoDelete: false,
				Exclusive:  cOpts.Exclusive,
				Passive:    false,
				NoWait:     cOpts.NoWait,
				Declare:    true,
				Args: rabbitmq.Table{ // queue args
					"delivery-limit": "3",
				},
			}
			options.CloseGracefully = true
		},
	)

	return consumer, err
}
