package controller

import (
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//评论操作
func CommentAction(c *gin.Context) {
	err := service.CommentActionService(c)
	if err != nil {
		logrus.Error(err)
		return
	}
}

//展示评论
func CommentList(c *gin.Context) {
	err := service.CommentListService(c)
	if err != nil {
		logrus.Error(err)
		return
	}
}
