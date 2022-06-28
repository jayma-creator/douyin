package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

//点赞列表
func FavoriteListService(c *gin.Context) (err error) {
	userId := c.Query("user_id")
	//查询出当前登录的用户点赞过的视频列表
	videoList, err := util.GetFavoriteListCache(userId)
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
		go util.SetRedisCache(fmt.Sprintf("favoriteList%v", userId), videoList)

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
	videoList, err := util.GetPublishListCache(userId)
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
		go util.SetRedisCache(fmt.Sprintf("publishList%v", userId), videoList)

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
