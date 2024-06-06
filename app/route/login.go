package route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"server.com/daren/app/store"
	"server.com/daren/config"
	"server.com/daren/def"
	"server.com/daren/middleware"
	"server.com/daren/pkg/crypto"
	"server.com/daren/util"
	"strconv"
	"strings"
	"time"
)

type LoginRsp struct {
	Uid        int    `json:"uid"`
	Name       string `json:"name"`
	VipTime    int64  `json:"vipTime"`
	CreateTime int64  `json:"createTime"`
	Sid        string `json:"sid"`
	VipExpire  bool   `json:"vip_expire"`
}

func setupLoginRouter(r *gin.Engine) {
	loginApi := r.Group("api", middleware.CheckAccessToken())
	{
		loginApi.POST("/login", Login)
		loginApi.POST("/register", Register)
		loginApi.GET("/config", ConfigData)

		loginApi.GET("/index", Index)
		loginApi.GET("/v_list", VDataList)
		loginApi.POST("/v_detail", VDataDetail)
		loginApi.POST("/v_episode", VDataEpisode)
		loginApi.GET("/search", Search)
	}
}

func Login(c *gin.Context) {
	var reqParams struct {
		OpenId    string `form:"open_id" binding:"required"`
		Password  string `form:"password" binding:"required"`
		Version   string `form:"version"`
		HeartBeat bool   `form:"heart_beat"` // 是否心跳检查
		Device    string `form:"device"`     // 唯一设备id，用来限制设备数量
	}

	if err := c.ShouldBind(&reqParams); err != nil {
		logrus.Errorln(err)
		c.JSON(http.StatusOK, util.Pack(def.CodeParamError, "参数不正确", err.Error()))
		return
	}

	userInfo, err := store.MC.GetUserByOpenId(reqParams.OpenId)
	if err != nil {
		logrus.Errorln(err)
		c.JSON(http.StatusOK, util.Pack(def.CodeParamError, "参数错误", nil))
		return
	}
	if userInfo == nil {
		c.JSON(http.StatusOK, util.Pack(def.CodeUserNotExist, "user not exist", nil))
		return
	}
	if strings.Compare(userInfo.Password, crypto.Md5(reqParams.Password)) != 0 {
		c.JSON(http.StatusOK, util.Pack(def.CodePasswordError, "password error", nil))
		return
	}

	// 非心跳的话，挤掉前面的设备 单设备登录
	if reqParams.HeartBeat && reqParams.Device != userInfo.Device {
		c.JSON(http.StatusOK, util.Pack(def.CodePasswordError, "又其他设备登录了该账号", nil))
		return
	}

	if !reqParams.HeartBeat && reqParams.Device != userInfo.Device {
		userInfo.Device = reqParams.Device
		store.MC.SaveUser(userInfo)
	}

	if userInfo.IsVipExpire() && reqParams.HeartBeat {
		c.JSON(http.StatusOK, util.Pack(def.CodePasswordError, "VIP已经过期了", nil))
		return
	}

	resp, err := loginHandler(c, userInfo, reqParams.Password)
	if err != nil {
		logrus.WithFields(logrus.Fields{"uid": userInfo.Uid, "err": err}).Error("login error")
		c.JSON(http.StatusOK, util.Pack(def.CodeFailed, "登录失败"+err.Error(), nil))
		return
	}

	c.JSON(http.StatusOK, util.Pack(def.CodeSuccess, "ok", resp))
}

func loginHandler(c *gin.Context, userInfo *store.User, password string) (LoginRsp, error) {
	var resp LoginRsp
	if c == nil || userInfo == nil {
		logrus.Errorln("nil pointer c or userinfo")
		return resp, errors.New("数据错误")
	}

	//创建session信息
	sid, err := store.MC.SaveSession(userInfo.Uid, c.ClientIP())
	if err != nil {
		logrus.WithFields(logrus.Fields{"uid": userInfo.Uid, "err": err}).Error("session error")
		return resp, errors.New("sid生成失败" + def.LoginHelpMessage)
	}
	if sid == "" {
		return resp, errors.New("sid生成失败" + def.LoginHelpMessage)
	}

	resp.Name = userInfo.Name
	resp.Uid = userInfo.Uid
	resp.VipTime = userInfo.VipTime
	resp.CreateTime = userInfo.CreateTime
	resp.Sid = sid
	resp.VipExpire = userInfo.IsVipExpire()

	return resp, nil
}

