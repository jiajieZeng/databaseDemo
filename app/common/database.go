package common

import (
	"databaseDemo/app/model"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	host := "127.0.0.1"
	port := "3306"
	database := "webtest"
	username := "root"
	password := "hjy021012"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=true",
		username,
		password,
		host,
		port,
		database,
		charset)

	db, err := gorm.Open("mysql", args)
	if err != nil {
		panic("failed to connect database, err:" + err.Error())
	}

	//迁移
	db.AutoMigrate(&model.User{})

	DB = db

	return db

}

func GetDB() *gorm.DB {
	return DB
}
