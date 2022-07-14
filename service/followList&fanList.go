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

type UserListResponse struct {
	common.Response
	UserList []common.User `json:"user_list"`
}
type UserLoginResponse struct {
	common.Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	common.Response
	User common.User `json:"user"`
}

//关注列表
func FollowListService(c *gin.Context) (err error) {
	userId := c.Query("user_id")
	key := fmt.Sprintf("followList%v", userId)
	var followList []common.User
	//先查询缓存
	if util.IsExistCache(key) == 1 {
		followList, err = util.GetFollowListCache(userId)
		if err != nil {
			logrus.Info("查询点赞列表缓存失败", err)
		}
	} else if util.IsExistCache(key) == 0 {
		//缓存不存在，从数据库查询
		var count int64
		lockNum := "1"
		if util.RedisLock(lockNum) == true {
			followList, count, err = dao.QueryFollowList(userId)
			if err != nil {
				logrus.Error("获取关注列表失败", err)
				return
			}
		}
		//如果数据库不存在，则缓存一个10秒的空值，防止缓存穿透
		if count == 0 {
			go util.SetNull(key)
		} else {
			//缓存到redis
			go util.SetRedisCache(key, followList)
		}
		util.RedisUnlock(lockNum)
	} else {
		time.Sleep(time.Millisecond * 100)
		followList, err = util.GetFollowListCache(userId)
		if err != nil {
			logrus.Info("查询点赞列表缓存失败", err)
		}
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: common.Response{
			StatusCode: 0,
		},
		UserList: followList,
	})
	return
}

//粉丝列表
func FanListService(c *gin.Context) (err error) {
	userId := c.Query("user_id")
	key := fmt.Sprintf("fansList%v", userId)
	var fansList []common.User
	//先查询缓存
	if util.IsExistCache(key) == 1 {
		fansList, err = util.GetFanListCache(userId)
		if err != nil {
			logrus.Info("查询点赞列表缓存失败", err)
		}
	} else if util.IsExistCache(key) == 0 {
		var count int64
		lockNum := "1"
		if util.RedisLock(lockNum) == true {
			//缓存不存在，从数据库查询
			fansList, count, err = dao.QueryFanList(userId)
			if err != nil {
				logrus.Error("获取粉丝列表失败", err)
				return
			}
			//如果数据库不存在，则缓存一个10秒的空值，防止缓存穿透
			if count == 0 {
				go util.SetNull(key)
			} else {
				//缓存到redis
				go util.SetRedisCache(key, fansList)
			}
			util.RedisUnlock(lockNum)
		} else {
			time.Sleep(time.Millisecond * 100)
			fansList, err = util.GetFanListCache(userId)
			if err != nil {
				logrus.Info("查询点赞列表缓存失败", err)
			}
		}
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: common.Response{
			StatusCode: 0,
		},
		UserList: fansList,
	})
	return
}
