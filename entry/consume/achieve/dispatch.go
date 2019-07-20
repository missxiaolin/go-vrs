package achieve

import (
	"github.com/go-redis/redis"
	"go-vrs/config"
	"go-vrs/consume"
	"go-vrs/lib/mong"
	"log"
)

var redisClient *redis.Client
var mongoClient *mong.Mgo

func init() {

	var err error
	cfg := config.Cfg

	// 初始化redis
	if redisClient == nil {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     cfg.Redis.Addr,
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})
	}

	// 初始化mongodb
	if mongoClient == nil {
		mgoCfg := &mong.Config{
			Host:       cfg.Mongodb.Host,
			Database:   cfg.Mongodb.Database,
			Collection: cfg.Mongodb.Collection,
		}

		mongoClient, err = mgoCfg.NewMgo()
		if err != nil {
			log.Fatal(err)
		}
		if err := mongoClient.Connect(); err != nil {
			log.Fatal(err)
		}
	}
}

// 创建数据调度器
func NewDispatch() *consume.Dispatch {

	redisAch := &RedisAch{
		Client: redisClient,
		Num: config.Cfg.Consume.BatchNum,
	}
	mongoAch := &MongoAch{
		Mgo: mongoClient,
	}

	return consume.NewDispatch(
		consume.NewCollection(redisAch),
		consume.NewBuilder(redisAch),
		consume.NewTransport(mongoAch),
	)
}