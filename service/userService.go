package service

import (
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}
type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

var count = 0

func RegisterService(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")
	token := username + password
	user := User{}
	dao.DB.Where("token = ?", token).Find(&user).Count(&count)
	//如果查询到已存在对应的token，返回错误信息“已存在”
	if count == 1 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
		//如果查询到不存在，则往数据库里添加对应的用户信息
	} else if count == 0 {
		newUser := User{
			//Id:       userIdSequence,
			Name:     username,
			Password: password,
			Token:    token,
		}
		//插入数据
		dao.DB.Create(&newUser)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	}
}

func LoginService(c *gin.Context) {
	user := User{}
	username := c.Query("username")
	password := c.Query("password")
	token := username + password
	dao.DB.Where("token = ?", token).Find(&user).Count(&count)
	//如果没有对应的token，返回错误信息“用户不存在”
	if count == 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
		return
		//如果有对应的token，返回用户信息
	} else if count == 1 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    user.Token,
		})
	}
	//匹配当前登录的账号是否已关注别的账号
	users := []User{}
	dao.DB.Table("users").
		Joins("join follow_fans_relations on follower_id = users.id and follow_id = ?", user.Id).
		Scan(&users)
	for i := 0; i < len(users); i++ {
		users[i].IsFollow = true
		dao.DB.Model(&User{}).Where("id = ?", users[i].Id).Update("is_follow", true)
	}

	//匹配当前登录的账号是否已点赞视频
	videos := []Video{}
	dao.DB.Table("videos").
		Joins("join user_favorite_relations on video_id = videos.id and user_id = ?", user.Id).
		Scan(&videos)
	for i := 0; i < len(videos); i++ {
		videos[i].IsFavorite = true
		dao.DB.Model(&Video{}).Where("id = ?", videos[i].Id).Update("is_favorite", true)
	}
}

//用户信息
func UserInfoService(c *gin.Context) {
	user := User{}
	token := c.Query("token")
	dao.DB.Where("token = ?", token).Find(&user).Count(&count)
	if count == 0 {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User doesn't exist"},
		})
		return
		//如果有对应的token，返回用户信息
	} else if count == 1 {
		c.JSON(http.StatusOK, UserResponse{
			Response: Response{StatusCode: 0},
			User:     user,
		})
	}
}

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

//关注列表
func FollowListService(c *gin.Context) {
	userId := c.Query("user_id")
	followList := []User{}
	//查询出当前用户关注的列表
	dao.DB.Table("users").
		Joins("join follow_fans_relations on follower_id = users.id and follow_id = ? and follow_fans_relations.deleted_at is null", userId).
		Scan(&followList)
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: followList,
	})
}

//粉丝列表
func FollowerListService(c *gin.Context) {
	userId := c.Query("user_id")
	fansList := []User{}
	//查询出当前用户的粉丝列表
	dao.DB.Table("users").
		Joins("join follow_fans_relations on follow_id = users.id and follower_id = ? and follow_fans_relations.deleted_at is null", userId).
		Scan(&fansList)
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: fansList,
	})
}
