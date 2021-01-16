package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/xuanbo/ohmydata/pkg/api/middleware"
	"github.com/xuanbo/ohmydata/pkg/api/util"
	"github.com/xuanbo/ohmydata/pkg/config"
	"github.com/xuanbo/ohmydata/pkg/log"
	"github.com/xuanbo/ohmydata/pkg/model"

	"github.com/labstack/echo/v4"
	mw "github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

// ServeHTTP 启动服务
func ServeHTTP() error {
	// Echo instance
	router := echo.New()

	// Middleware
	router.HTTPErrorHandler = customHTTPErrorHandler
	router.Use(middleware.ZapLogger(log.Logger()))
	router.Use(middleware.Recover(log.Logger()))
	router.Use(middleware.NewContext())
	router.Use(mw.JWTWithConfig(mw.JWTConfig{
		// 跳过登录
		Skipper: func(ctx echo.Context) bool {
			return ctx.Path() == "/v1/user/login"
		},
		SigningKey:  []byte("secret"),
		ContextKey:  "JWT_TOKEN",
		TokenLookup: "header:" + echo.HeaderAuthorization,
		ErrorHandler: func(err error) error {
			if he, ok := err.(*echo.HTTPError); ok {
				message := fmt.Sprintf("%s", he.Message)
				return echo.NewHTTPError(http.StatusUnauthorized, message)
			}
			return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
		},
		SuccessHandler: func(ctx echo.Context) {
			if token, ok := ctx.Get("JWT_TOKEN").(*jwt.Token); ok {
				userID := token.Header["userId"]
				userName := token.Header["userName"]
				username := token.Header["username"]
				ctx.Set("userId", userID)
				ctx.Set("userName", userName)
				ctx.Set("username", username)

				// 传递到context
				cctx := ctx.(*middleware.Context)
				c := cctx.Ctx()
				c = context.WithValue(c, util.UserID, userID)
				cctx.SetCtx(c)
			}
		},
	}))

	// 初始化路由
	if err := InitRoutes(router); err != nil {
		return err
	}

	addr := config.GetString("http.addr")
	if addr == "" {
		addr = ":9090"
	}

	log.Logger().Info("启动HTTP服务", zap.String("addr", addr))

	// Start server
	return router.Start(addr)
}

func customHTTPErrorHandler(err error, ctx echo.Context) {
	if he, ok := err.(*echo.HTTPError); ok {
		message := fmt.Sprintf("%s", he.Message)
		if err := ctx.JSON(he.Code, model.Fail(message)); err != nil {
			log.Logger().Error("统一异常处理响应错误", zap.Error(err))
		}
		return
	}
	if err := ctx.JSON(http.StatusOK, model.Fail(err.Error())); err != nil {
		log.Logger().Error("统一异常处理响应错误", zap.Error(err))
	}
}
