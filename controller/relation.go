package controller

import (
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

func RelationAction(c *gin.Context) {
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
		dao.DB.Model(&user).Where("id = ?", toUserId).Update("is_follow", true)
	} else {
		//2表示取消关注
		//把对方用户从关注列表里删除
		//把当前用户从对方用户的粉丝列表里删除
		dao.DB.Where("follow_id = ? and follower_id = ?", user.Id, toUserId).Delete(FollowFansRelation{})
		//修改对方用户的is_follow字段为false，表示未关注
		dao.DB.Model(&user).Where("id = ?", toUserId).Update("is_follow", false)
	}
}

func FollowList(c *gin.Context) {
	userId := c.Query("user_id")
	//找出当前id的结构体
	dao.DB.Where("follow_id = ?", userId).Find(&followSlice)
	//把当前用户关注的对方用户放到切片relations里
	//在关系表里查询
	//每一次都要重新清零，不然会一直在切片末尾追加，无法改变前面已添加的
	followIdSlice = []int64{}
	for i := 0; i < len(followSlice); i++ {
		followIdSlice = append(followIdSlice, followSlice[i].FollowerId)
	}
	//从User表里用对方的ID找出对应的结构体
	//在用户表里查询，不是关系表
	//根据切片toUserIds里的对方用户ID直接查询，放在UserList切片
	dao.DB.Where(followIdSlice).Find(&followList) // gorm写法
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: followList,
	})
}

func FollowerList(c *gin.Context) {
	userId := c.Query("user_id")
	//根据follower_id找出当前用户的结构体，和关注列表相反
	dao.DB.Where("follower_id = ?", userId).Find(&fansSlice)
	//把当前用户关注的followid放到切片relations里
	//每一次都要重新清零，不然会一直在切片末尾追加，无法改变前面已添加的
	fansIdSlice = []int64{}
	for i := 0; i < len(fansSlice); i++ {
		//粉丝列表是找FollowId，和关注列表相反
		//这里的toUserIds还是对方用户的Id，只不过现在对方用户Id在粉丝列表中
		fansIdSlice = append(fansIdSlice, fansSlice[i].FollowId)
	}
	//从User表里用对方的ID找出对应的结构体
	//在用户表里查询，不是关系表
	//根据切片toUserIds里的对方用户ID直接查询，放在UserList切片
	dao.DB.Where(fansIdSlice).Find(&fansList) // gorm写法
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: fansList,
	})
}
