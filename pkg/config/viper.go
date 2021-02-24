package config

import (
	"encoding/json"
	"strings"

	"github.com/xuanbo/ohmydata/pkg/log"

	"github.com/spf13/viper"
)

var internalViper *viper.Viper

// Init 初始化
func Init() error {
	log.Logger().Info("初始化配置")

	internalViper = viper.New()

	// 优先使用环境变量
	internalViper.AutomaticEnv()
	// database.mysql.dns => DATABASE_MYSQL_DNS
	internalViper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	internalViper.SetConfigName("config")
	internalViper.SetConfigType("yaml")
	internalViper.AddConfigPath("/etc/ohmydata/")
	internalViper.AddConfigPath(".")
	internalViper.AddConfigPath("./config")
	if err := internalViper.ReadInConfig(); err != nil {
		return err
	}

	b, err := json.MarshalIndent(internalViper.AllSettings(), "", "  ")
	if err != nil {
		return err
	}
	log.Logger().Info("配置: " + string(b))

	return nil
}

// GetInt 获取配置
func GetInt(key string) int {
	return internalViper.GetInt(key)
}

// GetString 获取配置
func GetString(key string) string {
	return internalViper.GetString(key)
}
