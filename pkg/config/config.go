package config

import (
	"context"
	"strings"

	"github.com/pperaltaisern/financing/internal/watermillzap"

	"github.com/ThreeDotsLabs/watermill-amqp/pkg/amqp"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ServerConfig struct {
	Network string
	Port    string
}

func LoadServerConfig() ServerConfig {
	return ServerConfig{
		Network: "tcp",
		Port:    viper.GetString("SERVER_PORT"),
	}
}

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

type PostgresConfig struct {
	ConnectionString string
}

func LoadPostgresConfig() PostgresConfig {
	return PostgresConfig{
		ConnectionString: viper.GetString("DB_CONNECTION_STRING"),
	}
}

func (c PostgresConfig) Build() (*pgxpool.Pool, error) {
	return pgxpool.Connect(context.Background(), c.ConnectionString)
}

type LoggerConfig struct {
	Level string
}

func LoadLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Level: viper.GetString("LOGGER_LEVEL"),
	}
}

func (c LoggerConfig) Build() (*zap.Logger, error) {
	zc := zap.NewDevelopmentConfig()

	if c.Level != "" {
		var lvl zapcore.Level
		switch strings.ToUpper(c.Level) {
		case "DEBUG":
			lvl = zapcore.DebugLevel
		case "INFO":
			lvl = zapcore.InfoLevel
		case "WARNING":
			lvl = zapcore.WarnLevel
		case "ERROR":
			lvl = zapcore.ErrorLevel
		case "FATAL":
			lvl = zapcore.FatalLevel
		}
		zc.Level = zap.NewAtomicLevelAt(lvl)
	}

	return zc.Build()
}
