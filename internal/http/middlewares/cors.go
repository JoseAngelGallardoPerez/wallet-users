package middlewares

import (
	"github.com/Confialink/wallet-users/internal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CorsMiddleware cors middleware
func CorsMiddleware() gin.HandlerFunc {
	// Retrieve config options.
	conf := config.GetConf()

	corsConfig := cors.DefaultConfig()

	corsConfig.AllowMethods = conf.Cors.Methods
	for _, origin := range conf.Cors.Origins {
		if origin == "*" {
			corsConfig.AllowAllOrigins = true
		}
	}
	if !corsConfig.AllowAllOrigins {
		corsConfig.AllowOrigins = conf.Cors.Origins
	}
	corsConfig.AllowHeaders = conf.Cors.Headers

	return cors.New(corsConfig)
}
