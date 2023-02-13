package conf

import "go.uber.org/zap"

// 初始化日志
func initLogging(debug bool, appName string) {
	var _conf zap.Config
	if debug {
		_conf = zap.NewDevelopmentConfig()
	} else {
		_conf = zap.NewProductionConfig()
	}
	logger, _ := _conf.Build()
	logger = logger.WithOptions(zap.AddCaller())
	logger = logger.With(zap.String("app", appName))
	zap.ReplaceGlobals(logger)
}
