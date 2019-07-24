package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type result struct {
	Success      bool
	Data         interface{}
	ErrorCode    uint32
	ErrorMessage string
}

func TokenVerify () gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Query("token") == "" {
			c.AbortWithStatusJSON(http.StatusForbidden, result{
				Success: true,
				ErrorCode: 1000,
				ErrorMessage: "not token",
			})
		}

		c.Next()
	}
}