package router

import (
	"github.com/RaymondCode/simple-demo/controller"
	"github.com/RaymondCode/simple-demo/middleware"
	"github.com/RaymondCode/simple-demo/setting"
	"github.com/gin-gonic/gin"
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
		apiRouter.GET("/feed/", middleware.FeedAuthMiddleware(), controller.Feed)
		apiRouter.GET("/user/", middleware.AuthMiddleware(), controller.UserInfo)
		apiRouter.POST("/user/register/", controller.Register)
		apiRouter.POST("/user/login/", controller.Login)
		apiRouter.POST("/publish/action/", middleware.PublishAuthMiddleware(), controller.Publish)
		apiRouter.GET("/publish/list/", middleware.AuthMiddleware(), controller.PublishList)

		// extra apis - I
		apiRouter.POST("/favorite/action/", middleware.AuthMiddleware(), controller.FavoriteAction)
		apiRouter.GET("/favorite/list/", middleware.AuthMiddleware(), controller.FavoriteList)
		apiRouter.POST("/comment/action/", middleware.AuthMiddleware(), controller.CommentAction)
		apiRouter.GET("/comment/list/", middleware.AuthMiddleware(), controller.CommentList)

		// extra apis - II
		apiRouter.POST("/relation/action/", middleware.AuthMiddleware(), controller.RelationAction)
		apiRouter.GET("/relation/follow/list/", controller.FollowList)
		apiRouter.GET("/relation/follower/list/", controller.FollowerList)
	}
}
