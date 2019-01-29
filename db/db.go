package db

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"github.com/tinklabs/golibs/cmd"
	"github.com/tinklabs/golibs/config"
)

var (
	DB *gorm.DB
)

func Init() {
	//in case no db service
	if config.TakeDbUrl() == "null" {
		return
	}
	db, err := gorm.Open("mysql", config.TakeDbUrl())
	if err != nil {
		panic(fmt.Sprintf("open db:%v", err))
	}

	DB = db

	if cmd.IsDebug() {
		DB.LogMode(true)
	}
}
