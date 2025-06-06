package mongod_uninstall

import (
	mongoduninstalloptions "DatabaseManage/internal/mongo-command-line/command/mongod-uninstall/options"
	"DatabaseManage/internal/mongo-command-line/options"
	"DatabaseManage/internal/pkg/log"
	"github.com/yuanbaopig/app"
)

const (
	basename    = "mongod-uninstall"
	description = "mongod instance for uninstall"
)

func New(opts *options.Options) *app.Command {
	o := mongoduninstalloptions.New()

	f := func(args []string) error {
		log.Register(opts.Log.ApplyTo().Build())
		defer log.Sync()
		return run(o)
	}

	return app.NewCommand(
		basename,
		description,
		app.WithCommandOptions(o),
		app.WithCommandRunFunc(f),
	)
}
