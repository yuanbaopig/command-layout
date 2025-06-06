package mongod_install

import (
	"DatabaseManage/internal/mongo-command-line/command/common"
	mongodinstalloptions "DatabaseManage/internal/mongo-command-line/command/mongod-install/options"
	"DatabaseManage/internal/mongo-command-line/contract"
	"DatabaseManage/internal/mongo-command-line/module/ini"
	"DatabaseManage/internal/mongo-command-line/module/yaml"
	installservice "DatabaseManage/internal/mongo-command-line/service"
	"DatabaseManage/internal/pkg/log"
	"context"
	"fmt"
	"path"
	"syscall"
)

type server struct {
	// 可选的 server 字段
	Port                int                   `json:"port"`
	BaseDataPath        string                `json:"dataPath"`
	WiredTigerCacheSize int                   `json:"wiredTigerCacheSize"`
	DbVersion           contract.DatabaseInfo `json:"dbVersion"`
	SysVersion          contract.SystemInfo   `json:"sysVersion"`
	DbPackage           contract.DbPackage    `json:"dbPackage"`
	FileSystemCheck     bool                  `json:"fileSystemCheck"`
	OplogSize           int                   `json:"oplogSize"`
	ReplName            string                `json:"replName"`
	ShardSvr            bool                  `json:"shardSvr"`
	ConfigSvr           bool                  `json:"configSvr"`
	LogDir              string                `json:"logDir"`
	DataDir             string                `json:"dataDir"`
}

func (s *server) GetDbPackage() contract.DbPackage {
	return s.DbPackage
}

func (s *server) GetDataPath() string {
	return s.BaseDataPath
}

func (s *server) GetPort() int {
	return s.Port
}

func (s *server) GetDataDir() string {
	return s.DataDir
}

func (s *server) GetLogDir() string {
	return s.LogDir
}

type preparedServer struct {
	*server
	service contract.InstallModule
}

func createServer(opts *mongodinstalloptions.Options) (server, error) {
	// 直接返回一个新的 server 实例
	log.Debug("create server for mongod-install")

	ips := opts.MongoDOpts.ApplyTo()
	cmc := opts.CommonConfigOpts.ApplyTo()

	dbInfo := contract.DatabaseInfo{
		Name:    contract.DBName,
		Version: ips.Version,
	}

	return server{
		Port:                ips.Port,
		BaseDataPath:        ips.DataPath,
		DbVersion:           dbInfo,
		FileSystemCheck:     cmc.FileSystemCheck,
		WiredTigerCacheSize: ips.CacheSizeGB,
		OplogSize:           ips.OplogSizeMB,
		ReplName:            ips.SetName,
		ShardSvr:            ips.ShardSvr,
		ConfigSvr:           ips.ConfigSvr,
		DataDir:             path.Join(ips.DataPath, fmt.Sprintf(contract.MongoDataDirFormat, ips.Port)),
		LogDir:              path.Join(ips.DataPath, fmt.Sprintf(contract.MongoLogDirFormat, ips.Port)),
	}, nil
}

func (s *server) PrepareRun() (preparedServer, error) {
	// 获取系统版本
	log.Debug("prepare run server")
	ps := preparedServer{}
	service := &installservice.Install{}

	dbPackage, err := service.GetPackage(s.DbVersion)
	if err != nil {
		log.Debug(err)
		return ps, err
	}
	s.DbPackage = dbPackage

	if err := common.PrepareDir(service, s.BaseDataPath); err != nil {
		return ps, err
	}

	//// 检查数据目录文件系统，默认要求xfs格式
	if !s.FileSystemCheck {
		if err := fileCheck(s.BaseDataPath); err != nil {
			return ps, err
		}
	}

	ps.server = s
	ps.service = service
	return ps, nil
}

