package util

import (
	"github.com/oschwald/geoip2-golang"
	"github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strings"
)

func GetLocalIp() string {
	addrSlice, err := net.InterfaceAddrs()
	if err == nil {
		for _, addr := range addrSlice {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if nil != ipnet.IP.To4() {
					return ipnet.IP.String()
				}
			}
		}
	}
	return ""
}

func IsPublicIP(IP net.IP) bool {
	if IP.IsLoopback() || IP.IsLinkLocalMulticast() || IP.IsLinkLocalUnicast() {
		return false
	}
	if ip4 := IP.To4(); ip4 != nil {
		switch true {
		case ip4[0] == 10:
			return false
		case ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31:
			return false
		case ip4[0] == 192 && ip4[1] == 168:
			return false
		default:
			return true
		}
	}
	return false
}

// 未获取到ip，则默认返回true，认为该ip来自审核人员
// 目前只屏蔽来自广州ip的访问
func IsReviewIp(ip string) bool {
	if len(ip) == 0 {
		return true
	}
	db, err := geoip2.Open("./GeoLite2-City.mmdb")
	if err != nil {
		logrus.WithError(err).Errorln("connect to geoip database error")
		return true
	}
	defer db.Close()
	// If you are using strings that may be invalid, check that ip is not nil
	ipParsed := net.ParseIP(ip)
	if !IsPublicIP(ipParsed) {
		//内网IP，也开启审核模式
		return true
	}
	record, err := db.City(ipParsed)
	if err != nil {
		logrus.WithFields(logrus.Fields{"err": err, "ip": ip}).Errorln("parse city from ip error")
	}
	city := record.City.Names["zh-CN"]
	return strings.Compare(city, "广州") == 0
}

func RealIp(r *http.Request) string {
	ips := proxy(r)
	if len(ips) > 0 && ips[0] != "" {
		realIP, _, err := net.SplitHostPort(ips[0])
		if err != nil {
			realIP = ips[0]
		}
		return realIP
	}
	xRealIp := r.Header.Get("X-Real-IP")
	if len(xRealIp) > 0 {
		return xRealIp
	}
	if ip, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return ip
	}
	return r.RemoteAddr
}

func proxy(r *http.Request) []string {
	if ips := r.Header.Get("X-Forwarded-For"); ips != "" {
		return strings.Split(ips, ",")
	}
	return []string{}
}
