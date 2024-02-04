package main

import (
	"github.com/spf13/viper"
	"gitlab.mvalley.com/wind/rime-utils/internal/pkg/config"
	"gitlab.mvalley.com/wind/rime-utils/internal/server"
)

func main() {
	err := InitConfiguration("config", "./configs/", &config.CONFIG)
	if err != nil {
		panic(err)
	}

	server.New().Run()
}

// InitConfiguration ...
func InitConfiguration(configName string, configPath string, config interface{}) error {
	vp := viper.New()
	vp.SetConfigName(configName)
	vp.AutomaticEnv()
	vp.AddConfigPath(configPath)

	if err := vp.ReadInConfig(); err != nil {
		return err
	}

	err := vp.Unmarshal(config)
	if err != nil {
		return err
	}

	return nil
}
