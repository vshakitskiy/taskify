package config

import (
	shconfig "app.shared/config"
	"app.shared/pkg/env"
)

var Config *shconfig.Config

func InitConfig() error {
	if err := env.LoadEnv(); err != nil {
		return err
	}
	conf, err := shconfig.ReadConfig()
	Config = conf
	return err
}
