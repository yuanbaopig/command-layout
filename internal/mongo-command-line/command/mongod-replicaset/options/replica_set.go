package mongo_replicaset_options

import (
	"DatabaseManage/internal/mongo-command-line/contract"
	"fmt"
	"github.com/spf13/pflag"
)

const (
	flagOptionType    = "opt"
	flagConfig        = "config"
	flagAddress       = "host"
	flagPort          = "port"
	flagPasswd        = "passwd"
	flagAuthDb        = "authdb"
	flagUser          = "user"
	flagReplSetName   = "replSet"
	flagNotDirect     = "direct"
	flagTimeout       = "timeout"
	flagNotPortActive = "port-available"

	AddOption   = "add"
	InitOption  = "init"
	CheckOption = "check"
)

var optionList = []string{AddOption, InitOption, CheckOption}
var _ OptsInterfaces[*MongoReplicaSetOptions] = (*MongoReplicaSetOptions)(nil)

type MongoReplicaSetOptions struct {
	Address          string                        `json:"address"`
	Port             int                           `json:"port"`
	User             string                        `json:"user"`
	Passwd           string                        `json:"passwd"`
	AuthDatabase     string                        `json:"authDatabase"`
	ReplSetName      string                        `json:"replSetName"`
	ConfigInit       *contract.MongoReplInitConfig `json:"configInit"`
	OpType           string                        `json:"initType"`
	Timeout          int                           `json:"timeout"`
	NotDirect        bool                          `json:"notDirect"`
	NotPortAvailable bool                          `json:"notPortAvailable"`
}

func NewMongoReplicaSetInitOptions() *MongoReplicaSetOptions {
	return &MongoReplicaSetOptions{
		Address:    "127.0.0.1",
		Port:       27017,
		OpType:     InitOption,
		Timeout:    30,
		ConfigInit: &contract.MongoReplInitConfig{},
	}
}

func (m *MongoReplicaSetOptions) ApplyTo() *MongoReplicaSetOptions {
	return m
}

func (m *MongoReplicaSetOptions) AddFlags(set *pflag.FlagSet) {
	set.Var(m.ConfigInit, flagConfig, ""+
		"The following document provides a representation of a replica set configuration document.")

	set.StringVar(&m.OpType, flagOptionType, m.OpType, ""+
		"options type.")

	set.StringVar(&m.Address, flagAddress, m.Address, ""+
		"Mongo service host address.")

	set.StringVar(&m.User, flagUser, m.User, ""+
		"Username for access to mongo service.")

	set.StringVar(&m.Passwd, flagPasswd, m.Passwd, ""+
		"Password for access to mongo, should be used pair with auth database.")

	set.StringVar(&m.AuthDatabase, flagAuthDb, m.AuthDatabase, ""+
		"Database name for the server to use auth.")

	set.IntVar(&m.Port, flagPort, m.Port, ""+
		"Mongo service port address.")

	set.StringVar(&m.ReplSetName, flagReplSetName, m.ReplSetName, ""+
		"Mongo service replica set name.")

	set.BoolVar(&m.NotDirect, flagNotDirect, m.NotDirect, ""+
		"Mongo Not direct connect mode.")

	set.IntVar(&m.Timeout, flagTimeout, m.Timeout, ""+
		"Mongo connection and execution timeout.")

	set.BoolVar(&m.NotPortAvailable, flagNotPortActive, m.NotPortAvailable, ""+
		"Not check port connect available.")
}

func (m *MongoReplicaSetOptions) Validate() []error {
	var errs []error

	if !validType(m.OpType) {
		errs = append(errs, fmt.Errorf("%s option: flag [%s] only supports %v", m.Name(), flagOptionType, optionList))
	}

	if len(m.User) > 0 && len(m.Passwd) > 0 {
		if 0 == len(m.ReplSetName) {
			errs = append(errs, fmt.Errorf("%s option: flag [%s] must not be empty", m.Name(), m.ReplSetName))
		}
	}

	if len(m.ConfigInit.Members) == 0 {
		errs = append(errs, fmt.Errorf("%s option: flag [%s] config parse error", m.Name(), flagConfig))
	}

	return errs
}

func (m *MongoReplicaSetOptions) Complete() error {

	if len(m.User) > 0 && len(m.Passwd) > 0 {
		if 0 == len(m.AuthDatabase) {
			m.AuthDatabase = "admin"
		}
	}

	return nil
}

func (m *MongoReplicaSetOptions) Name() string {
	return "mongo_replicaset"
}

func validType(op string) bool {

	for _, validOption := range optionList {
		if op == validOption {
			return true
		}
	}
	return false
}
