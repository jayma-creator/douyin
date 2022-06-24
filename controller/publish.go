package controller

import (
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//发布视频
func Publish(c *gin.Context) {
	err := service.PublishService(c)
	if err != nil {
		logrus.Error(err)
		return
	}
}

//展示发布视频的列表
func PublishList(c *gin.Context) {
	err := service.PublishListService(c)
	if err != nil {
		logrus.Error(err)
		return
	}
}
