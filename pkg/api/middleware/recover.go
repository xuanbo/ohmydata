package middleware

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// Recover 恢复中间件
func Recover(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					logger.Error("PANIC RECOVER", zap.Error(err), zap.Stack("stack"))
					ctx.Error(err)
				}
			}()
			return next(ctx)
		}
	}
}
