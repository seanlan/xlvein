package conf

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
)

func NewDB(dns string) *gorm.DB {
	var (
		err error
		db  *gorm.DB
	)
	if dns == "" {
		log.Fatalf("get mysql config error")
	}
	db, err = gorm.Open(mysql.Open(dns),
		&gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				TablePrefix: "", SingularTable: true,
			},
		})
	if err != nil {
		log.Fatalf("connect db error: %#v", err)
	}
	sqlDB, err := db.DB()
	if err != nil || sqlDB == nil {
		panic("db connect error")
	}
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(100)
	return db
}
