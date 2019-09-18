package middleware

import (
	"github.com/gin-gonic/gin"
	rscors "github.com/rs/cors/wrapper/gin"
)

var (
	// DefaultAllowedOrigins are the default allowed origins for the api, allowing for access
	// via both of our internet uplinks when using the web interface.
	DefaultAllowedOrigins = []string{"https://temporal.cloud", "https://backup.temporal.cloud"}
)

// CORSMiddleware is used to load our CORS handling logic
func CORSMiddleware(devMode bool, allowedOrigins []string) gin.HandlerFunc {
	opts := rscors.Options{}
	if devMode {
		opts.AllowedOrigins = []string{"*"}
		opts.AllowCredentials = true
	} else {
		opts.AllowedOrigins = allowedOrigins
	}
	opts.AllowedMethods = []string{"GET", "POST", "OPTIONS", "DELETE", "PUT"}
	opts.AllowedHeaders = []string{"cache-control", "Authorization", "Content-Type", "X-Request-ID"}
	opts.OptionsPassthrough = true
	return rscors.New(opts)
}
