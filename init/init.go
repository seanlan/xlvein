package init

import (
	"github.com/seanlan/goutils/xlconfig"
	"github.com/seanlan/xlvein/app/dao"
	"github.com/seanlan/xlvein/conf"
	"math/rand"
	"time"
)

func init() {
	// 初始化随机种子
	rand.Seed(time.Now().Unix())
	// 初始化时区
	local, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	time.Local = local
	// 初始化配置
	xlconfig.Setup()
	conf.Setup()
	dao.DB = conf.NewDB(xlconfig.GetString("mysql"))
}
