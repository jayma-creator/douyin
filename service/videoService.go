package service

import (
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

func FavoriteActionService(c *gin.Context) {
	user := User{}
	token := c.Query("token")
	videoIdStr := c.Query("video_id")
	actionType := c.Query("action_type")
	videoId, _ := strconv.Atoi(videoIdStr)
	count := 0
	dao.DB.Where("token = ?", token).Find(&user).Count(&count)
	if count == 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	} else if count == 1 {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	}
	if actionType == "1" {
		fr := UserFavoriteRelation{
			UserId:  user.Id,
			VideoId: int64(videoId),
		}
		dao.DB.Create(&fr)
		//把video结构体里的IsFavorite改为true
		dao.DB.Model(&Video{}).Where("id = ?", videoId).Update("is_favorite", true)
		//video的favorite_count+1
		dao.DB.Model(&Video{}).Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count + ?", "1"))

	} else if actionType == "2" {
		dao.DB.Where("user_id = ? and video_id = ?", user.Id, videoId).Delete(&UserFavoriteRelation{})
		//把video结构体里的IsFavorite改为false
		dao.DB.Model(&Video{}).Where("id = ?", videoId).Update("is_favorite", false)
		//video的favorite_count-1
		dao.DB.Model(&Video{}).Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count - ?", "1"))
	}
}

func FavoriteListService(c *gin.Context) {
	userId := c.Query("user_id")
	//查询出当前登录的用户点赞过的视频列表
	videoList := []Video{}
	dao.DB.Table("videos").
		Joins("join user_favorite_relations on video_id = videos.id and user_id = ? and videos.deleted_at is null", userId).
		Scan(&videoList)
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}

func PublishListService(c *gin.Context) {
	userId := c.Query("user_id")
	//查询出当前用户发布过的视频列表
	videoList := []Video{}
	dao.DB.Table("videos").
		Joins("join users on publisher_token = token and users.id = ? and videos.deleted_at is null", userId).
		Scan(&videoList)

	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}
