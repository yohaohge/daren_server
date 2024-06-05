package route

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"server.com/daren/app/store"
	"server.com/daren/def"
	"server.com/daren/middleware"
	"server.com/daren/util"
)

func setupUserRouter(r *gin.Engine) {
	userApi := r.Group("user", middleware.CheckLogin())
	{
		userApi.POST("/history", History)
		userApi.POST("/collect", Collect)
		userApi.POST("/refresh", Refresh)
		userApi.POST("/report_event", ReportEvent)
		userApi.POST("/user_info", GetUserInfo)
		userApi.POST("/use_cdkey", UseCdKey)
	}
}

func History(c *gin.Context) {
	c.JSON(http.StatusOK, util.Pack(def.CodeSuccess, "ok", nil))
}

func Collect(c *gin.Context) {
	var reqParams struct {
		Sid string `form:"sid" binding:"required"`
		Uid string `form:"uid" binding:"required"`
		Vid int    `form:"vid" binding:"required"`
		Num int    `form:"Num" binding:"required"`
	}
	if err := c.ShouldBind(&reqParams); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, util.Pack(def.CodeParamError, "参数不正确", err.Error()))
		return
	}

	// 逻辑

	c.JSON(http.StatusOK, util.Pack(def.CodeSuccess, "ok", nil))
}

func Refresh(c *gin.Context) {
	var reqParams struct {
		Sid string `form:"sid" binding:"required"`
		Uid string `form:"uid" binding:"required"`
	}

	if err := c.ShouldBind(&reqParams); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, util.Pack(def.CodeParamError, "参数不正确", err.Error()))
		return
	}

	if store.MC.UpdateSession(reqParams.Sid, reqParams.Uid) {
		c.JSON(http.StatusOK, util.Pack(def.CodeSuccess, "ok", nil))
	} else {
		c.JSON(http.StatusOK, util.Pack(def.CodeFailed, "failed", nil))
	}
}

func ReportEvent(c *gin.Context) {
	var reqParams struct {
		Sid     string `form:"sid" binding:"required"`
		Uid     string `form:"uid" binding:"required"`
		EvtId   string `form:"event_id" binding:"required"`
		EvtData string `form:"event_data" binding:"required"`
	}

	if err := c.ShouldBind(&reqParams); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, util.Pack(def.CodeParamError, "参数不正确", err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Pack(def.CodeSuccess, "ok", nil))
}

func GetUserInfo(c *gin.Context) {
	var reqParams struct {
		Sid string `form:"sid" binding:"required"`
		Uid int    `form:"uid" binding:"required"`
	}

	if err := c.ShouldBind(&reqParams); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, util.Pack(def.CodeParamError, "参数不正确", err.Error()))
		return
	}

	user, err := store.MC.GetUserByUid(reqParams.Uid)
	if err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, util.Pack(def.CodeFailed, "内部错误", err.Error()))
		return
	}
	if user == nil {
		c.JSON(http.StatusOK, util.Pack(def.CodeUserNotExist, "用户不存在", err.Error()))
		return
	}
	c.JSON(http.StatusOK, util.Pack(def.CodeSuccess, "ok", user))
}

func UseCdKey(c *gin.Context) {
	var reqParams struct {
		Uid   int    `form:"uid" binding:"required"`
		Cdkey string `form:"cdkey" binding:"required"`
	}

	if err := c.ShouldBind(&reqParams); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusOK, util.Pack(def.CodeParamError, "参数不正确", err.Error()))
		return
	}

	if len(reqParams.Cdkey) < 2 {
		c.JSON(http.StatusOK, util.Pack(def.CodeParamError, "参数不正确", nil))
		return
	}

	user, err := store.MC.GetUserByUid(reqParams.Uid)
	if err != nil {
		logrus.Errorln(err)
		c.JSON(http.StatusOK, util.Pack(def.CodeFailed, "兑换失败", err.Error()))
		return
	}
	if user == nil {
		c.JSON(http.StatusOK, util.Pack(def.CodeUserNotExist, "用户不存在", err.Error()))
		return
	}
	// 校验CDKEY
	b, cdkeyInfo := store.MC.GetCDKEYInfo(reqParams.Cdkey)
	if !b {
		c.JSON(http.StatusOK, util.Pack(def.CodeInvalidCDKEY, "兑换码无效或已经使用", nil))
		return
	}
	// 校验玩家
	isUsed := store.MC.IsUserUsed(reqParams.Uid, reqParams.Cdkey, cdkeyInfo.CdkeyType)
	if isUsed {
		c.JSON(http.StatusOK, util.Pack(def.CodeFailed, "玩家已经使用", nil))
		return
	}
	// 使用CDKEY
	if !store.UseCDKEY(reqParams.Uid, reqParams.Cdkey) {
		c.JSON(http.StatusOK, util.Pack(def.CodeInvalidCDKEY, "兑换失败,请稍后重试", nil))
		return
	}
	rewards := []def.ItemOpe{}
	if err = json.Unmarshal([]byte(cdkeyInfo.Items), &rewards); err != nil {
		logrus.WithFields(logrus.Fields{"uid": reqParams.Uid, "CDKEY": reqParams.Cdkey}).Errorln(err, "items parse error")
	}
	store.AddUserRewards(reqParams.Uid, rewards)

	c.JSON(http.StatusOK, util.Pack(def.CodeSuccess, "ok", rewards))
}
