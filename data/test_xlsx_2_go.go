package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"log"
	"os"
	"text/template"
)

type Field struct {
	Name string
	Type string
}

type TemplateData struct {
	PackageName string
	StructName  string
	Fields      []Field
}

func TypeConvert(srcType string) string {
	if srcType == "double" || srcType == "float" {
		return "float64"
	}
	if srcType == "[]double" || srcType == "[]float" {
		return "[]float64"
	}
	return srcType
}

const templateStr = `package {{.PackageName}}

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

type {{.StructName}} struct {
{{range .Fields}}
	{{.Name}} {{TypeConvert .Type}}
{{end}}
}

var myConfig map[int]*{{.StructName}}

func LoadConfigFromExcel(filePath string) (map[int]*{{.StructName}}, error) {
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

	config := make(map[int]*{{.StructName}})
	for r := 5; r < len(sheet.Rows); r++ {
		rconfig := &{{.StructName}}{}
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

func GetConfig_{{.StructName}}(id int) *{{.StructName}} {
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
`

func main() {
	// 打开 Excel 文件
	file, err := xlsx.OpenFile("./data/test.xlsx")
	if err != nil {
		log.Fatal(err)
	}

	// 选择第一个工作表
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

	// 准备模板数据
	data := TemplateData{
		PackageName: "main",
		StructName:  sheet.Name,
		Fields:      fields,
	}

	// 解析模板
	//tmpl, err := template.New("struct").Parse(templateStr)
	//if err != nil {
	//	log.Fatal(err)
	//}
	// 绑定处理函数
	tmpl := template.Must(template.New("struct").Funcs(template.FuncMap{"TypeConvert": TypeConvert}).Parse(templateStr))

	// 生成代码文件
	fileName := "./data/test_xlsx.go"
	outputFile, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// 执行模板渲染
	err = tmpl.Execute(outputFile, data)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Go code exported to %s\n", fileName)
}
