package main

import (
	"github.com/xuanbo/ohmydata/pkg/api"
	"github.com/xuanbo/ohmydata/pkg/cache"
	"github.com/xuanbo/ohmydata/pkg/config"
	"github.com/xuanbo/ohmydata/pkg/db"
	"github.com/xuanbo/ohmydata/pkg/db/elastic"
	"github.com/xuanbo/ohmydata/pkg/db/mysql"
	"github.com/xuanbo/ohmydata/pkg/db/postgres"
	"github.com/xuanbo/ohmydata/pkg/log"

	"go.uber.org/zap"
)

func main() {
	// 日志
	if err := log.Init(); err != nil {
		panic(err)
	}
	defer log.Flush()

	// 配置
	if err := config.Init(); err != nil {
		log.Logger().Panic("初始化配置错误", zap.Error(err))
	}

	// 驱动适配层
	if err := mysql.Register(); err != nil {
		log.Logger().Panic("注册mysql驱动错误", zap.Error(err))
	}
	if err := postgres.Register(); err != nil {
		log.Logger().Panic("注册postgres驱动错误", zap.Error(err))
	}
	if err := elastic.Register(); err != nil {
		log.Logger().Panic("注册elastic驱动错误", zap.Error(err))
	}

	// 初始化redis
	if err := cache.Init(); err != nil {
		log.Logger().Panic("初始化redis错误", zap.Error(err))
	}

	// 同步数据库
	if err := db.Init(); err != nil {
		log.Logger().Panic("同步数据库表结构错误", zap.Error(err))
	}

	// 启动服务
	if err := api.ServeHTTP(); err != nil {
		log.Logger().Panic("启动HTTP服务错误", zap.Error(err))
	}
}
