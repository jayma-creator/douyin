package controller

import (
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//未登录状态下显示倒序显示视频，登录状态会匹配每个视频和作者
func Feed(c *gin.Context) {
	err := service.FeedService(c)
	if err != nil {
		logrus.Error(err)
		return
	}
}
