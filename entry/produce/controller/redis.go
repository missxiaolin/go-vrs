package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go-vrs/Model"
	"go-vrs/config"
	"go-vrs/lib/pool"
	"log"
	"net/http"
	"time"
)

var p *pool.Pool
var redix *redis.Client

func init()  {
	if redix == nil {
		redix = redis.NewClient(&redis.Options{
			Addr:     config.Cfg.Redis.Addr,
			Password: config.Cfg.Redis.Password,
			DB:       config.Cfg.Redis.DB,
		})
	}
	if p == nil {
		var err error
		p, err = pool.NewPool(config.Cfg.Produce.PoolMaxCapacity)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// 协程池测试
func RedisSend(c *gin.Context) {
	for i := 0; i < 50000; i++ {
		_ = p.Submit(func() {
			data := Model.Data{
				Uid:        0,
				IP:         "127.0.0.1",
				Content:    "Hello World!",
				CreateTime: time.Now().UnixNano() / 1e3, // 微秒
			}

			b, err := json.Marshal(data)
			if err != nil {
				log.Fatal(err)
			}
			redix.LPush(config.Cfg.Redis.ListKey, b)
		})
	}
	c.String(http.StatusOK, "ok")
}
