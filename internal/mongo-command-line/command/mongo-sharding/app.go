package mongo_sharding

import (
	commandOptions "DatabaseManage/internal/mongo-command-line/command/mongo-sharding/options"
	appOptions "DatabaseManage/internal/mongo-command-line/options"
	"DatabaseManage/internal/pkg/log"
	"github.com/yuanbaopig/app"
)

const (
	basename    = "mongo-sharding"
	description = "mongo sharding cluster for add sharding node"
)

// appOptions 主command options，通用配置
func New(opts *appOptions.Options) *app.Command {
	// commandOptions 子命令的options，本地目录中自身的options
	o := commandOptions.New()

	f := func(args []string) error {
		log.Register(opts.Log.ApplyTo().Build())
		defer log.Sync()

		return Run(o)
	}

	return app.NewCommand(
		basename,
		description,
		app.WithCommandOptions(o),
		app.WithCommandRunFunc(f),
	)
}
