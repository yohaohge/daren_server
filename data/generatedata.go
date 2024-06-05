package main

import (
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"io/fs"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

var MysqlConnStr string = "root:lVkvB1*m@tcp(127.0.0.1:3306)/mini_video?charset=utf8mb4"

// var MysqlConnStr string = "root:lVkvB1*m@tcp(139.9.50.17:9997)/mini_video?charset=utf8mb4"
var MySqlIns *sqlx.DB

var LZJUrl string = "https://122.228.19.15/user/data/"
var DuanJu5Url string = "http://cdn63.nzdd.cn/"

type LZJItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Cover    string `json:"cover"`
	Dir      string `json:"dir,omitempty"`
	TotalNum int    `json:"total_num,omitempty"`
}
type DuanJu5Item struct {
	Name     string `json:"name"`
	Cover    string `json:"cover"`
	TotalNum int    `json:"total_num"`
	Url      string `json:"url"`
}

func main() {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	IsUseMysql := false
	// init mysql
	var InitMysql = func() {
		db, err := sqlx.Connect("mysql", MysqlConnStr)
		if err != nil {
			log.Fatal("mysql connect error: " + err.Error())
			panic(err)
		}
		db.SetMaxOpenConns(100)
		MySqlIns = db
	}
	if IsUseMysql {
		InitMysql()
	}
	//ProcLzj(IsUseMysql)
	ProcDuanju5Com(IsUseMysql)
}

