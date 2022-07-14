package service

import (
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
		err = dao.UpdateInitLikeInfo(tx)
		if err != nil {
			logrus.Error("修改失败", err)
			tx.Rollback()
			return
		}
		err = dao.UpdateInitFollowInfo(tx)
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
	users := []common.User{}
	//匹配当前登录的账号是否已关注别的账号
	users, err = dao.QueryIsFollow(tx, user)
	if err != nil {
		logrus.Error("修改失败", err)
		tx.Rollback()
		return err
	}
	err = dao.UpdateIsFollowInfo(tx, users)
	if err != nil {
		logrus.Error("修改失败", err)
		tx.Rollback()
		return err
	}

	//匹配当前登录的账号是否已点赞视频
	videos := []common.Video{}
	videos, err = dao.QueryIsLike(tx, user)
	if err != nil {
		logrus.Error("修改失败", err)
		tx.Rollback()
		return err
	}

	err = dao.UpdateIsLike(tx, videos)
	if err != nil {
		logrus.Error("修改失败", err)
		tx.Rollback()
		return err
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
		var count int64
		lockNum := "1"
		if util.RedisLock(lockNum) == true {
			videoList, count, err = dao.QueryFeed()
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
