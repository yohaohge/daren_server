package crypto

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

func Md5(data string) string {
	md5 := md5.New()
	md5.Write([]byte(data))
	return hex.EncodeToString(md5.Sum(nil))
}

func Hmac(key, data string) string {
	hmac := hmac.New(sha1.New, []byte(key))
	hmac.Write([]byte(data))
	return hex.EncodeToString(hmac.Sum(nil))
}

func Sha1(data string) string {
	sha1 := sha1.New()
	sha1.Write([]byte(data))
	return hex.EncodeToString(sha1.Sum(nil))
}