func ProcLzj(IsDirectMysql bool) {

	buf, err := ioutil.ReadFile("./data/item.json")
	if err != nil {
		panic(err)
	}

	dataList := make([]LZJItem, 0)
	if err := json.Unmarshal(buf, &dataList); err != nil {
		panic(err)
	}

	// "https://122.228.19.15/user/data/20231117/_id_3/1/1.m3u8"
	// query := "INSERT INTO video_content (name, data, cover, total, desc, label) VALUES (?,?,?,?,?,?)"

	totalInsert := 0
	counter := 0
	if IsDirectMysql {

		sb := strings.Builder{}
		sb.WriteString("INSERT INTO video_content (name, data, cover, total, `desc`, label) VALUES")
		params := make([]interface{}, 0)

		for _, v := range dataList {
			content := make(map[int]interface{})
			for i := 1; i <= v.TotalNum; i++ {
				path := LZJUrl + v.Dir + "/_id_" + v.ID + "/" + strconv.Itoa(i) + "/" + strconv.Itoa(i) + ".m3u8"
				subCover := LZJUrl + v.Dir + "/_id_" + v.ID + "/" + strconv.Itoa(i) + ".jpg"

				// 每集的结构
				episodeNode := struct {
					Num      int    `json:"num""`
					SubCover string `json:"sub_cover""`
					PlayUrl  string `json:"play_url"`
				}{
					Num:      i,
					SubCover: subCover,
					PlayUrl:  path,
				}
				content[i] = episodeNode
			}
			b, err := json.Marshal(content)
			if err != nil {
				panic(err)
			}
			contentStr := string(b)

			sb.WriteString("(?,?,?,?,?,?)")
			params = append(params, v.Name)
			params = append(params, contentStr)
			params = append(params, v.Cover)
			params = append(params, v.TotalNum)
			params = append(params, v.Name)
			params = append(params, "label")

			counter++
			if counter%100 == 0 {
				// 每次写入100个
				if err != nil {
					panic(err)
				}
				result, err := MySqlIns.Exec(sb.String(), params...)
				if err != nil {
					panic(err)
				}

				n, err := result.RowsAffected()
				totalInsert += int(n)

				counter = 0
				params = params[:0]

				sb.Reset()
				sb.WriteString("INSERT INTO video_content (name, data, cover, total, `desc`, label) VALUES")

				fmt.Println("total insert", totalInsert)

			} else {
				sb.WriteString(",")
			}
		}

		sLen := len(sb.String())
		if sLen > 0 && len(params) > 0 {
			result, err := MySqlIns.Exec(sb.String()[:sLen-1], params...)
			if err != nil {
				panic(err)
			}
			n, err := result.RowsAffected()
			totalInsert += int(n)

			fmt.Println("total insert", totalInsert)
		}
	} else {
		// 写sql
		counter = 0
		sbSql := strings.Builder{}
		sbSql.WriteString("use mini_video;")
		sbSql.WriteString("INSERT INTO video_content (name, data, cover, total, `desc`, label) VALUES")

		for _, v := range dataList {
			content := make(map[int]interface{})
			for i := 1; i <= v.TotalNum; i++ {
				path := LZJUrl + v.Dir + "/_id_" + v.ID + "/" + strconv.Itoa(i) + "/" + strconv.Itoa(i) + ".m3u8"
				subCover := LZJUrl + v.Dir + "/_id_" + v.ID + "/" + strconv.Itoa(i) + ".jpg"

				// 每集的结构
				episodeNode := struct {
					Num      int    `json:"num""`
					SubCover string `json:"sub_cover""`
					PlayUrl  string `json:"play_url"`
				}{
					Num:      i,
					SubCover: subCover,
					PlayUrl:  path,
				}
				content[i] = episodeNode
			}
			b, err := json.Marshal(content)
			if err != nil {
				panic(err)
			}
			contentStr := string(b)

			sbSql.WriteString("(")
			sbSql.WriteString("'" + v.Name + "',")
			sbSql.WriteString("'" + contentStr + "',")
			sbSql.WriteString("'" + v.Cover + "',")
			sbSql.WriteString("'" + strconv.Itoa(v.TotalNum) + "',")
			sbSql.WriteString("'" + v.Name + "',")
			sbSql.WriteString("'" + "label" + "'")
			sbSql.WriteString("),")

			counter++
			if counter%100 == 0 {
				fmt.Println("total insert", counter)
			}
		}
		sqlLen := len(sbSql.String())
		if sqlLen > 0 {
			sqlStr := sbSql.String()[:sqlLen-1]

			if err := ioutil.WriteFile("./data/lzj_data.sql", []byte(sqlStr), fs.ModePerm); err != nil {
				panic(err)
			}

			fmt.Println("total insert", counter)
		}
	}
}
func ProcDuanju5Com(IsDirectMysql bool) {

	buf, err := ioutil.ReadFile("./data/item_duanju5.json")
	if err != nil {
		panic(err)
	}

	dataList := make([]DuanJu5Item, 0)
	if err := json.Unmarshal(buf, &dataList); err != nil {
		panic(err)
	}

	// "http://cdn63.nzdd.cn/6312-%E6%81%B6%E9%AD%94%E8%80%81%E5%85%AC%E4%B8%8D%E5%A5%BD%E6%83%B9%EF%BC%8880%E9%9B%86%EF%BC%89/1.mp4"
	// query := "INSERT INTO video_content (name, data, cover, total, desc, label) VALUES (?,?,?,?,?,?)"

	totalInsert := 0
	counter := 0
	if IsDirectMysql {

		sb := strings.Builder{}
		sb.WriteString("INSERT INTO video_content (name, data, cover, total, `desc`, label) VALUES")
		params := make([]interface{}, 0)

		for _, v := range dataList {
			content := make(map[int]interface{})
			lastIdx := strings.LastIndex(v.Url, "/")
			if lastIdx < 0 {
				panic(errors.New(fmt.Sprint(v.Url, v.TotalNum, v.Name)))
			}
			pathPre := v.Url[0:lastIdx] + "/"
			for i := 1; i <= v.TotalNum; i++ {
				path := pathPre + strconv.Itoa(i) + ".mp4"
				subCover := ""

				// 每集的结构
				episodeNode := struct {
					Num      int    `json:"num""`
					SubCover string `json:"sub_cover""`
					PlayUrl  string `json:"play_url"`
				}{
					Num:      i,
					SubCover: subCover,
					PlayUrl:  path,
				}
				content[i] = episodeNode
			}
			b, err := json.Marshal(content)
			if err != nil {
				panic(err)
			}
			contentStr := string(b)

			sb.WriteString("(?,?,?,?,?,?)")
			params = append(params, v.Name)
			params = append(params, contentStr)
			params = append(params, v.Cover)
			params = append(params, v.TotalNum)
			params = append(params, v.Name)
			params = append(params, "label")

			counter++
			if counter%100 == 0 {
				// 每次写入100个
				if err != nil {
					panic(err)
				}
				result, err := MySqlIns.Exec(sb.String(), params...)
				if err != nil {
					panic(err)
				}

				n, err := result.RowsAffected()
				totalInsert += int(n)

				counter = 0
				params = params[:0]

				sb.Reset()
				sb.WriteString("INSERT INTO video_content (name, data, cover, total, `desc`, label) VALUES")

				fmt.Println("total insert", totalInsert)

			} else {
				sb.WriteString(",")
			}
		}

		sLen := len(sb.String())
		if sLen > 0 && len(params) > 0 {
			result, err := MySqlIns.Exec(sb.String()[:sLen-1], params...)
			if err != nil {
				panic(err)
			}
			n, err := result.RowsAffected()
			totalInsert += int(n)

			fmt.Println("total insert", totalInsert)
		}
	} else {
		// 写sql
		counter = 0
		sbSql := strings.Builder{}
		sbSql.WriteString("use mini_video;")
		sbSql.WriteString("INSERT INTO video_content (name, data, cover, total, `desc`, label) VALUES")

		for _, v := range dataList {
			content := make(map[int]interface{})
			lastIdx := strings.LastIndex(v.Url, "/")
			if lastIdx < 0 {
				panic(errors.New(fmt.Sprint(v.Url, v.TotalNum, v.Name)))
			}
			pathPre := v.Url[0:lastIdx] + "/"
			for i := 1; i <= v.TotalNum; i++ {
				path := pathPre + strconv.Itoa(i) + ".mp4"
				subCover := ""

				// 每集的结构
				episodeNode := struct {
					Num      int    `json:"num""`
					SubCover string `json:"sub_cover""`
					PlayUrl  string `json:"play_url"`
				}{
					Num:      i,
					SubCover: subCover,
					PlayUrl:  path,
				}
				content[i] = episodeNode
			}
			b, err := json.Marshal(content)
			if err != nil {
				panic(err)
			}
			contentStr := string(b)

			sbSql.WriteString("(")
			sbSql.WriteString("'" + v.Name + "',")
			sbSql.WriteString("'" + contentStr + "',")
			sbSql.WriteString("'" + v.Cover + "',")
			sbSql.WriteString("'" + strconv.Itoa(v.TotalNum) + "',")
			sbSql.WriteString("'" + v.Name + "',")
			sbSql.WriteString("'" + "label" + "'")
			sbSql.WriteString("),")

			counter++
			if counter%100 == 0 {
				fmt.Println("total insert", counter)
			}
		}

		sqlLen := len(sbSql.String())
		if sqlLen > 0 {
			sqlStr := sbSql.String()[:sqlLen-1]

			sqlStr = strings.ReplaceAll(sqlStr, "\\u003c", "<")
			sqlStr = strings.ReplaceAll(sqlStr, "\\u003e", ">")
			sqlStr = strings.ReplaceAll(sqlStr, "\\u0026", "&")

			if err := ioutil.WriteFile("./data/duanju5_data.sql", []byte(sqlStr), fs.ModePerm); err != nil {
				panic(err)
			}

			fmt.Println("total insert", counter)
		}
	}
}
