package main

import (
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

var U1 service.User
var U2 service.User

func main() {
	r := gin.Default()
	initRouter(r)
	dao.InitDB()

	dao.DB.AutoMigrate(&service.User{})
	dao.DB.AutoMigrate(&service.Comment{})
	dao.DB.AutoMigrate(&service.Video{})
	dao.DB.AutoMigrate(&service.FollowFansRelation{})
	dao.DB.AutoMigrate(&service.UserFavoriteRelation{})

	service.InitDemo() //初始化测试数据

	defer dao.DB.Close()
	r.Run()

}
