package mongos_install

import (
	"DatabaseManage/internal/mongo-command-line/command/common"
	mongosinstalloptions "DatabaseManage/internal/mongo-command-line/command/mongos-install/options"
	"DatabaseManage/internal/mongo-command-line/contract"
	"DatabaseManage/internal/mongo-command-line/module/ini"
	"DatabaseManage/internal/mongo-command-line/module/yaml"
	installservice "DatabaseManage/internal/mongo-command-line/service"
	"DatabaseManage/internal/pkg/log"
	"context"
	"fmt"
	"path"
)

type server struct {
	Port       int                   `json:"port"`
	DataPath   string                `json:"dataPath"`
	DbVersion  contract.DatabaseInfo `json:"dbVersion"`
	SysVersion contract.SystemInfo   `json:"sysVersion"`
	DbPackage  contract.DbPackage    `json:"dbPackage"`
	ConfigDB   string                `json:"configDB"`
	LogDir     string                `json:"logDir"`
	DataDir    string                `json:"dataDir"`
}

func (s *server) GetDataDir() string {
	return s.DataDir
}

func (s *server) GetLogDir() string {
	return s.LogDir
}

func (s *server) GetDbPackage() contract.DbPackage {
	return s.DbPackage
}

func (s *server) GetDataPath() string {
	return s.DataPath
}

func (s *server) GetPort() int {
	return s.Port
}

type preparedServer struct {
	*server
	service contract.InstallModule
}

func createServer(opts *mongosinstalloptions.Options) (server, error) {
	log.Debug("create server for mongos-install")

	ips := opts.MongoSOpts.ApplyTo()

	dbInfo := contract.DatabaseInfo{
		Name:    contract.DBName,
		Version: ips.Version,
	}

	return server{
		Port:      ips.Port,
		DataPath:  ips.DataPath,
		DbVersion: dbInfo,
		ConfigDB:  ips.ConfigDB,
		DataDir:   path.Join(ips.DataPath, fmt.Sprintf(contract.MongoDataDirFormat, ips.Port)),
		LogDir:    path.Join(ips.DataPath, fmt.Sprintf(contract.MongoLogDirFormat, ips.Port)),
	}, nil

}

func (s *server) PrepareRun() (preparedServer, error) {
	log.Debug("prepare run server")
	ps := preparedServer{}
	service := &installservice.Install{}

	//sName, SVersion, err := common.GetLinuxVersion()
	//if err != nil {
	//	log.Error(err)
	//	return ps, err
	//}
	//
	//s.SysVersion = common.SystemInfo{
	//	Name:    sName,
	//	Version: SVersion,
	//}
	//
	//if pg, ok := common.GetDbPackageMapping().GetPackage(s.SysVersion, s.DbVersion); ok {
	//	s.DbPackage = pg
	//} else {
	//	pgInfo := common.DbPackageInfoFormat(s.SysVersion, s.DbVersion)
	//	supportInfo := common.PrintPackageInfo()
	//	return ps, fmt.Errorf("%s package not matching, only support: %v", pgInfo, supportInfo)
	//}

	dbPackage, err := service.GetPackage(s.DbVersion)
	if err != nil {
		log.Debug(err)
		return ps, err
	}
	s.DbPackage = dbPackage

	//// 检查数据目录
	//if err := ps.service.CheckDirEmpty(s.DataPath); err != nil {
	//	log.Debug(err)
	//	return ps, err
	//}
	//
	//// 创建目录
	//if err := ps.service.CreateDir(s.DataPath); err != nil {
	//	log.Debug(err)
	//	return ps, err
	//}
	//
	//if err := ps.service.CreateDir(contract.PackagePath); err != nil {
	//	log.Debug(err)
	//	return ps, err
	//}

	if err := common.PrepareDir(service, s.DataPath); err != nil {
		return ps, err
	}

	ps.server = s
	ps.service = service
	return ps, nil
}

func (s preparedServer) Run() error {
	log.Debug("run mongos-install server")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
		yaml.WithConfigDB(s.server.ConfigDB),
		yaml.WithSetParameter(),
	)

	switch s.DbVersion.Version {
	case contract.Mongo4_2_10:
		//opts = append(opts, config.WithWiredTigerMaxOverflowFileSize(0))
		opts = append(opts, yaml.WithOperation(yaml.MongoOperationNull{}, 100, 1))
	case contract.Mongo3_4_24:
	}

	// 根据参数生成的配置
	//opts = append(opts, s.makeMongoConfig()...)

	DbCfg := yaml.NewMongoYamlConfig(opts...)
	DbCfgName := path.Join(s.DataPath, fmt.Sprintf("mg%d.conf", s.Port))
	if err := yaml.CreateYamlConfig(DbCfgName, DbCfg); err != nil {
		log.Debug(err)
		return err
	}

	// 获取指定用户的信息，并且对目录设置属主、组
	log.Debug("data directory set user owner")
	if err := s.service.SetUserOwner(contract.MongoUser, s.DataPath); err != nil {
		log.Debug(err)
		return err
	}

	log.Debug("create systemd config")
	systemdCfg := common.SystemdServiceConfig{}
	systemdCfg.ServiceName = fmt.Sprintf(contract.SystemdMongoSServiceName, s.Port)
	systemdCfg.Config = ini.NewMongoSSystemIniConfig(ini.WithMongoSConfig(DbCfgName, s.server.Port))
	systemdCfg.FileName = path.Join(contract.SystemdConfigPath, systemdCfg.ServiceName)

	if err := common.ServiceStartAndEnable(ctx, systemdCfg, s.service); err != nil {
		return err
	}

	return nil
}
