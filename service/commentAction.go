package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	if u != nil && e != nil {
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
			}
		}
	}
	return
}

//评论列表
func CommentListService(c *gin.Context) (err error) {
	videoId := c.Query("video_id")
	key := fmt.Sprintf("commentList%v", videoId)
	var commentList []common.Comment
	var count int64
	//先查询缓存
	if util.IsExistCache(key) == 1 {
		commentList, err = util.GetCommentCache(videoId)
		if err != nil {
			logrus.Info("查询评论列表缓存失败", err)
		}
	} else if util.IsExistCache(key) == 0 {
		lockNum := "1"
		if util.RedisLock(lockNum) == true {
			commentList, count, err = dao.QueryCommentList(videoId)
			if err != nil {
				logrus.Error(err)
				return
			}
			if count == 0 {
				//如果数据库不存在，则缓存一个10秒的空值，防止缓存穿透
				go util.SetNull(key)
			} else {
				//缓存到redis
				go util.SetRedisCache(key, commentList)
			}
			util.RedisUnlock(lockNum)
		} else {
			time.Sleep(time.Millisecond * 100)
			commentList, err = util.GetCommentCache(videoId)
			if err != nil {
				logrus.Info("查询评论列表缓存失败", err)
				return
			}
		}
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
	err = dao.UpdateCommentAdd(tx, videoId)
	if err != nil {
		logrus.Error("修改评论数失败", err)
		tx.Rollback()
		return
	}

	//删除redis缓存
	err = util.DelCache("feed")
	err = util.DelCache(fmt.Sprintf("commentList%v", videoId))
	if err != nil {
		return
	}

	tx.Commit()

	//延时双删
	time.Sleep(time.Millisecond * 50)
	err = util.DelCache("feed")
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
	err = dao.DeleteComment(tx, commentId)
	if err != nil {
		logrus.Error("删除评论信息失败", err)
		tx.Rollback()
		return
	}
	//video的comment_count-1
	err = dao.UpdateCommentDel(tx, videoId)
	if err != nil {
		logrus.Error("修改评论信息失败", err)
		tx.Rollback()
		return
	}
	//删除redis缓存
	err = util.DelCache("feed")
	err = util.DelCache(fmt.Sprintf("commentList%v", videoId))
	if err != nil {
		return
	}

	tx.Commit()
	//延时双删
	time.Sleep(time.Millisecond * 50)
	err = util.DelCache("feed")
	err = util.DelCache(fmt.Sprintf("commentList%v", videoId))

	c.JSON(http.StatusOK, CommentActionResponse{Response: common.Response{StatusCode: 0, StatusMsg: "删除评论成功"}})
	return
}
