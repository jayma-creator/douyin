package service

import (
	"crypto/md5"
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync/atomic"
)

var userIdSequence = int64(2)

func RegisterService(c *gin.Context) (err error) {
	username := c.Query("username")
	password := GetMD5(c.Query("password"))
	token, err := GetToken(username, password)
	if err != nil {
		logrus.Error("获取token失败", err)
		return
	}
	user := common.User{}
	err = dao.DB.Where("name = ? ", username).Find(&user).Count(&count).Error
	if err != nil {
		logrus.Error("查询name失败", err)
		return
	} //如果查询到已存在对应的name，返回错误信息“已存在”
	if count == 1 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
		//如果查询到不存在，则往数据库里添加对应的用户信息
	} else if count == 0 {
		atomic.AddInt64(&userIdSequence, 1)
		newUser := common.User{
			Id:       userIdSequence,
			Name:     username,
			Password: password,
			Token:    token,
		}
		//插入数据
		err = dao.DB.Create(&newUser).Error
		if err != nil {
			logrus.Error("新增用户失败", err)
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

func GetMD5(str string) string {
	data := []byte(str)
	strMD5 := fmt.Sprintf("%x", md5.Sum(data))
	return strMD5
}
