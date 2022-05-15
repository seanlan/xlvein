package conf

import (
	"github.com/seanlan/goutils/xlconfig"
	"github.com/seanlan/xlvein/app/dao"
	"github.com/seanlan/xlvein/pkg/xlredis"
)

func initRedis() {
	hosts := xlconfig.GetString("redis", "hosts")
	username := xlconfig.GetString("redis", "username")
	password := xlconfig.GetString("redis", "password")
	prefix := xlconfig.GetString("redis", "prefix")
	db := xlconfig.GetInt("redis", "db")
	r, err := xlredis.NewClient(hosts, username, password, prefix, int(db))
	if err != nil {
		panic(err)
	}
	dao.Redis = r
}
