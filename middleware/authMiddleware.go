package middleware

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 1,
				StatusMsg:  "请登录账号",
			})
			c.Abort()
			return
		}
		user, exist, err := CheckToken(token)
		if err != nil {
			logrus.Error("鉴权失败", err)
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 1,
				StatusMsg:  "token已过期，请重新登陆",
			})
			c.Abort()
		}
		c.Set("user", user)
		c.Set("exist", exist)
		return
	}
}

func PublishAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.PostForm("token")
		if token == "" {
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 1,
				StatusMsg:  "请登录账号",
			})
			c.Abort()
			return
		}
		user, exist, err := CheckToken(token)
		if err != nil {
			logrus.Error("鉴权失败", err)
			c.JSON(http.StatusOK, common.Response{
				StatusCode: 1,
				StatusMsg:  "token已过期，请重新登陆",
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
		user, exist, err := CheckToken(token)
		if err != nil {
			logrus.Error("鉴权失败", err)
			c.Next()
			return
		}
		c.Set("user", user)
		c.Set("exist", exist)
		return
	}
}

func CheckToken(token string) (user common.User, bool bool, err error) {
	conn := dao.Pool.Get()
	defer conn.Close()
	claims, err := util.ParseToken(token)
	exist, _ := conn.Do("exists", token)
	if exist.(int64) == 1 {
		user, err = util.GetUserCache(claims.Username)
		if err != nil {
			logrus.Info("查询用户信息缓存失败", err)
		}

		//说明redis没有缓存，改为从数据库读取,并缓存到redis
		if user == (common.User{}) {
			err = dao.DB.Where("name = ?", claims.Username).Find(&user).Error
			//把user信息缓存到redis
			go util.SetRedisCache(fmt.Sprintf("user%v", claims.Username), user)
		}
		//每次请求都会刷新token
		util.RefreshToken(token)
	} else {
		logrus.Info("jwt设定的token已过期", err)
		return user, false, err
	}
	return user, true, err
}
