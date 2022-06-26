package dao

import (
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

var Pool *redis.Pool

func InitRedisPool() {
	Pool = &redis.Pool{
		MaxIdle:     100,               //最大空闲链接数，表示即使没有redis链接事依然可以保持N个空闲链接，而不被清除
		MaxActive:   200,               //最大激活连接数，表示同时最多有多少个链接
		IdleTimeout: 240 * time.Second, //最大空闲链接等待时间，超过此时间，空闲将被关闭
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "localhost:6379")
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}
}
