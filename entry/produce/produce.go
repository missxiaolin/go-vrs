package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zctod/go-tool/common/util_server"
	"go-vrs/config"
	"net/http"
	"time"
)

func main()  {
	g := gin.Default()

	// 程序测试
	g.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello produce!")
	})

	startServer(g)
}

func startServer (g *gin.Engine)  {
	server := &http.Server{
		Addr:           ":" + config.PRODUCE_PORT,
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
