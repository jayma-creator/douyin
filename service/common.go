package service

import (
	"time"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

//添加了PublisherToken字段，来判别视频属于谁发布的
type Video struct {
	Id             int64      `json:"id,omitempty"`
	Author         User       `json:"author"`
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
	VideoId    int64      `json:"video_id" gorm:"comment:'发表评论的视频id'"`
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
	Token         string     `json:"token,omitempty" gorm:"unique_index"`
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
