package mongos_install_options

import (
	"DatabaseManage/internal/mongo-command-line/contract"
	"fmt"
	"github.com/spf13/pflag"
)

const (
	flagDbPath       = "dbpath"
	flagPort         = "port"
	flagMongoVersion = "version"
	flagConfigDB     = "configdb"
)

var _ OptsInterfaces[*MongoSOptions] = (*MongoSOptions)(nil)

type MongoSOptions struct {
	//CacheSizeGB int    `json:"cacheSizeGB"`
	DataPath string `json:"dataPath"`
	Port     int    `json:"port"`
	//DeployMode  string `json:"deployMode"`
	//SetName     string `json:"flagSetName"`
	//OplogSizeMB int    `json:"oplogSizeMB"`
	Version string `json:"version"`
	//ShardSvr  bool   `json:"sharding"`
	//ConfigSvr bool   `json:"configSvr"`
	ConfigDB string `yaml:"configDB"`
	//FileCheck   bool   `json:"fileCheck"`
}

func NewMongoSInstallOpts() *MongoSOptions {
	return &MongoSOptions{
		//CacheSizeGB: 1,
		Port: 27017,
		//OplogSizeMB: 20480,
	}
}

func (m *MongoSOptions) ApplyTo() *MongoSOptions {
	return m
}

func (m *MongoSOptions) AddFlags(set *pflag.FlagSet) {
	//set.IntVar(&m.CacheSizeGB, flagWgCache, m.CacheSizeGB, ""+
	//	"Maximum amount of memory to allocate for cache; Defaults to 1GB of physical RAM.")
	set.StringVar(&m.DataPath, flagDbPath, m.DataPath, ""+
		"Directory for data files.")
	set.IntVar(&m.Port, flagPort, m.Port, ""+
		"Specify flagPort number.")

	set.StringVar(&m.Version, flagMongoVersion, m.Version, ""+
		"Specified mongos version.")
	//set.StringVar(&m.SetName, flagSetName, m.SetName, ""+
	//	"Mongo set name.")
	//set.BoolVar(&m.ShardSvr, flagShardSvr, m.ShardSvr, ""+
	//	"Declare this is a shard db of a cluster.")
	//set.BoolVar(&m.ConfigSvr, flagConfigSvr, m.ConfigSvr, ""+
	//	"Declare this is a config db of a cluster.")
	set.StringVar(&m.ConfigDB, flagConfigDB, m.ConfigDB, ""+
		"Connection string for communicating with config servers: <config replset name>/<host1:port>,<host2:port>")
	//set.BoolVar(&m.FileCheck, fCheck, m.FileCheck, ""+
	//	"file system check.")
}

func (m *MongoSOptions) Validate() []error {
	var errs []error

	if len(m.Version) == 0 {
		errs = append(errs, fmt.Errorf("%s option: flag [%s] must not be empty", m.Name(), flagMongoVersion))
	}

	if len(m.ConfigDB) == 0 {
		errs = append(errs, fmt.Errorf("%s option: flag [%s] connection string for communicating with config servers: <config replset name>/<host1:port>,<host2:port>,[...]", m.Name(), flagConfigDB))
	}

	return errs
}

func (m *MongoSOptions) Complete() error {
	if len(m.DataPath) == 0 {
		m.DataPath = fmt.Sprintf(contract.BaseDataDirFormat, m.Port)
	}

	//if len(m.SetName) > 0 {
	//	if 0 == m.OplogSizeMB {
	//		m.OplogSizeMB = 20480
	//	}
	//}

	return nil
}

func (m *MongoSOptions) Name() string {
	return "mongos_install"
}
