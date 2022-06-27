package util

import (
	logr "github.com/sirupsen/logrus"
	"time"
)

func InitLogRecord() (err error) {
	//now := getTime()
	logr.SetReportCaller(true)                                                     //设置文件和行数
	logr.SetFormatter(&logr.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05"}) //json格式,时间格式
	logr.SetLevel(logr.ErrorLevel)                                                 //级别
	//file, err := os.OpenFile(now+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	//if err != nil {
	//	logr.Info("Failed to log to file, using default stderr")
	//}
	//logr.SetOutput(file)
	logr.Info("init log done")
	return
}

func getTime() string {
	return time.Now().Format("2006-01-02")
}
