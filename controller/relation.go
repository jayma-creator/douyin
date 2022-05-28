package controller

import (
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

//关注操作
func RelationAction(c *gin.Context) {
	service.RelationActionService(c)

}

//关注列表
func FollowList(c *gin.Context) {
	service.FollowListService(c)
}

//粉丝列表
func FollowerList(c *gin.Context) {
	service.FollowerListService(c)
}
