package log

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

// Init 初始化
func Init() error {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:    "time",
		LevelKey:   "level",
		NameKey:    "logger",
		CallerKey:  "caller",
		MessageKey: "msg",
		// warn级别以上不显示堆栈
		// StacktraceKey: "stacktrace",
		LineEnding:  zapcore.DefaultLineEnding,
		EncodeLevel: zapcore.CapitalColorLevelEncoder, // 大写编码器
		// 时间
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
		},
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 日志级别
	var atom zap.AtomicLevel
	level := os.Getenv("OH_MY_DATA_LOGGER_LEVEL")
	switch level {
	case "debug":
		atom = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		atom = zap.NewAtomicLevelAt(zap.InfoLevel)
	default:
		atom = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	config := zap.Config{
		Level:            atom,               // 日志级别
		Development:      true,               // 开发模式，堆栈跟踪
		Encoding:         "console",          // 输出格式 console 或 json
		EncoderConfig:    encoderConfig,      // 编码器配置
		OutputPaths:      []string{"stdout"}, // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
		ErrorOutputPaths: []string{"stderr"},
	}

	// 构建日志
	var err error
	logger, err = config.Build()
	if err != nil {
		return err
	}
	return nil
}

// Logger 日志实例
func Logger() *zap.Logger {
	return logger
}

// Flush 刷新缓存，程序退出前调用
func Flush() error {
	return logger.Sync()
}
