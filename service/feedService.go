package service

import (
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

func FeedService(c *gin.Context) (err error) {
	token := c.Query("token")
	videoList := []Video{}
	//把数据库里所有视频放在videoList内,且按照创建时间降序排列
	err = dao.DB.Order("created_at desc").Preload("Author").Find(&videoList).Error
	if err != nil {
		logrus.Error("获取视频列表失败", err)
		return
	}
	//无用户登录
	if token == "" {
		//每次获取先把默认点赞标识改为false
		for i := 0; i < len(videoList); i++ {
			dao.DB.Model(&Video{}).Update("is_favorite", false)
			videoList[i].IsFavorite = false
		}
		//每次获取先把关注图标标识改为false
		users := []User{}
		dao.DB.Find(&users)
		for i := 0; i < len(users); i++ {
			users[i].IsFollow = false
			dao.DB.Model(&User{}).Update("is_follow", false)
		}
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  time.Now().Unix(),
	})
	return
}
