package conf

import "github.com/seanlan/goutils/xlconfig"

var (
	DebugMode = false
	AppName   = ""
)

func Setup() {
	xlconfig.Setup()
	DebugMode = xlconfig.GetBool("debug")
	AppName = xlconfig.GetString("app")
	initLogging(DebugMode, AppName)
	initRedis()
	initDB()
}
