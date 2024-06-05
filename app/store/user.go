package store

import (
	"LittleVideo/def"
	"github.com/sirupsen/logrus"
	"time"
)

type User struct {
	OpenId     string      `json:"openId" db:"openId"`
	Uid        int         `json:"uid" db:"uid"`
	Name       string      `json:"name" db:"name"`
	VipTime    int64       `json:"vipTime" db:"vipTime"`
	Version    string      `json:"version" db:"version"`
	Password   string      `json:"password" db:"password"`
	Data       interface{} `json:"data" db:"data"`
	ClientIP   string      `json:"clientIP" db:"clientIP"`
	CreateTime int64       `json:"createTime" db:"createTime"`
}

type Session struct {
	SessionId string `json:"sessionId" db:"sessionId"`
	Uid       int    `json:"uid" db:"uid"`
	Expire    int64  `json:"expire" db:"expire"`
}

type CDKEYInfo struct {
	Cdkey      string `json:"cdkey" db:"cdkey"`
	Num        int    `json:"num" db:"num"`
	CdkeyType  int    `json:"cdkeyType" db:"cdkeyType"`
	Items      string `json:"items" db:"items"`
	CreateTime string `json:"createTime" db:"createTime"`
	ExpireTime string `json:"expireTime" db:"expireTime"`
}

func (user *User) IsVipExpire() bool {
	return user.VipTime < time.Now().Unix()
}

func AddVipTime(uid int, num int) bool {
	now := time.Now().Unix()
	query := "UPDATE user SET vipTime = ? WHERE uid = ? AND vipTime < ?"
	_, err := MC.c.Exec(query, now, uid, now)
	if err != nil {
		logrus.Errorln(err)
	}
	query = "UPDATE user SET vipTime = vipTime + ? WHERE uid = ?"
	result, err := MC.c.Exec(query, num, uid)

	n, err := result.RowsAffected()
	if err != nil {
		logrus.Errorln(err)
	}
	return n > 0
}

func AddUserRewards(uid int, items []def.ItemOpe) {
	for _, v := range items {
		switch v.GetProp() {
		case def.PropVipTime:
			AddVipTime(uid, v.Num)
		default:
			// nothing
			logrus.Warn("not found the item", v.Id, v.Num, v.Extra)
		}
	}
}
