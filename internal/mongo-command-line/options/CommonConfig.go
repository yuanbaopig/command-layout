package options

import (
	"github.com/spf13/pflag"
)

const (
	FileSystemCheckFlag = "fscheck"
)

var _ OptsInterfaces[*CommonConfig] = (*CommonConfig)(nil)

type CommonConfig struct {
	FileSystemCheck bool `json:"fileSystemCheck"`
}

func NewCommonConfigOpts() *CommonConfig {
	return &CommonConfig{}

}

func (c *CommonConfig) ApplyTo() *CommonConfig {
	return c
}

func (c *CommonConfig) AddFlags(set *pflag.FlagSet) {
	set.BoolVar(&c.FileSystemCheck, FileSystemCheckFlag, c.FileSystemCheck, ""+
		"Close file system check.")

}

func (c *CommonConfig) Validate() []error {
	var errs []error

	return errs
}

func (c *CommonConfig) Complete() error {
	return nil
}

func (c *CommonConfig) Name() string {
	return "common-config"
}
