package util

import (
	"crypto/md5"
	password1 "github.com/anaskhan96/go-password-encoder"
)

//type Options struct {
//	SaltLen      int              //用户生成的长度，默认256
//	Iterations   int              //PBKDF2函数中的迭代次数，默认10000
//	KeyLen       int              //BKDF2函数中编码密钥的长度，默认512
//	HashFunction func() hash.Hash //使用的哈希算法，默认sha512
//}

func GetMD5WithSalted(password string) (string, string) {
	// 方式一：使用默认选项
	//salt, encodedPwd := password.Encode("generic password", nil)
	//check := password.Verify("generic password", salt, encodedPwd, nil)
	//fmt.Println(check) // true

	// 方式二：使用自定义选项
	options := &password1.Options{10, 5000, 25, md5.New}
	salt, encodedPwd := password1.Encode(password, options)
	//check := password.Verify(str, salt, encodedPwd, options)
	return encodedPwd, salt
}

func VerifyPassword(password, encodePwd, salt string) bool {
	options := &password1.Options{10, 5000, 25, md5.New}
	check := password1.Verify(password, salt, encodePwd, options)
	return check
}
