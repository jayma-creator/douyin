package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/util"
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
	videoList := []common.Video{}
	key := "feed"
	u, _ := c.Get("user")
	e, _ := c.Get("exist")
	//把数据库里所有视频放在videoList内,且按照创建时间降序排列
	//无用户登录
	if u == nil && e == nil {
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
		videoList, _ = feedList(key)
	} else {
		user := u.(common.User)
		exist := e.(bool)
		if exist {
			err = checkUserSetting(user)
			if err != nil {
				return err
			}
			videoList, _ = feedList(key)

		}
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

func feedList(key string) (videoList []common.Video, err error) {
	//先从redis查询
	if util.IsExistCache(key) == 1 {
		videoList, err = util.GetFeed()
		if err != nil {
			logrus.Info("查询feed列表缓存失败", err)
		}
	} else if util.IsExistCache(key) == 0 {
		//倒序播放30条视频
		var count int64
		lockNum := "1"
		if util.RedisLock(lockNum) == true {
			err = dao.DB.Where("created_at <= ?", time.Now().Format("2006-01-02 15:04:05")).Order("created_at desc").Preload("Author").Find(&videoList).Limit(30).Count(&count).Error
			if err != nil {
				logrus.Error("获取视频列表失败", err)
				return
			}
			if count == 0 {
				//如果数据库不存在，则缓存一个10秒的空值，防止缓存穿透
				go util.SetNull(key)
			} else {
				//缓存到redis
				go util.SetRedisCache(key, videoList)
				fmt.Println(videoList, 111111111111)
			}
			util.RedisUnlock(lockNum)
		} else {
			time.Sleep(time.Millisecond * 100)
			videoList, err = util.GetFeed()
			if err != nil {
				logrus.Info("查询feed列表缓存失败", err)
			}
		}
	}
	return
}
