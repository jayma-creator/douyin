package util

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/sirupsen/logrus"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"os"
)

//截图做封面
func exampleReadFrameAsJpeg(inFileName string, frameNum int) io.Reader {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(inFileName).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		logrus.Error("获取封面失败", err)
		return nil
	}
	return buf
}

//保存截图
func GetSnapShot(snapShotName string, videoFilePath string) {
	reader := exampleReadFrameAsJpeg(videoFilePath, 48)
	img, err := imaging.Decode(reader)
	if err != nil {
		logrus.Error("保存截图失败", err)
		return
	}
	err = imaging.Save(img, "./public/"+snapShotName)
	if err != nil {
		logrus.Error("保存截图失败", err)
		return
	}
	return
}
