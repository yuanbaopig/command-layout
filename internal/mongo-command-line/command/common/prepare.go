package common

import (
	"DatabaseManage/internal/mongo-command-line/contract"
	"DatabaseManage/internal/pkg/log"
)

func PrepareDir(service contract.Prepare, DataPath string) error {
	// 检查数据目录
	if err := service.CheckDirEmpty(DataPath); err != nil {
		log.Debug(err)
		return err
	}

	// 创建目录
	if err := service.CreateDir(DataPath); err != nil {
		log.Debug(err)
		return err
	}

	if err := service.CreateDir(contract.PackagePath); err != nil {
		log.Debug(err)
		return err
	}

	return nil
}
