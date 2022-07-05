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
	conn.Do("AUTH", "123456")
	if err != nil {
		logrus.Error("密码错误", err)
	}
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
	conn.Do("AUTH", "123456")
	if err != nil {
		logrus.Error("密码错误", err)
	}
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
	conn.Do("AUTH", "123456")
	if err != nil {
		logrus.Error("密码错误", err)
	}
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

func GetFeed() (videoList []common.Video, err error) {
	//从连接池当中获取链接
	conn := dao.Pool.Get()
	//先查看redis中是否有数据
	defer conn.Close()
	conn.Do("AUTH", "123456")
	if err != nil {
		logrus.Error("密码错误", err)
	}
	//redis读取缓存
	rebytes, err := redis.Bytes(conn.Do("get", "feed"))
	if err != nil {
		logrus.Info("读取feed%v缓存失败", err)
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
	conn.Do("AUTH", "123456")
	if err != nil {
		logrus.Error("密码错误", err)
	}
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
	conn.Do("AUTH", "123456")
	if err != nil {
		logrus.Error("密码错误", err)
	}
	//redis读取缓存
	rebytes, err := redis.Bytes(conn.Do("get", fmt.Sprintf("fansList%v", userId)))
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
	conn.Do("AUTH", "123456")
	if err != nil {
		logrus.Error("密码错误", err)
	}
	//redis读取缓存
	rebytes, err := redis.Bytes(conn.Do("get", username))
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
	conn.Do("AUTH", "123456")
	if err != nil {
		logrus.Error("密码错误", err)
	}
	//将数据进行gob序列化
	var buffer bytes.Buffer
	ecoder := gob.NewEncoder(&buffer)
	err = ecoder.Encode(data)
	if err != nil {
		logrus.Error(err)
		return
	}
	//加上随机数，防止同时过期造成缓存雪崩
	randNum := rangeRand(30*60, 60*60)
	time := 10*60*60 + randNum //10小时
	//redis缓存数据
	_, err = conn.Do("setex", key, time, buffer.Bytes())
	if err != nil {
		logrus.Infof("写入%s缓存失败,err:%v", key, err)
	}
	return
}

func SetRedisNum(key, value string) {
	conn := dao.Pool.Get()
	defer conn.Close()
	conn.Do("AUTH", "123456")

	randNum := rangeRand(30*60, 60*60)
	time := 10*60*60 + randNum //10小时
	_, err := conn.Do("setex", key, time, value)
	if err != nil {
		logrus.Error("设置缓存失败", err)
	}

}

//删除缓存
func DelCache(key string) (err error) {
	conn := dao.Pool.Get()
	defer conn.Close()
	conn.Do("AUTH", "123456")

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
	conn.Do("AUTH", "123456")

	//redis缓存数据
	time := 10 //单位秒
	_, err = conn.Do("setex", key, time, "")
	if err != nil {
		logrus.Infof("缓存空值到%s失败,err:%v", key, err)
	}
	return
}

func RefreshToken(token string) (err error) {
	conn := dao.Pool.Get()
	defer conn.Close()
	conn.Do("AUTH", "123456")

	time := 1 * 60 * 60 * 10 //单位秒
	_, err = conn.Do("setex", token, time, 5)
	if err != nil {
		logrus.Error("刷新token失败", err)
	}
	return
}

func IsExistCache(key string) (exists int64) {
	conn := dao.Pool.Get()
	defer conn.Close()
	conn.Do("AUTH", "123456")

	exist, err := conn.Do("exists", key)
	if err != nil {
		logrus.Error("查询缓存是否存在失败", err)
	}
	exists = exist.(int64)
	return exists
}

func RedisLock(key string) (isLock bool) {
	conn := dao.Pool.Get()
	defer conn.Close()
	conn.Do("AUTH", "123456")

	redisLockTimeout := 10
	//这里需要redis.String包一下，才能返回redis.ErrNil
	_, err := redis.String(conn.Do("set", key, 1, "ex", redisLockTimeout, "nx"))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		logrus.Error("加锁失败", err)
		return
	}
	return true
}

func RedisUnlock(key string) (err error) {
	conn := dao.Pool.Get()
	defer conn.Close()
	conn.Do("AUTH", "123456")
	if err != nil {
		logrus.Error("密码错误", err)
	}
	_, err = conn.Do("del", key)
	if err != nil {
		logrus.Error("解锁失败", err)
		return
	}
	return
}
