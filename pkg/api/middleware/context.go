package middleware

import (
	"context"

	"github.com/labstack/echo/v4"
)

// Context 上下文
type Context struct {
	echo.Context
}

// Ctx 返回context.Context
func (c Context) Ctx() context.Context {
	return c.Get("CTX").(context.Context)
}

// SetCtx 设置context.Context
func (c Context) SetCtx(ctx context.Context) {
	c.Set("CTX", ctx)
}

// NewContext 包装echo.Context
func NewContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			ctx.Set("CTX", context.Background())
			return next(&Context{ctx})
		}
	}
}
