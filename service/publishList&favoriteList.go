package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

//点赞列表
func FavoriteListService(c *gin.Context) (err error) {
	userId := c.Query("user_id")
	key := fmt.Sprintf("favoriteList%v", userId)
	var videoList []common.Video
	//先查询缓存
	if util.IsExistCache(key) == 1 {
		videoList, err = util.GetFavoriteListCache(userId)
		if err != nil {
			logrus.Info("查询点赞列表缓存失败", err)
		}
	} else if util.IsExistCache(key) == 0 { //没有缓存，从数据库取，并缓存到redis
		var count int64
		lockNum := "1"
		if util.RedisLock(lockNum) == true {
			err = dao.DB.Table("videos").
				Joins("join user_favorite_relations on video_id = videos.id and user_id = ? and user_favorite_relations.deleted_at is null", userId).Preload("Author").Find(&videoList).Count(&count).Error
			if err != nil {
				logrus.Error("获取点赞列表失败", err)
				return
			}
			if count == 0 {
				go util.SetNull(key)
			} else {
				//缓存到redis
				go util.SetRedisCache(key, videoList)
			}
			util.RedisUnlock(lockNum)
		} else {
			time.Sleep(time.Millisecond * 100)
			videoList, err = util.GetFavoriteListCache(userId)
			if err != nil {
				logrus.Info("查询点赞列表缓存失败", err)
			}
		}

		//如果数据库不存在，则缓存一个10秒的空值，防止缓存穿透

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
	key := fmt.Sprintf("publishList%v", userId)
	var videoList []common.Video
	//先查询缓存
	if util.IsExistCache(key) == 1 {
		videoList, err = util.GetPublishListCache(userId)
		if err != nil {
			logrus.Info("查询发布列表缓存失败", err)
		}
	} else if util.IsExistCache(key) == 0 { //没有缓存，从数据库取，并缓存到redis
		var count int64
		lockNum := "1"
		if util.RedisLock(lockNum) == true {
			err = dao.DB.Where("author_id = ?", userId).Preload("Author").Find(&videoList).Count(&count).Error
			if err != nil {
				logrus.Error("获取发布列表失败", err)
				return
			}
			//如果数据库不存在，则缓存一个10秒的空值，防止缓存穿透
			if count == 0 {
				go util.SetNull(key)
			} else {
				//缓存到redis
				go util.SetRedisCache(key, videoList)
			}
			util.RedisUnlock(lockNum)
		} else {
			time.Sleep(time.Millisecond * 100)
			videoList, err = util.GetPublishListCache(userId)
			if err != nil {
				logrus.Info("查询发布列表缓存失败", err)
			}
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
