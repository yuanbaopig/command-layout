package module

import (
	"fmt"
	"os"
)

//func CreateSymlink(target, linkName string) error {
//	// 创建软链接
//	err := os.Symlink(target, linkName)
//	if err != nil {
//		return fmt.Errorf("failed to create symlink: %w", err)
//	}
//	return nil
//}

func CreateSymlink(target, linkName string) error {
	// 如果 linkName 已存在，先删除旧的链接或文件
	if _, err := os.Lstat(linkName); err == nil {
		// 只在链接存在时删除
		if err := os.Remove(linkName); err != nil {
			return fmt.Errorf("failed to remove existing link: %w", err)
		}
	}

	// 创建新的软链接
	if err := os.Symlink(target, linkName); err != nil {
		return fmt.Errorf("failed to create symlink: %w", err)
	}
	return nil
}

func GetSymlinkTarget(symlink string) (string, error) {
	// 检查文件是否存在
	info, err := os.Lstat(symlink)
	if err != nil {
		return "", err
	}

	// 检查是否是符号链接
	if info.Mode()&os.ModeSymlink != 0 {
		// 获取符号链接指向的目标
		target, err := os.Readlink(symlink)
		if err != nil {
			return "", fmt.Errorf("error reading symlink: %w", err)
		}

		//fmt.Printf("Symbolic link %s points to: %s\n", symlink, target)

		return target, nil

	}
	//fmt.Printf("%s is not a symbolic link.\n", symlink)

	return "", nil
}
