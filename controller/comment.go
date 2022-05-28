package controller

import (
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

//评论操作
func CommentAction(c *gin.Context) {
	service.CommentActionService(c)

}

//展示评论
func CommentList(c *gin.Context) {
	service.CommentListService(c)
}
