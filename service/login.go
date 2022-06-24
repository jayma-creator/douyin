package service

import (
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func LoginService(c *gin.Context) (err error) {
	user := User{}
	username := c.Query("username")
	password := GetMD5(c.Query("password"))
	token := username + password
	err = dao.DB.Where("name = ?", username).Find(&user).Count(&count).Error
	if err != nil {
		logrus.Error("查询name失败", err)
		return
	}
	//如果没有对应的token，返回错误信息“用户不存在”
	if count == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户不存在"},
		})
		return
		//如果token不匹配，提示密码错误
	} else if count == 1 && token != user.Token {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "密码错误"},
		})
		return
	} else if count == 1 && token == user.Token {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    user.Token,
		})
	}
	//匹配当前登录的账号是否已关注别的账号
	users := []User{}
	dao.DB.Table("users").
		Joins("join follow_fans_relations on follower_id = users.id and follow_id = ?", user.Id).
		Scan(&users)
	for i := 0; i < len(users); i++ {
		users[i].IsFollow = true
		dao.DB.Model(&User{}).Where("id = ?", users[i].Id).Update("is_follow", true)
	}

	//匹配当前登录的账号是否已点赞视频
	videos := []Video{}
	dao.DB.Table("videos").
		Joins("join user_favorite_relations on video_id = videos.id and user_id = ?", user.Id).
		Scan(&videos)
	for i := 0; i < len(videos); i++ {
		videos[i].IsFavorite = true
		dao.DB.Model(&Video{}).Where("id = ?", videos[i].Id).Update("is_favorite", true)
	}
	return
}
