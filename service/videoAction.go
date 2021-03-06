package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type VideoListResponse struct {
	common.Response
	VideoList []common.Video `json:"video_list"`
}

const (
	like   = "1"
	unLike = "2"
)

//点赞与取消
func FavoriteActionService(c *gin.Context) (err error) {
	videoIdStr := c.Query("video_id")
	actionType := c.Query("action_type")
	videoId, _ := strconv.Atoi(videoIdStr)
	var count int64
	u, _ := c.Get("user")
	e, _ := c.Get("exist")
	if u != nil && e != nil {
		user := u.(common.User)
		exist := e.(bool)
		key := strconv.Itoa(int(user.Id)) + strconv.Itoa(videoId) + "favorite"
		if exist {
			if actionType == like {
				//先查询缓存对应的ID有没有点赞该视频
				exist := util.IsExistCache(key)
				//如果有，则直接返回已经关注
				if exist == 1 {
					c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: "已经点赞该视频，请刷新视频查看"})
					return
				} else {
					//如果缓存没有，则查询数据库
					count, err = dao.QueryLike(user, videoId)
					if err != nil {
						return
					}
					//如果数据库没有，则执行关注操作，并把关注信息缓存到redis
					if count == 0 {
						err = likeAct(c, user, videoId)
						if err != nil {
							return err
						}
					} else {
						c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: "已经点赞该视频，请刷新视频查看"})
						return err
					}
					go util.SetRedisNum(key, key)
				}
			} else if actionType == unLike {
				//查询数据库
				count, err = dao.QueryLike(user, videoId)
				if err != nil {
					logrus.Error(err)
				}
				if count == 1 {
					err = unlikeAct(c, user, videoId)
					if err != nil {
						return err
					}
				} else {
					c.JSON(http.StatusOK, common.Response{StatusCode: 1, StatusMsg: "已经取消点赞该视频，请刷新视频查看"})
					return err
				}
				go util.DelCache(key)
			}
		}
	}
	return
}

//点赞
func likeAct(c *gin.Context, user common.User, videoId int) (err error) {
	tx := dao.DB.Begin()
	ufr := common.UserFavoriteRelation{
		UserId:  user.Id,
		VideoId: int64(videoId),
	}
	err = tx.Create(&ufr).Error
	if err != nil {
		logrus.Error("新增视频信息失败", err)
		tx.Rollback()
		return
	}
	//把video结构体里的IsFavorite改为true
	//video的favorite_count+1
	err = dao.UpdateLikeAdd(tx, videoId)
	if err != nil {
		logrus.Error("修改视频信息失败", err)
		tx.Rollback()
		return
	}

	//删除redis缓存
	err = util.DelCache("feed")
	err = util.DelCache(fmt.Sprintf("favoriteList%v", user.Id))
	if err != nil {
		return
	}
	tx.Commit()
	//延时双删
	time.Sleep(time.Millisecond * 50)
	err = util.DelCache("feed")
	err = util.DelCache(fmt.Sprintf("favoriteList%v", user.Id))
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, common.Response{StatusCode: 0, StatusMsg: "点赞成功"})
	return
}

//取消赞
func unlikeAct(c *gin.Context, user common.User, videoId int) (err error) {
	tx := dao.DB.Begin()
	err = dao.DeleteLike(tx, user, videoId)
	if err != nil {
		logrus.Error("删除视频信息失败", err)
		tx.Rollback()
		return
	}
	err = dao.UpdateLikeDel(tx, videoId)
	if err != nil {
		logrus.Error(err)
		tx.Rollback()
		return
	}

	//删除redis缓存
	err = util.DelCache("feed")
	err = util.DelCache(fmt.Sprintf("favoriteList%v", user.Id))
	if err != nil {
		return
	}
	tx.Commit()
	//延时双删
	time.Sleep(time.Millisecond * 50)
	err = util.DelCache("feed")
	err = util.DelCache(fmt.Sprintf("favoriteList%v", user.Id))
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, common.Response{StatusCode: 0, StatusMsg: "取消赞成功"})
	return
}
