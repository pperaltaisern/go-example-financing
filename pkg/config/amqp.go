package config

import (
	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
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

var newCommandConfig = amqp.NewDurableQueueConfig
var newEventConfig = amqp.NewDurablePubSubConfig