func (s preparedServer) Run() error {
	log.Debug("run mongod-install server")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	/*

		配置ulimit		通过systemd上面的约束配置
		关闭big page		通过systemd service 启动顺序配置完
		创建mongod用户 (done)
		下载安装包	(done)
		解压缩安装包	(done)
		配置安装包软连接 (done)
		添加环境变量	(done)

		创建数据目录 (done)
		创建日志目录 (done)
		生成配置文件	(done)
		文件及目录授权 (done)
		生成systemctl配置文件	(done)

		启动mongod服务 (done)
		启动hugepage服务 (done)
	*/

	err := common.BaseInstall(ctx, s.service, s.server)
	if err != nil {
		return err
	}

	// 创建配置文件
	log.Debug("make config file")
	var opts []yaml.MongoYamlConfigOption
	// 设置目录相关配置
	logFile := path.Join(s.LogDir, fmt.Sprintf("mongodb%d.log", s.Port))
	opts = append(opts,
		yaml.WithLogAndPidPath(s.DataDir, logFile),
		yaml.WithPort(s.Port),
		yaml.WithStorage(s.DataDir),
	)
	// 根据参数生成的配置
	opts = append(opts, s.makeMongoConfig()...)

	DbCfg := yaml.NewMongoYamlConfig(opts...)
	DbCfgName := path.Join(s.BaseDataPath, fmt.Sprintf("mg%d.conf", s.Port))
	if err := yaml.CreateYamlConfig(DbCfgName, DbCfg); err != nil {
		log.Debug(err)
		return err
	}

	// 获取指定用户的信息，并且对目录设置属主、组
	log.Debug("data directory set user owner")
	if err := s.service.SetUserOwner(contract.MongoUser, s.BaseDataPath); err != nil {
		log.Debug(err)
		//return fmt.Errorf("directory user owner set failed, %v", err)
		return err
	}

	log.Debug("create systemd config")
	systemdDbCfg := common.SystemdServiceConfig{}
	systemdDbCfg.ServiceName = fmt.Sprintf(contract.SystemdMongoDServiceName, s.Port)
	systemdDbCfg.Config = ini.NewMongoDSystemIniConfig(ini.WithMongoDConfig(DbCfgName))
	systemdDbCfg.FileName = path.Join(contract.SystemdConfigPath, systemdDbCfg.ServiceName)

	systemdHugePageCfg := common.SystemdServiceConfig{}
	systemdHugePageCfg.ServiceName = contract.SystemdHugepageServiceName
	systemdHugePageCfg.Config = ini.NewHugePageSystemdIniConfig(ini.WithService(systemdDbCfg.ServiceName))
	systemdHugePageCfg.FileName = path.Join(contract.SystemdConfigPath, contract.SystemdHugepageServiceName)

	if err := common.ServiceStartAndEnable(ctx, systemdHugePageCfg, s.service); err != nil {
		return err
	}

	if err := common.ServiceStartAndEnable(ctx, systemdDbCfg, s.service); err != nil {
		return err
	}

	log.Debug("mongod install done")

	return nil
}

func fileCheck(path string) error {
	var stat syscall.Statfs_t

	if err := syscall.Statfs(path, &stat); err != nil {
		log.Error(err)
		return fmt.Errorf("failed to get file system type: %v", err)
	}

	if stat.Type != contract.Xfs {
		// 检查文件系统类型（不同系统上类型值不同）
		log.Debugf("File system type: %d", stat.Type)
		return fmt.Errorf("file system is not xfs format")
	}

	return nil
}

func (s *server) makeMongoConfig() []yaml.MongoYamlConfigOption {
	var opts []yaml.MongoYamlConfigOption

	if s.OplogSize > 0 && len(s.ReplName) > 0 {
		opts = append(opts, yaml.WithReplication(s.OplogSize, s.ReplName))
	}

	opts = append(opts, yaml.WithWiredTigerCacheSize(s.WiredTigerCacheSize))

	if s.ShardSvr {
		opts = append(opts, yaml.WithShardSvr())
	}

	if s.ConfigSvr {
		opts = append(opts, yaml.WithConfigSvr())
	}

	// 版本差异预留
	switch s.DbVersion.Version {
	case contract.Mongo4_2_10:
		//opts = append(opts, config.WithWiredTigerMaxOverflowFileSize(0))
		opts = append(opts, yaml.WithOperation(yaml.MongoOperationSlowOp{}, 100, 1))
	case contract.Mongo3_4_24:
		opts = append(opts, yaml.WithOperation(yaml.MongoOperationSlowOp{}, 100, 0))
	}

	return opts

}
