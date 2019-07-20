package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type configure struct {
	Redis   RedisConfig   `yaml:"redis"`
	Mongodb MongodbConfig `yaml:"mongodb"`
	Produce Produce       `yaml:"produce"`
	Consume Consume       `yaml:"consume"`
}

type RedisConfig struct {
	Addr      string `yaml:"addr"`       // 地址
	Password  string `yaml:"password"`   // 密码
	DB        int    `yaml:"db"`         // 库名
	MaxMemory int64  `yaml:"max_memory"` // 最大内存
	ListKey   string `yaml:"list_key"`   // list键名
}

type MongodbConfig struct {
	Host       string `yaml:"host"`       // 地址
	Username   string `yaml:"username"`   // 用户名
	Password   string `yaml:"password"`   // 密码
	Database   string `yaml:"database"`   // 数据库名称
	Collection string `yaml:"collection"` // 数据表名称
}

type Produce struct {
	Port            string `yaml:"port"`              // 监听端口
	PoolMaxCapacity int    `yaml:"pool_max_capacity"` // 协程池最大容量
}

type Consume struct {
	BatchNum        int   `yaml:"batch_num"`         // 批量处理数量
	SleepInterval   int64 `yaml:"sleep_interval"`    // 休眠间隔
	PoolMaxCapacity int   `yaml:"pool_max_capacity"` // 协程池最大容量
	TriggerCount    int   `yaml:"trigger_count"`     // 协程触发分裂的计数
}

var Cfg = &configure{}

func init() {

	var err error
	b, err := ioutil.ReadFile(CONFIG_PATH)
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(b, Cfg); err != nil {
		panic(err)
	}

	if Cfg.Redis.MaxMemory == 0 {
		Cfg.Redis.MaxMemory = REDIS_MAX_MEMORY
	}
	if Cfg.Redis.ListKey == "" {
		Cfg.Redis.ListKey = REDIS_LIST_KEY
	}
	if Cfg.Produce.Port == "" {
		Cfg.Produce.Port = PRODUCE_PORT
	}
	if Cfg.Consume.SleepInterval == 0 {
		Cfg.Consume.SleepInterval = CONSUME_SLEEP_INTERVAL
	}
}