package config

import (
	"time"

	"github.com/spf13/viper"
)

func Wait() {
	ms := viper.GetDuration("WAIT")
	if ms > 0 {
		time.Sleep(time.Duration(ms) * time.Millisecond)
	}
}
