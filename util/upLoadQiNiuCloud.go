package util

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"mime/multipart"
)

//上传视频到七牛云
func upLoadQiniuCloud(file *multipart.FileHeader, fileName string) {
	accessKey := "rql9AX6S4i0vl5nR6N0CxEQD08bmOIh7vbwKIr4w"
	secretKey := "HZM9mxeNud7AjWZDlCzlreZRcs7qd4TgxLBjlN5t"
	bucket := "douyin123456"
	mac := qbox.NewMac(accessKey, secretKey)
	src, err := file.Open()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer src.Close()
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
	err = formUploader.Put(context.Background(), &ret, upToken, key, src, file.Size, &putExtra)
}
