package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//用户信息
func UserInfoService(c *gin.Context) (err error) {
	token := c.Query("token")
	user, exist, err := CheckToken(token)
	if exist {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	} else {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "token已过期，请重新登录"},
		})
		return
	}
	return
}
