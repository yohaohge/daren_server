package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server.com/daren/config"
	"server.com/daren/def"
	"server.com/daren/util"
)

func setupGmRouter(r *gin.Engine) {
	backendApi := r.Group("gm")
	{
		backendApi.GET("/reload_config", ReloadConfig)
	}
}

func ReloadConfig(c *gin.Context) {
	var reqParam struct {
		Sig string `form:"sig" binding:"required"`
	}
	//if "127.0.0.1" != c.ClientIP() {
	//	return
	//}
	if err := c.ShouldBindQuery(&reqParam); err != nil {
		c.JSON(http.StatusOK, util.Pack(def.CodeParamError, "参数不正确", err.Error()))
		return
	}
	if reqParam.Sig != "haoju223_com" {
		c.JSON(http.StatusOK, util.Pack(def.CodeFailed, "签名错误", nil))
		return
	}
	config.LoadConfig()
	c.JSON(http.StatusOK, util.Pack(def.CodeSuccess, "ok", nil))
}
