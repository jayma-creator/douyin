package controller

import (
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//注册
func Register(c *gin.Context) {
	err := service.RegisterService(c)
	if err != nil {
		logrus.Error(err)
		return
	}
}

//登录
func Login(c *gin.Context) {
	err := service.LoginService(c)
	if err != nil {
		logrus.Error(err)
		return
	}
}

//用户信息
func UserInfo(c *gin.Context) {
	err := service.UserInfoService(c)
	if err != nil {
		logrus.Error(err)
		return
	}
}
