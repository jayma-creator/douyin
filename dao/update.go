package dao

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func UpdateFollowUserAdd(tx *gorm.DB, user common.User) (err error) {
	err = tx.Model(&common.User{}).Where("id = ?", user.Id).Update("follow_count", gorm.Expr("follow_count + ?", "1")).Error
	return
}

func UpdateFollowFansAdd(tx *gorm.DB, toUserId int) (err error) {
	err = tx.Model(&common.User{}).Where("id = ?", toUserId).Updates(map[string]interface{}{"follower_count": gorm.Expr("follower_count + ?", "1"), "is_follow": true}).Error
	return
}

func UpdateFollowUserDel(tx *gorm.DB, user common.User) (err error) {
	err = tx.Model(&common.User{}).Where("id = ?", user.Id).Update("follow_count", gorm.Expr("follow_count - ?", "1")).Error
	return
}

func UpdateFollowFansDel(tx *gorm.DB, toUserId int) (err error) {
	err = tx.Model(&common.User{}).Where("id = ?", toUserId).Updates(map[string]interface{}{"follower_count": gorm.Expr("follower_count - ?", "1"), "is_follow": false}).Error
	return
}

func UpdateLikeAdd(tx *gorm.DB, videoId int) (err error) {
	err = tx.Model(&common.Video{}).Where("id = ?", videoId).Updates(map[string]interface{}{"is_favorite": true, "favorite_count": gorm.Expr("favorite_count + ?", "1")}).Error
	if err != nil {
		logrus.Error("修改视频信息失败", err)
	}
	return
}

func UpdateLikeDel(tx *gorm.DB, videoId int) (err error) {
	err = tx.Model(&common.Video{}).Where("id = ?", videoId).Updates(map[string]interface{}{"is_favorite": false, "favorite_count": gorm.Expr("favorite_count - ?", "1")}).Error
	if err != nil {
		logrus.Error("修改视频信息失败", err)
	}
	return
}

func UpdateCommentAdd(tx *gorm.DB, videoId int) (err error) {
	err = tx.Model(&common.Video{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + ?", "1")).Error
	return
}

func UpdateCommentDel(tx *gorm.DB, videoId int) (err error) {
	err = tx.Model(&common.Video{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count - ?", "1")).Error
	return
}

func UpdateInitLikeInfo(tx *gorm.DB) (err error) {
	err = tx.Model(common.Video{}).Where("is_favorite = ?", true).Update("is_favorite", false).Error
	return
}

func UpdateInitFollowInfo(tx *gorm.DB) (err error) {
	err = tx.Model(common.User{}).Where("is_follow = ?", true).Update("is_follow", false).Error
	return
}

func UpdateIsFollowInfo(tx *gorm.DB, users []common.User) (err error) {
	for i := 0; i < len(users); i++ {
		err = tx.Model(&common.User{}).Where("id = ?", users[i].Id).Update("is_follow", true).Error
	}
	return
}

func UpdateIsLike(tx *gorm.DB, videos []common.Video) (err error) {
	for i := 0; i < len(videos); i++ {
		err = tx.Model(&common.Video{}).Where("id = ?", videos[i].Id).Update("is_favorite", true).Error
	}
	return
}
