package config

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"time"
)

type PayConfig struct {
	Id      int    `json:"id"`      // 商品Id
	Price   int    `json:"price"`   // 单位:分
	Desc    string `json:"desc"`    // 描述
	PayCode string `json:"payCode"` // 支付二维码base64数据
	PayLink string `json:"payLink"` // 支付链接
}

type BannerConfig struct {
	Src  string `json:"src"`
	Link string `json:"link"`
}

type Config struct {
	PayCfg      []PayConfig    `json:"payConfig"`
	BanCfg      []BannerConfig `json:"bannerConfig"`
	Notice      string         `json:"notice"`
	Support     string         `json:"support"`
	XYUrl       string         `json:"xyUrl"`
	ShareUrl    string         `json:"shareUrl"`
	WxUrl       string         `json:"wxUrl"`
	DownloadUrl string         `json:"downloadUrl"`
}

var (
	BaseConfigFile = "./config/config.json"
	baseConf       Config
	lastLoadTime   = int64(0)
)

func GetConfig() *Config {
	now := time.Now().Unix()
	if now-lastLoadTime > 3600 {
		LoadConfig()
	}
	return &baseConf
}

func LoadConfig() {
	succ := true
	for {
		buf, err := ioutil.ReadFile(BaseConfigFile)
		if err != nil {
			logrus.Errorln(err)
			succ = false
			break
		}
		if err = json.Unmarshal(buf, &baseConf); err != nil {
			logrus.Errorln(err)
			succ = false
			break
		}
		break
	}
	if !succ {
		payCfg := make([]PayConfig, 0)
		banCfg := make([]BannerConfig, 0)

		payCfg = append(payCfg, PayConfig{Id: 1, Price: 990, Desc: "9.9 VIP1个月", PayCode: "", PayLink: "http://haoju223.cc/tmp/9_9.jpg"})

		banCfg = append(banCfg, BannerConfig{Src: "", Link: ""})
		banCfg = append(banCfg, BannerConfig{Src: "", Link: ""})

		baseConf = Config{
			PayCfg:      payCfg,
			BanCfg:      banCfg,
			Notice:      "好剧天天看!",
			Support:     "Banlangen54321",
			XYUrl:       "【闲鱼】https://m.tb.cn/h.5xyWNly?tk=OytgWMw4tCn HU7632 「我在闲鱼发布了【短剧30天无限畅享】」\n点击链接直接打开",
			ShareUrl:    "http://haoju223.cc/",
			WxUrl:       "http://haoju223.cc/tmp/qr.jpg",
			DownloadUrl: "",
		}
	}
	oldLastLoadTime := lastLoadTime
	lastLoadTime = time.Now().Unix()
	logrus.Infoln("reload config", oldLastLoadTime, lastLoadTime)
}
