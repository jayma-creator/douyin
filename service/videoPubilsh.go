package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"
)

var videoIdSequence = int64(0)

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
		finalName := fmt.Sprintf("%d__%s", user.Id, data.Filename)
		//设定路径public文件夹下
		saveFile := filepath.Join("./", finalName)
		//保存文件
		if err = c.SaveUploadedFile(data, saveFile); err != nil {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			})
			return err
		}

		snapShotName := finalName + "-cover.jpeg"
		if err != nil {
			logrus.Error("获取ip失败", err)
			return err
		}
		util.GetSnapShot(snapShotName, saveFile)
		go util.Consumer()
		go util.Producer(snapShotName)
		go util.Producer(finalName)
		atomic.AddInt64(&videoIdSequence, 1)
		video := common.Video{
			Id:       videoIdSequence,
			Author:   user,
			PlayUrl:  "http://rd1qd4izf.hn-bkt.clouddn.com/" + finalName,
			CoverUrl: "http://rd1qd4izf.hn-bkt.clouddn.com/" + snapShotName,
			Title:    title,
		}

		//删除redis缓存
		err = util.DelCache("feed")
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
		err = util.DelCache("feed")
		err = util.DelCache(fmt.Sprintf("publishList%v", user.Id))
		if err != nil {
			return err
		}

		//大于10秒客户端会显示超时
		timeSleep(finalName)

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

func timeSleep(fileName string) {
	file, _ := os.Stat("./" + fileName)
	if file.Size() >= 41943040 { //文件大于40兆，睡眠9秒
		time.Sleep(time.Second * 9)
	} else if file.Size() < 41943040 && file.Size() >= 20971520 { //如果文件大于20兆，睡眠7秒
		time.Sleep(time.Second * 7)
	} else if file.Size() < 20971520 && file.Size() >= 10485760 { //如果文件大于10兆，睡眠5秒
		time.Sleep(time.Second * 5)
	} else if file.Size() < 10485760 && file.Size() >= 5242880 { //如果文件大于5兆，睡眠3秒
		time.Sleep(time.Second * 3)
	}
}
