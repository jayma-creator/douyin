package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"path/filepath"
	"sync/atomic"
	"time"
)

var videoIdSequence = int64(2)

func PublishService(c *gin.Context) (err error) {
	u, _ := c.Get("user")
	e, _ := c.Get("exist")
	user := u.(common.User)
	exist := e.(bool)
	if exist {
		//获取视频文件数据
		title := c.PostForm("title")
		data, err := c.FormFile("data")
		if err != nil {
			c.JSON(http.StatusOK, common.Response{
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
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return err
		}
		snapShotName := finalName + "-cover.jpeg"
		ip := util.GetIp()
		if err != nil {
			logrus.Error("获取ip失败", err)
			return err
		}
		util.GetSnapShot(snapShotName, saveFile)
		atomic.AddInt64(&videoIdSequence, 1)
		video := common.Video{
			Id:            videoIdSequence,
			Author:        user,
			PlayUrl:       "http://" + ip + ":8080/static/" + finalName,
			CoverUrl:      "http://" + ip + ":8080/static/" + snapShotName,
			FavoriteCount: 0,
			CommentCount:  0,
			IsFavorite:    false,
			Title:         title,
		}

		//删除redis缓存
		err = util.DelCache(fmt.Sprintf("publishList%v", user.Id))
		if err != nil {
			return err
		}
		err = dao.DB.Create(&video).Error
		if err != nil {
			logrus.Error("插入视频失败", err)
			return err
		}
		//延时双删
		time.Sleep(time.Millisecond * 50)
		err = util.DelCache(fmt.Sprintf("publishList%v", user.Id))
		if err != nil {
			return err
		}

		c.JSON(http.StatusOK, common.Response{
			StatusCode: 0,
			StatusMsg:  finalName + " uploaded successfully",
		})
	} else {
		c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: "token已过期，请重新登录"})
		return
	}
	return
}
