package controller

import (
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	token := username + password

	//查找数据库有没有对应的token
	dao.DB.Where("token = ?", token).Find(&user).Count(&count)
	//如果查询到已存在对应的token，返回错误信息“已存在”
	if count == 1 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
		//如果查询到不存在，则往数据库里添加对应的用户信息
	} else if count == 0 {
		//atomic.AddInt64(&userIdSequence, 1)
		newUser := User{
			//Id:       userIdSequence,
			Name:     username,
			Password: password,
			Token:    token,
		}
		//往数据库添加一行数据
		dao.DB.Create(&newUser)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	}
}

func Login(c *gin.Context) {
	//要添加user := User{} 才能重置count数
	user := User{}
	username := c.Query("username")
	password := c.Query("password")
	token := username + password
	//查找数据库有没有对应的token
	dao.DB.Where("token = ?", token).Find(&user).Count(&count)
	//如果没有对应的token，返回错误信息“用户不存在”
	if count == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
		return
		//如果有对应的token，返回用户信息
	} else if count == 1 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    user.Token,
		})
	}
}

func UserInfo(c *gin.Context) {
	user := User{}
	token := c.Query("token")
	dao.DB.Where("token = ?", token).Find(&user).Count(&count)
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
}
