package gorm

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
)

// ZapLogger zap实现日志
type ZapLogger struct {
	SlowThreshold time.Duration
	level         logger.LogLevel
	logger        *zap.Logger
	prefix        string
}

// NewZapLogger 创建实例
func NewZapLogger(logger *zap.Logger, slowThreshold time.Duration, prefix string) *ZapLogger {
	return &ZapLogger{SlowThreshold: slowThreshold, logger: logger, prefix: prefix}
}

// LogMode 实现LogMode接口
func (zl *ZapLogger) LogMode(level logger.LogLevel) logger.Interface {
	zl.level = level
	return zl
}

// Info 实现Info接口
func (zl *ZapLogger) Info(ctx context.Context, msg string, v ...interface{}) {
	zl.logger.Debug(fmt.Sprintf(msg, v...))
}

// Warn 实现Warn接口
func (zl *ZapLogger) Warn(ctx context.Context, msg string, v ...interface{}) {
	zl.logger.Warn(fmt.Sprintf(msg, v...))
}

// Error 实现Error接口
func (zl *ZapLogger) Error(ctx context.Context, msg string, v ...interface{}) {
	zl.logger.Error(fmt.Sprintf(msg, v...))
}

// Trace 实现Trace接口
func (zl *ZapLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if zl.level > 0 {
		elapsed := time.Since(begin)
		msg := zl.prefix + "Execute SQL"
		switch {
		case err != nil && zl.level >= logger.Error:
			sql, rows := fc()
			zl.logger.Error(msg, zap.String("sql", sql), zap.String("cost", elapsed.String()), zap.Int64("rows", rows), zap.Error(err))
		case zl.level >= logger.Info:
			sql, rows := fc()
			if elapsed >= zl.SlowThreshold {
				zl.logger.Warn(msg, zap.String("sql", sql), zap.String("cost", elapsed.String()), zap.Int64("rows", rows))
			} else {
				zl.logger.Debug(msg, zap.String("sql", sql), zap.String("cost", elapsed.String()), zap.Int64("rows", rows))
			}
		}
	}
}
