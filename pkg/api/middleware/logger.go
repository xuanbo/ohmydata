package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// ZapLogger 日志中间件
func ZapLogger(logger *zap.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) (err error) {
			req := ctx.Request()
			resp := ctx.Response()
			start := time.Now()
			if err = next(ctx); err != nil {
				ctx.Error(err)
			}
			cost := time.Since(start)
			if err == nil || resp.Status == http.StatusNotFound {
				logger.Debug("API", zap.String("uri", req.RequestURI), zap.String("method", req.Method),
					zap.Int("status", resp.Status), zap.String("cost", cost.String()), zap.Any("username", ctx.Get("username")))
			} else {
				logger.Warn("API", zap.String("uri", req.RequestURI), zap.String("method", req.Method),
					zap.Int("status", resp.Status), zap.String("cost", cost.String()), zap.Any("username", ctx.Get("username")),
					zap.Error(err))
			}
			return nil
		}
	}
}
