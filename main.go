package main

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/RaymondCode/simple-demo/router"
	"github.com/RaymondCode/simple-demo/setting"
	"github.com/RaymondCode/simple-demo/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	// 加载配置文件
	if err := setting.Init("conf/config.ini"); err != nil {
		logrus.Error("加载配置文件失败", err)
		return
	}

	// 连接数据库
	err := dao.InitMySQL(setting.Conf.MySQLConfig)
	if err != nil {
		logrus.Error("初始化mysql失败", err)
		return
	}

	//连接redis
	err = dao.InitRedisPool(setting.Conf.RedisConfig)
	if err != nil {
		logrus.Error("初始化redis失败", err)
		return
	}

	//初始化日志配置
	err = util.InitLogRecord()
	if err != nil {
		logrus.Error("初始化日志失败", err)
		return
	}

	util.InitDemo() //初始化测试数据

	r := gin.Default()
	router.InitRouter(r)
	err = r.Run(fmt.Sprintf(":%d", setting.Conf.GinConfig.Port))
	if err != nil {
		logrus.Errorf("server startup failed, err:%v\n", err)
	}

}
