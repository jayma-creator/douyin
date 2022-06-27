package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/dao"
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

//点赞列表
func FavoriteListService(c *gin.Context) (err error) {
	userId := c.Query("user_id")
	//查询出当前登录的用户点赞过的视频列表
	videoList, err := getFavoriteListCache(userId)
	if err != nil {
		logrus.Info("查询点赞列表缓存失败", err)
	}
	//没有缓存，从数据库取，并缓存到redis
	if len(videoList) == 0 {
		err = dao.DB.Table("videos").
			Joins("join user_favorite_relations on video_id = videos.id and user_id = ? and user_favorite_relations.deleted_at is null", userId).Preload("Author").
			Find(&videoList).Error
		if err != nil {
			logrus.Error("获取点赞列表失败", err)
			return
		}
		//缓存到redis
		err = setRedisCache(fmt.Sprintf("favoriteList%v", userId), videoList)
		if err != nil {
			logrus.Error("缓存失败")
		}
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: common.Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
	return
}

//发布列表
func PublishListService(c *gin.Context) (err error) {
	userId := c.Query("user_id")
	//查询出当前用户发布过的视频列表
	videoList, err := getPublishListCache(userId)
	if err != nil {
		logrus.Info("查询发布列表缓存失败", err)
	}
	fmt.Println("从redis读取缓存")
	//没有缓存，从数据库取，并缓存到redis
	if len(videoList) == 0 {
		err = dao.DB.Where("author_id = ?", userId).Preload("Author").Find(&videoList).Error
		if err != nil {
			logrus.Error("获取发布列表失败", err)
			return
		}
		fmt.Println("mysql")
		err = setRedisCache(fmt.Sprintf("publishList%v", userId), videoList)
		if err != nil {
			logrus.Error("缓存失败")
		}
		fmt.Println("缓存到redis")
	}

	c.JSON(http.StatusOK, VideoListResponse{
		Response: common.Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
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
	err = delCache(fmt.Sprintf("favoriteList%v", user.Id))
	if err != nil {
		return
	}
	tx.Commit()
	//延时双删
	time.Sleep(time.Millisecond * 50)
	err = delCache(fmt.Sprintf("favoriteList%v", user.Id))
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
	err = delCache(fmt.Sprintf("favoriteList%v", user.Id))
	if err != nil {
		return
	}
	tx.Commit()
	//延时双删
	time.Sleep(time.Millisecond * 50)
	err = delCache(fmt.Sprintf("favoriteList%v", user.Id))
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, common.Response{StatusCode: 0, StatusMsg: "unlike success"})
	return
}
