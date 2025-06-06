package mongo_command_line

import (
	mongosharding "DatabaseManage/internal/mongo-command-line/command/mongo-sharding"
	mongodinstall "DatabaseManage/internal/mongo-command-line/command/mongod-install"
	mongoreplicaset "DatabaseManage/internal/mongo-command-line/command/mongod-replicaset"
	mongoduninstall "DatabaseManage/internal/mongo-command-line/command/mongod-uninstall"
	mongosinstall "DatabaseManage/internal/mongo-command-line/command/mongos-install"
	"DatabaseManage/internal/mongo-command-line/options"
	"github.com/yuanbaopig/app"
)

const description = "mongo command line manage tools"

func New(basename string) *app.App {
	opts := options.New()

	return app.NewApp("mongo-command-line", basename,
		app.WithDescription(description),
		app.WithNoConfig(),
		app.WithNoVersion(),
		app.WithAddCommands(
			mongodinstall.New(opts),
			mongoduninstall.New(opts),
			mongosinstall.New(opts),
			mongoreplicaset.New(opts),
			mongosharding.New(opts),
		),
		//app.WithRunFunc(run(opts)),
		app.WithOptions(opts),
	)
}
