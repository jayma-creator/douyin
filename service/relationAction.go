package service

import (
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

func RelationActionService(c *gin.Context) {
	user := User{}
	token := c.Query("token")
	actionType := c.Query("action_type")
	toUserIdStr := c.Query("to_user_id")
	dao.DB.Where("token = ?", token).Find(&user).Count(&count)
	if count == 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "User doesn't exist"})
		return
	} else if count == 1 {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	}
	toUserId, _ := strconv.Atoi(toUserIdStr)
	if actionType == "1" {
		//如果当前用户点击关注自己，返回错误提示
		if user.Id == int64(toUserId) {
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "不能关注自己"})
			return
		}
		//把当前用户添加到对方用户的粉丝列表
		r := FollowFansRelation{
			FollowId:   user.Id,
			FollowerId: int64(toUserId),
		}
		dao.DB.Create(&r)
		//修改对方用户的is_follow字段为true，表示已关注
		dao.DB.Model(&User{}).Where("id = ?", toUserId).Update("is_follow", true)
		//当前ID的user结构体里的关注数follow_count+1，对方ID的粉丝数follower_count+1
		dao.DB.Model(&User{}).Where("id = ?", user.Id).Update("follow_count", gorm.Expr("follow_count + ?", "1"))
		dao.DB.Model(&User{}).Where("id = ?", toUserId).Update("follower_count", gorm.Expr("follower_count + ?", "1"))

	} else {
		//把当前用户从对方用户的粉丝列表里删除
		dao.DB.Where("follow_id = ? and follower_id = ?", user.Id, toUserId).Delete(FollowFansRelation{})
		//修改对方用户的is_follow字段为false，表示未关注
		dao.DB.Model(&User{}).Where("id = ?", toUserId).Update("is_follow", false)
		//当前ID的user结构体里的关注数follow_count-1，对方ID的粉丝数follower_count-1
		dao.DB.Model(&User{}).Where("id = ?", user.Id).Update("follow_count", gorm.Expr("follow_count - ?", "1"))
		dao.DB.Model(&User{}).Where("id = ?", toUserId).Update("follower_count", gorm.Expr("follower_count - ?", "1"))
	}

}
