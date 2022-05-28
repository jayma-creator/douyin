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
		FollowCount:   1,
		FollowerCount: 550,
		IsFollow:      false,
		Password:      "123123",
		Token:         "lisi123123123",
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
		DeletedAt:     nil,
	}
	dao.DB.Create(&U1)
	dao.DB.Create(&U2)

	////视频v1,v2
	//v1 := Video{
	//	Id: 1,
	//	//Author:         u1, 因为Author属性不能写入数据库，所以要在后面手动赋值给Author
	//	PlayUrl:        "http://192.168.220.1:8080/static/1_test1.mp4",
	//	CoverUrl:       "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
	//	FavoriteCount:  0,
	//	CommentCount:   0,
	//	IsFavorite:     false,
	//	PublisherToken: "zhangsan123123123",
	//	CreatedAt:      time.Time{},
	//	UpdatedAt:      time.Time{},
	//	DeletedAt:      nil,
	//}
	//v2 := Video{
	//	Id: 2,
	//	//Author:         u2,
	//	PlayUrl:        "https://www.w3schools.com/html/movie.mp4",
	//	CoverUrl:       "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
	//	FavoriteCount:  30,
	//	CommentCount:   50,
	//	IsFavorite:     false,
	//	PublisherToken: "lisi123123123",
	//	CreatedAt:      time.Time{},
	//	UpdatedAt:      time.Time{},
	//	DeletedAt:      nil,
	//}
	//
	//dao.DB.Create(&v1)
	//dao.DB.Create(&v2)
	////测试用例直接手动赋值
	//v1.Author = U1
	//v2.Author = U2
}
