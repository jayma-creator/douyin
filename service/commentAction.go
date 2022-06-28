package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/util"
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
	common.Response
	CommentList []common.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	common.Response
	Comment common.Comment `json:"comment,omitempty"`
}

//评论和删除评论
func CommentActionService(c *gin.Context) (err error) {
	u, _ := c.Get("user")
	e, _ := c.Get("exist")
	user := u.(common.User)
	exist := e.(bool)

	if exist {
		actionType := c.Query("action_type")
		videoIdStr := c.Query("video_id")
		videoId, _ := strconv.Atoi(videoIdStr)
		if actionType == createComment {
			err = comment(c, user, videoId)
			if err != nil {
				return
			}
		} else if actionType == delComment {
			err = deleteComment(c, videoId)
			if err != nil {
				return
			}
		} else {
			c.JSON(http.StatusOK, CommentActionResponse{Response: common.Response{StatusCode: 1, StatusMsg: "错误操作"}})
			return
		}
	} else {
		c.JSON(http.StatusOK, CommentActionResponse{Response: common.Response{StatusCode: 1, StatusMsg: "token已过期，请重新登录"}})
		return
	}
	return
}

//评论列表
func CommentListService(c *gin.Context) (err error) {
	videoId := c.Query("video_id")
	commentList, err := util.GetCommentCache(videoId)
	if err != nil {
		logrus.Info("查询评论列表缓存失败", err)
	}
	//说明redis没有缓存，改为从数据库读取,并缓存到redis
	if len(commentList) == 0 {
		err = dao.DB.Where("video_id = ?", videoId).Preload("User").Preload("Video").Preload("Video.Author").Order("created_at desc").Find(&commentList).Error
		if err != nil {
			return
		}
		//缓存到redis
		go util.SetRedisCache(fmt.Sprintf("commentList%v", videoId), commentList)
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    common.Response{StatusCode: 0},
		CommentList: commentList,
	})
	return
}

//新增评论
func comment(c *gin.Context, user common.User, videoId int) (err error) {
	userId := user.Id
	tx := dao.DB.Begin()
	text := c.Query("comment_text")
	//新增评论
	comment := common.Comment{
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
	err = tx.Model(&common.Video{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + ?", "1")).Error
	if err != nil {
		logrus.Error("修改评论数失败", err)
		tx.Rollback()
		return
	}

	//删除redis缓存
	err = util.DelCache(fmt.Sprintf("commentList%v", videoId))
	if err != nil {
		return
	}

	tx.Commit()

	//延时双删
	time.Sleep(time.Millisecond * 50)
	err = util.DelCache(fmt.Sprintf("commentList%v", videoId))

	c.JSON(http.StatusOK, CommentActionResponse{
		Response: common.Response{StatusCode: 0, StatusMsg: "评论成功"},
		Comment:  comment,
	})
	return
}

//删除评论
func deleteComment(c *gin.Context, videoId int) (err error) {
	tx := dao.DB.Begin()
	commentId := c.Query("comment_id")
	err = tx.Where("id = ?", commentId).Delete(&common.Comment{}).Error
	if err != nil {
		logrus.Error("删除评论信息失败", err)
		tx.Rollback()
		return
	}
	//video的comment_count-1
	err = tx.Model(&common.Video{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count - ?", "1")).Error
	if err != nil {
		logrus.Error("修改评论信息失败", err)
		tx.Rollback()
		return
	}
	//删除redis缓存
	err = util.DelCache(fmt.Sprintf("commentList%v", videoId))
	if err != nil {
		return
	}

	tx.Commit()
	//延时双删
	time.Sleep(time.Millisecond * 50)
	err = util.DelCache(fmt.Sprintf("commentList%v", videoId))

	c.JSON(http.StatusOK, CommentActionResponse{Response: common.Response{StatusCode: 0, StatusMsg: "删除评论成功"}})
	return
}
