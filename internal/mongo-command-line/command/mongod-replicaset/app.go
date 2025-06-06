package mongod_replicaset

import (
	mongoreplicaset "DatabaseManage/internal/mongo-command-line/command/mongod-replicaset/options"
	"DatabaseManage/internal/mongo-command-line/options"
	"DatabaseManage/internal/pkg/log"
	"github.com/yuanbaopig/app"
)

const (
	basename    = "mongod-replicaset"
	description = "mongod replica set for set initialize or add node"
)

func New(opts *options.Options) *app.Command {
	o := mongoreplicaset.New()

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
