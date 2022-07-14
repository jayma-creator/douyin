package dao

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

func QueryFollow(user common.User, toUserId int) (count int64, err error) {
	err = DB.Where("follow_id = ? and follower_id = ?", user.Id, toUserId).Find(&common.FollowFansRelation{}).Count(&count).Error
	return
}

func QueryLike(user common.User, videoId int) (count int64, err error) {
	err = DB.Where("user_id = ? and video_id = ?", user.Id, videoId).Find(&common.UserFavoriteRelation{}).Count(&count).Error
	if err != nil {
		logrus.Error("查询点赞信息失败", err)
		return
	}
	return
}

func QueryCommentList(videoId string) (commentList []common.Comment, count int64, err error) {
	err = DB.Where("video_id = ?", videoId).Preload("User").Preload("Video").Preload("Video.Author").Order("created_at desc").Find(&commentList).Count(&count).Error
	return
}

func QueryIsFollow(tx *gorm.DB, user common.User) (users []common.User, err error) {
	err = tx.Table("users").
		Joins("join follow_fans_relations on follower_id = users.id and follow_id = ? and follow_fans_relations.deleted_at is null", user.Id).
		Find(&users).Error
	return
}

func QueryIsLike(tx *gorm.DB, user common.User) (videos []common.Video, err error) {
	err = tx.Table("videos").
		Joins("join user_favorite_relations on video_id = videos.id and user_id = ? and user_favorite_relations.deleted_at is null", user.Id).
		Find(&videos).Error
	return
}

func QueryFeed() (videoList []common.Video, count int64, err error) {
	err = DB.Where("created_at <= ?", time.Now().Format("2006-01-02 15:04:05")).Order("created_at desc").Preload("Author").Find(&videoList).Limit(30).Count(&count).Error
	return
}

func QueryFollowList(userId string) (followList []common.User, count int64, err error) {
	err = DB.Table("users").
		Joins("join follow_fans_relations on follower_id = users.id and follow_id = ? and follow_fans_relations.deleted_at is null", userId).
		Find(&followList).Count(&count).Error
	return
}

func QueryFanList(userId string) (fansList []common.User, count int64, err error) {
	err = DB.Table("users").
		Joins("join follow_fans_relations on follow_id = users.id and follower_id = ? and follow_fans_relations.deleted_at is null", userId).
		Find(&fansList).Count(&count).Error
	return
}

func QueryUsernameIsExit(username string) (user common.User, count int64, err error) {
	err = DB.Where("name = ? ", username).Find(&user).Count(&count).Error
	return
}

func QueryPublishList(userId string) (videoList []common.Video, count int64, err error) {
	err = DB.Where("author_id = ?", userId).Preload("Author").Find(&videoList).Count(&count).Error
	return
}

func QueryLikeList(userId string) (videoList []common.Video, count int64, err error) {
	err = DB.Table("videos").
		Joins("join user_favorite_relations on video_id = videos.id and user_id = ? and user_favorite_relations.deleted_at is null", userId).Preload("Author").Find(&videoList).Count(&count).Error
	return
}
