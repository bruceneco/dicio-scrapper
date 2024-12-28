package settings

import "github.com/rabbitmq/amqp091-go"

const PrefetchCount = 10

type ExchangeName string

const (
	DefaultExchangeName ExchangeName = "events"
)

type ExchangeType string

const (
	DefaultExchangeType ExchangeType = amqp091.ExchangeHeaders
)

func (n *ExchangeName) String() string {
	return string(*n)
}

func (t *ExchangeType) String() string {
	return string(*t)
}
