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

type VideoListResponse struct {
	common.Response
	VideoList []common.Video `json:"video_list"`
}

const (
	like   = "1"
	unLike = "2"
)

//点赞与取消
func FavoriteActionService(c *gin.Context) (err error) {
	videoIdStr := c.Query("video_id")
	actionType := c.Query("action_type")
	videoId, _ := strconv.Atoi(videoIdStr)
	u, _ := c.Get("user")
	e, _ := c.Get("exist")
	user := u.(common.User)
	exist := e.(bool)
	if exist {
		if actionType == like {
			err = likeAct(c, user, videoId)
			if err != nil {
				return
			}
		} else if actionType == unLike {
			err = unlikeAct(c, user, videoId)
			if err != nil {
				return
			}
		} else {
			c.JSON(http.StatusOK, CommentActionResponse{Response: common.Response{StatusCode: 1, StatusMsg: "错误操作"}})
			return
		}
	} else {
		c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: "token已过期，请重新登录"})
		return
	}
	return
}

//点赞
func likeAct(c *gin.Context, user common.User, videoId int) (err error) {
	tx := dao.DB.Begin()
	ufr := common.UserFavoriteRelation{
		UserId:  user.Id,
		VideoId: int64(videoId),
	}
	err = tx.Create(&ufr).Error
	if err != nil {
		logrus.Error("新增视频信息失败", err)
		tx.Rollback()
		return
	}
	//把video结构体里的IsFavorite改为true
	//video的favorite_count+1
	err = tx.Model(&common.Video{}).Where("id = ?", videoId).Updates(map[string]interface{}{"is_favorite": true, "favorite_count": gorm.Expr("favorite_count + ?", "1")}).Error
	if err != nil {
		logrus.Error("修改视频信息失败", err)
		tx.Rollback()
		return
	}

	//删除redis缓存
	err = util.DelCache(fmt.Sprintf("favoriteList%v", user.Id))
	if err != nil {
		return
	}
	tx.Commit()
	//延时双删
	time.Sleep(time.Millisecond * 50)
	err = util.DelCache(fmt.Sprintf("favoriteList%v", user.Id))
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, common.Response{StatusCode: 0, StatusMsg: "like success"})
	return
}

//取消赞
func unlikeAct(c *gin.Context, user common.User, videoId int) (err error) {
	tx := dao.DB.Begin()
	err = tx.Where("user_id = ? and video_id = ?", user.Id, videoId).Delete(&common.UserFavoriteRelation{}).Error
	if err != nil {
		logrus.Error("删除视频信息失败", err)
		tx.Rollback()
		return
	}
	//把video结构体里的IsFavorite改为false
	//video的favorite_count-1
	err = tx.Model(&common.Video{}).Where("id = ?", videoId).Updates(map[string]interface{}{"is_favorite": false, "favorite_count": gorm.Expr("favorite_count - ?", "1")}).Error
	if err != nil {
		logrus.Error("修改视频信息失败", err)
		tx.Rollback()
		return
	}
	//删除redis缓存
	err = util.DelCache(fmt.Sprintf("favoriteList%v", user.Id))
	if err != nil {
		return
	}
	tx.Commit()
	//延时双删
	time.Sleep(time.Millisecond * 50)
	err = util.DelCache(fmt.Sprintf("favoriteList%v", user.Id))
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, common.Response{StatusCode: 0, StatusMsg: "unlike success"})
	return
}