package install_service

import (
	"DatabaseManage/internal/mongo-command-line/contract"
	"DatabaseManage/internal/mongo-command-line/install-package"
	"DatabaseManage/internal/mongo-command-line/module"
	"DatabaseManage/internal/pkg/filecheck"
	"context"
	"fmt"
)

var _ contract.InstallModule = (*Install)(nil)

type Install struct {
}

func (i *Install) GetPackage(databaseInfo contract.DatabaseInfo) (contract.DbPackage, error) {
	var (
		ok bool
		pg contract.DbPackage
	)

	for _, d := range contract.MongoVersionList {
		install_package.RegistryDbPackages(d)
	}

	systemInfo, err := module.GetLinuxVersion()
	if err != nil {
		return pg, fmt.Errorf("linux version get failed: %w", err)
	}

	SysVersion := contract.SystemInfo{
		Name:    systemInfo.SystemName,
		Version: systemInfo.Version,
	}

	if pg, ok = install_package.GetDbPackageMapping().GetPackage(SysVersion, databaseInfo); ok {
		return pg, nil
	} else {
		pgInfo := install_package.DbPackageInfoFormat(SysVersion, databaseInfo)
		supportInfo := install_package.PrintPackageInfo()
		return pg, fmt.Errorf("%s package not matching, only support: %v", pgInfo, supportInfo)
	}

}

func (i *Install) CheckDirEmpty(path string) error {
	err := filecheck.CheckDirEmpty(path)
	if err != nil {
		return fmt.Errorf("directory %s check failed: %w", path, err)
	}
	return nil
}

func (i *Install) CreateDir(path string) error {
	err := module.CreateDir(path)
	if err != nil {
		return fmt.Errorf("directory %s create failed: %w", path, err)
	}
	return nil
}

func (i *Install) AddNoLoginUser(username, shell string) error {
	err := module.AddNoLoginUser(username, shell)
	if err != nil {
		return fmt.Errorf("user %s create fialed: %w", username, err)
	}
	return nil
}

func (i *Install) CheckFileExist(path string) error {
	err := filecheck.CheckFileOrDirExist(path)
	if err != nil {
		return fmt.Errorf("file %s check failed: %w", path, err)
	}
	return nil
}

func (i *Install) DownloadFile(ctx context.Context, url string, descFilepath string) error {
	err := module.DownloadFile(ctx, url, descFilepath)
	if err != nil {
		return fmt.Errorf("url %s download file: %w", url, err)
	}
	return nil
}

func (i *Install) Extract(ctx context.Context, path, descDir string) error {
	err := module.Extract(ctx, path, descDir)
	if err != nil {
		return fmt.Errorf("file %s extract failed: %w", path, err)
	}
	return nil
}

func (i *Install) CreateSymlink(target, linkName string) error {
	err := module.CreateSymlink(target, linkName)
	if err != nil {
		return fmt.Errorf("symlink %s create failed: %w", linkName, err)
	}
	return nil
}

func (i *Install) AddEnvToFile(profileFile, variable, value string) error {
	err := module.AddEnvToFile(profileFile, variable, value)
	if err != nil {
		return fmt.Errorf("profile add failed: %w", err)
	}
	return nil
}

func (i *Install) SetUserOwner(uName, path string) error {
	err := module.SetUserOwner(uName, path)
	if err != nil {
		return fmt.Errorf("user owner grant failed: %w", err)
	}
	return nil
}

func (i *Install) StartService(ctx context.Context, serviceName string) error {
	err := module.StartService(ctx, serviceName)
	if err != nil {
		return fmt.Errorf("%s service start failed: %w", serviceName, err)
	}
	return nil
}

func (i *Install) SystemdReload(ctx context.Context) error {
	err := module.SystemdReload(ctx)
	if err != nil {
		return fmt.Errorf("systemd reload failed: %w", err)
	}
	return nil
}

func (i *Install) EnableService(ctx context.Context, unitFiles []string) error {
	err := module.EnableService(ctx, unitFiles)
	if err != nil {
		return fmt.Errorf("%v units service enable failed: %w", unitFiles, err)
	}
	return nil
}
