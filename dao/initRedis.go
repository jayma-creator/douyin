package dao

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/setting"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"time"
)

var Pool *redis.Pool

func InitRedisPool(cfg *setting.RedisConfig) (err error) {
	Pool = &redis.Pool{
		MaxIdle:     100,               //最大空闲链接数，表示即使没有redis链接事依然可以保持N个空闲链接，而不被清除
		MaxActive:   200,               //最大激活连接数，表示同时最多有多少个链接
		IdleTimeout: 240 * time.Second, //最大空闲链接等待时间，超过此时间，空闲将被关闭
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Address, cfg.Port))
			if err != nil {
				logrus.Error("连接redis失败", err)
				return nil, err
			}
			return c, err
		}}
	return
}
