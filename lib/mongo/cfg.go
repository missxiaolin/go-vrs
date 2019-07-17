package mongo

import "time"

type Config struct {
	Host       string        // 连接地址
	Username   string        // 用户名
	Password   string        // 密码
	Database   string        // 数据库名称
	Collection string        // 数据库集合名称
	Timeout    time.Duration // 请求超时时间
}

func (c *Config) NewMgo() (*Mgo, error) {
	return &Mgo{cfg: *c}, nil
}