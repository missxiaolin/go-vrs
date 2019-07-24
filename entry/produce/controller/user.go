package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func UserInfo(c *gin.Context)  {
	c.JSON(http.StatusOK, "userInfo")
}
