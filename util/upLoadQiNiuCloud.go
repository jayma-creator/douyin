package util

import (
	"context"
	"github.com/RaymondCode/simple-demo/setting"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/sirupsen/logrus"
	"path/filepath"
)

//上传视频到七牛云
func UpLoadQiniuCloud(fileName string) {
	accessKey := setting.Conf.QiNiuCloud.AccessKey
	secretKey := setting.Conf.QiNiuCloud.SecretKey
	bucket := "douyin123456"
	mac := qbox.NewMac(accessKey, secretKey)
	path := filepath.Join("./", fileName)

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
	resumeUploader := storage.NewResumeUploaderV2(&cfg)

	ret := storage.PutRet{}           // 上传后返回的结果
	putExtra := storage.RputV2Extra{} // 额外参数

	//key为上传的文件名
	key := fileName // 上传路径，如果当前目录中已存在相同文件，则返回上传失败错误
	err := resumeUploader.PutFile(context.Background(), &ret, upToken, key, path, &putExtra)
	if err != nil {
		logrus.Error(err)
		return
	}
	return
}
