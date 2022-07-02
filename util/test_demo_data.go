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
	encodePwd, _ := GetMD5WithSalted("123123")

	U1 = common.User{
		Id:            1,
		Name:          "zhangsan",
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
		Password:      encodePwd,
	}
	U2 = common.User{
		Id:            2,
		Name:          "lisi",
		FollowCount:   0,
		FollowerCount: 0,
		IsFollow:      false,
		Password:      encodePwd,
	}
	dao.DB.Create(&U1)
	dao.DB.Create(&U2)

}
