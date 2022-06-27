package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}
type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

//关注列表
func FollowListService(c *gin.Context) (err error) {
	userId := c.Query("user_id")
	//查询出当前用户关注的列表
	followList, err := getFollowListCache(userId)
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
		err = setRedisCache(fmt.Sprintf("followList%v", userId), followList)
		if err != nil {
			logrus.Error("缓存失败")
		}
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: followList,
	})
	return
}

//粉丝列表
func FollowerListService(c *gin.Context) (err error) {
	userId := c.Query("user_id")
	fansList, err := getFanListCache(userId)
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
		err = setRedisCache(fmt.Sprintf("fansList%v", userId), fansList)
		if err != nil {
			logrus.Error("缓存失败")
		}
	}

	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: fansList,
	})
	return
}
