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

func FeedService(c *gin.Context) {
	//判断有没有用户登录
	token := c.Query("token")
	user := User{}
	videoList := []Video{}
	count := 0
	dao.DB.Where("token = ?", token).Find(&user).Count(&count)
	//把数据库里所有视频放在videoList内,且按照创建时间降序排列
	dao.DB.Order("created_at desc").Find(&videoList)
	//每次获取先把默认点赞标识改为false
	//每次都把是否已点赞默认改为false
	for i := 0; i < len(videoList); i++ {
		dao.DB.Model(&Video{}).Update("is_favorite", false)
	}

	//匹配视频与作者
	//发布的视频有publish_token，和当前用户的token对应起来
	//根据publish_token取出对应的user结构体，赋值给video的Author
	videoTokenSlice := []string{}
	for i := 0; i < len(videoList); i++ {
		videoTokenSlice = append(videoTokenSlice, videoList[i].PublisherToken)
	}
	//循环取出User结构体,赋给相对应的video，通俗点说就是发布者匹配
	//才能点进头像后正常显示用户信息
	for i := 0; i < len(videoTokenSlice); i++ {
		user := User{}
		dao.DB.Where("token = ?", videoTokenSlice[i]).Find(&user)
		videoList[i].Author = user
	}

	//无用户登录的逻辑，直接展示
	if count == 0 {
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0},
			VideoList: videoList,
			NextTime:  time.Now().Unix(),
		})
	} else if count == 1 { //如果有用户登录，根据用户点赞关系表正确显示当前用户是否已点赞该视频
		//拿出当前用户的Id
		userId := user.Id
		userFavoriteRelations := []UserFavoriteRelation{} //存放当前用户点过赞的视频关系表的结构体
		//在用户点赞关系表里找出当前用户点赞了的视频的ID
		dao.DB.Where("user_id = ?", userId).Find(&userFavoriteRelations)
		//从上面的结构体取出点过赞的视频ID
		favoriteVideoIdSlice := []int64{}
		for i := 0; i < len(userFavoriteRelations); i++ {
			favoriteVideoIdSlice = append(favoriteVideoIdSlice, userFavoriteRelations[i].VideoId)
		}
		//根据点过赞的视频ID找到对应的视频，把该视频的is_favorite改为true，默认为false
		for i := 0; i < len(favoriteVideoIdSlice); i++ {
			dao.DB.Model(&Video{}).Where("id = ?", favoriteVideoIdSlice[i]).Update("is_favorite", true)
		}
		c.JSON(http.StatusOK, FeedResponse{
			Response:  Response{StatusCode: 0},
			VideoList: videoList,
			NextTime:  time.Now().Unix(),
		})

	}

}
