package util

import (
	"regexp"
	"strconv"
	"strings"
	"time"
	"LittleVideo/def"
)

func CheckDeviceId(deviceId string) bool {
	len := len(deviceId)
	return len > 10 && len < 100
}

func CheckChannel(channel int32) bool {
	if channel == def.ChannelIOSAppStoreMIPushId || channel == def.ChannelAndroidMIPushId ||
		channel == def.ChannelAndroidHuaweiPushId || channel == def.ChannelAndroidBaiduPushId {
		return true
	}
	return false
}

func CheckPushId(pushId string) bool {
	if len(pushId) == 0 {
		return true
	}
	if len(pushId) < 15 {
		return false
	}
	return true
}

func CheckVmDevice(deviceId string) bool {
	emilator := []string{"qemu", "vbox", "memu", "emulator", "shEmulator"}
	for _, v := range emilator {
		if strings.Contains(deviceId, v) {
			return true
		}
	}

	return false
}

// only match english and chinese
func CheckSpecialCharacter(text string) bool {
	reg := regexp.MustCompile("^[a-zA-Z\u4e00-\u9fa5]+$")
	return reg.MatchString(text)
}

func CheckPunishFinish(freezeInfo map[string]interface{}) bool {
	punishValue := 0
	punishStartTime := 0
	switch freezeInfo["punish_value"].(type) {
	case string:
		punishValue, _ = strconv.Atoi(freezeInfo["punish_value"].(string))
	case float64:
		punishValue = int(freezeInfo["punish_value"].(float64))
	}

	switch freezeInfo["start_time"].(type) {
	case string:
		punishStartTime, _ = strconv.Atoi(freezeInfo["start_time"].(string))
	case float64:
		punishStartTime = int(freezeInfo["start_time"].(float64))
	}
	if punishStartTime == 0 {
		return true
	}
	return int64(punishValue*def.HourSeconds+punishStartTime) < time.Now().Unix()
}

func GetPunishMsgByType(freezeInfo map[string]interface{}) string {
	determineType := 0
	switch freezeInfo["determine_type"].(type) {
	case string:
		determineType, _ = strconv.Atoi(freezeInfo["determine_type"].(string))
	case float64:
		determineType = int(freezeInfo["determine_type"].(float64))
	}

	punishValue := 0
	switch freezeInfo["punish_value"].(type) {
	case string:
		punishValue, _ = strconv.Atoi(freezeInfo["punish_value"].(string))
	case float64:
		punishValue = int(freezeInfo["punish_value"].(float64))
	}
	if determineType == 0 || punishValue == 0 {
		return ""
	}
	//
	//switch determineType {
	//case def.DetermineTypePornSex:
	//	return "你上传的照片严重影响了社区环境，你已被封号。出来混迟早要还的，请文明游戏！"
	//case def.DetermineTypePornBreast:
	//	return "你上传的照片违反了游戏规定，已被删除。同时你已被封号" + strconv.Itoa(punishValue) + "小时。你先静静吧，我也不问你静静是谁。"
	//case def.DetermineTypePornVulgar:
	//	return "你上传的照片违反了游戏规定，已被删除。同时你已被封号" + strconv.Itoa(punishValue) + "小时。你先静静吧，我也不问你静静是谁。"
	//case def.DetermineTypePornDisturb:
	//	return "骚扰他人毕竟是不对的，你已被封号" + strconv.Itoa(punishValue) + "小时。这里不是陌陌，还能好好游戏吗？你先静静吧，我也不问你静静是谁。"
	//case def.DetermineTypeAbuse:
	//	return "在游戏里骂人毕竟是不好的，你已被封号" + strconv.Itoa(punishValue) + "小时。你先静静吧，我也不问你静静是谁。"
	//case def.DetermineTypeTrouble:
	//	return "在游戏里捣乱毕竟是不好的，你已被封号" + strconv.Itoa(punishValue) + "小时。你先静静吧，我也不问你静静是谁。"
	//case def.DetermineTypeCheat:
	//	return "游戏作弊毕竟是不好的，你已被封号" + strconv.Itoa(punishValue) + "小时。你先静静吧，我也不问你静静是谁。"
	//case def.DetermineTypeFreezeDevice:
	//	return "你被多名用户举报严重破坏社区环境，你的设备已被封禁。"
	//case def.DetermineTypeFreezeAccount:
	//	return "你被多名用户举报严重破坏社区环境，目前已被封号。出来混迟早要还的，请文明游戏！"
	//case def.DetermineTypeTwoHour:
	//	return "你被举报破坏社区环境，被封号" + strconv.Itoa(punishValue) + "小时。 如果法官再给一次机会，你一定会做个好人对吧？冷静一下，再来游戏吧！"
	//case def.DetermineTypeHangup:
	//	return "在游戏里挂机毕竟是不好的，你已被封号" + strconv.Itoa(punishValue) + "小时。你先静静吧，我也不问你静静是谁。"
	//}
	return ""
}

func IsBlockMediaName(mediaName string) bool {
	reg := regexp.MustCompile("^.*(讯|报|闻|条|网|文|号|吧|社|刊|军|市|省|院|中央|纪委|国家|中国|检察|长江|新华|青年|媒体|频道|热线|在线|信息|电视台|百度|新浪|搜狐|腾讯|CCTV|36氪|创业邦|艾瑞|i黑马|简书|小红书|喜马拉雅FM|电视猫|芒果TV|荔枝FM|参考消息|人民教育出版社|蜻蜓fm|TVB|成都全搜索|抽屉新热榜|中国蓝TV官方网站|世界经理人|电子产品世界|电台之家|摄影之友|深圳之窗|第一视频|河南|安徽|福建|甘肃|贵州|海南|河北|黑龙江|湖北|湖南|吉林|江苏|江西|辽宁|青海|山东|山西|陕西|四川|云南|浙江|台湾|香港|澳门|西藏|南海|钓鱼岛|美国|日本|印度|广东|体育|NBA|球).*$")
	match := reg.MatchString(mediaName)
	return match
}
