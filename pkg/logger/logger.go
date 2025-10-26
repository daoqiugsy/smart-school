package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log 是一个全局的、可被项目其他地方直接使用的 zap Logger 实例
var Log *zap.Logger

// InitLogger 初始化 Logger
// mode 可以是 "dev" (开发模式) 或 "prod" (生产模式)
func InitLogger(mode string) error {
	var err error
	var logger *zap.Logger

	// 根据不同的模式，使用不同的配置
	if mode == "dev" {
		// 开发模式下，使用更易读的 ConsoleEncoder
		config := zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder // 带颜色的日志级别
		logger, err = config.Build()
	} else {
		// 生产模式下，使用高性能的 JSONEncoder
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return err
	}

	// 将我们创建的 logger 赋值给全局变量 Log
	Log = logger
	return nil
}
