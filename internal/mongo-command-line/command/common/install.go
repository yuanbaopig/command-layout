package common

import (
	"DatabaseManage/internal/mongo-command-line/contract"
	"DatabaseManage/internal/mongo-command-line/module"
	"DatabaseManage/internal/pkg/log"
	"context"
	"path"
)

type InstallInfo interface {
	GetDbPackage() contract.DbPackage
	GetDataPath() string
	GetPort() int
	GetDataDir() string
	GetLogDir() string
}

func BaseInstall(ctx context.Context, service contract.BaseInstall, s InstallInfo) error {
	log.Debug("create process user")
	if err := service.AddNoLoginUser(contract.MongoUser, module.NoLogin); err != nil {
		log.Debug(err)
		return err
	}

	// 检查安装包是否存在，不存在则下载
	log.Debug("down load install package")
	installPackage := path.Join(contract.PackagePath, s.GetDbPackage().PackageName)
	if err := service.CheckFileExist(installPackage); err != nil {
		// 文件不存在或者打开异常
		if err := service.DownloadFile(ctx, s.GetDbPackage().DownLoadURL, path.Join(contract.PackagePath, s.GetDbPackage().PackageName)); err != nil {
			log.Debug(err)
			return err
		}
	}

	// 解压缩
	log.Debug("extract install package")
	baseDir := path.Join(contract.PackagePath, s.GetDbPackage().DirName)
	if err := service.CheckFileExist(baseDir); err != nil {
		// 文件不存在或者打开异常
		if err := service.Extract(ctx, installPackage, contract.PackagePath); err != nil {
			log.Debug(err)
			return err
		}
	}

	// 创建软链接
	log.Debug("create link file")
	if err := service.CreateSymlink(baseDir, contract.LinkPath); err != nil {
		log.Debug(err)
		return err
	}

	// 创建环境变量
	log.Debug("set env variables")
	if err := service.AddEnvToFile(contract.Profile, "PATH", contract.ProfileEnv); err != nil {
		log.Debug(err)
		return err
	}

	// 创建数据目录
	log.Debug("create data directory")
	//dataDir := path.Join(s.GetDataPath(), fmt.Sprintf("data%d", s.GetPort()))
	if err := service.CreateDir(s.GetDataDir()); err != nil {
		log.Debug(err)
		return err
	}

	// 创建日志目录
	log.Debug("create log directory")
	//logDir := path.Join(s.GetDataPath(), fmt.Sprintf("log%d", s.GetPort()))
	if err := service.CreateDir(s.GetLogDir()); err != nil {
		log.Debug(err)
		return err
	}

	return nil
}
