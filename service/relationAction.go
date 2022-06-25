package service

import (
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

const (
	follow   = "1"
	unfollow = "2"
)

func RelationActionService(c *gin.Context) (err error) {
	token := c.Query("token")
	user, exist, err := CheckToken(token)
	if exist {
		actionType := c.Query("action_type")
		toUserIdStr := c.Query("to_user_id")
		err = dao.DB.Where("token = ?", token).Find(&user).Count(&count).Error
		if err != nil {
			logrus.Error("查询token失败", err)
			return
		}

		toUserId, _ := strconv.Atoi(toUserIdStr)
		if actionType == follow {
			tx := dao.DB.Begin()
			//如果当前用户点击关注自己，返回错误提示
			if user.Id == int64(toUserId) {
				c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "不能关注自己"})
				return
			}
			//把当前用户添加到对方用户的粉丝列表
			r := FollowFansRelation{
				FollowId:   user.Id,
				FollowerId: int64(toUserId),
			}
			err = tx.Create(&r).Error
			if err != nil {
				logrus.Error("插入关注信息失败", err)
				tx.Rollback()
				return
			}
			//修改对方用户的is_follow字段为true，表示已关注
			//修改当前ID的user结构体里的关注数follow_count+1，对方ID的粉丝数follower_count+1
			err = tx.Model(&User{}).Where("id = ?", user.Id).Update("follow_count", gorm.Expr("follow_count + ?", "1")).Error
			if err != nil {
				logrus.Error("修改关注信息失败", err)
				tx.Rollback()
				return
			}
			err = tx.Model(&User{}).Where("id = ?", toUserId).Updates(map[string]interface{}{"follower_count": gorm.Expr("follower_count + ?", "1"), "is_follow": true}).Error
			if err != nil {
				logrus.Error("修改关注信息失败", err)
				tx.Rollback()
				return
			}
			tx.Commit()
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "关注成功"})

		} else if actionType == unfollow {
			tx := dao.DB.Begin()
			err = tx.Where("follow_id = ? and follower_id = ?", user.Id, toUserId).Delete(&FollowFansRelation{}).Error
			if err != nil {
				logrus.Error("删除关注信息失败", err)
				tx.Rollback()
				return
			}
			err = tx.Model(&User{}).Where("id = ?", user.Id).Update("follow_count", gorm.Expr("follow_count - ?", "1")).Error
			if err != nil {
				logrus.Error("修改关注信息失败", err)
				tx.Rollback()
				return
			}
			err = tx.Model(&User{}).Where("id = ?", toUserId).Updates(map[string]interface{}{"follower_count": gorm.Expr("follower_count - ?", "1"), "is_follow": false}).Error
			if err != nil {
				logrus.Error("修改关注信息失败", err)
				tx.Rollback()
				return
			}
			tx.Commit()
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "取关成功"})
		} else {
			c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 1, StatusMsg: "错误操作"}})
			return
		}
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "当前用户不存在"})
		return
	}
	return
}
