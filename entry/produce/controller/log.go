package controller

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"go-vrs/Model"
	"go-vrs/config"
	"go-vrs/lib/logger"
	"net/http"
	"time"
)

var logs = logger.New()

func LogFile (c *gin.Context)  {
	logs.Println("12312")
	logs.Warningln("test")

	c.String(http.StatusOK, "ok")
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

// 日志推送接收接口
func LogPush (c *gin.Context)  {
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