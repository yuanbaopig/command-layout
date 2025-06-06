package mongos_install

import (
	mongosinstalloptions "DatabaseManage/internal/mongo-command-line/command/mongos-install/options"
	"DatabaseManage/internal/mongo-command-line/options"
	"DatabaseManage/internal/pkg/log"
	"github.com/yuanbaopig/app"
)

const (
	basename    = "mongos-install"
	description = "mongos instance for install"
)

func New(opts *options.Options) *app.Command {

	o := mongosinstalloptions.New()

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
