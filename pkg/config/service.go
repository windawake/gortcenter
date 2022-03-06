package config

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	viperlib "github.com/spf13/viper" // 自定义包名，避免与内置 viper 实例冲突
)

// viper 库实例
var Service *viperlib.Viper

type ServiceConfig struct{}

func init() {
	Service = viperlib.New()
	Service.SetConfigType("json")
	Service.AddConfigPath("./config")
	Service.SetConfigName("service")

	if err := Service.ReadInConfig(); err != nil {
		log.Fatalf("read config failed: %v", err)
	}
}

func (c *ServiceConfig) GetDb(name string) (*sqlx.DB, error) {
	map1 := Service.GetStringMapString("db." + name)
	sdn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", map1["username"], map1["password"], map1["host"], map1["port"], map1["database"])

	db, err := sqlx.Connect("mysql", sdn)
	if err != nil {
		return nil, err
	}

	return db, nil
}
