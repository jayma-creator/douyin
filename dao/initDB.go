package dao

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/setting"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitMySQL(cfg *setting.MySQLConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		// gorm日志模式：Warn
		Logger: logger.Default.LogMode(logger.Warn),
		// 禁用默认事务（提高运行速度）
		SkipDefaultTransaction: true,
	})
	if err != nil {
		return
	}
	return
}
