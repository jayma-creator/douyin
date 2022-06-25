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
	followList := []User{}
	//查询出当前用户关注的列表
	err = dao.DB.Table("users").
		Joins("join follow_fans_relations on follower_id = users.id and follow_id = ? and follow_fans_relations.deleted_at is null", userId).
		Find(&followList).Error
	if err != nil {
		logrus.Error("获取关注列表失败", err)
		return
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
	fansList := []User{}
	//查询出当前用户的粉丝列表
	err = dao.DB.Table("users").
		Joins("join follow_fans_relations on follow_id = users.id and follower_id = ? and follow_fans_relations.deleted_at is null", userId).
		Find(&fansList).Error
	if err != nil {
		logrus.Error("获取粉丝列表失败", err)
		return
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: fansList,
	})
	return
}
