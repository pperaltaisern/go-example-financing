package config

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	viper.AutomaticEnv()
	viper.AddConfigPath(directoryPath())
	viper.SetConfigName("config_local")
	viper.ReadInConfig()
}

func directoryPath() string {
	_, fileName, _, _ := runtime.Caller(0)
	prefixPath := filepath.Dir(fileName)
	return prefixPath + "/"
}

type ServerConfig struct {
	Network string
	Address string
}

func LoadCommandServerConfig() ServerConfig {
	return ServerConfig{
		Network: "tcp",
		Address: viper.GetString("COMMAND_SERVER_ADDRESS"),
	}
}

func LoadQueryServerConfig() ServerConfig {
	return ServerConfig{
		Network: "tcp",
		Address: viper.GetString("QUERY_SERVER_ADDRESS"),
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
