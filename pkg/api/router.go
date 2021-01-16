package api

import (
	v1 "github.com/xuanbo/ohmydata/pkg/api/v1"
	"github.com/xuanbo/ohmydata/pkg/srv"

	"github.com/labstack/echo/v4"
)

type controller interface {
	// 初始化
	Init() error
	// 添加路由
	AddRoutes(router *echo.Echo)
}

// InitRoutes 初始化路由
func InitRoutes(router *echo.Echo) error {
	if err := addRoutes(router, v1.NewUser(srv.NewUser())); err != nil {
		return err
	}
	if err := addRoutes(router, v1.NewDict()); err != nil {
		return err
	}
	if err := addRoutes(router, v1.NewDataSource(srv.NewDataSource())); err != nil {
		return err
	}
	if err := addRoutes(router, v1.NewDataSet(srv.NewDataSet())); err != nil {
		return err
	}
	return nil
}

func addRoutes(router *echo.Echo, controllers ...controller) error {
	for _, controller := range controllers {
		if err := controller.Init(); err != nil {
			return err
		}
		controller.AddRoutes(router)
	}
	return nil
}
