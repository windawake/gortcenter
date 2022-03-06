package config

import (
	"log"

	viperlib "github.com/spf13/viper" // 自定义包名，避免与内置 viper 实例冲突
)

// viper 库实例
var App *viperlib.Viper

func init() {
	App = viperlib.New()
	App.SetConfigType("json")
	App.AddConfigPath("./config")
	App.SetConfigName("app")

	if err := App.ReadInConfig(); err != nil {
		log.Fatalf("read config failed: %v", err)
	}
}
