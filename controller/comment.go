package controller

import (
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

func CommentAction(c *gin.Context) {
	user := User{}
	token := c.Query("token")
	actionType := c.Query("action_type")
	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.Atoi(videoIdStr)
	dao.DB.Where("token = ?", token).Find(&user).Count(&count)
	if count == 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	if actionType == "1" {
		text := c.Query("comment_text")
		//新增评论
		comment := Comment{
			//User:       user, //User是结构体类型，该字段不会在数据库里创建，所以这里可以省略
			Content:    text,
			CreateDate: time.Now().Format("2006-01-02 15:04:05")[5:10], //按格式输出日期，5:10表示月-日
			UserToken:  token,
			VideoId:    int64(videoId),
		}
		dao.DB.Create(&comment)
		//为了能评论后即时显示用户名，这里手动赋值comment的User
		comment.User = user

		c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
			Comment: comment,
		})

	} else if actionType == "2" {
		commentId := c.Query("comment_id")
		dao.DB.Where("id = ?", commentId).Delete(&Comment{})
	}

}

func CommentList(c *gin.Context) {
	//把数据库里所有评论放在commentList内
	commentList := []Comment{}
	user := User{}
	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.Atoi(videoIdStr)
	token := c.Query("token")
	dao.DB.Where("token = ?", token).Find(&user)

	//根据当前播放的视频ID匹配对应的评论
	//comment里的videoId要等于当前播放的Id
	dao.DB.Where("video_id = ?", videoId).Find(&commentList)

	//匹配评论与作者
	//发表的评论有user_token，和当前用户的token对应起来
	//根据user_token取出对应的user结构体，赋值给comment的User
	commentTokenSlice := []string{}
	for i := 0; i < len(commentList); i++ {
		commentTokenSlice = append(commentTokenSlice, commentList[i].UserToken)
	}
	//循环取出User结构体,赋给相对应的comment，通俗点说就是发布者匹配
	for i := 0; i < len(commentTokenSlice); i++ {
		user := User{}
		dao.DB.Where("token = ?", commentTokenSlice[i]).Find(&user)
		commentList[i].User = user
	}
	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0},
		CommentList: commentList,
	})
}
