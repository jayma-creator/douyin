package service

import (
	"gorm.io/gorm"
	"time"
)

var count int64

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

//添加了PublisherToken字段，来判别视频属于谁发布的
type Video struct {
	Id            int64          `json:"id,omitempty"`
	AuthorId      int64          `json:"author_id"`
	Author        User           `json:"author" gorm:"foreignKey:AuthorId"`
	PlayUrl       string         `json:"play_url,omitempty"`
	CoverUrl      string         `json:"cover_url,omitempty"`
	FavoriteCount int64          `json:"favorite_count,omitempty"`
	CommentCount  int64          `json:"comment_count,omitempty" `
	IsFavorite    bool           `json:"is_favorite,omitempty"`
	Title         string         `json:"title"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at"`
}

type Comment struct {
	Id         int64          `json:"id,omitempty"`
	UserId     int64          `json:"user_id,omitempty"`
	User       User           `json:"user" gorm:"foreignKey:UserId;"`
	VideoId    int64          `json:"video_id" `
	Video      Video          `json:"video" gorm:"foreignKey:VideoId"`
	Content    string         `json:"content,omitempty"`
	CreateDate string         `json:"create_date,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`
}

type User struct {
	Id            int64          `json:"id,omitempty"`
	Name          string         `json:"name,omitempty"`
	FollowCount   int64          `json:"follow_count,omitempty" `
	FollowerCount int64          `json:"follower_count,omitempty" `
	IsFollow      bool           `json:"is_follow,omitempty" `
	Password      string         `json:"password,omitempty"`
	Token         string         `json:"token,omitempty" gorm:"unique_index"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at"`
}

//粉丝和关注的关系表
type FollowFansRelation struct {
	Id         int64          `json:"id,omitempty"`
	FollowId   int64          `json:"follow_id"`
	FollowerId int64          `json:"follower_id"`
	Follow     User           `json:"follow" gorm:"foreignKey:FollowId"`
	Follower   User           `json:"follower" gorm:"foreignKey:FollowerId"`
	CreatedAt  time.Time      `json:"created_at" `
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" `
}

//用户点赞视频关系表
type UserFavoriteRelation struct {
	Id        int64          `json:"id,omitempty"`
	UserId    int64          `json:"user_id"`
	VideoId   int64          `json:"video_id"`
	User      User           `json:"follow,omitempty" gorm:"foreignKey:UserId"`
	Video     Video          `json:"video" gorm:"foreignKey:VideoId"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
