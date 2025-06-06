package test_use

import "github.com/spf13/viper"

func GetConfig(cf, dir string, config interface{}) error {
	//// 设置配置文件名（不带扩展名）
	viper.SetConfigName(cf)

	// 设置配置文件路径
	viper.AddConfigPath(dir) // 当前目录

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	// 将配置文件内容解析到结构体中
	if err := viper.Unmarshal(config); err != nil {
		return err
	}

	return nil
}
