package service

import (
	"bytes"
	"fmt"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func PublishService(c *gin.Context) {
	user := User{}
	token := c.PostForm("token")
	title := c.PostForm("title")
	dao.DB.Where("token = ?", token).Find(&user)
	dao.DB.Where("token = ?", token).Find(&user).Count(&count)
	if count == 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	//获取视频文件数据
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
	snapShotName := finalName + "-cover.jpeg"
	ip := getIp()
	getSnapShot(snapShotName, saveFile)
	video := Video{
		PlayUrl:        "http://" + ip + ":8080/static/" + finalName,
		CoverUrl:       "http://" + ip + ":8080/static/" + snapShotName,
		FavoriteCount:  0,
		CommentCount:   0,
		IsFavorite:     false,
		Title:          title,
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

//获取当前主机IP
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

//截图做封面
func ExampleReadFrameAsJpeg(inFileName string, frameNum int) io.Reader {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		log.Fatal(err)
	}
	return buf
}

//保存截图
func getSnapShot(snapShotName string, videoFilePath string) error {
	reader := ExampleReadFrameAsJpeg(videoFilePath, 1)
	img, err := imaging.Decode(reader)
	if err != nil {
		return err
	}
	err = imaging.Save(img, "./public/"+snapShotName)
	if err != nil {
		return err
	}
	return nil
}
