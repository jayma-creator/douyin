package controller

import (
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
)

//注册
func Register(c *gin.Context) {
	service.RegisterService(c)
}

//登录
func Login(c *gin.Context) {
	service.LoginService(c)
}

//用户信息
func UserInfo(c *gin.Context) {
	service.UserInfoService(c)
}
