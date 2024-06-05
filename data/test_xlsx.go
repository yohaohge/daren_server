package main

import (
	"fmt"
	"reflect"
	"strings"
	"strconv"
	"github.com/tealeg/xlsx"
)

type Field struct {
	Name string
	Type string
}

type T_Compose struct {

	KEY_id int

	Val1 int

	Val2 float64

	Val3 []int

	Val4 int

}

var myConfig map[int]*T_Compose

func LoadConfigFromExcel(filePath string) (map[int]*T_Compose, error) {
	file, err := xlsx.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	sheet := file.Sheets[0]

	// 获取列头信息
	headerRow := sheet.Rows[0]
	headers := make([]string, len(headerRow.Cells))
	for i, cell := range headerRow.Cells {
		headers[i] = cell.String()
	}

	// 列头对应的类型
	typeRow := sheet.Rows[3]
	types := make([]string, len(headerRow.Cells))
	for i, cell := range typeRow.Cells {
		types[i] = cell.String()
	}

	// 获取字段信息
	fields := make([]Field, len(headers))
	for i, header := range headers {
		tye := "interface{}"
		if i < len(types) {
			tye = types[i]
		}
		fields[i] = Field{Name: header, Type: tye} // 默认为 string 类型
	}

	config := make(map[int]*T_Compose)
	for r := 5; r < len(sheet.Rows); r++ {
		rconfig := &T_Compose{}
		dataRow := sheet.Rows[r] // 数据行从第五行开始
		for i, cell := range dataRow.Cells {
			value := cell.String()

			// 根据字段名设置字段值
			fieldName := fields[i].Name
			fieldType := fields[i].Type

			switch fieldType {
			case "int":
				// 解析为 int 类型
				fieldValue, err := strconv.Atoi(value)
				if err != nil {
					return nil, err
				}
				reflect.ValueOf(rconfig).Elem().FieldByName(fieldName).SetInt(int64(fieldValue))
			case "float", "double":
				// 解析为 float64 类型
				fieldValue, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return nil, err
				}
				reflect.ValueOf(rconfig).Elem().FieldByName(fieldName).SetFloat(fieldValue)
			case "[]int":
				// 解析为 int 数组
				newVal := value[1 : len(value)-1]
				valArr := []string{}
				if len(newVal) > 0 {
					valArr = strings.Split(newVal, ",")
				}
				fieldValue := []int{}
				if len(valArr) > 0 {
					fieldValue = make([]int, len(valArr), len(valArr))
				}
				for i, v := range valArr {
					n, err := strconv.Atoi(v)
					if err != nil {
						return nil, err
					}
					fieldValue[i] = n
				}
				reflect.ValueOf(rconfig).Elem().FieldByName(fieldName).Set(reflect.ValueOf(fieldValue))
			case "[]float", "[]double":
				// 解析为 float64 数组
				newVal := value[1 : len(value)-1]
				valArr := []string{}
				if len(newVal) > 0 {
					valArr = strings.Split(newVal, ",")
				}
				fieldValue := []float64{}
				if len(valArr) > 0 {
					fieldValue = make([]float64, len(valArr), len(valArr))
				}
				for i, v := range valArr {
					f, err := strconv.ParseFloat(v, 64)
					if err != nil {
						return nil, err
					}
					fieldValue[i] = f
				}
				reflect.ValueOf(rconfig).Elem().FieldByName(fieldName).Set(reflect.ValueOf(fieldValue))
			default:
				reflect.ValueOf(rconfig).Elem().FieldByName(fieldName).SetString(value)
			}
		}
		config[rconfig.KEY_id] = rconfig
	}

	return config, nil
}

func GetConfig_T_Compose(id int) *T_Compose {
	if myConfig == nil {
		return nil
	}
	return myConfig[id]
}

func main() {
	cfg, err := LoadConfigFromExcel("./data/test.xlsx")
	if err != nil {
		panic(err)
	}
	for k, v := range cfg {
		fmt.Println(k, v)
	}
	myConfig = cfg
}
