package util

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
)

func GetCommentCache(videoId string) (commentList []common.Comment, err error) {
	//从连接池当中获取链接
	conn := dao.Pool.Get()
	//先查看redis中是否有数据
	defer conn.Close()
	//redis读取缓存
	rebytes, err := redis.Bytes(conn.Do("get", fmt.Sprintf("commentList%v", videoId)))
	if err != nil {
		logrus.Infof("读取commentList%v缓存失败,err:%v", videoId, err)
	}
	//进行gob序列化
	reader := bytes.NewReader(rebytes)
	dec := gob.NewDecoder(reader)
	err = dec.Decode(&commentList)
	return
}

func GetPublishListCache(userId string) (videoList []common.Video, err error) {
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

func GetFavoriteListCache(userId string) (videoList []common.Video, err error) {
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

func GetFollowListCache(userId string) (followList []common.User, err error) {
	//从连接池当中获取链接
	conn := dao.Pool.Get()
	//先查看redis中是否有数据
	defer conn.Close()
	//redis读取缓存
	rebytes, err := redis.Bytes(conn.Do("get", fmt.Sprintf("followList%v", userId)))
	if err != nil {
		logrus.Infof("读取followList%v缓存失败,err:%v", userId, err)
	}
	//进行gob序列化
	reader := bytes.NewReader(rebytes)
	dec := gob.NewDecoder(reader)
	err = dec.Decode(&followList)
	return
}

func GetFanListCache(userId string) (fanList []common.User, err error) {
	//从连接池当中获取链接
	conn := dao.Pool.Get()
	//先查看redis中是否有数据
	defer conn.Close()
	//redis读取缓存
	rebytes, err := redis.Bytes(conn.Do("get", fmt.Sprintf("fanList%v", userId)))
	if err != nil {
		logrus.Infof("读取fanList%v缓存失败,err:%v", userId, err)
	}
	//进行gob序列化
	reader := bytes.NewReader(rebytes)
	dec := gob.NewDecoder(reader)
	err = dec.Decode(&fanList)
	return
}

func GetUserCache(username string) (user common.User, err error) {
	//从连接池当中获取链接
	conn := dao.Pool.Get()
	//先查看redis中是否有数据
	defer conn.Close()
	//redis读取缓存
	rebytes, err := redis.Bytes(conn.Do("get", fmt.Sprintf("user%v", username)))
	if err != nil {
		logrus.Infof("读取user%v缓存失败,err:%v", username, err)
	}
	//进行gob序列化
	reader := bytes.NewReader(rebytes)
	dec := gob.NewDecoder(reader)
	err = dec.Decode(&user)
	return
}

//设置缓存
func SetRedisCache(key string, data interface{}) (err error) {
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
	time := 10 * 60 * 60 //10小时
	_, err = conn.Do("setex", key, time, buffer.Bytes())
	if err != nil {
		logrus.Infof("写入%s缓存失败,err:%v", key, err)
	}
	return
}

//删除缓存
func DelCache(key string) (err error) {
	conn := dao.Pool.Get()
	defer conn.Close()
	_, err = conn.Do("del", key)
	if err != nil {
		logrus.Infof("删除%s缓存失败,err:%v", key, err)
	}
	return
}

func SetNull(key string) (err error) {
	//缓存到redis
	conn := dao.Pool.Get()
	defer conn.Close()

	//redis缓存数据
	time := 5 * 60 //单位秒
	_, err = conn.Do("setex", key, time, "")
	if err != nil {
		logrus.Infof("缓存空值到%s失败,err:%v", key, err)
	}
	return
}

func RefreshToken(token string) (err error) {
	conn := dao.Pool.Get()
	defer conn.Close()
	time := 1 * 60 * 60 * 10 //单位秒
	_, err = conn.Do("setex", token, time, 5)
	if err != nil {
		logrus.Error("刷新token失败", err)
	}
	return
}
