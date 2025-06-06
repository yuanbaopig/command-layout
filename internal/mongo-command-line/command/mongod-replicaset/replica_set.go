package mongod_replicaset

import (
	"DatabaseManage/internal/mongo-command-line/command/common"
	mongoreplicaset "DatabaseManage/internal/mongo-command-line/command/mongod-replicaset/options"
	"DatabaseManage/internal/mongo-command-line/contract"
	"DatabaseManage/internal/mongo-command-line/module"
	"DatabaseManage/internal/pkg/log"
	"DatabaseManage/internal/pkg/mongodb"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

type server struct {
	Address          string                        `json:"address"`
	Port             int                           `json:"port"`
	User             string                        `json:"user"`
	Passwd           string                        `json:"passwd"`
	AuthDatabase     string                        `json:"authDatabase"`
	ReplSetName      string                        `json:"replSetName"`
	Timeout          int                           `json:"timeout"`
	NotDirect        bool                          `json:"direct"`
	NotPortAvailable bool                          `json:"notPortAvailable"`
	OpType           string                        `json:"opType"`
	Config           *contract.MongoReplInitConfig `json:"config"`
	MongoService     *mongodb.MongoService         `json:"mongoService"`
}

type prepareServer struct {
	*server
}

func createServer(opts *mongoreplicaset.Options) (*server, error) {
	log.Debug("create server for mongo replica set")
	rps := opts.MgReplicaSetOpts.ApplyTo()

	return &server{
		Address:          rps.Address,
		Port:             rps.Port,
		User:             rps.User,
		Passwd:           rps.Passwd,
		AuthDatabase:     rps.AuthDatabase,
		ReplSetName:      rps.ReplSetName,
		Timeout:          rps.Timeout,
		NotDirect:        rps.NotDirect,
		OpType:           rps.OpType,
		Config:           rps.ConfigInit,
		NotPortAvailable: rps.NotPortAvailable,
	}, nil

}

func (s *server) PrepareRun() (*prepareServer, error) {
	log.Debug("prepare run server")
	var ps = &prepareServer{}
	mgOpts := &mongodb.MongoOptions{}

	mgOpts.Host = s.Address
	mgOpts.Port = s.Port
	mgOpts.Direct = !s.NotDirect
	mgOpts.ConnectTimeout = s.Timeout

	// 如果设置了用户密码，按照URL链接方式进行登陆
	if len(s.User) > 0 && len(s.Passwd) > 0 {
		mgOpts.Username = s.User
		mgOpts.Password = s.Passwd
		mgOpts.Database = s.AuthDatabase
		mgOpts.ReplicaSet = s.ReplSetName
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

func (p *prepareServer) Run() error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(p.Timeout))
	defer cancel()

	// 探测对应节点的端口是否通
	// 用于一些非节点本地执行的场景，例如控制机进行副本集节点添加
	if !p.server.NotPortAvailable {
		for _, member := range p.server.Config.Members {
			if !module.PortAvailable(ctx, member.Host) {
				return fmt.Errorf("host %s connect failed", member.Host)
			}
		}
	}

	if p.server.OpType == mongoreplicaset.CheckOption {
		log.Debug("port available check done")
		return nil
	}

	// 数据库ping
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

	// 判断option 类型
	switch p.server.OpType {
	case mongoreplicaset.InitOption:
		// 初始化操作
		command := bson.D{
			{"replSetInitiate", p.server.Config},
		}

		if err := client.Database("admin").RunCommand(ctx, command).Decode(&result); err != nil {
			log.Debug(err)
			return fmt.Errorf("RunCommand %s execute failed: %w", p.server.OpType, err)
		}

	case mongoreplicaset.AddOption:
		return common.AddMembersToReplicaSet(ctx, client, p.server.Config.Members)

	}

	return nil
}
