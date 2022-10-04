package common

import (
	"gorm.io/gorm"
	"time"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// 添加了PublisherToken字段，来判别视频属于谁发布的
type Video struct {
	Id            int64          `json:"id,omitempty"`
	AuthorId      int64          `json:"author_id"`
	Author        User           `json:"author" gorm:"foreignKey:AuthorId"`
	PlayUrl       string         `json:"play_url,omitempty" gorm:"NOT NLL"`
	CoverUrl      string         `json:"cover_url,omitempty" gorm:"NOT NULL"`
	FavoriteCount int64          `json:"favorite_count,omitempty" gorm:"NOT NULL"`
	CommentCount  int64          `json:"comment_count,omitempty" gorm:"NOT NULL"`
	IsFavorite    bool           `json:"is_favorite,omitempty" gorm:"NOT NULL"`
	Title         string         `json:"title" gorm:"NOT NULL"`
	CreatedAt     time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at"`
}

type Comment struct {
	Id         int64          `json:"id,omitempty"`
	UserId     int64          `json:"user_id,omitempty"`
	User       User           `json:"user" gorm:"foreignKey:UserId"`
	VideoId    int64          `json:"video_id" gorm:"index:idx"`
	Video      Video          `json:"video" gorm:"foreignKey:VideoId"`
	Content    string         `json:"content,omitempty" gorm:"NOT NULL"`
	CreateDate string         `json:"create_date,omitempty" gorm:"NOT NULL"`
	CreatedAt  time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at"`
}

type User struct {
	Id            int64          `json:"id,omitempty"`
	Name          string         `json:"name,omitempty" gorm:"index:idx;unique;NOT NUll"`
	FollowCount   int64          `json:"follow_count,omitempty" gorm:"NOT NULL"`
	FollowerCount int64          `json:"follower_count,omitempty" gorm:"NOT NULL"`
	IsFollow      bool           `json:"is_follow,omitempty" gorm:"NOT NULL"`
	Password      string         `json:"password,omitempty" gorm:"NOT NULL"`
	CreatedAt     time.Time      `json:"created_at" gorm:"index"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at"`
}

// 粉丝和关注的关系表
type FollowFansRelation struct {
	Id         int64          `json:"id,omitempty"`
	FollowId   int64          `json:"follow_id" gorm:"index:follow_fan"`
	FollowerId int64          `json:"follower_id" gorm:"index:follow_fan"`
	Follow     User           `json:"follow" gorm:"foreignKey:FollowId"`
	Follower   User           `json:"follower" gorm:"foreignKey:FollowerId"`
	CreatedAt  time.Time      `json:"created_at" `
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" `
}

// 用户点赞视频关系表
type UserFavoriteRelation struct {
	Id        int64          `json:"id,omitempty"`
	UserId    int64          `json:"user_id" gorm:"index:user_video"`
	VideoId   int64          `json:"video_id" gorm:"index:user_video"`
	User      User           `json:"follow,omitempty" gorm:"foreignKey:UserId"`
	Video     Video          `json:"video" gorm:"foreignKey:VideoId"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at"`
}
