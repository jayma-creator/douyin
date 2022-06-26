package service

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
)

func getCommentCache() (commentList []Comment, err error) {
	//从连接池当中获取链接
	conn := dao.Pool.Get()
	//先查看redis中是否有数据
	defer conn.Close()
	//redis读取缓存
	rebytes, err := redis.Bytes(conn.Do("get", "commentList"))
	if err != nil {
		logrus.Infof("读取commentList缓存失败,err:%v", err)
	}
	//进行gob序列化
	reader := bytes.NewReader(rebytes)
	dec := gob.NewDecoder(reader)
	err = dec.Decode(&commentList)
	return
}

func getPublishListCache(userId string) (videoList []Video, err error) {
	//从连接池当中获取链接
	conn := dao.Pool.Get()
	//先查看redis中是否有数据
	defer conn.Close()
	//redis读取缓存
	rebytes, err := redis.Bytes(conn.Do("get", fmt.Sprintf("publishList%v", userId)))
	if err != nil {
		logrus.Infof("读取videoList%v缓存失败,err:%v", userId, err)
	}
	//进行gob序列化
	reader := bytes.NewReader(rebytes)
	dec := gob.NewDecoder(reader)
	err = dec.Decode(&videoList)
	return
}

func getFavoriteListCache(userId string) (videoList []Video, err error) {
	//从连接池当中获取链接
	conn := dao.Pool.Get()
	//先查看redis中是否有数据
	defer conn.Close()
	//redis读取缓存
	rebytes, err := redis.Bytes(conn.Do("get", fmt.Sprintf("favoriteList%v", userId)))
	if err != nil {
		logrus.Infof("读取favoriteList%v缓存失败,err:%v", userId, err)
	}
	//进行gob序列化
	reader := bytes.NewReader(rebytes)
	dec := gob.NewDecoder(reader)
	err = dec.Decode(&videoList)
	return
}

func setRedisCache(key string, data interface{}) (err error) {
	//缓存到redis
	conn := dao.Pool.Get()
	defer conn.Close()

	//将数据进行gob序列化
	var buffer bytes.Buffer
	ecoder := gob.NewEncoder(&buffer)
	err = ecoder.Encode(data)
	if err != nil {
		logrus.Error(err)
		return
	}
	//redis缓存数据
	_, err = conn.Do("set", key, buffer.Bytes())
	conn.Do("expire", key, 1*60*60)
	if err != nil {
		logrus.Infof("写入%s缓存失败,err:%v", key, err)
	}

	return
}

func delCache(key string) (err error) {
	conn := dao.Pool.Get()
	defer conn.Close()
	_, err = conn.Do("del", key)
	if err != nil {
		logrus.Infof("删除%s缓存失败,err:%v", key, err)
	}
	return
}
