package controller

import (
	"time"
)

//关注和粉丝部分
var count int                        //查询数据库的记录条数标记
var user User                        //定义一个结构体修改用户数据，gorm写法&user会自动去user表里找
var followSlice []FollowFansRelation //放当前用户的结构体，根据follower来找，方便取出对方用户的ID，也就是关注对象的ID
var fansSlice []FollowFansRelation   //放当前用户的结构体，根据followerid来找，方便取出对方用户的ID，也就是粉丝对象的ID
var followList []User                //根据对方用户的ID，从用户表里找出来的结构体，用来返回给软件关注列表展示
var fansList []User                  //根据对方用户的ID，从用户表里找出来的结构体，用来返回给软件粉丝列表展示
var followIdSlice []int64            //放对方用户的ID，对方用户是当前ID的关注
var fansIdSlice []int64              //放对方用户的ID，对方用户是当前ID的粉丝

//用户点赞部分
var favoriteRelationVideoIdList []UserFavoriteRelation
var videoList []Video

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

//添加了PublisherToken字段，来判别视频属于谁发布的
type Video struct {
	Id             int64      `json:"id,omitempty"`
	Author         User       `json:"author"` //注意这里Author属性是不会导入数据库的
	PlayUrl        string     `json:"play_url,omitempty"`
	CoverUrl       string     `json:"cover_url,omitempty"`
	FavoriteCount  int64      `json:"favorite_count,omitempty" gorm:"default:'0'"`
	CommentCount   int64      `json:"comment_count,omitempty" gorm:"default:'0'"`
	IsFavorite     bool       `json:"is_favorite,omitempty"`
	PublisherToken string     `json:"publisher_token"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
	DeletedAt      *time.Time `json:"deleted_at"`
}

type Comment struct {
	Id         int64      `json:"id,omitempty"`
	User       User       `json:"user"`
	Content    string     `json:"content,omitempty"`
	CreateDate string     `json:"create_date,omitempty"`
	UserToken  string     `json:"user_token" gorm:"comment:'发表评论用户的token'"`
	VideoId    int64      `json:"video_id" gorm:"'发表评论的视频id'"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}

type User struct {
	Id            int64      `json:"id,omitempty"`
	Name          string     `json:"name,omitempty"`
	FollowCount   int64      `json:"follow_count,omitempty" gorm:"default:'0'"`
	FollowerCount int64      `json:"follower_count,omitempty" gorm:"default:'0'"`
	IsFollow      bool       `json:"is_follow,omitempty" gorm:"default:'0'"`
	Password      string     `json:"password,omitempty"`
	Token         string     `json:"token,omitempty" gorm:"unique"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}

//粉丝和关注的关系表
type FollowFansRelation struct {
	Id         int64      `json:"id,omitempty"`
	FollowId   int64      `json:"follow_id"`
	FollowerId int64      `json:"follower_id"`
	CreatedAt  time.Time  `json:"created_at" gorm:"comment:'关注时间'"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at" gorm:"comment:'取关时间'"`
}

//用户点赞关系表
type UserFavoriteRelation struct {
	Id        int64      `json:"id,omitempty"`
	UserId    int64      `json:"user_id"`
	VideoId   int64      `json:"video_id"`
	CreatedAt time.Time  `json:"created_at" gorm:"comment:'like时间'"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"comment:'unlike时间'"`
}
