package controller

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"time"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

func Publish(c *gin.Context) {
	user := User{}
	token := c.PostForm("token")
	//在user结构体里查找token=客户端传来的token，count计数表示获取条数
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

func PublishList(c *gin.Context) {
	//封面问题还未解决
	token := c.Query("token")
	dao.DB.Where("publisher_token = ?", token).Find(&videoList)
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}
