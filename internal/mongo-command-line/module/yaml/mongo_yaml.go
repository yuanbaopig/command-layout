package yaml

import (
	"path"
)

// OperationMode 接口
type OperationMode interface {
	Mode() string
}

type MongoOperationOff struct{}
type MongoOperationSlowOp struct{}
type MongoOperationAll struct{}
type MongoOperationNull struct {
}

func (MongoOperationOff) Mode() string    { return "off" }
func (MongoOperationSlowOp) Mode() string { return "slowOp" }
func (MongoOperationAll) Mode() string    { return "all" }
func (MongoOperationNull) Mode() string {
	return ""
}

type MongoYamlConfigOption func(config *MongoYamlConfig)

type MongoYamlConfig struct {
	CPU                bool               `json:"cpu" yaml:"cpu,omitempty"`
	ProcessManagement  ProcessManagement  `json:"processManagement" yaml:"processManagement"`
	SystemLog          SystemLog          `json:"systemLog" yaml:"systemLog"`
	Net                Net                `json:"net" yaml:"net"`
	Storage            Storage            `json:"storage" yaml:"storage,omitempty"`
	OperationProfiling OperationProfiling `json:"operationProfiling" yaml:"operationProfiling,omitempty"`
	Replication        Replication        `json:"replication" yaml:"replication,omitempty"`
	Sharding           Sharding           `json:"sharding" yaml:"sharding,omitempty"`
	SetParameter       SetParameter       `json:"setParameter" yaml:"setParameter,omitempty"`
}

type ProcessManagement struct {
	Fork        bool   `json:"fork" yaml:"fork"`
	PidFilePath string `json:"pidFilePath" yaml:"pidFilePath"`
}

type SystemLog struct {
	Verbosity   int    `json:"verbosity" yaml:"verbosity"`
	Quiet       bool   `json:"quiet" yaml:"quiet"`
	Path        string `json:"path" yaml:"path"`
	LogAppend   bool   `json:"logAppend" yaml:"logAppend"`
	LogRotate   string `json:"logRotate" yaml:"logRotate"`
	Destination string `json:"destination" yaml:"destination"`
}

type Compression struct {
	Compressors string `json:"compressors" yaml:"compressors"`
}

type Net struct {
	Port                   int         `json:"port" yaml:"port"`
	BindIP                 string      `json:"bindIp" yaml:"bindIp"`
	MaxIncomingConnections int         `json:"maxIncomingConnections" yaml:"maxIncomingConnections"`
	WireObjectCheck        bool        `json:"wireObjectCheck" yaml:"wireObjectCheck"`
	Compression            Compression `json:"compression" yaml:"compression"`
}

type Journal struct {
	Enabled          bool `json:"enabled" yaml:"enabled"`
	CommitIntervalMs int  `json:"commitIntervalMs" yaml:"commitIntervalMs"`
}

type EngineConfig struct {
	CacheSizeGB                int    `json:"cacheSizeGB" yaml:"cacheSizeGB"`
	JournalCompressor          string `json:"journalCompressor" yaml:"journalCompressor"`
	DirectoryForIndexes        bool   `json:"directoryForIndexes" yaml:"directoryForIndexes"`
	MaxCacheOverflowFileSizeGB int    `json:"maxCacheOverflowFileSizeGB" yaml:"maxCacheOverflowFileSizeGB,omitempty"`
}

type CollectionConfig struct {
	BlockCompressor string `json:"blockCompressor" yaml:"blockCompressor"`
}

type IndexConfig struct {
	PrefixCompression bool `json:"prefixCompression" yaml:"prefixCompression"`
}

type WiredTiger struct {
	EngineConfig     EngineConfig     `json:"engineConfig" yaml:"engineConfig"`
	CollectionConfig CollectionConfig `json:"collectionConfig" yaml:"collectionConfig"`
	IndexConfig      IndexConfig      `json:"indexConfig" yaml:"indexConfig"`
}

