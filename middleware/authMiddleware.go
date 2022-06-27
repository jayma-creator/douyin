package middleware

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

var count int64

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		fmt.Println(111, token, 222)
		fmt.Println(token == "")
		if token == "" {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 1,
				StatusMsg:  "请登录账号",
			})
			c.Abort()
			return
		}
		user, exist, err := checkToken(token)
		if err != nil {
			logrus.Error("鉴权失败", err)
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 1,
				StatusMsg:  "token超时，请重新登陆",
			})
			c.Abort()
		}
		c.Set("user", user)
		c.Set("exist", exist)

		return
	}
}

func FeedAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		//无用户登录也可以返回视频流
		if token == "" {
			return
		}
		user, exist, err := checkToken(token)
		if err != nil {
			logrus.Error("鉴权失败", err)
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 1,
				StatusMsg:  "token超时，请重新登陆",
			})
			c.Abort()
		}
		c.Set("user", user)
		c.Set("exist", exist)
		return
	}
}

func checkToken(token string) (common.User, bool, error) {
	user := common.User{}
	claims, err := service.ParseToken(token)
	if err != nil {
		logrus.Error(err)
		return user, false, err
	}
	err = dao.DB.Where("name = ? and password = ?", claims.Username, claims.Password).Find(&user).Count(&count).Error
	fmt.Println(claims.Username, claims.Password)
	if err != nil {
		logrus.Error("token is invalid", err)
		return user, false, err
	}
	if count == 0 {
		logrus.Error("token已过期", err)
		return user, false, err
	}
	return user, true, err
}
