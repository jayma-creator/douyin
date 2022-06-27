package util

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/dao"
)

var U1 common.User
var U2 common.User

func InitDemo() {
	//测试用例，启动直接写在数据库
	//用户u1，u2
	U1 = common.User{
		Id:            1,
		Name:          "zhangsan",
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
		Password:      GetMD5("123123"),
	}
	U2 = common.User{
		Id:            2,
		Name:          "lisi",
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
		Password:      GetMD5("123123"),
	}
	dao.DB.Create(&U1)
	dao.DB.Create(&U2)

}
