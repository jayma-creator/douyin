package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
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
	//查询出当前用户关注的列表
	followList, err := util.GetFollowListCache(userId)
	if err != nil {
		logrus.Info("查询点赞列表缓存失败", err)
	}
	//没有缓存，从数据库取，并缓存到redis
	if len(followList) == 0 {
		err = dao.DB.Table("users").
			Joins("join follow_fans_relations on follower_id = users.id and follow_id = ? and follow_fans_relations.deleted_at is null", userId).
			Find(&followList).Error
		if err != nil {
			logrus.Error("获取关注列表失败", err)
			return
		}
		//缓存到redis
		go util.SetRedisCache(fmt.Sprintf("followList%v", userId), followList)
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
	fansList, err := util.GetFanListCache(userId)
	if err != nil {
		logrus.Info("查询点赞列表缓存失败", err)
	}
	if len(fansList) == 0 {
		//查询出当前用户的粉丝列表
		err = dao.DB.Table("users").
			Joins("join follow_fans_relations on follow_id = users.id and follower_id = ? and follow_fans_relations.deleted_at is null", userId).
			Find(&fansList).Error
		if err != nil {
			logrus.Error("获取粉丝列表失败", err)
			return
		}
		//缓存到redis
		go util.SetRedisCache(fmt.Sprintf("fansList%v", userId), fansList)

	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: common.Response{
			StatusCode: 0,
		},
		UserList: fansList,
	})
	return
}
