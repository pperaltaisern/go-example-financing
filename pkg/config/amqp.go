package config

import (
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/pperaltaisern/financing/internal/watermillzap"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type AMQPConfig struct {
	Address string
}

func LoadAMQPConfig() AMQPConfig {
	return AMQPConfig{
		Address: viper.GetString("AMQP_ADDRESS"),
	}
}

func (c AMQPConfig) BuildCommandPublisher(log *zap.Logger) (*amqp.Publisher, error) {
	return amqp.NewPublisher(
		newCommandConfig(c.Address),
		watermillzap.NewLogger(log))
}

func (c AMQPConfig) BuildCommandSubscriber(log *zap.Logger) (*amqp.Subscriber, error) {
	return amqp.NewSubscriber(
		newCommandConfig(c.Address),
		watermillzap.NewLogger(log))
}

func (c AMQPConfig) BuildEventPublisher(log *zap.Logger) (*amqp.Publisher, error) {
	return amqp.NewPublisher(
		newEventConfig(c.Address, nil),
		watermillzap.NewLogger(log))
}

func (c AMQPConfig) BuildEventSubscriber(log *zap.Logger, handlerName string) (*amqp.Subscriber, error) {
	return amqp.NewSubscriber(
		newEventConfig(c.Address, amqp.GenerateQueueNameTopicNameWithSuffix(handlerName)),
		watermillzap.NewLogger(log))
}

func (c AMQPConfig) BuildCommandBus(log *zap.Logger) (*cqrs.EventBus, error) {
	amqpConfig := LoadAMQPConfig()

	publisher, err := amqpConfig.BuildEventPublisher(log)
	if err != nil {
		return nil, err
	}

	return cqrs.NewEventBus(
		publisher,
		generateEventsTopic,
		CommandEventMarshaler)
}

func (c AMQPConfig) BuildEventBus(log *zap.Logger) (*cqrs.EventBus, error) {
	pub, err := c.BuildCommandPublisher(log)
	if err != nil {
		return nil, err
	}
	return cqrs.NewEventBus(
		pub,
		generateEventsTopic,
		CommandEventMarshaler)
}

var newCommandConfig = amqp.NewDurableQueueConfig
var newEventConfig = amqp.NewDurablePubSubConfig
