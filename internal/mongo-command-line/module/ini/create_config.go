package ini

import (
	"DatabaseManage/internal/pkg/bufferpool"
	"fmt"
	"gopkg.in/ini.v1"
	"os"
	"reflect"
	"regexp"
	"strings"
)

func CreateIniConfig(fileName string, cfg interface{}) error {
	iniFile := ToIni(cfg)

	iniStr, err := ToString(iniFile)
	if err != nil {
		return err
	}

	//file, err := os.Create(fileName)
	//if err != nil {
	//	return err
	//}
	//defer file.Close()
	//
	//_, err = file.WriteString(iniStr)
	if err := os.WriteFile(fileName, []byte(iniStr), 0666); err != nil {
		return err
	}

	return nil
}

func ToString(f *ini.File) (string, error) {
	bf := bufferpool.GetBytesBuffer()
	defer bufferpool.PutBytesBuffer(bf)

	_, err := f.WriteTo(bf)
	if err != nil {
		return "", err
	}

	// 定义正则表达式
	re := regexp.MustCompile(`ExecStart\d+\s+=`)
	// 替换 ExecStartExtend 为 ExecStart
	result := re.ReplaceAllString(bf.String(), "ExecStart  =")

	return result, nil
}

func ToIni(config interface{}) *ini.File {
	// 创建一个新的 ini 文件对象
	file := ini.Empty()

	// 使用反射来遍历主配置结构体的字段
	v := reflect.ValueOf(config)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		sectionName := field.Tag.Get("ini") // 获取 section 标签

		// 跳过没有定义 ini 标签的字段
		if sectionName == "" {
			continue
		}

		section := file.Section(sectionName)
		sectionStruct := v.Field(i)

		// 遍历子结构体的字段
		for j := 0; j < sectionStruct.NumField(); j++ {

			subField := sectionStruct.Type().Field(j)
			tagStr := subField.Tag.Get("ini")

			if len(tagStr) == 0 {
				continue
			}

			tagList := strings.Split(tagStr, ",")

			/*
				keyName := tagList[0]
				omitempty := len(tagList) > 1 && tagList[1] == "omitempty"
			*/
			var keyName string
			var omitempty bool
			for _, tag := range tagList {
				if tag == "omitempty" {
					omitempty = true
					continue
				}
				keyName = tag
			}

			// 获取字段值并检查零值
			value := sectionStruct.Field(j).Interface()
			if reflect.DeepEqual(value, reflect.Zero(sectionStruct.Field(j).Type()).Interface()) && omitempty {
				continue // 忽略零值字段
			}

			// 将字段写入 ini 文件
			section.Key(keyName).SetValue(fmt.Sprintf("%v", value))
		}
	}

	return file
}
