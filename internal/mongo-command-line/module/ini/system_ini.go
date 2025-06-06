package ini

import (
	"DatabaseManage/internal/mongo-command-line/contract"
	"fmt"
)

type Unit struct {
	Description string `ini:"Description,omitempty"`
	Requires    string `ini:"Requires,omitempty"`
	After       string `ini:"After,omitempty"`
	Before      string `ini:"Before,omitempty"`
}

type Service struct {
	User         string `ini:"User,omitempty"`
	Group        string `ini:"Group,omitempty"`
	Type         string `ini:"Type,omitempty"`
	ExecStart    string `ini:"ExecStart,omitempty"`
	ExecStart1   string `ini:"ExecStart1,omitempty"`
	ExecReload   string `ini:"ExecReload,omitempty"`
	ExecStop     string `ini:"ExecStop,omitempty"`
	PrivateTmp   bool   `ini:"PrivateTmp,omitempty"`
	LimitFSIZE   string `ini:"LimitFSIZE,omitempty"`
	LimitCPU     string `ini:"LimitCPU,omitempty"`
	LimitAS      string `ini:"LimitAS,omitempty"`
	LimitMEMLOCK string `ini:"LimitMEMLOCK,omitempty"`
	LimitNOFILE  int    `ini:"LimitNOFILE,omitempty"`
	LimitNPROC   int    `ini:"LimitNPROC,omitempty"`
}

type Install struct {
	WantedBy   string `ini:"WantedBy,omitempty"`
	RequiredBy string `ini:"RequiredBy,omitempty"`
}

type SystemdConfig struct {
	Unit `ini:"Unit"`

	Service `ini:"Service"`

	Install `ini:"Install"`
}

func NewHugePageSystemdIniConfig(opts ...SystemdIniConfigOption) *SystemdConfig {
	cfg := &SystemdConfig{
		Unit: Unit{
			Description: "Disable Transparent Hugepage before MongoDB boots",
			Before:      "mongodb.service",
		},
		Service: Service{
			Type:       "oneshot",
			ExecStart:  "/bin/bash -c 'echo never > /sys/kernel/mm/transparent_hugepage/enabled'",
			ExecStart1: "/bin/bash -c 'echo never > /sys/kernel/mm/transparent_hugepage/defrag'",
		},
		Install: Install{
			RequiredBy: "mongodb.service",
		},
	}
	for _, op := range opts {
		op(cfg)
	}

	return cfg

}

func NewMongoDSystemIniConfig(opts ...SystemdIniConfigOption) *SystemdConfig {

	iniCfg := &SystemdConfig{
		Unit: Unit{
			Description: "mongod",
			Requires:    "network.target",
			After:       "network.target remote-fs.target nss-lookup.target",
		},
		Service: Service{
			User:         contract.MongoUser,
			Group:        contract.MongoUser,
			Type:         "forking",
			ExecStart:    "/usr/local/mongodb/bin/mongod --config /data1/mg27017/mg27017.conf",
			ExecReload:   "/bin/kill -s HUP $MAINPID",
			ExecStop:     "/usr/local/mongodb/bin/mongod --shutdown --config /data1/mg27017/mg27017.conf",
			PrivateTmp:   true,
			LimitFSIZE:   "infinity",
			LimitCPU:     "infinity",
			LimitAS:      "infinity",
			LimitMEMLOCK: "infinity",
			LimitNOFILE:  64000,
			LimitNPROC:   64000,
		},
		Install: Install{
			WantedBy: "multi-user.target",
		},
	}

	for _, op := range opts {
		op(iniCfg)
	}

	return iniCfg
}

func NewMongoSSystemIniConfig(opts ...SystemdIniConfigOption) *SystemdConfig {
	iniCfg := &SystemdConfig{
		Unit: Unit{
			Description: "mongos",
			Requires:    "network.target",
			After:       "network.target remote-fs.target nss-lookup.target",
		},
		Service: Service{
			User:         contract.MongoUser,
			Group:        contract.MongoUser,
			Type:         "forking",
			ExecStart:    "/usr/local/mongodb/bin/mongos --config /data1/mg27017/mg27017.conf",
			ExecReload:   "/bin/kill -s HUP $MAINPID",
			ExecStop:     "ExecStop=/usr/local/mongodb/bin/mongo admin --port 27017 --eval 'db.shutdownServer()'",
			PrivateTmp:   true,
			LimitFSIZE:   "infinity",
			LimitCPU:     "infinity",
			LimitAS:      "infinity",
			LimitMEMLOCK: "infinity",
			LimitNOFILE:  64000,
			LimitNPROC:   64000,
		},
		Install: Install{
			WantedBy: "multi-user.target",
		},
	}

	for _, op := range opts {
		op(iniCfg)
	}

	return iniCfg
}

type SystemdIniConfigOption func(conf *SystemdConfig)

func WithMongoDConfig(conf string) SystemdIniConfigOption {
	return func(c *SystemdConfig) {
		//conf := fmt.Sprintf("/data1/mg%d/mg%d.conf", port, port)
		c.ExecStart = contract.MongoDPath + " --config " + conf
		c.ExecStop = contract.MongoDPath + " --shutdown --config " + conf
	}
}

func WithMongoSConfig(conf string, port int) SystemdIniConfigOption {
	return func(c *SystemdConfig) {
		c.ExecStart = contract.MongoSPath + " --config " + conf
		c.ExecStop = contract.MongoPath + fmt.Sprintf(" admin --port %d --eval 'db.shutdownServer()'", port)
	}
}

func WithService(serviceName string) SystemdIniConfigOption {
	return func(c *SystemdConfig) {
		c.RequiredBy = serviceName
		c.Before = serviceName
	}
}
