package routes

import (
	"go-vrs/bootstrap"
	"go-vrs/entry/produce/controller"
	"go-vrs/entry/produce/middleware"
)

func Configure(b *bootstrap.Bootstrapper)  {
	user := b.Group("user", middleware.TokenVerify())
	{
		user.GET("/info", controller.UserInfo)
	}

	b.GET("/", controller.Index)
	b.GET("/log/file", controller.LogFile)
	b.GET("/redis/send", controller.RedisSend)
	b.POST("/push", controller.LogPush)
}