type Storage struct {
	DbPath         string     `json:"dbPath" yaml:"dbPath"`
	Journal        Journal    `json:"journal" yaml:"journal"`
	DirectoryPerDB bool       `json:"directoryPerDB" yaml:"directoryPerDB"`
	SyncPeriodSecs int        `json:"syncPeriodSecs" yaml:"syncPeriodSecs"`
	Engine         string     `json:"engine" yaml:"engine"`
	WiredTiger     WiredTiger `json:"wiredTiger" yaml:"wiredTiger"`
}

type OperationProfiling struct {
	Mode              string  `json:"mode" yaml:"mode,omitempty"`
	SlowOpThresholdMs int     `json:"slowOpThresholdMs" yaml:"slowOpThresholdMs"`
	SlowOpSampleRate  float64 `json:"slowOpSampleRate" yaml:"slowOpSampleRate,omitempty"`
}

type Replication struct {
	OplogSizeMB int    `json:"oplogSizeMB" yaml:"oplogSizeMB"`
	ReplSetName string `json:"replSetName" yaml:"replSetName"`
}

type Sharding struct {
	ClusterRole string `json:"clusterRole" yaml:"clusterRole,omitempty"`
	ConfigDB    string `json:"configDB" yaml:"configDB,omitempty"`
}

type SetParameter struct {
	ShardingTaskExecutorPoolMaxSize int `json:"ShardingTaskExecutorPoolMaxSize" yaml:"ShardingTaskExecutorPoolMaxSize"`
	ShardingTaskExecutorPoolMinSize int `json:"ShardingTaskExecutorPoolMinSize" yaml:"ShardingTaskExecutorPoolMinSize"`
	TaskExecutorPoolSize            int `json:"taskExecutorPoolSize" yaml:"taskExecutorPoolSize"`
}

func NewMongoYamlConfig(opts ...MongoYamlConfigOption) *MongoYamlConfig {
	cfg := MongoYamlConfig{
		ProcessManagement: ProcessManagement{
			Fork:        true,
			PidFilePath: "/data1/mg27017/data27017/mongo.pid",
		},
		SystemLog: SystemLog{
			Verbosity:   0,
			Quiet:       false,
			Path:        "/data1/mg27017/log27017/mongodb.log",
			LogAppend:   false,
			LogRotate:   "rename",
			Destination: "file",
		},

		Net: Net{
			BindIP:                 "0.0.0.0",
			MaxIncomingConnections: 2000,
			WireObjectCheck:        true,
			Compression: Compression{
				Compressors: "snappy",
			},
			Port: 27017,
		},

		//Storage: Storage{
		//	DbPath: "/data1/mg27017/data27017/",
		//	Journal: Journal{
		//		Enabled:          true,
		//		CommitIntervalMs: 100,
		//	},
		//	DirectoryPerDB: true,
		//	SyncPeriodSecs: 60,
		//	Engine:         "wiredTiger",
		//	WiredTiger: WiredTiger{
		//		EngineConfig: EngineConfig{
		//			CacheSizeGB:         1,
		//			JournalCompressor:   "snappy",
		//			DirectoryForIndexes: true,
		//			// Available starting in MongoDB 4.2.1 (and 4.0.12)
		//			// db.serverStatus().wiredTiger.cache["cache overflow table max on-disk size"]
		//			//MaxCacheOverflowFileSizeGB: 0,
		//		},
		//		CollectionConfig: CollectionConfig{
		//			BlockCompressor: "snappy",
		//		},
		//		IndexConfig: IndexConfig{
		//			PrefixCompression: true,
		//		},
		//	},
		//},
		//OperationProfiling: OperationProfiling{
		//	Mode:              "slowOp",
		//	SlowOpThresholdMs: 100,
		// 3.4 版本不支持
		//SlowOpSampleRate:  1,
		//},
	}

	for _, op := range opts {
		op(&cfg)
	}

	return &cfg
}

