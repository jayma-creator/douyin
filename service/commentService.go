package service

import (
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"time"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

//评论和删除评论
func CommentActionService(c *gin.Context) {
	user := User{}
	token := c.Query("token")
	actionType := c.Query("action_type")
	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.Atoi(videoIdStr)
	count := 0
	dao.DB.Where("token = ?", token).Find(&user).Count(&count)
	if count == 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	if actionType == "1" {
		text := c.Query("comment_text")
		//新增评论
		comment := Comment{
			Content:    text,
			CreateDate: time.Now().Format("2006-01-02 15:04:05")[5:10], //按格式输出日期，5:10表示月-日
			UserToken:  token,
			VideoId:    int64(videoId),
		}
		dao.DB.Create(&comment)
		//为了能评论后即时显示用户名，这里手动赋值comment的User
		comment.User = user
		//video的comment_count+1
		dao.DB.Model(&Video{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + ?", "1"))
		c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
			Comment: comment,
		})
	} else if actionType == "2" {
		commentId := c.Query("comment_id")
		dao.DB.Where("id = ?", commentId).Delete(&Comment{})
		//video的comment_count-1
		dao.DB.Model(&Video{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count - ?", "1"))
		c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0}})
	}
}

//评论列表
func CommentListService(c *gin.Context) {
	videoId := c.Query("video_id")
	//取出所有当前视频的评论
	commentList := []Comment{}
	dao.DB.Where("video_id = ?", videoId).Find(&commentList)
	//匹配评论作者
	for i := 0; i < len(commentList); i++ {
		user := User{}
		dao.DB.Where("token = ?", commentList[i].UserToken).Find(&user)
		commentList[i].User = user
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: commentList,
	})
}
