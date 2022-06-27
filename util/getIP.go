package util

import (
	"github.com/sirupsen/logrus"
	"net"
	"strings"
)

//获取当前主机IP
func GetIp() string {
	conn, err := net.Dial("udp", "google.com:80")
	if err != nil {
		logrus.Error("获取ip失败", err)
		return ""
	}
	defer conn.Close()
	ip := strings.Split(conn.LocalAddr().String(), ":")[0]
	return ip
}
