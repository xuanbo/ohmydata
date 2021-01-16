package config

import (
	"encoding/json"

	"github.com/xuanbo/ohmydata/pkg/log"

	"github.com/spf13/viper"
)

// Init 初始化
func Init() error {
	log.Logger().Info("初始化配置")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/ohmydata/")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	b, err := json.MarshalIndent(viper.AllSettings(), "", "  ")
	if err != nil {
		return err
	}
	log.Logger().Info("配置: " + string(b))

	return nil
}

// GetInt 获取配置
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetString 获取配置
func GetString(key string) string {
	return viper.GetString(key)
}
