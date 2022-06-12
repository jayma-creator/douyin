package service

import (
	"crypto/md5"
	"fmt"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterService(c *gin.Context) {
	username := c.Query("username")
	password := GetMD5(c.Query("password"))
	token := username + password

	user := User{}
	dao.DB.Where("token = ?", token).Find(&user).Count(&count)
	//如果查询到已存在对应的token，返回错误信息“已存在”
	if count == 1 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
		//如果查询到不存在，则往数据库里添加对应的用户信息
	} else if count == 0 {
		newUser := User{
			//Id:       userIdSequence,
			Name:     username,
			Password: password,
			Token:    token,
		}
		//插入数据
		dao.DB.Create(&newUser)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	}
}

func GetMD5(str string) string {
	data := []byte(str)
	strMD5 := fmt.Sprintf("%x", md5.Sum(data))
	return strMD5
}
