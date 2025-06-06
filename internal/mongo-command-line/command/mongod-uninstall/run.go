package mongod_uninstall

import (
	mongod_uninstall_options "DatabaseManage/internal/mongo-command-line/command/mongod-uninstall/options"
	"DatabaseManage/internal/pkg/log"
	"github.com/fatih/color"
)

func run(opts *mongod_uninstall_options.Options) error {

	log.Debug("start mongod uninstall service")

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
