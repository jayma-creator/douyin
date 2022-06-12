package service

import (
	"github.com/RaymondCode/simple-demo/dao"
	"time"
)

var U1 User
var U2 User

func InitDemo() {
	//测试用例，启动直接写在数据库
	//用户u1，u2
	U1 = User{
		Id:            1,
		Name:          "张三",
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
		Password:      "123123",
		Token:         "zhangsan123123123",
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
		DeletedAt:     nil,
	}
	U2 = User{
		Id:            2,
		Name:          "李四",
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
		Password:      "123123",
		Token:         "lisi123123123",
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
		DeletedAt:     nil,
	}
	dao.DB.Create(&U1)
	dao.DB.Create(&U2)

}
