package mongo_sharding

import (
	commandOptions "DatabaseManage/internal/mongo-command-line/command/mongo-sharding/options"
	"DatabaseManage/internal/mongo-command-line/contract"
	"DatabaseManage/internal/pkg/log"
	"DatabaseManage/internal/pkg/mongodb"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type server struct {
	Address      string                        `json:"address"`
	Port         int                           `json:"port"`
	User         string                        `json:"user"`
	Passwd       string                        `json:"passwd"`
	AuthDatabase string                        `json:"authDatabase"`
	Config       *contract.MongoShardingConfig `json:"config"`
	Timeout      int                           `json:"timeout"`
	MongoService *mongodb.MongoService         `json:"mongoService"`
}

type preparedServer struct {
	*server
}

func createServer(options *commandOptions.Options) (*server, error) {
	log.Debug("create server for mongo-sharding")

	mongoSharding := options.MongoSharding.ApplyTo()

	return &server{
		Address:      mongoSharding.Address,
		Port:         mongoSharding.Port,
		User:         mongoSharding.User,
		Passwd:       mongoSharding.Passwd,
		AuthDatabase: mongoSharding.AuthDatabase,
		Config:       mongoSharding.Config,
		Timeout:      mongoSharding.Timeout,
	}, nil
}

func (s *server) PrepareRun() (*preparedServer, error) {
	log.Debug("prepare run server")

	var ps = &preparedServer{}

	mgOpts := &mongodb.MongoOptions{}

	mgOpts.Host = s.Address
	mgOpts.Port = s.Port
	mgOpts.ConnectTimeout = s.Timeout

	// 如果设置了用户密码，按照URL链接方式进行登陆
	if len(s.User) > 0 && len(s.Passwd) > 0 {
		mgOpts.Username = s.User
		mgOpts.Password = s.Passwd
		mgOpts.Database = s.AuthDatabase
	}

	service, err := mongodb.NewClient(mgOpts)
	if err != nil {
		log.Debug(err)
		return ps, err
	}

	s.MongoService = service

	ps.server = s
	return ps, nil
}

func (p *preparedServer) Run() error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(p.Timeout))
	defer cancel()

	if err := p.server.MongoService.Ping(); err != nil {
		if err != nil {
			log.Debug(err)
			return fmt.Errorf("mongodb ping failed: %w", err)
		}
	}

	defer p.server.MongoService.Close()

	log.Debug(p.server.Config.String())

	var result bson.M

	client := p.server.MongoService.GetClient()

	for _, shardNode := range p.server.Config.Sharding {

		sharding := contract.ConvertToShardString(shardNode)

		command := bson.D{
			{
				Key:   "addShard",
				Value: sharding,
			},
			{
				Key:   "name",
				Value: shardNode.ReplicaSet,
			},
		}

		if err := client.Database("admin").RunCommand(ctx, command).Decode(&result); err != nil {
			log.Debug(err)
			return fmt.Errorf("RunCommand addShard excute failed: %w", err)
		}

		log.Debug(result)
	}

	return nil
}
