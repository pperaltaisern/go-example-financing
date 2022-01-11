package main

import (
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	CommandServer CommandServerConfig
	AMQP          AMQPConfig
	Postgres      PostgresConfig
}

type CommandServerConfig struct {
	Network string
	Address string
}

type AMQPConfig struct {
	Address string
}

type PostgresConfig struct {
	ConnectionString string
}

func NewConfig(directory, configFileName string) (Config, error) {
	vs := strings.Split(configFileName, ".json")
	configFileName = vs[0]

	v := viper.New()
	v.AddConfigPath(directoryPath())
	v.AddConfigPath(directory)

	v.SetConfigName(configFileName)
	err := v.ReadInConfig()
	if err != nil {
		return Config{}, err
	}

	var c Config
	err = v.Unmarshal(&c)
	return c, err
}

func directoryPath() string {
	_, fileName, _, _ := runtime.Caller(0)
	prefixPath := filepath.Dir(fileName)
	return prefixPath + "/"
}
