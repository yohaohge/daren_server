package def

const (
	PropLevel   = 0
	PropExp     = 1
	PropVipTime = 2
)

const (
	CDKEY_TYPE_NONE   = 0
	CDKEY_TYPE_NORMAL = 1 // 普通码,一人一码
	CDKEY_TYPE_ANYONE = 2 // 万能码,每个人都可用一次
)

type ItemOpe struct {
	Id    int         `json:"id"`
	Num   int         `json:"num"`
	Extra interface{} `json:"extra"`
}

func (item *ItemOpe) GetProp() int {
	// 目前只有VIP
	return PropVipTime
}
