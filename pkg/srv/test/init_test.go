package srv_test

import (
	"github.com/xuanbo/ohmydata/pkg/config"
	"github.com/xuanbo/ohmydata/pkg/db"
	"github.com/xuanbo/ohmydata/pkg/db/mysql"
	"github.com/xuanbo/ohmydata/pkg/log"
	"github.com/xuanbo/ohmydata/pkg/srv"
)

func init() {
	// 日志
	if err := log.Init(); err != nil {
		panic(err)
	}

	// 配置
	if err := config.Init(); err != nil {
		panic(err)
	}

	// 驱动
	if err := mysql.Register(); err != nil {
		panic(err)
	}

	// 数据库
	if err := db.Init(); err != nil {
		panic(err)
	}

	// 同步数据源
	if err := srv.SyncDataSource(srv.NewDataSource()); err != nil {
		panic(err)
	}
}
