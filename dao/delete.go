package dao

import (
	"github.com/RaymondCode/simple-demo/common"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func DeleteLike(tx *gorm.DB, user common.User, videoId int) (err error) {
	err = tx.Where("user_id = ? and video_id = ?", user.Id, videoId).Delete(&common.UserFavoriteRelation{}).Error
	if err != nil {
		logrus.Error("删除视频信息失败", err)
		return
	}
	return
}

func DeleteComment(tx *gorm.DB, commentId string) (err error) {
	err = tx.Where("id = ?", commentId).Delete(&common.Comment{}).Error
	return err
}

func DeleteFollow(tx *gorm.DB, user common.User, toUserId int) (err error) {
	err = tx.Where("follow_id = ? and follower_id = ?", user.Id, toUserId).Delete(&common.FollowFansRelation{}).Error
	return
}
