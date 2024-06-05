package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				logrus.WithFields(logrus.Fields{
					"request": string(httpRequest),
					"err":     err,
				}).Error("panic recovered")
				path := c.Request.URL.Path
				msg := fmt.Sprintf("服务器出现异常panic, path: %s，err: %v", path, err)
				// 通知
				logrus.Errorln(msg)
				//
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
