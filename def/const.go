package def

import (
	"math"
)

const (
	CodeSuccess          = 0
	CodeParamError       = 1
	CodeFailed           = 2
	CodeNeedLogin        = 3
	CodeSignError        = 4
	CodeUserNotExist     = 5
	CodeUserAlreadyExist = 6
	CodePasswordError    = 7
	CodeCanNotWatch      = 8
	CodeNotExist         = 9
	CodeInvalidCDKEY     = 10
)

// user相关
const (
	DefaultAvatarUrl = "https://img.qukankeji.com/avatar/1.jpg"
	GenderUnknown    = 0
	GenderMale       = 1
	GenderFemale     = 2
	DefaultUserCoin  = 100
	UserDeviceExceed = 11

	SessionExpireTime = 7 * 24 * 3600 //session过期时间
	NotLoginWatchMin  = 1
	NotLoginWatchMax  = 3
	LoginWatchMin     = 1
	LoginWatchMax     = 5
	VipWatchMin       = 1
	VipWatchMax       = math.MaxInt32

	BaseUrl = "http://haoju223.cc"
)

const (
	ChannelIOSAppStoreMIPushId = 1
	ChannelAndroidMIPushId     = 2
	ChannelAndroidHuaweiPushId = 3
	ChannelAndroidBaiduPushId  = 4
)

const (
	DaySeconds        = 24 * 60 * 60
	WeekSeconds       = 24 * 60 * 60 * 7
	HourSeconds       = 60 * 60
	MonthMilliSeconds = 30 * 24 * 60 * 60 * 1000
	WeekMilliSeconds  = 7 * 24 * 60 * 60 * 1000
)

const (
	PlatformIOS     = "ios"
	PlatformAndroid = "Android"
)
