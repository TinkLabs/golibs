package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"github.com/tinklabs/golibs/cmd"
	"github.com/tinklabs/golibs/config"
)

var (
	DB *gorm.DB
)

func Init() {
	db, err := gorm.Open("mysql", config.TakeDbUrl())
	if err != nil {
		panic(err)
	}

	DB = db

	if cmd.IsDebug() {
		DB.LogMode(true)
	}
}
