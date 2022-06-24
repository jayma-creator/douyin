package service

import (
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

//用户信息
func UserInfoService(c *gin.Context) (err error) {
	user := User{}
	token := c.Query("token")
	err = dao.DB.Where("token = ?", token).Find(&user).Count(&count).Error
	if err != nil {
		logrus.Error("查询token失败", err)
		return
	}
	if count == 0 {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
		return
		//如果有对应的token，返回用户信息
	} else if count == 1 {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	}
	return
}
