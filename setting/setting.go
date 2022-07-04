package setting

import "gopkg.in/ini.v1"

var Conf = new(Config)

// Config 应用程序配置
type Config struct {
	*GinConfig   `ini:"gin"`
	*MySQLConfig `ini:"mysql"`
	*RedisConfig `ini:"redis"`
	*QiNiuCloud  `ini:"qiniucloud"`
}

// GinConfig 配置
type GinConfig struct {
	Release bool `ini:"release"`
	Port    int  `ini:"port"`
}

// MySQLConfig 数据库配置
type MySQLConfig struct {
	User     string `ini:"user"`
	Password string `ini:"password"`
	DB       string `ini:"db"`
	Address  string `ini:"address"`
	Port     int    `ini:"port"`
}

// RedisConfig 数据库配置
type RedisConfig struct {
	Address string `ini:"address"`
	Port    int    `ini:"port"`
}

type QiNiuCloud struct {
	AccessKey string `ini:"access_key"`
	SecretKey string `ini:"secret_key"`
}

func Init(file string) error {
	return ini.MapTo(Conf, file)
}