func Register(c *gin.Context) {
	var reqParams struct {
		OpenId   string `form:"open_id" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	if err := c.ShouldBind(&reqParams); err != nil {
		logrus.Errorln(err)
		c.JSON(http.StatusOK, util.Pack(def.CodeParamError, "参数不正确", err.Error()))
		return
	}

	userInfo, err := store.MC.GetUserByOpenId(reqParams.OpenId)
	if err != nil {
		logrus.Errorln(err)
		c.JSON(http.StatusOK, util.Pack(def.CodeParamError, "参数错误", nil))
		return
	}
	if userInfo != nil {
		c.JSON(http.StatusOK, util.Pack(def.CodeUserAlreadyExist, "user already exist", nil))
		return
	}

	userInfo, err = newUser(reqParams.OpenId, reqParams.Password)
	if err != nil {
		logrus.WithFields(logrus.Fields{"openId": reqParams.OpenId, "err": err}).Error("register error")
		c.JSON(http.StatusOK, util.Pack(def.CodeFailed, "注册失败", nil))
		return
	}

	userInfo.ClientIP = c.ClientIP()
	_, err = store.MC.AddUser(userInfo)
	if err != nil {
		logrus.WithFields(logrus.Fields{"openId": reqParams.OpenId, "err": err}).Error("register error")
		c.JSON(http.StatusOK, util.Pack(def.CodeFailed, "注册失败", nil))
		return
	}

	c.JSON(http.StatusOK, util.Pack(def.CodeSuccess, "ok", nil))
}

func newUser(openId string, password string) (*store.User, error) {
	uid, err := store.MC.GetNewUid()
	if err != nil {
		return nil, err
	}
	pwdHash := crypto.Md5(password)
	data := ""
	user := &store.User{
		openId,
		int(uid),
		fmt.Sprintf("user_%d", uid),
		0,
		"1.0.0",
		pwdHash,
		data,
		"",
		time.Now().Unix(),
		"",
	}
	return user, nil
}

func Index(c *gin.Context) {

}

func VDataList(c *gin.Context) {
	var reqParams struct {
		PageNo    int `form:"page_no" binding:"required"`
		PageCount int `form:"page_count" binding:"required"`
	}

	if err := c.ShouldBindQuery(&reqParams); err != nil {
		logrus.Errorln(err)
		c.JSON(http.StatusOK, util.Pack(def.CodeParamError, "参数不正确", err.Error()))
		return
	}

	total := store.MC.GetVideListTotal()
	list := store.MC.GetVideoListWithOutData(reqParams.PageNo, reqParams.PageCount)

	// response
	resp := struct {
		Total int                `json:"total"`
		List  []*store.VideoInfo `json:"list"`
	}{
		Total: total,
		List:  list,
	}

	c.JSON(http.StatusOK, util.Pack(def.CodeSuccess, "ok", resp))
}

type VDetailResp struct {
	*store.VideoInfo
	VEpisodeInfo []store.VideoEpisodeInfo `json:"episodes"`
	WatchMin     int                      `json:"w_min"`
	WatchMax     int                      `json:"w_max"`
}

func VDataDetail(c *gin.Context) {
	var reqParams struct {
		Id int `form:"id" binding:"required"`
	}

	if err := c.ShouldBind(&reqParams); err != nil {
		logrus.Errorln(err)
		c.JSON(http.StatusOK, util.Pack(def.CodeParamError, "参数不正确", err.Error()))
		return
	}

	isLogin := false
	isVip := false
	sid := c.PostForm("sid")
	uid := c.PostForm("uid")

	iUid, _ := strconv.Atoi(uid)
	if len(sid) > 0 && iUid > 0 {
		isLogin = store.MC.IsSessionValid(sid, iUid)
	}
	if isLogin {
		user, err := store.MC.GetUserByUid(iUid)
		if err != nil {
			logrus.WithFields(logrus.Fields{"uid": uid}).Error(err)
		} else {
			if user != nil {
				isVip = !user.IsVipExpire()
			}
		}
	}

	resp := dataDetailHandler(reqParams.Id, isLogin, isVip)
	c.JSON(http.StatusOK, util.Pack(def.CodeSuccess, "ok", resp))
}

func dataDetailHandler(id int, isLogin, isVip bool) *VDetailResp {
	resp := &VDetailResp{}

	watchMin := def.NotLoginWatchMin
	watchMax := def.NotLoginWatchMax
	vd := store.GetVideoDetail(id)
	if vd != nil {
		if isLogin {
			watchMin = def.LoginWatchMin
			watchMax = def.LoginWatchMax
		}
		if isLogin && isVip {
			watchMin = def.VipWatchMin
			watchMax = vd.Total
		}

		episodeInfo := make([]store.VideoEpisodeInfo, 0)
		// 第一集详细信息
		val, b := vd.VData[1]
		if b {
			episodeInfo = append(episodeInfo, val)
		}
		// 后面集的基本信息
		for i := 2; i <= vd.Total; i++ {
			// 每集信息封面
			val, b = vd.VData[i]
			if b {
				episodeItem := store.VideoEpisodeInfo{Num: val.Num, SubCover: val.SubCover, PlayUrl: ""}
				if i >= watchMin && i <= watchMax {
					episodeItem.PlayUrl = val.PlayUrl
				}
				episodeInfo = append(episodeInfo, episodeItem)
			}
		}
		resp.VEpisodeInfo = episodeInfo
		resp.VideoInfo = vd.VideoInfo
	}
	resp.WatchMin = watchMin
	resp.WatchMax = watchMax

	return resp
}

func VDataEpisode(c *gin.Context) {
	var reqParams struct {
		Id    int `form:"id" binding:"required"`
		Index int `form:"index" binding:"required"`
	}

	if err := c.ShouldBind(&reqParams); err != nil {
		logrus.Errorln(err)
		c.JSON(http.StatusOK, util.Pack(def.CodeParamError, "参数不正确", err.Error()))
		return
	}

	isLogin := false
	isVip := false
	sid := c.PostForm("sid")
	uid := c.PostForm("uid")
	iUid, _ := strconv.Atoi(uid)
	if len(sid) > 0 && iUid > 0 {
		isLogin = store.MC.IsSessionValid(sid, iUid)
	}
	if isLogin {
		user, err := store.MC.GetUserByUid(iUid)
		if err != nil {
			logrus.WithFields(logrus.Fields{"uid": uid}).Error(err)
		} else {
			if user != nil {
				isVip = !user.IsVipExpire()
			}
		}
	}

	watchMin := def.NotLoginWatchMin
	watchMax := def.NotLoginWatchMax
	if isVip {
		watchMin = def.VipWatchMin
		watchMax = def.VipWatchMax
	} else {
		if isLogin {
			watchMin = def.LoginWatchMin
			watchMax = def.LoginWatchMax
		}
	}
	if reqParams.Index < watchMin || reqParams.Index > watchMax {
		c.JSON(http.StatusOK, util.Pack(def.CodeCanNotWatch, "条件不满足无法观看", nil))
		return
	}
	vd := store.GetVideoDetail(reqParams.Id)
	if vd == nil {
		c.JSON(http.StatusOK, util.Pack(def.CodeNotExist, "资源不存在~", nil))
		return
	}
	val, b := vd.VData[reqParams.Index]
	if !b {
		c.JSON(http.StatusOK, util.Pack(def.CodeNotExist, "该资源不存在~", nil))
		return
	}
	c.JSON(http.StatusOK, util.Pack(def.CodeSuccess, "ok", val))
}

func Search(c *gin.Context) {
	var reqParams struct {
		Name      string `form:"name" binding:"required"`
		PageNo    int    `form:"page_no" binding:"required"`
		PageCount int    `form:"page_count" binding:"required"`
	}

	if err := c.ShouldBindQuery(&reqParams); err != nil {
		logrus.Errorln(err)
		c.JSON(http.StatusOK, util.Pack(def.CodeParamError, "参数不正确", err.Error()))
		return
	}

	total, list := store.MC.SearchVideo(reqParams.Name, reqParams.PageNo, reqParams.PageCount)

	// response
	resp := struct {
		Total int                `json:"total"`
		List  []*store.VideoInfo `json:"list"`
	}{
		Total: total,
		List:  list,
	}

	c.JSON(http.StatusOK, util.Pack(def.CodeSuccess, "ok", resp))
}

func ConfigData(c *gin.Context) {
	resp := config.GetConfig()
	c.JSON(http.StatusOK, util.Pack(def.CodeSuccess, "ok", resp))
}
