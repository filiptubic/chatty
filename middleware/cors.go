package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowCredentials: true,
		AllowAllOrigins:  true,
		AllowHeaders:     []string{"Origin", "Authorization"},
	})
}
