package mongod_uninstall

import (
	mongoduninstalloptions "DatabaseManage/internal/mongo-command-line/command/mongod-uninstall/options"
	"DatabaseManage/internal/mongo-command-line/contract"
	"DatabaseManage/internal/mongo-command-line/module"
	"DatabaseManage/internal/pkg/log"
	"context"
	"fmt"
	"os"
	"path"
	"time"
)

type server struct {
	Port         int    `json:"port"`
	DataPath     string `json:"dataPath"`
	Uninstall    bool   `json:"uninstall"`
	InstanceType string `json:"instanceType"`
}

type preparedServer struct {
	*server
}

func createServer(options *mongoduninstalloptions.Options) (server, error) {
	log.Debug("create server for mongod-uninstall")
	ups := options.MgUninstallOpts.ApplyTo()

	return server{
		Port:         ups.Port,
		DataPath:     ups.DataPath,
		Uninstall:    ups.Uninstall,
		InstanceType: ups.InstanceType,
	}, nil
}

func (s *server) PrepareRun() (preparedServer, error) {
	log.Debug("prepare run server")

	return preparedServer{s}, nil
}

func (p *preparedServer) Run() error {
	// 实例卸载
	// 检查服务是否运行
	// 移动目录
	// 删除systemd 配置文件

	log.Debug("instance uninstall prepare")

	var service string
	ctx := context.Background()

	switch p.server.InstanceType {
	case contract.MongoD:
		service = fmt.Sprintf(contract.SystemdMongoDServiceName, p.server.Port)
	case contract.MongoS:
		service = fmt.Sprintf(contract.SystemdMongoSServiceName, p.server.Port)
	}

	log.Debug("service status check")
	property, err := module.StatusService(ctx, service, module.PropertySubState)
	if err != nil {
		log.Debug(err)
		return fmt.Errorf("service substate check fialed: %w", err)
	}

	if module.SubSateRunning == property.Value.Value() {
		return fmt.Errorf("service %s substate %s", service, module.SubSateRunning)
	}

	log.Debug("data directory rename")
	rmDir := fmt.Sprintf(contract.DropDataDirFormat, p.server.Port, time.Now().Format("20060102150405"))
	if err := os.Rename(p.server.DataPath, rmDir); err != nil {
		log.Debug(err)
		return fmt.Errorf("data directory rename failed: %w", err)
	}

	log.Debug("drop systemd config")
	systemdConf := path.Join(contract.SystemdConfigPath, service)
	if err := os.Remove(systemdConf); err != nil {
		log.Debug(err)
		return fmt.Errorf("systemd config file drop fialed: %w", err)
	}

	log.Debug("systemd reload")
	if err := module.SystemdReload(ctx); err != nil {
		log.Debug(err)
		return fmt.Errorf("systemd reload fialed: %w", err)
	}

	// 判断是否清理软件包
	if !p.Uninstall {
		return nil
	}

	// 部署卸载
	// 删除安装目录
	// 删除软连接

	log.Debug("package directory drop")
	targetDir, err := module.GetSymlinkTarget(contract.LinkPath)
	if err != nil {
		log.Debug(err)
		return fmt.Errorf("target of symlink get failed: %w", err)
	}

	if err := os.RemoveAll(targetDir); err != nil {
		log.Debug(err)
		return fmt.Errorf("package directory drop failed: %w", err)
	}

	log.Debug("symlink file drop")
	if err := os.Remove(contract.LinkPath); err != nil {
		log.Debug(err)
		return fmt.Errorf("symlink file drop fialed: %w", err)
	}

	return nil
}
