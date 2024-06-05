package middleware

import (
	"LittleVideo/app/store"
	"LittleVideo/config"
	"LittleVideo/def"
	"LittleVideo/util"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.IsDev() {
			return
		}

		uid := c.PostForm("uid")
		sid := c.PostForm("sid")

		iUid, _ := strconv.Atoi(uid)
		if iUid <= 0 || !store.MC.IsSessionValid(sid, iUid) {
			logrus.WithFields(logrus.Fields{
				"uid": uid,
				"sid": sid,
			}).Warn("未登录")
			c.AbortWithStatusJSON(http.StatusOK, util.Pack(def.CodeNeedLogin, "未登录", nil))
			return
		}
		c.Next()
	}
}

func CheckAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.IsDev() {
			return
		}
		if config.IsLocalDev() {
			return
		}
		accessToken := c.PostForm("access_token")
		if len(accessToken) == 0 {
			return
		}

		//hashSign, _ := hex.DecodeString(crypto.Hmac(appKey, sortedQueryStr))
		//serverAccessToken := base64.StdEncoding.EncodeToString(hashSign)
		serverAccessToken := accessToken

		if strings.Compare(serverAccessToken, accessToken) != 0 {
			logrus.WithFields(logrus.Fields{
				"uri":               c.FullPath(),
				"accessToken":       accessToken,
				"serverAccessToken": serverAccessToken,
				"sortedQueryStr":    "",
			}).Error("签名不正确")
			c.AbortWithStatusJSON(http.StatusOK, util.Pack(def.CodeSignError, "签名不正确", nil))
			return
		}
		c.Next()
	}
}
