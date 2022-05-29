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
	//查找数据库有没有对应的token
	dao.DB.Where("token = ?", token).Find(&user).Count(&count)
	//如果查询到已存在对应的token，返回错误信息“已存在”
	if count == 1 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
		//如果查询到不存在，则往数据库里添加对应的用户信息
	} else if count == 0 {
		//atomic.AddInt64(&userIdSequence, 1)
		newUser := User{
			//Id:       userIdSequence,
			Name:     username,
			Password: password,
			Token:    token,
		}
		//往数据库添加一行数据
		dao.DB.Create(&newUser)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   user.Id,
			Token:    token,
		})
	}
}

func LoginService(c *gin.Context) {
	//要添加user := User{} 才能重置count数
	user := User{}
	username := c.Query("username")
	password := c.Query("password")
	token := username + password
	//查找数据库有没有对应的token
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
	//如果当前登录的为ID1,1关注了2，那么从关注粉丝表里找出当前ID所关注了的对方ID
	//然后根据对方ID在users表里更改对应id的ISfollow为true
	//每次更改账号后都要运行一次
	//每次换账号登录后都要重新运行一次的有:更新关注的显示，点赞的显示

	//首先每次登陆账号,或切换账号要做的事：
	//1.把所有用户的is_follow改为false  2.把所有视频的is_favorite改为false
	//目前关掉软件重开可以正确显示，不关软件切换账号显示异常，但数据库有改动，客户端问题
	dao.DB.Model(&User{}).Update("is_follow", false)
	dao.DB.Model(&Video{}).Update("is_favorite", false)

	//匹配当前登录的账号是否已关注别的账号
	//拿出当前用户的关注粉丝表结构体
	followFansRelations := []FollowFansRelation{}
	dao.DB.Where("follow_id = ?", user.Id).Find(&followFansRelations)
	//拿出当前用户关注的对方ID
	temp := []int64{} //放对方用户的ID
	for i := 0; i < len(followFansRelations); i++ {
		temp = append(temp, followFansRelations[i].FollowerId)
	}
	//根据对方用户的ID找到相应的user，把is_follow改为true
	dao.DB.Model(&User{}).Where(temp).Update("is_follow", true)

	//匹配当前登录的账号是否已点赞视频
	userFavoriteRelations := []UserFavoriteRelation{}
	dao.DB.Where("user_id = ?", user.Id).Find(&userFavoriteRelations)
	temp2 := []int64{}
	for i := 0; i < len(userFavoriteRelations); i++ {
		temp2 = append(temp2, userFavoriteRelations[i].VideoId)
	}
	dao.DB.Model(&Video{}).Where(temp2).Update("is_favorite", true)
}

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
		//1表示关注
		//如果当前用户点击关注自己，返回错误提示
		if user.Id == int64(toUserId) {
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "不能关注自己"}) //实际看不到文字，客户端问题
			return
		}
		//把对方用户添加到关注列表里
		//把当前用户添加到对方用户的粉丝列表
		r := FollowFansRelation{
			FollowId:   user.Id,
			FollowerId: int64(toUserId),
		}
		//gorm 增加一行
		dao.DB.Create(&r)
		//修改对方用户的is_follow字段为true，表示已关注
		dao.DB.Model(&User{}).Where("id = ?", toUserId).Update("is_follow", true)
		//当前ID的user结构体里的关注数follow_count+1，对方ID的粉丝数follower_count+1
		dao.DB.Model(&User{}).Where("id = ?", user.Id).Update("follow_count", gorm.Expr("follow_count + ?", "1"))
		dao.DB.Model(&User{}).Where("id = ?", toUserId).Update("follower_count", gorm.Expr("follower_count + ?", "1"))

	} else {
		//2表示取消关注
		//把对方用户从关注列表里删除
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
