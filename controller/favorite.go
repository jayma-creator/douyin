package controller

import (
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func FavoriteAction(c *gin.Context) {
	user := User{}
	token := c.Query("token")
	videoIdStr := c.Query("video_id")
	actionType := c.Query("action_type")
	videoId, _ := strconv.Atoi(videoIdStr)
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
	} else if actionType == "2" {
		dao.DB.Where("user_id = ? and video_id = ?", user.Id, videoId).Delete(&UserFavoriteRelation{})
		//把video结构体里的IsFavorite改为false
		dao.DB.Model(&Video{}).Where("id = ?", videoId).Update("is_favorite", false)
	}
}

func FavoriteList(c *gin.Context) {
	userId := c.Query("user_id")
	//从点赞关系表中取出当前id的结构体
	dao.DB.Where("user_id = ?", userId).Find(&favoriteRelationVideoIdList)
	//从当前id的结构体中取出video_id字段，保存在切片中
	videoIdSlice := []int64{}
	for i := 0; i < len(videoIdSlice); i++ {
		videoIdSlice = append(videoIdSlice, favoriteRelationVideoIdList[i].VideoId)
	}
	//根据video_id找出对应的video结构体放在结构体切片中，并返回前端显示
	dao.DB.Where(videoIdSlice).Find(&videoList)
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}
