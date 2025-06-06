package options

import (
	"DatabaseManage/internal/mongo-command-line/contract"
	"fmt"
	"github.com/spf13/pflag"
)

const (
	flagConfig  = "config"
	flagAddress = "host"
	flagPort    = "port"
	flagPasswd  = "passwd"
	flagAuthDb  = "authdb"
	flagUser    = "user"
)

var _ OptsInterfaces[*MongoShardingOptions] = (*MongoShardingOptions)(nil)

type MongoShardingOptions struct {
	Address      string                        `json:"address"`
	Port         int                           `json:"port"`
	User         string                        `json:"user"`
	Passwd       string                        `json:"passwd"`
	AuthDatabase string                        `json:"authDatabase"`
	Config       *contract.MongoShardingConfig `json:"config"`
	Timeout      int                           `json:"timeout"`
}

func NewMongoShardingOptions() *MongoShardingOptions {
	return &MongoShardingOptions{
		Address: "127.0.0.1",
		Port:    27017,
		Timeout: 30,
		Config:  &contract.MongoShardingConfig{},
	}
}

func (m *MongoShardingOptions) ApplyTo() *MongoShardingOptions {
	return m
}

func (m *MongoShardingOptions) AddFlags(set *pflag.FlagSet) {
	set.Var(m.Config, flagConfig, ""+
		"The command takes for add a shard replica set.")

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
}

func (m *MongoShardingOptions) Validate() []error {
	var errs []error
	if 0 == len(m.Config.Sharding) {
		errs = append(errs, fmt.Errorf("%s option: flag [%s] config parse error", m.Name(), flagConfig))
	}
	return errs
}

func (m *MongoShardingOptions) Complete() error {
	if len(m.User) > 0 && len(m.Passwd) > 0 {
		if 0 == len(m.AuthDatabase) {
			m.AuthDatabase = "admin"
		}
	}

	return nil
}

func (m *MongoShardingOptions) Name() string {
	return "mongo-sharding"
}
