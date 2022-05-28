package service

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}
type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

func FavoriteActionService(c *gin.Context) {
	user := User{}
	token := c.Query("token")
	videoIdStr := c.Query("video_id")
	actionType := c.Query("action_type")
	videoId, _ := strconv.Atoi(videoIdStr)
	count := 0
	dao.DB.Where("token = ?", token).Find(&user).Count(&count)
	if count == 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	} else if count == 1 {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	}
	if actionType == "1" {
		fr := UserFavoriteRelation{
			UserId:  user.Id,
			VideoId: int64(videoId),
		}
		dao.DB.Create(&fr)
		//把video结构体里的IsFavorite改为true
		dao.DB.Model(&Video{}).Where("id = ?", videoId).Update("is_favorite", true)
		//video的favorite_count+1
		dao.DB.Model(&Video{}).Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count + ?", "1"))

	} else if actionType == "2" {
		dao.DB.Where("user_id = ? and video_id = ?", user.Id, videoId).Delete(&UserFavoriteRelation{})
		//把video结构体里的IsFavorite改为false
		dao.DB.Model(&Video{}).Where("id = ?", videoId).Update("is_favorite", false)
		//video的favorite_count-1
		dao.DB.Model(&Video{}).Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count - ?", "1"))
	}
}

func FavoriteListService(c *gin.Context) {
	userId := c.Query("user_id")
	//从点赞关系表中取出当前id的结构体
	userFavoriteRelations := []UserFavoriteRelation{}
	dao.DB.Where("user_id = ?", userId).Find(&userFavoriteRelations)
	//从当前id的结构体中取出video_id字段，保存在切片中
	videoIdSlice := []int64{}
	for i := 0; i < len(videoIdSlice); i++ {
		videoIdSlice = append(videoIdSlice, userFavoriteRelations[i].VideoId)
	}
	//根据video_id找出对应的video结构体放在结构体切片中，并返回前端显示
	videoList := []Video{}
	dao.DB.Where(videoIdSlice).Find(&videoList)
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}

func CommentActionService(c *gin.Context) {
	user := User{}
	token := c.Query("token")
	actionType := c.Query("action_type")
	videoIdStr := c.Query("video_id")
	videoId, _ := strconv.Atoi(videoIdStr)
	count := 0
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
			CreateDate: time.Now().Format("2006-01-02 15:04:05")[5:10], //按格式输出日期，5:10表示月-日  2006-01-02 15:04:05是官方定义的规定格式
			UserToken:  token,
			VideoId:    int64(videoId),
		}
		dao.DB.Create(&comment)
		//为了能评论后即时显示用户名，这里手动赋值comment的User
		comment.User = user
		//video的comment_count+1
		dao.DB.Model(&Video{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + ?", "1"))

		c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
			Comment: comment,
		})

	} else if actionType == "2" {
		commentId := c.Query("comment_id")
		dao.DB.Where("id = ?", commentId).Delete(&Comment{})
		//video的comment_count-1
		dao.DB.Model(&Video{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count - ?", "1"))
	}

}
func CommentListService(c *gin.Context) {
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
func PublishService(c *gin.Context) {
	user := User{}
	token := c.PostForm("token")
	//在user结构体里查找token=客户端传来的token，count计数表示获取条数
	count := 0
	dao.DB.Where("token = ?", token).Find(&user).Count(&count)
	if count == 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}
	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	//文件名
	//filename := filepath.Base(data.Filename)
	finalName := fmt.Sprintf("%d_%s", user.Id, data.Filename)
	//保存在public文件夹下
	saveFile := filepath.Join("./public/", finalName)
	fmt.Println(saveFile)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	video := Video{
		//Author:         user, //Author是User结构体类型，该字段不会在数据库里创建，所以这里可以省略
		PlayUrl:        "http://192.168.220.1:8080/static/" + finalName,
		CoverUrl:       "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
		FavoriteCount:  0,
		CommentCount:   0,
		IsFavorite:     false,
		PublisherToken: token,
		CreatedAt:      time.Time{},
		UpdatedAt:      time.Time{},
		DeletedAt:      nil,
	}
	dao.DB.Create(&video)
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}
func PublishListService(c *gin.Context) {
	//封面问题还未解决
	token := c.Query("token")
	videoList := []Video{}
	dao.DB.Where("publisher_token = ?", token).Find(&videoList)
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}
