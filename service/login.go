package service

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func LoginService(c *gin.Context) (err error) {
	user := common.User{}
	username := c.Query("username")
	password := util.GetMD5(c.Query("password"))
	token, err := util.GetToken(username, password)
	if err != nil {
		logrus.Error("获取token失败", err)
		return
	}
	var count int64
	err = dao.DB.Where("name = ? ", username).Find(&user).Count(&count).Error
	if err != nil {
		logrus.Error("查询name失败", err)
		return
	}
	//如果没有对应的name，返回错误信息“用户不存在”
	if count == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: "用户不存在"},
		})
		return
		//如果token不匹配，提示密码错误
	} else if count == 1 && password != user.Password {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: "密码错误"},
		})
		return
	} else if count == 1 && password == user.Password {
		err = dao.DB.Model(&user).Where("id = ?", user.Id).Update("token", token).Error
		if err != nil {
			logrus.Error("更新token失败", err)
			return
		}
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: common.Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	}

	return
}
