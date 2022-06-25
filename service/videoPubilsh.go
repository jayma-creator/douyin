package service

import (
	"bytes"
	"fmt"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
)

var videoIdSequence = int64(0)

func PublishService(c *gin.Context) (err error) {
	token := c.PostForm("token")
	user, exist, err := CheckToken(token)
	if exist {
		//获取视频文件数据
		title := c.PostForm("title")
		data, err := c.FormFile("data")
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return err
		}
		//设定文件名
		finalName := fmt.Sprintf("%d_%s", user.Id, data.Filename)
		//设定路径public文件夹下
		saveFile := filepath.Join("./public/", finalName)
		//保存文件
		if err = c.SaveUploadedFile(data, saveFile); err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return err
		}
		snapShotName := finalName + "-cover.jpeg"
		ip := getIp()
		if err != nil {
			logrus.Error("获取ip失败", err)
			return err
		}
		getSnapShot(snapShotName, saveFile)
		atomic.AddInt64(&videoIdSequence, 1)
		video := Video{
			Id:            videoIdSequence,
			Author:        user,
			PlayUrl:       "http://" + ip + ":8080/static/" + finalName,
			CoverUrl:      "http://" + ip + ":8080/static/" + snapShotName,
			FavoriteCount: 0,
			CommentCount:  0,
			IsFavorite:    false,
			Title:         title,
		}
		err = dao.DB.Create(&video).Error
		if err != nil {
			logrus.Error("插入视频失败", err)
			return err
		}
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  finalName + " uploaded successfully",
		})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	return
}

//获取当前主机IP
func getIp() string {
	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		logrus.Error("获取ip失败", err)
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
		logrus.Error("获取封面失败", err)
		return nil
	}
	return buf
}

//保存截图
func getSnapShot(snapShotName string, videoFilePath string) {
	reader := ExampleReadFrameAsJpeg(videoFilePath, 1)
	img, err := imaging.Decode(reader)
	if err != nil {
		logrus.Error("保存截图失败", err)
		return
	}
	err = imaging.Save(img, "./public/"+snapShotName)
	if err != nil {
		logrus.Error("保存截图失败", err)
		return
	}
	return
}
