package mongod_install

import (
	mongodinstalloptions "DatabaseManage/internal/mongo-command-line/command/mongod-install/options"
	"DatabaseManage/internal/mongo-command-line/options"
	"DatabaseManage/internal/pkg/log"
	"github.com/yuanbaopig/app"
)

const (
	basename    = "mongod-install"
	description = "mongod instance for install"
)

func New(opts *options.Options) *app.Command {

	o := mongodinstalloptions.New()

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
