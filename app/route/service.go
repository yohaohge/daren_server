package route

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"server.com/daren/middleware"
)

func SetupRouter(r *gin.Engine) {

	r.Use(gin.Recovery())
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(middleware.Cors())
	setupLoginRouter(r)
	setupUserRouter(r)
	setupGmRouter(r)
}
