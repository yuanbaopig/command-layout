package yaml

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

func CreateYamlConfig(fName string, config interface{}) error {
	// 将结构体编码成 YAML 格式
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config to YAML: %v", err)
	}

	// 将 YAML 写入文件
	//file, err := os.OpenFile(fName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	file, err := os.Create(fName)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {

		return fmt.Errorf("failed to write data to file: %v\n", err)
	}
	//fmt.Println("YAML configuration written to config.yaml")
	return nil
}
