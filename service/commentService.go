package service

import (
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

const (
	createComment = "1"
	delComment    = "2"
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
func CommentActionService(c *gin.Context) (err error) {
	token := c.Query("token")
	user, exist, err := CheckToken(token)
	if exist {
		actionType := c.Query("action_type")
		videoIdStr := c.Query("video_id")
		videoId, _ := strconv.Atoi(videoIdStr)
		userId := user.Id
		if actionType == createComment {
			tx := dao.DB.Begin()
			text := c.Query("comment_text")
			//新增评论
			comment := Comment{
				Content:    text,
				UserId:     userId,
				User:       user,
				CreateDate: time.Now().Format("2006-01-02 15:04:05")[5:10], //按格式输出日期，5:10表示月-日
				VideoId:    int64(videoId),
			}
			err = tx.Create(&comment).Error
			if err != nil {
				logrus.Error("插入评论信息失败", err)
				tx.Rollback()
				return
			}
			//video的comment_count+1
			err = tx.Model(&Video{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + ?", "1")).Error
			if err != nil {
				logrus.Error("修改评论数失败", err)
				tx.Rollback()
				return
			}
			tx.Commit()
			c.JSON(http.StatusOK, CommentActionResponse{
				Response: Response{StatusCode: 0, StatusMsg: "评论成功"},
				Comment:  comment,
			})
		} else if actionType == delComment {
			tx := dao.DB.Begin()
			commentId := c.Query("comment_id")
			err = tx.Where("id = ?", commentId).Delete(&Comment{}).Error
			if err != nil {
				logrus.Error("删除评论信息失败", err)
				tx.Rollback()
				return
			}
			//video的comment_count-1
			err = tx.Model(&Video{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count - ?", "1")).Error
			if err != nil {
				logrus.Error("修改评论信息失败", err)
				tx.Rollback()
				return
			}
			tx.Commit()
			c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0, StatusMsg: "删除评论成功"}})
		} else {
			c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 1, StatusMsg: "错误操作"}})
			return
		}
	} else {
		c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 1, StatusMsg: "当前用户不存在"}})
		return
	}
	return
}

//评论列表
func CommentListService(c *gin.Context) (err error) {
	videoId := c.Query("video_id")
	//取出所有当前视频的评论
	commentList := []Comment{}
	err = dao.DB.Where("video_id = ?", videoId).Preload("User").Find(&commentList).Error
	if err != nil {
		logrus.Error("获取评论列表失败", err)
		return
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: commentList,
	})
	return
}
