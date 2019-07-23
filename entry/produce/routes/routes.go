package routes

import (
	"go-vrs/bootstrap"
	"go-vrs/entry/produce/controller"
)

func Configure(b *bootstrap.Bootstrapper)  {
	b.GET("/", controller.Index)
	b.GET("/log/file", controller.LogFile)
	b.GET("/redis/send", controller.RedisSend)
	b.POST("/push", controller.LogPush)
}
