package conf

import "github.com/seanlan/goutils/xlconfig"

var (
	DebugMode = false
	AppName   = ""
)

func Setup() {
	DebugMode = xlconfig.GetBool("debug")
	AppName = xlconfig.GetString("app_name")
	initLogging(DebugMode, AppName)
}
