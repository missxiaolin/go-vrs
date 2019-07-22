package routes

import (
	"go-vrs/bootstrap"
	"go-vrs/entry/produce/controller"
)

func Configure(b *bootstrap.Bootstrapper)  {
	b.GET("/", controller.Index)
}
