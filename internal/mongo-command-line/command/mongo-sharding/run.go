package mongo_sharding

import (
	commandOptions "DatabaseManage/internal/mongo-command-line/command/mongo-sharding/options"
	"DatabaseManage/internal/pkg/log"
	"github.com/fatih/color"
)

func Run(opts *commandOptions.Options) error {
	log.Debug("mongo-sharding server start running")

	// 返回第一个error就行
	for _, err := range opts.Validate() {
		return err
	}

	err := opts.Complete()
	if err != nil {
		return err
	}

	log.Debugf("Config: %s", color.GreenString(opts.String()))

	s, err := createServer(opts)
	if err != nil {
		return err
	}

	ps, err := s.PrepareRun()
	if err != nil {
		return err
	}
	return ps.Run()
}
