package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/zctod/go-tool/common/util_server"
	"go-vrs/Model"
	"go-vrs/bootstrap"
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

func newApp() *bootstrap.Bootstrapper {
	// 初始化应用
	app := bootstrap.New("日志收集", "xiaolin")
	app.Bootstrap()
	//app.Configure(identity.Configure, routes.Configure)

	return app
}

func main()  {

	app := newApp()


	app.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello produce!")
	})

	startServer(app)




	//g := gin.Default()
	//
	//// 程序测试
	//g.GET("/", func(c *gin.Context) {
	//	c.String(http.StatusOK, "Hello produce!")
	//})
	//
	//g.GET("/log/file", func(c *gin.Context) {
	//	logs.Println("12312")
	//	logs.Warningln("test")
	//
	//	c.String(http.StatusOK, "ok")
	//})
	//
	//// 协程池测试
	//g.GET("/redis/send", func(c *gin.Context) {
	//	for i := 0; i < 50000; i++ {
	//		_ = p.Submit(func() {
	//			data := Model.Data{
	//				Uid:        0,
	//				IP:         "127.0.0.1",
	//				Content:    "Hello World!",
	//				CreateTime: time.Now().UnixNano() / 1e3, // 微秒
	//			}
	//
	//			b, err := json.Marshal(data)
	//			if err != nil {
	//				log.Fatal(err)
	//			}
	//			redix.LPush(config.Cfg.Redis.ListKey, b)
	//		})
	//	}
	//	c.String(http.StatusOK, "ok")
	//})
	//
	//// 日志推送接收接口
	//g.POST("/push", LogPush)
	//
	//startServer(g)
}

func startServer (b *bootstrap.Bootstrapper)  {
	server := &http.Server{
		Addr:           ":" + config.Cfg.Produce.Port,
		Handler:        b,
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

// 检查redis list占用内存
func checkRedisMemory() error {

	cmd := redix.MemoryUsage(config.Cfg.Redis.ListKey)
	size, err := cmd.Result()
	if err != nil {
		return err
	}

	if size >= config.Cfg.Redis.MaxMemory {
		return errors.New("redis memory usage reaches the upper limit")
	}
	return nil
}

func LogPush(c *gin.Context) {

	if err := checkRedisMemory(); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, err.Error())
		return
	}

	_ = c.Request.ParseForm()

	jsonData := c.PostForm("logData")
	if jsonData == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, "内容为空")
		return
	}

	var data Model.Data
	if err := json.Unmarshal([]byte(jsonData), &data); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, "内容字段异常")
		return
	}

	if data.Uid == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, "用户id异常")
		return
	}
	if data.ReceiveTime == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, "接收时间异常")
		return
	}

	data.CreateTime = time.Now().UnixNano() / 1e3

	go func() {
		b, err := json.Marshal(data)
		if err != nil {
			logs.Fatal(err)
		}
		redix.LPush(config.Cfg.Redis.ListKey, b)
	}()

	c.JSON(http.StatusOK, "ok")
}
