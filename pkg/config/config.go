package config

import (
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	viper.AutomaticEnv()
}

type ServerConfig struct {
	Network string
	Port    string
}

func LoadServerConfig() ServerConfig {
	return ServerConfig{
		Network: "tcp",
		Port:    viper.GetString("SERVER_ADDRESS"),
	}
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
