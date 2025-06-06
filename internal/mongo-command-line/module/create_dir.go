package module

import (
	"fmt"
	"os"
)

func CreateDir(path string) error {
	// 使用 os.MkdirAll 递归创建目录
	err := os.MkdirAll(path, 0755) // os.ModePerm = 0777
	if err != nil {
		//log.Error(err)
		return fmt.Errorf("%s directory create failed: %v", path, err)
	}
	return nil
}
