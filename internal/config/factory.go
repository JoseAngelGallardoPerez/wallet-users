package config

import "github.com/inconshreveable/log15"

func Factory(logger log15.Logger) *Configuration {
	if err := InitConfig(logger.New("service", "configReader")); err != nil {
		panic("cannot init settings: " + err.Error())
	}
	return GetConf()
}
