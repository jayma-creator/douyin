package service

import (
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

//初始化视频流
//拉取视频之前，未登录状态要取消点赞图标，登录状态要正确显示点赞图标
//未登录状态点进头像显示未关注，登录后点进头像要正确的显示是否已关注
func initVideo(videoList []Video) {

}

func FeedService(c *gin.Context) {
	//initVideo(videoList )
	//判断有没有用户登录
	token := c.Query("token")
	videoList := []Video{}
	//把数据库里所有视频放在videoList内,且按照创建时间降序排列
	dao.DB.Order("created_at desc").Find(&videoList)

	//匹配视频与作者
	for i := 0; i < len(videoList); i++ {
		user := User{}
		dao.DB.Where("token = ?", videoList[i].PublisherToken).Find(&user)
		videoList[i].Author = user
	}

	//如果无用户登录，则把点赞图标取消
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
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0},
			VideoList: videoList,
			NextTime:  time.Now().Unix(),
		})

	} else {

		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0},
			VideoList: videoList,
			NextTime:  time.Now().Unix(),
		})

	}

}
