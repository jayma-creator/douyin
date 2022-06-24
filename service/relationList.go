package service

import (
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
	userList := []User{}
	followList := []FollowFansRelation{}
	//查询出当前用户关注的列表
	err = dao.DB.Where("follow_id = ?", userId).Preload("Follower").Find(&followList).Error
	if err != nil {
		logrus.Error("获取关注列表失败", err)
		return
	}
	for i := 0; i < len(followList); i++ {
		userList = append(userList, followList[i].Follower)
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: userList,
	})
	return
}

//粉丝列表
func FollowerListService(c *gin.Context) (err error) {
	userId := c.Query("user_id")
	fansList := []User{}
	followList := []FollowFansRelation{}
	//查询出当前用户的粉丝列表
	err = dao.DB.Where("follower_id = ?", userId).Preload("Follow").Find(&followList).Error
	if err != nil {
		logrus.Error("获取粉丝列表失败", err)
		return
	}
	for i := 0; i < len(followList); i++ {
		fansList = append(fansList, followList[i].Follow)
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: fansList,
	})
	return
}
