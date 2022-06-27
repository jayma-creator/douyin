package util

import (
	"crypto/md5"
	"fmt"
)

func GetMD5(str string) string {
	data := []byte(str)
	strMD5 := fmt.Sprintf("%x", md5.Sum(data))
	return strMD5
}
