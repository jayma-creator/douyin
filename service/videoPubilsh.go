package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

func PublishService(c *gin.Context) {
	user := User{}
	token := c.PostForm("token")
	//在user结构体里查找token=客户端传来的token
	dao.DB.Where("token = ?", token).Find(&user)
	count := 0
	dao.DB.Where("token = ?", token).Find(&user).Count(&count)
	if count == 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	//获取视频文件数据，前端传来data
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	//设定文件名
	finalName := fmt.Sprintf("%d_%s", user.Id, data.Filename)
	//设定路径public文件夹下
	saveFile := filepath.Join("./public/", finalName)
	//保存文件
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	ip := getIp()
	video := Video{
		//Author:         user, //Author是User结构体类型，该字段不会在数据库里创建，所以这里可以省略
		PlayUrl:        "http://" + ip + ":8080/static/" + finalName,
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

func getIp() string {
	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	defer conn.Close()
	ip := strings.Split(conn.LocalAddr().String(), ":")[0]
	return ip
}
