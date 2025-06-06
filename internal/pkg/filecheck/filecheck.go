package filecheck

import (
	"errors"
	"fmt"
	"io"
	"os"
)

// CheckFileOrDirExist 检查目录和文件是否存在，返回错误为不存在
func CheckFileOrDirExist(path string) error {
	_, err := os.Stat(path)
	// 不存在
	if os.IsNotExist(err) {
		return fmt.Errorf("file or directory %s does not exist: %w", path, err)
	}
	// 其他错误
	if err != nil {
		return fmt.Errorf("file or directory %s doesn't access, error %w", path, err)
	}

	return nil // 文件或目录存在
}

// CheckDirEmpty 检查目录是否为空，或目录不存在
func CheckDirEmpty(path string) error {
	// 打开目录
	dir, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("directory %s doesn't open, error %v", path, err)
	}
	defer dir.Close()

	// 尝试读取目录中的文件名列表，如果存在任何文件，返回错误
	_, err = dir.Readdirnames(1)
	if err == nil {
		return fmt.Errorf("directory %s isn't empty", path)
	} else if err != io.EOF {
		// 非EOF的错误
		return fmt.Errorf("read directory error %v", err)
	}

	return nil // 目录存在且为空
}
