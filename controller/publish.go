package controller

import (
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

//发布视频
func Publish(c *gin.Context) {
	service.PublishService(c)
}

//展示发布视频的列表
func PublishList(c *gin.Context) {

	service.PublishListService(c)
}
