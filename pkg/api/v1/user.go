package v1

import (
	"net/http"

	"github.com/xuanbo/ohmydata/pkg/api/middleware"
	"github.com/xuanbo/ohmydata/pkg/entity"
	"github.com/xuanbo/ohmydata/pkg/model"
	"github.com/xuanbo/ohmydata/pkg/srv"

	"github.com/labstack/echo/v4"
)

// User 用户API管理
type User struct {
	srv *srv.User
}

// NewUser 创建
func NewUser(srv *srv.User) *User {
	return &User{srv}
}

// Init 初始化
func (u *User) Init() error {
	return nil
}

// AddRoutes 添加路由
func (u *User) AddRoutes(e *echo.Echo) {
	g := e.Group("/v1")
	{
		// 用户管理
		g.POST("/user/login", u.Login)
	}
}

// Login 登录
func (u *User) Login(ctx echo.Context) error {
	var user entity.User
	if err := ctx.Bind(&user); err != nil {
		return err
	}
	c := ctx.(*middleware.Context).Ctx()
	s, err := u.srv.Login(c, &user)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(s))
}
