package service

import (
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

const (
	like   = "1"
	unLike = "2"
)

//点赞与取消
func FavoriteActionService(c *gin.Context) (err error) {
	user := User{}
	token := c.Query("token")
	videoIdStr := c.Query("video_id")
	actionType := c.Query("action_type")
	videoId, _ := strconv.Atoi(videoIdStr)
	err = dao.DB.Where("token = ?", token).Find(&user).Count(&count).Error
	if err != nil {
		logrus.Error("查询token失败", err)
		return
	}
	if count == 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	} else if count == 1 {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	}
	if actionType == like {
		tx := dao.DB.Begin()
		ufr := UserFavoriteRelation{
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
		err = tx.Model(&Video{}).Where("id = ?", videoId).Updates(map[string]interface{}{"is_favorite": true, "favorite_count": gorm.Expr("favorite_count + ?", "1")}).Error
		if err != nil {
			logrus.Error("修改视频信息失败", err)
			tx.Rollback()
			return
		}
		tx.Commit()

	} else if actionType == unLike {
		tx := dao.DB.Begin()
		err = tx.Where("user_id = ? and video_id = ?", user.Id, videoId).Delete(&UserFavoriteRelation{}).Error
		if err != nil {
			logrus.Error("删除视频信息失败", err)
			tx.Rollback()
			return
		}
		//把video结构体里的IsFavorite改为false
		//video的favorite_count-1
		err = tx.Model(&Video{}).Where("id = ?", videoId).Updates(map[string]interface{}{"is_favorite": false, "favorite_count": gorm.Expr("favorite_count - ?", "1")}).Error
		if err != nil {
			logrus.Error("修改视频信息失败", err)
			tx.Rollback()
			return
		}
		tx.Commit()
	}
	return
}

//点赞列表
func FavoriteListService(c *gin.Context) (err error) {
	userId := c.Query("user_id")
	//查询出当前登录的用户点赞过的视频列表
	videoList := []Video{}
	userFavoriteRelation := []UserFavoriteRelation{}
	err = dao.DB.Where("user_id = ?", userId).Preload("Video").Preload("User").Preload("Video.Author").Find(&userFavoriteRelation).Error
	if err != nil {
		logrus.Error("获取点赞列表失败", err)
		return
	}
	for i := 0; i < len(userFavoriteRelation); i++ {
		videoList = append(videoList, userFavoriteRelation[i].Video)
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
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
	videoList := []Video{}
	err = dao.DB.Where("author_id = ?", userId).Preload("Author").Find(&videoList).Error
	if err != nil {
		logrus.Error("获取发布列表失败", err)
		return
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
	return
}
