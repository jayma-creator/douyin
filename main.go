package main

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/router"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/RaymondCode/simple-demo/setting"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置文件
	if err := setting.Init("conf/config.ini"); err != nil {
		fmt.Printf("load config from file failed, err:%v\n", err)
		return
	}
	// 连接数据库
	err := dao.InitMySQL(setting.Conf.MySQLConfig)
	if err != nil {
		fmt.Printf("init mysql failed, err:%v\n", err)
		return
	}
	r := gin.Default()
	router.InitRouter(r)

	dao.DB.AutoMigrate(&service.User{})
	dao.DB.AutoMigrate(&service.Comment{})
	dao.DB.AutoMigrate(&service.Video{})
	dao.DB.AutoMigrate(&service.FollowFansRelation{})
	dao.DB.AutoMigrate(&service.UserFavoriteRelation{})

	service.InitDemo() //初始化测试数据

	defer dao.DB.Close()
	if err := r.Run(fmt.Sprintf(":%d", setting.Conf.Port)); err != nil {
		fmt.Printf("server startup failed, err:%v\n", err)
	}
}
