package mongos_install

import (
	mongosinstalloptions "DatabaseManage/internal/mongo-command-line/command/mongos-install/options"
	"DatabaseManage/internal/pkg/log"
	"github.com/fatih/color"
)

func run(opts *mongosinstalloptions.Options) error {
	log.Debug("start mongos install service")

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
