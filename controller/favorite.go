package controller

import (
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//点赞
func FavoriteAction(c *gin.Context) {
	err := service.FavoriteActionService(c)
	if err != nil {
		logrus.Error(err)
		return
	}
}

//展示点赞过的视频
func FavoriteList(c *gin.Context) {
	err := service.FavoriteListService(c)
	if err != nil {
		logrus.Error(err)
		return
	}
}
