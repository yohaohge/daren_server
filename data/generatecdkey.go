package main

import (
	"LittleVideo/def"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

type CdKeyContent struct {
	Cdkey      string `json:"cdkey"`
	Num        int    `json:"num"`
	CdkeyType  int    `json:"cdkeyType"`
	Items      string `json:"items"`
	CreateTime int64  `json:"createTime"`
}

func main() {
	f, err := os.Create("./data/1.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// 判重
	recMap := make(map[string]string)

	// 普通码
	// type 1-普通码、2-万能码
	count := 0
	sbSql := strings.Builder{}
	sbSql.WriteString("use mini_video;")
	sbSql.WriteString("INSERT INTO cdkey_config (cdkey, num, cdkeyType, items, createTime, expireTime) VALUES")
	now := time.Now().Unix()

	items := [][]*def.ItemOpe{
		[]*def.ItemOpe{&def.ItemOpe{10000, 2 * def.DaySeconds, ""}},
		[]*def.ItemOpe{&def.ItemOpe{10000, 5 * def.DaySeconds, ""}},
		[]*def.ItemOpe{&def.ItemOpe{10000, 7 * def.DaySeconds, ""}},
	}
	for _, item := range items {
		itemJson, err := json.Marshal(item)
		if err != nil {
			panic(err)
		}
		for i := 0; i < 1000; i++ {
			key := strconv.Itoa(rand.Intn(8999)+1000) + strconv.Itoa(rand.Intn(8999)+1000) + fmt.Sprintf("%02d", rand.Intn(100))
			_, b := recMap[key]
			if !b {
				count++
				// 写文件
				f.WriteString(key + "\t" + string(itemJson))
				f.WriteString("\n")

				// 写入mysql
				sbSql.WriteString("(")
				sbSql.WriteString("'" + key + "',")
				sbSql.WriteString("'" + "1" + "',")
				sbSql.WriteString("'" + "1" + "',")
				sbSql.WriteString("'" + string(itemJson) + "',")
				sbSql.WriteString("'" + strconv.Itoa(int(now)) + "',")
				sbSql.WriteString("'" + strconv.Itoa(math.MaxUint32) + "'")
				sbSql.WriteString("),")

				recMap[key] = "1"
			}
		}
		// 下一批
		f.WriteString("\n")
	}
	sqlLen := len(sbSql.String())
	if sqlLen > 0 {
		sqlStr := sbSql.String()[:sqlLen-1]

		if err = ioutil.WriteFile("./data/cdkey_data.sql", []byte(sqlStr), fs.ModePerm); err != nil {
			panic(err)
		}
	}
	fmt.Println("total:", count)
}
