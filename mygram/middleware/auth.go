package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/refandas/scalable-web-service/mygram/helper"
	"net/http"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		userData, err := helper.VerifyToken(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid token",
			})
			return
		}

		c.Set("userData", userData)
		c.Next()
	}
}
