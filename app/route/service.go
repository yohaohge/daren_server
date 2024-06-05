package route

import (
	"LittleVideo/middleware"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {

	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(middleware.Cors())
	setupLoginRouter(r)
	setupUserRouter(r)
	setupGmRouter(r)
}
