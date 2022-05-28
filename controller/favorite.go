package controller

import (
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

//点赞
func FavoriteAction(c *gin.Context) {
	service.FavoriteActionService(c)
}

//展示点赞过的视频
func FavoriteList(c *gin.Context) {
	service.FavoriteListService(c)
}