func WithPort(port int) MongoYamlConfigOption {
	return func(c *MongoYamlConfig) {
		c.Net.Port = port
	}
}

func WithMaxConnect(maxConnect int) MongoYamlConfigOption {
	return func(c *MongoYamlConfig) {
		c.Net.MaxIncomingConnections = maxConnect
	}
}

func WithWiredTigerCacheSize(GB int) MongoYamlConfigOption {
	return func(c *MongoYamlConfig) {
		c.Storage.WiredTiger.EngineConfig.CacheSizeGB = GB
	}
}

func WithReplication(oplogSize int, replName string) MongoYamlConfigOption {
	return func(c *MongoYamlConfig) {
		c.Replication.OplogSizeMB = oplogSize
		c.Replication.ReplSetName = replName
	}
}

func WithLogAndPidPath(dataPath, logPath string) MongoYamlConfigOption {
	return func(c *MongoYamlConfig) {
		//c.Storage.DbPath = dataPath
		//c.SystemLog.Path = path.Join(logPath, "mongodb.log")
		c.SystemLog.Path = logPath
		c.ProcessManagement.PidFilePath = path.Join(dataPath, "mongo.pid")
	}
}

func WithOperation[T OperationMode](mode T, timeMs int, rate float64) MongoYamlConfigOption {
	return func(c *MongoYamlConfig) {
		c.OperationProfiling.Mode = mode.Mode()
		c.OperationProfiling.SlowOpThresholdMs = timeMs
		c.OperationProfiling.SlowOpSampleRate = rate
	}
}

func WithCPU(set bool) MongoYamlConfigOption {
	return func(c *MongoYamlConfig) {
		c.CPU = set
	}
}

func WithWiredTigerMaxOverflowFileSize(GB int) MongoYamlConfigOption {
	return func(c *MongoYamlConfig) {
		c.Storage.WiredTiger.EngineConfig.MaxCacheOverflowFileSizeGB = GB
	}
}

func WithShardSvr() MongoYamlConfigOption {
	return func(c *MongoYamlConfig) {
		c.Sharding.ClusterRole = "shardsvr"
	}
}

func WithConfigSvr() MongoYamlConfigOption {
	return func(c *MongoYamlConfig) {
		c.Sharding.ClusterRole = "configsvr"
	}

}

func WithConfigDB(configDB string) MongoYamlConfigOption {
	return func(c *MongoYamlConfig) {

		//configDB := fmt.Sprintf("cfg/%s", strings.Join(addressList, ","))

		c.Sharding.ConfigDB = configDB
	}

}

func WithSetParameter() MongoYamlConfigOption {
	return func(c *MongoYamlConfig) {
		c.SetParameter.ShardingTaskExecutorPoolMaxSize = 30
		c.SetParameter.ShardingTaskExecutorPoolMinSize = 10
		c.SetParameter.TaskExecutorPoolSize = 8
	}
}

func WithStorage(DbPath string) MongoYamlConfigOption {
	return func(c *MongoYamlConfig) {
		c.Storage = Storage{
			DbPath: DbPath,
			Journal: Journal{
				Enabled:          true,
				CommitIntervalMs: 100,
			},
			DirectoryPerDB: true,
			SyncPeriodSecs: 60,
			Engine:         "wiredTiger",
			WiredTiger: WiredTiger{
				EngineConfig: EngineConfig{
					CacheSizeGB:         1,
					JournalCompressor:   "snappy",
					DirectoryForIndexes: true,
					// Available starting in MongoDB 4.2.1 (and 4.0.12)
					// db.serverStatus().wiredTiger.cache["cache overflow table max on-disk size"]
					//MaxCacheOverflowFileSizeGB: 0,
				},
				CollectionConfig: CollectionConfig{
					BlockCompressor: "snappy",
				},
				IndexConfig: IndexConfig{
					PrefixCompression: true,
				},
			},
		}
	}
}
