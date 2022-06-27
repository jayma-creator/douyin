package service

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type FeedResponse struct {
	common.Response
	VideoList []common.Video `json:"video_list,omitempty"`
	NextTime  int64          `json:"next_time,omitempty"`
}

func FeedService(c *gin.Context) (err error) {
	token := c.Query("token")
	videoList := []common.Video{}
	//把数据库里所有视频放在videoList内,且按照创建时间降序排列
	//无用户登录
	if token == "" {
		tx := dao.DB.Begin()
		//每次获取先把点赞图标和用户关注改为false
		err = tx.Model(common.Video{}).Where("is_favorite = ?", true).Update("is_favorite", false).Error
		if err != nil {
			logrus.Error("修改失败", err)
			tx.Rollback()
			return
		}
		err = tx.Model(common.User{}).Where("is_follow = ?", true).Update("is_follow", false).Error
		if err != nil {
			logrus.Error("修改失败", err)
			tx.Rollback()
			return
		}
		tx.Commit()
	} else {
		u, _ := c.Get("user")
		e, _ := c.Get("exist")
		user := u.(common.User)
		exist := e.(bool)
		if exist {
			err = checkUserSetting(user)
			if err != nil {
				return err
			}
		}
	}
	err = dao.DB.Order("created_at desc").Preload("Author").Find(&videoList).Error
	if err != nil {
		logrus.Error("获取视频列表失败", err)
		return
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  common.Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  time.Now().Unix(),
	})
	return
}

//匹配当前登录的账号是否已关注别的账号，是否点赞视频
func checkUserSetting(user common.User) (err error) {
	tx := dao.DB.Begin()
	//匹配当前登录的账号是否已关注别的账号
	users := []common.User{}
	err = tx.Table("users").
		Joins("join follow_fans_relations on follower_id = users.id and follow_id = ? and follow_fans_relations.deleted_at is null", user.Id).
		Find(&users).Error
	if err != nil {
		logrus.Error("修改失败", err)
		tx.Rollback()
		return err
	}
	for i := 0; i < len(users); i++ {
		err = tx.Model(&common.User{}).Where("id = ?", users[i].Id).Update("is_follow", true).Error
		if err != nil {
			logrus.Error("修改失败", err)
			tx.Rollback()
			return err
		}
	}

	//匹配当前登录的账号是否已点赞视频
	videos := []common.Video{}
	err = tx.Table("videos").
		Joins("join user_favorite_relations on video_id = videos.id and user_id = ? and user_favorite_relations.deleted_at is null", user.Id).
		Find(&videos).Error
	if err != nil {
		logrus.Error("修改失败", err)
		tx.Rollback()
		return err
	}
	for i := 0; i < len(videos); i++ {
		//videos[i].IsFavorite = true
		err = tx.Model(&common.Video{}).Where("id = ?", videos[i].Id).Update("is_favorite", true).Error
		if err != nil {
			logrus.Error("修改失败", err)
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return err
}
