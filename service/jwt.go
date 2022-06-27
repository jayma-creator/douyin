package service

import (
	"errors"
	"fmt"
	"github.com/RaymondCode/simple-demo/dao"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"time"
)

type MyClaims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.RegisteredClaims
}

// MySecret 定义Secret
var MySecret = []byte("ma")

// GetToken 生成JWT
func GetToken(username string, password string) (tokenString string, err error) {
	// 创建一个我们自己的声明
	claim := MyClaims{
		Username: username, // 自定义字段
		Password: password,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * time.Duration(1))), // 过期时间24小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ma", // 签发人
		}}
	// 使用指定的签名方法创建签名对象,HS256算法
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	tokenString, err = token.SignedString(MySecret)
	return
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{},
		func(token *jwt.Token) (i interface{}, err error) {
			return MySecret, nil
		})
	if err != nil {
		ve, ok := err.(*jwt.ValidationError)
		if ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				logrus.Error("that is not even a token")
				return nil, err
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				logrus.Error("token is expired")
				return nil, err
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				logrus.Error("token not active yet")
				return nil, err
			} else {
				logrus.Error("could not handle this token")
				return nil, err
			}
		}
	}
	claims, ok := token.Claims.(*MyClaims)
	if ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func CheckToken(token string) (User, bool, error) {
	user := User{}
	claims, err := ParseToken(token)
	if err != nil {
		logrus.Error(err)
		return user, false, err
	}
	err = dao.DB.Where("name = ? and password = ?", claims.Username, claims.Password).Find(&user).Count(&count).Error
	fmt.Println(claims.Username, claims.Password)
	if err != nil {
		logrus.Error("token is invalid", err)
		return user, false, err
	}
	if count == 0 {
		logrus.Error("token已过期", err)
		return user, false, err
	}
	return user, true, err
}
