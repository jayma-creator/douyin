package router

import (
	"github.com/RaymondCode/simple-demo/controller"
	"github.com/RaymondCode/simple-demo/setting"
	"github.com/gin-gonic/gin"
	"net/http"
)

func InitRouter(r *gin.Engine) {
	//配置文件里如果Release为true则为生产环境模式
	if setting.Conf.Release {
		gin.SetMode(gin.ReleaseMode)
	}
	r.Static("/static", "./public")
	{
		apiRouter := r.Group("/douyin") //提取公用的前缀，下面的就省略前缀

		// basic apis
		apiRouter.GET("/feed/", controller.Feed)
		apiRouter.GET("/user/", controller.UserInfo)
		apiRouter.POST("/user/register/", controller.Register)
		apiRouter.POST("/user/login/", controller.Login)
		apiRouter.POST("/publish/action/", controller.Publish)
		apiRouter.GET("/publish/list/", controller.PublishList)
		apiRouter.GET("/hello", func(c *gin.Context) {
			c.JSON(http.StatusOK, 3)
		})
		// extra apis - I
		apiRouter.POST("/favorite/action/", controller.FavoriteAction)
		apiRouter.GET("/favorite/list/", controller.FavoriteList)
		apiRouter.POST("/comment/action/", controller.CommentAction)
		apiRouter.GET("/comment/list/", controller.CommentList)

		// extra apis - II
		apiRouter.POST("/relation/action/", controller.RelationAction)
		apiRouter.GET("/relation/follow/list/", controller.FollowList)
		apiRouter.GET("/relation/follower/list/", controller.FollowerList)
	}
}
