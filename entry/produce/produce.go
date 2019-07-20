package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/zctod/go-tool/common/util_server"
	"go-vrs/Model"
	"go-vrs/config"
	"go-vrs/lib/logger"
	"go-vrs/lib/pool"
	"log"
	"net/http"
	"time"
)

var logs = logger.New()
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

func main()  {
	g := gin.Default()

	// 程序测试
	g.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello produce!")
	})

	g.GET("/log/file", func(c *gin.Context) {
		logs.Println("12312")
		logs.Warningln("test")

		c.String(http.StatusOK, "ok")
	})

	// 协程池测试
	g.GET("/redis/send", func(c *gin.Context) {
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
	})

	startServer(g)
}

func startServer (g *gin.Engine)  {
	server := &http.Server{
		Addr:           ":" + config.Cfg.Produce.Port,
		Handler:        g,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	// 平滑退出，先结束所有在执行的任务
	util_server.GracefulExitWeb(server)
}
