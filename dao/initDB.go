package dao

import (
	"github.com/jinzhu/gorm"
)

var DB *gorm.DB

//初始化数据库
func InitDB() (err error) {
	dsn := "root:333@(localhost:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	return nil
}
