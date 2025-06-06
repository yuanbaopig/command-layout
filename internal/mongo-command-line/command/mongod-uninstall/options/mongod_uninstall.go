package mongod_uninstall_options

import (
	"DatabaseManage/internal/mongo-command-line/contract"
	"fmt"
	"github.com/spf13/pflag"
)

const (
	flagDbPath       = "datapath"
	flagPort         = "port"
	flagUninstall    = "uninstall"
	flagInstanceType = "type"
)

var _ OptsInterfaces[*MongoUninstallOpts] = (*MongoUninstallOpts)(nil)

type MongoUninstallOpts struct {
	DataPath     string `json:"dataPath"`
	Port         int    `json:"port"`
	Uninstall    bool   `json:"uninstall"`
	InstanceType string `json:"instanceType"`
}

func NewMongoUninstallOpts() *MongoUninstallOpts {
	return &MongoUninstallOpts{}
}

func (u *MongoUninstallOpts) ApplyTo() *MongoUninstallOpts {
	return u
}

func (u *MongoUninstallOpts) AddFlags(set *pflag.FlagSet) {
	set.StringVar(&u.DataPath, flagDbPath, u.DataPath, ""+
		"Directory for data files.")
	set.IntVar(&u.Port, flagPort, u.Port, ""+
		"Specify flagPort number.")
	set.BoolVar(&u.Uninstall, flagUninstall, u.Uninstall, ""+
		"Mongodb binary package uninstall.")
	set.StringVar(&u.InstanceType, flagInstanceType, u.InstanceType, ""+
		"Mongodb instance type, mongod or mongos.")
}

func (u *MongoUninstallOpts) Validate() []error {
	var errs []error

	if u.Port <= 0 {
		errs = append(errs, fmt.Errorf("%s option: flag [%s] must not be empty or less than 0", u.Name(), flagPort))
	}

	if u.InstanceType != contract.MongoD && u.InstanceType != contract.MongoS {
		errs = append(errs, fmt.Errorf("%s option: flag [%s] instance type error", u.Name(), flagInstanceType))
	}

	return errs
}

func (u *MongoUninstallOpts) Complete() error {
	if len(u.DataPath) == 0 {
		u.DataPath = fmt.Sprintf(contract.BaseDataDirFormat, u.Port)
	}

	return nil
}

func (u *MongoUninstallOpts) Name() string {
	return "mongod_uninstall"
}
