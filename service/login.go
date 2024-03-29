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
	password := c.Query("password")
	encodePwd, salt := util.GetMD5WithSalted(password)
	token, err := util.GetToken(username, encodePwd)
	if err != nil {
		logrus.Error("获取token失败", err)
		return
	}
	var count int64
	user, count, err = dao.QueryUsernameIsExit(username)
	if err != nil {
		logrus.Error("查询username失败", err)
		return
	}
	//如果没有对应的账号，返回错误信息“用户不存在”
	if count == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: "用户不存在"},
		})
		return
		//如果token不匹配，提示密码错误
	} else if count == 1 && !util.VerifyPassword(password, encodePwd, salt) {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: "密码错误"},
		})
		return
	} else if count == 1 && util.VerifyPassword(password, encodePwd, salt) {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: common.Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	}
	util.DelCache("feed")
	return
}
