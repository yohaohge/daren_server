package util

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

//map key sort ,return sorted value array
func MapKeySort(params map[string]string) []string {
	var keys = make([]string, 0)
	var values = make([]string, 0)
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, v := range keys {
		values = append(values, v)
	}
	return values
}

// 单位是秒
func GetDayStartEndTime() (int64, int64) {
	t := time.Now()
	todayStartTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local).Unix()
	todayEndTime := todayStartTime + 24*3600
	return todayStartTime, todayEndTime
}

func Unmarshal(data []byte, v interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	//first unmarshal using old way
	err = json.Unmarshal(data, v)
	if err == nil {
		return nil
	} else if _, ok := err.(*json.UnmarshalTypeError); !ok {
		return err
	}

	oldType := reflect.TypeOf(v).Elem()
	//first check whether field tag match corresponding type
	mbStringCount := 0
	for i := 0; i < oldType.NumField(); i++ {
		if oldType.Field(i).Tag.Get("mbString") == "true" {
			mbStringCount++
		}
		if oldType.Field(i).Tag.Get("mbString") == "true" &&
			!reflect.TypeOf(int(1)).ConvertibleTo(oldType.Field(i).Type) {
			return errors.New(oldType.Field(i).Name + " tag not match its type")
		}
	}
	if mbStringCount == 0 {
		return err
	}

	//construct new type and new value for unmarshal
	newFields := make([]reflect.StructField, oldType.NumField())
	for i := 0; i < oldType.NumField(); i++ {
		fieldType := oldType.Field(i)
		newFields[i] = fieldType
		if fieldType.Tag.Get("mbString") == "true" {
			newFields[i].Type = reflect.TypeOf("")
		}
	}
	newType := reflect.StructOf(newFields)
	newValue := reflect.New(newType)
	err = json.Unmarshal(data, newValue.Interface())
	if err != nil {
		return err
	}

	//assign new value to old value
	oldValue := reflect.ValueOf(v).Elem()
	for i := 0; i < oldType.NumField(); i++ {
		fieldType := oldType.Field(i)
		if fieldType.Tag.Get("mbString") == "true" {
			value, err := strconv.Atoi(newValue.Elem().Field(i).String())
			if err != nil {
				return err
			}
			oldValue.Field(i).SetInt(int64(value))
		}
	}
	return nil
}

func CheckPhoneNum(phone string) bool {
	reg := regexp.MustCompile(`^(13[0-9]|14[57]|15[0-35-9]|18[07-9])\d{8}$`)
	return reg.MatchString(phone)
}

func GenerateVcode() string {
	rand.Seed(time.Now().UnixNano())
	codeNum := rand.Int31n(1000000)
	return fmt.Sprintf("%06v", codeNum)
}

func UrlEncode(urlstr string) string {
	u := url.Values{}
	u.Set(urlstr, "")
	resstr := u.Encode()
	resstr = resstr[0 : len(resstr)-1]
	resbyte := []byte(resstr)
	resbyte = regexp.MustCompile(`\+`).ReplaceAll(resbyte, []byte("%20"))
	resbyte = regexp.MustCompile(`\*`).ReplaceAll(resbyte, []byte("%2A"))
	resbyte = regexp.MustCompile(`\%7E`).ReplaceAll(resbyte, []byte("~"))
	return string(resbyte)
}

func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\".`)
	}
	return string(path[0 : i+1]), nil
}

func CheckImgUrl(imgurl string) bool {
	reg := regexp.MustCompile(`^(http|https):\/\/.*$`)
	return reg.MatchString(imgurl)
}

func GetSliceFromStr(sourcestr string) []int32 {
	resSlice := make([]int32, 0)
	sourcestr = strings.TrimLeft(sourcestr, "[")
	sourcestr = strings.TrimRight(sourcestr, "]")
	if len(sourcestr) == 0 {
		return resSlice
	}
	sourceSlice := strings.Split(sourcestr, ",")
	if len(sourceSlice) > 0 {
		for _, idstr := range sourceSlice {
			num, err := strconv.Atoi(idstr)
			if err != nil {
				logrus.Error("parse source str error ", err.Error())
				continue
			}
			resSlice = append(resSlice, int32(num))
		}
	}
	return resSlice
}

func GetSliceFromStrMulti(sourcestr string, multi int32) []int32 {
	resSlice := GetSliceFromStr(sourcestr)
	for k, v := range resSlice {
		resSlice[k] = v * multi
	}
	return resSlice
}

func FillStrParam(sourceStr string, params []int32) string {
	reg := regexp.MustCompile(`{[0-9]}`)
	strParams := reg.FindAllString(sourceStr, -1)
	for k, repParam := range strParams {
		if k+1 < len(params) {
			continue
		}
		sourceStr = strings.Replace(sourceStr, repParam, strconv.Itoa(int(params[k])), 1)
	}
	return sourceStr
}

func SliceIntContains(dest []int32, target int32) (bool, int) {
	if len(dest) == 0 {
		return false, 0
	}
	for k, v := range dest {
		if reflect.TypeOf(v) != reflect.TypeOf(target) {
			return false, 0
		}
		if v == target {
			return true, k
		}
	}
	return false, 0
}

func SliceStringContains(dest []string, target string) (bool, int) {
	if len(dest) == 0 {
		return false, 0
	}
	for k, v := range dest {
		if reflect.TypeOf(v) != reflect.TypeOf(target) {
			return false, 0
		}
		if strings.Compare(v, target) == 0 {
			return true, k
		}
	}
	return false, 0
}

func VersionCompare(v1, v2 string) int {
	v1s := strings.Split(v1, ".")
	v2s := strings.Split(v2, ".")
	var length int
	if len(v1s) > len(v2s) {
		length = len(v2s)
	} else {
		length = len(v1s)
	}
	for i := 0; i < length; i++ {
		n1, err := strconv.Atoi(v1s[i])
		if err != nil {
			n1 = 0
		}
		n2, err := strconv.Atoi(v2s[i])
		if err != nil {
			n2 = 0
		}
		if n1 > n2 {
			return 1
		} else if n1 < n2 {
			return -1
		}
	}
	if len(v1s) > length {
		for i := length; i < len(v1s); i++ {
			n1, err := strconv.Atoi(v1s[i])
			if err != nil {
				n1 = 0
			}
			if n1 > 0 {
				return 1
			}
		}
	} else if len(v2s) > length {
		for i := length; i < len(v2s); i++ {
			n2, err := strconv.Atoi(v2s[i])
			if err != nil {
				n2 = 0
			}
			if n2 > 0 {
				return -1
			}
		}
	}
	return 0
}

func BuildQueryStr(data map[string]string) string {
	var queryStr string
	for k, v := range data {
		queryStr += k + "=" + UrlEncode(v) + "&"
	}
	return strings.TrimRight(queryStr, "&")
}
