package util

import (
	"bytes"
	"context"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

//上传视频到七牛云
func UpLoadQiniuCloud(fileName string) error {
	accessKey := "XXXXXXXXXXXXXXXXX"
	secretKey := "XXXXXXXXXXXXXXXXX"
	bucket := "douyin123456"
	mac := qbox.NewMac(accessKey, secretKey)
	path := filepath.Join("./", fileName)
	a, err := os.ReadFile(path)
	if err != nil {
		logrus.Error(err)
		return err
	}
	//defer a.Close()
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	upToken := putPolicy.UploadToken(mac)
	// 配置参数
	cfg := storage.Config{
		Zone:          &storage.ZoneHuanan, // 华南区
		UseCdnDomains: false,
		UseHTTPS:      false, // 非https
	}
	formUploader := storage.NewFormUploader(&cfg)

	ret := storage.PutRet{}        // 上传后返回的结果
	putExtra := storage.PutExtra{} // 额外参数

	//key为上传的文件名
	key := fileName // 上传路径，如果当前目录中已存在相同文件，则返回上传失败错误
	err = formUploader.Put(context.Background(), &ret, upToken, key, bytes.NewReader(a), int64(len(a)), &putExtra)
	return err
}
