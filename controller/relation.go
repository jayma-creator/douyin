package controller

import (
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
)

// 关注
func RelationAction(c *gin.Context) {
	err := service.RelationActionService(c)
	if err != nil {
		logrus.Error(err)
		return
	}
}

// 关注
func FollowList(c *gin.Context) {
	err := service.FollowListService(c)
	if err != nil {
		logrus.Error(err)
		return
	}
}

// 粉丝
func FollowerList(c *gin.Context) {
	err := service.FanListService(c)
	if err != nil {
		logrus.Error(err)
		return
	}
}
