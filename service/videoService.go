package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
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
	//取出当前用户的所有发布列表
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

func PublishService(c *gin.Context) {
	user := User{}
	token := c.PostForm("token")
	//在user结构体里查找token=客户端传来的token，count计数表示获取条数
	count := 0
	dao.DB.Where("token = ?", token).Find(&user).Count(&count)
	if count == 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	//文件名
	//filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", user.Id, data.Filename)
	//保存在public文件夹下
	saveFile := filepath.Join("./public/", finalName)
	fmt.Println(saveFile)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	video := Video{
		//Author:         user, //Author是User结构体类型，该字段不会在数据库里创建，所以这里可以省略
		PlayUrl:        "http://192.168.220.1:8080/static/" + finalName,
		CoverUrl:       "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
		FavoriteCount:  0,
		CommentCount:   0,
		IsFavorite:     false,
		PublisherToken: token,
		CreatedAt:      time.Time{},
		UpdatedAt:      time.Time{},
		DeletedAt:      nil,
	}
	dao.DB.Create(&video)
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

func PublishListService(c *gin.Context) {
	//封面问题还未解决
	userId := c.Query("user_id")
	//取出当前用户的所有发布列表
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
