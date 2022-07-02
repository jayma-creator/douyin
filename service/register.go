package service

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"sync/atomic"
)

var userIdSequence = int64(2)

func RegisterService(c *gin.Context) (err error) {
	username := c.Query("username")
	password := c.Query("password")
	encodePwd, _ := util.GetMD5WithSalted(password)
	token, err := util.GetToken(username, encodePwd)
	if err != nil {
		logrus.Error("获取token失败", err)
		return
	}
	user := common.User{}
	var count int64
	err = dao.DB.Where("name = ? ", username).Find(&user).Count(&count).Error
	if err != nil {
		logrus.Error("查询name失败", err)
		return
	} //如果查询到已存在对应的name，返回错误信息“已存在”
	if count == 1 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: common.Response{StatusCode: 1, StatusMsg: "已存在用户，请更换用户名"},
		})
		//如果查询到不存在，则往数据库里添加对应的用户信息
	} else if count == 0 {
		atomic.AddInt64(&userIdSequence, 1)
		newUser := common.User{
			Id:       userIdSequence,
			Name:     username,
			Password: encodePwd,
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
