package config

import (
	"go.uber.org/zap"
)

// 初始化日志
func initLogging() {
	var _conf zap.Config
	if C.Debug {
		_conf = zap.NewDevelopmentConfig()
	} else {
		_conf = zap.NewProductionConfig()
	}
	logger, _ := _conf.Build()
	logger.WithOptions(zap.AddCaller())
	logger = logger.With(zap.String("app", C.App))
	zap.ReplaceGlobals(logger)
}
