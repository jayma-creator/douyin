package dao

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/setting"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitMySQL(cfg *setting.MySQLConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Address, cfg.Port, cfg.DB)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// gorm日志模式：Warn
		Logger: logger.Default.LogMode(logger.Warn),
		// 禁用默认事务（提高运行速度）
		SkipDefaultTransaction: true,
	})
	if err != nil {
		logrus.Error("连接mysql失败", err)
		return
	}

	err = createTable()
	if err != nil {
		logrus.Error("建表失败", err)
		return
	}
	return

}

func createTable() (err error) {
	//创建表
	err = DB.AutoMigrate(&common.User{})
	if err != nil {
		logrus.Error(err)
		return
	}
	err = DB.AutoMigrate(&common.Comment{})
	if err != nil {
		logrus.Error(err)
		return
	}
	err = DB.AutoMigrate(&common.Video{})
	if err != nil {
		logrus.Error(err)
		return
	}
	err = DB.AutoMigrate(&common.FollowFansRelation{})
	if err != nil {
		logrus.Error(err)
		return
	}
	err = DB.AutoMigrate(&common.UserFavoriteRelation{})
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}
