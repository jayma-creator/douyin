package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

const (
	follow   = "1"
	unfollow = "2"
)

func RelationActionService(c *gin.Context) (err error) {
	u, _ := c.Get("user")
	e, _ := c.Get("exist")
	if u != nil && e != nil {
		user := u.(common.User)
		exist := e.(bool)
		if exist {
			actionType := c.Query("action_type")
			toUserIdStr := c.Query("to_user_id")
			toUserId, _ := strconv.Atoi(toUserIdStr)
			//关注操作
			if actionType == follow {
				var count int64
				key := strconv.Itoa(int(user.Id)) + strconv.Itoa((toUserId)) + "follow"
				//先查询缓存对应的ID有没有关注对方
				exist := util.IsExistCache(key)
				//如果有，则直接返回已经关注
				if exist == 1 {
					c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: "已经关注对方，请刷新视频查看"})
					return
				} else { //如果缓存没有，则查询数据库
					err = dao.DB.Where("follow_id = ? and follower_id = ?", user.Id, toUserId).Find(&common.FollowFansRelation{}).Count(&count).Error
					if err != nil {
						logrus.Error("查询关注信息失败", err)
						return
					}
				}
				//如果数据库没有，则执行关注操作，并把关注信息缓存到redis
				if count == 0 {
					err = followAct(c, user, toUserId, key)
					if err != nil {
						return
					}
				} else {
					c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: "已经关注对方，请刷新视频查看"})
					return
				}
			} else if actionType == unfollow {
				err = unFollow(c, user, toUserId)
				if err != nil {
					return
				}
			}
		}
	}
	return
}

//关注操作
func followAct(c *gin.Context, user common.User, toUserId int, key string) (err error) {
	tx := dao.DB.Begin()
	//如果当前用户点击关注自己，返回错误提示
	if user.Id == int64(toUserId) {
		c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: "不能关注自己"})
		return
	}
	//把当前用户添加到对方用户的粉丝列表
	r := common.FollowFansRelation{
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
	err = tx.Model(&common.User{}).Where("id = ?", user.Id).Update("follow_count", gorm.Expr("follow_count + ?", "1")).Error
	if err != nil {
		logrus.Error("修改关注信息失败", err)
		tx.Rollback()
		return
	}
	err = tx.Model(&common.User{}).Where("id = ?", toUserId).Updates(map[string]interface{}{"follower_count": gorm.Expr("follower_count + ?", "1"), "is_follow": true}).Error
	if err != nil {
		logrus.Error("修改关注信息失败", err)
		tx.Rollback()
		return
	}
	//删除redis缓存
	err = util.DelCache(fmt.Sprintf("followList%v", user.Id))
	err = util.DelCache(fmt.Sprintf("fanList%v", user.Id))
	if err != nil {
		return
	}
	tx.Commit()
	//延时双删
	time.Sleep(time.Millisecond * 50)
	err = util.DelCache(fmt.Sprintf("followList%v", user.Id))
	err = util.DelCache(fmt.Sprintf("fanList%v", user.Id))
	c.JSON(http.StatusOK, common.Response{StatusCode: 0, StatusMsg: "关注成功"})
	go util.SetRedisNum(key, key)
	return
}

//取关操作
func unFollow(c *gin.Context, user common.User, toUserId int) (err error) {
	tx := dao.DB.Begin()
	err = tx.Where("follow_id = ? and follower_id = ?", user.Id, toUserId).Delete(&common.FollowFansRelation{}).Error
	if err != nil {
		logrus.Error("删除关注信息失败", err)
		tx.Rollback()
		return
	}
	err = tx.Model(&common.User{}).Where("id = ?", user.Id).Update("follow_count", gorm.Expr("follow_count - ?", "1")).Error
	if err != nil {
		logrus.Error("修改关注信息失败", err)
		tx.Rollback()
		return
	}
	err = tx.Model(&common.User{}).Where("id = ?", toUserId).Updates(map[string]interface{}{"follower_count": gorm.Expr("follower_count - ?", "1"), "is_follow": false}).Error
	if err != nil {
		logrus.Error("修改关注信息失败", err)
		tx.Rollback()
		return
	}
	err = util.DelCache(fmt.Sprintf("followList%v", user.Id))
	err = util.DelCache(fmt.Sprintf("fanList%v", user.Id))
	if err != nil {
		return
	}
	tx.Commit()
	//延时双删
	time.Sleep(time.Millisecond * 50)
	err = util.DelCache(fmt.Sprintf("followList%v", user.Id))
	err = util.DelCache(fmt.Sprintf("fanList%v", user.Id))
	c.JSON(http.StatusOK, common.Response{StatusCode: 0, StatusMsg: "取关成功"})
	return
}
