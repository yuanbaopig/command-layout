package mongod_install_options

import (
	"DatabaseManage/internal/mongo-command-line/contract"
	"fmt"
	"github.com/spf13/pflag"
)

const (
	flagWgCache      = "wiredTigerCacheSizeGB"
	flagDbPath       = "dbpath"
	flagPort         = "port"
	flagMongoType    = "type"
	flagMongoVersion = "version"
	flagOpSize       = "oplogSize"
	flagSetName      = "replSet"
	flagShardSvr     = "shardsvr"
	flagConfigSvr    = "configsvr"
	//flagConfigDB     = "configdb"
)

var _ OptsInterfaces[*MongodInstallOpts] = (*MongodInstallOpts)(nil)

var MongoDeployModeList = []string{contract.Repl, contract.MongoD, contract.Cluster, contract.MongoS}

type MongodInstallOpts struct {
	CacheSizeGB int    `json:"cacheSizeGB"`
	DataPath    string `json:"dataPath"`
	Port        int    `json:"port"`
	DeployMode  string `json:"deployMode"`
	SetName     string `json:"flagSetName"`
	OplogSizeMB int    `json:"oplogSizeMB"`
	Version     string `json:"version"`
	ShardSvr    bool   `json:"sharding"`
	ConfigSvr   bool   `json:"configSvr"`
	//ConfigDB    string `yaml:"configDB"`
	//FileCheck   bool   `json:"fileCheck"`
}

func NewMongodInstallOpts() *MongodInstallOpts {
	return &MongodInstallOpts{
		CacheSizeGB: 1,
		Port:        27017,
		//OplogSizeMB: 20480,
	}
}

func (m *MongodInstallOpts) ApplyTo() *MongodInstallOpts {
	return m
}

func (m *MongodInstallOpts) AddFlags(set *pflag.FlagSet) {
	set.IntVar(&m.CacheSizeGB, flagWgCache, m.CacheSizeGB, ""+
		"Maximum amount of memory to allocate for cache; Defaults to 1GB of physical RAM.")
	set.StringVar(&m.DataPath, flagDbPath, m.DataPath, ""+
		"Directory for data files.")
	set.IntVar(&m.Port, flagPort, m.Port, ""+
		"Specify flagPort number.")
	set.StringVar(&m.DeployMode, flagMongoType, m.DeployMode, ""+
		"Mongod install model, single mongod, repl set or cluster.")
	set.IntVar(&m.OplogSizeMB, flagOpSize, m.OplogSizeMB, ""+
		"Size to use (in MB) for replication op log.")
	set.StringVar(&m.Version, flagMongoVersion, m.Version, ""+
		"Specified mongod version.")
	set.StringVar(&m.SetName, flagSetName, m.SetName, ""+
		"Mongo set name.")
	set.BoolVar(&m.ShardSvr, flagShardSvr, m.ShardSvr, ""+
		"Declare this is a shard db of a cluster.")
	set.BoolVar(&m.ConfigSvr, flagConfigSvr, m.ConfigSvr, ""+
		"Declare this is a config db of a cluster.")
	//set.StringVar(&m.ConfigDB, flagConfigDB, m.ConfigDB, ""+
	//	"Connection string for communicating with config servers: <config replset name>/<host1:port>,<host2:port>")
	//set.BoolVar(&m.FileCheck, fCheck, m.FileCheck, ""+
	//	"file system check.")
}

func (m *MongodInstallOpts) Validate() []error {
	var errs []error

	if len(m.Version) == 0 {
		errs = append(errs, fmt.Errorf("%s option: flag [%s] must not be empty", m.Name(), flagMongoVersion))
	}

	// 参数是否在指定范围内
	if len(m.DeployMode) == 0 || !Contains(m.DeployMode) {
		errs = append(errs, fmt.Errorf("%s option: flag [%s] must not be empty or model error, options list %v", m.Name(), flagMongoType, MongoDeployModeList))
	}

	switch m.DeployMode {
	case contract.MongoS:
		errs = append(errs, fmt.Errorf("%s option: flag [%s] use mongos-install command"))
	//	// 部署mongod时，无法设置setName
	//	if len(m.SetName) > 0 {
	//		errs = append(errs, fmt.Errorf("%s option: flag [%s] must be empty when deploy mode(%s) is mongod", m.Name(), flagSetName, flagMongoType))
	//	}
	//
	//	if len(m.ConfigDB) == 0 {
	//		errs = append(errs, fmt.Errorf("%s option: flag [%s] connection string for communicating with config servers: <config replset name>/<host1:port>,<host2:port>,[...]", m.Name(), flagConfigDB))
	//	}

	case contract.MongoD:
		if len(m.SetName) > 0 {
			errs = append(errs, fmt.Errorf("%s option: flag [%s] must be empty when deploy mode(%s) is mongod", m.Name(), flagSetName, flagMongoType))
		}

	case contract.Cluster:
		if m.ShardSvr && m.ConfigSvr {
			errs = append(errs, fmt.Errorf("%s option: flag [%s] [%s] select only one", m.Name(), flagShardSvr, flagConfigSvr))
		}

		if !m.ShardSvr && !m.ConfigSvr {
			errs = append(errs, fmt.Errorf("%s option: flag [%s] [%s] must be select only", m.Name(), flagShardSvr, flagConfigSvr))

		}

		// 部署setRepl，必须设置setName
		if len(m.SetName) == 0 {
			//errs = append(errs, fmt.Errorf("%s option: flag [%s] must not be empty for replSet or cluster install", m.Name(), flagSetName))
			errs = append(errs, fmt.Errorf("%s option: not running with --%s", m.Name(), flagSetName))
		}

	case contract.Repl:
		// 部署setRepl，必须设置setName
		if len(m.SetName) == 0 {
			//errs = append(errs, fmt.Errorf("%s option: flag [%s] must not be empty for replSet or cluster install", m.Name(), flagSetName))
			errs = append(errs, fmt.Errorf("%s option: not running with --%s", m.Name(), flagSetName))
		}

	}

	if m.ShardSvr || m.ConfigSvr {
		if m.DeployMode != contract.Cluster {
			errs = append(errs, fmt.Errorf("%s option: not running with --%s %s", m.Name(), flagMongoType, contract.Cluster))
		}
	}

	return errs
}

func (m *MongodInstallOpts) Complete() error {
	if len(m.DataPath) == 0 {
		m.DataPath = fmt.Sprintf(contract.BaseDataDirFormat, m.Port)
	}

	if len(m.SetName) > 0 {
		if 0 == m.OplogSizeMB {
			m.OplogSizeMB = 20480
		}
	}

	return nil
}

func (m *MongodInstallOpts) Name() string {
	return "mongod_install"
}

func Contains(item string) bool {
	lookup := make(map[string]bool)
	for _, v := range MongoDeployModeList {
		lookup[v] = true
	}
	return lookup[item]
}
