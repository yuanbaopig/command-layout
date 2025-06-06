package mongodb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

// MongoOptions defines options for mongo database.
type MongoOptions struct {
	Host                   string        `json:"host"`
	Username               string        `json:"username"`
	Password               string        `json:"password"`
	Database               string        `json:"database"`
	Direct                 bool          `json:"direct"`
	Port                   int           `json:"port"`
	ReplicaSet             string        `json:"replica_set"` // 添加副本集名字段
	MaxPoolSize            uint64        `json:"max_pool_size"`
	MinPoolSize            uint64        `json:"min_pool_size"`
	ReadPreference         string        `json:"read_preference"`
	RetryWrites            bool          `json:"retry_writes"`
	RetryReads             bool          `json:"retry_reads"`
	WriteConcernLevel      int           `json:"write_concern_level"`
	ReadConcernLevel       string        `json:"read_concern_level"`
	MaxConnIdleTime        time.Duration `json:"max_conn_idle_time"`
	ServerSelectionTimeout time.Duration `json:"server_selection_timeout"`
	ConnectTimeout         int           `json:"connect_timeout"`
}

type MongoService struct {
	options *MongoOptions
	client  *mongo.Client
}

// NewClient 实例化一个mongo连接池.
func NewClient(opts *MongoOptions) (*MongoService, error) {
	var URI string
	if opts.ReplicaSet != "" {
		if opts.Database == "" {
			opts.Database = "admin"
		}
		// 使用副本集的URI格式
		URI = fmt.Sprintf("mongodb://%s:%d/%s?replicaSet=%s",
			opts.Host, opts.Port, opts.Database, opts.ReplicaSet)
	} else {
		// 单节点的URI格式
		URI = fmt.Sprintf("mongodb://%s:%d", opts.Host, opts.Port)
	}

	// 客户端连接参数
	clientOptions := options.Client().ApplyURI(URI)

	if opts.Password != "" && opts.Username != "" {
		clientOptions.SetAuth(options.Credential{
			Username: opts.Username,
			Password: opts.Password,
		}) // 设置数据库认证信息
	}

	// connect pool set
	if opts.Direct {
		clientOptions.Direct = &opts.Direct // 将连接视为单实例连接
	} else {
		clientOptions.SetMaxPoolSize(opts.MaxPoolSize)                                    // 设置连接池最大连接数
		clientOptions.SetMinPoolSize(opts.MinPoolSize)                                    // 设置连接池最小连接数
		clientOptions.SetConnectTimeout(time.Duration(opts.ConnectTimeout) * time.Second) // 设置连接超时为3秒，首次连接测试无效
		clientOptions.SetMaxConnIdleTime(opts.MaxConnIdleTime * time.Second)              // 设置最大空闲时间,单位为秒

		clientOptions.SetServerSelectionTimeout(opts.ServerSelectionTimeout * time.Second) // 设置服务器选择超时为5秒
		clientOptions.SetRetryWrites(opts.RetryWrites)                                     // 开启写操作的自动重试
		clientOptions.SetRetryReads(opts.RetryReads)                                       // 开启读操作的自动重试

		var readPreference *readpref.ReadPref
		switch strings.ToLower(opts.ReadPreference) {
		case "primary":
			readPreference = readpref.Primary()
		case "primarypreferred":
			readPreference = readpref.PrimaryPreferred()
		case "secondary":
			readPreference = readpref.Secondary()
		case "secondarypreferred":
			readPreference = readpref.SecondaryPreferred()
		case "nearest":
			readPreference = readpref.Nearest()
		default:
			// 默认模式或错误处理
			readPreference = readpref.Primary() // 或其他默认行为
		}
		clientOptions.SetReadPreference(readPreference) // 设置读取偏好为Primary

		// 设置写关注
		writeConcern := writeconcern.New(writeconcern.W(opts.WriteConcernLevel))
		clientOptions.SetWriteConcern(writeConcern)

		// 设置读关注
		var readConcern *readconcern.ReadConcern
		switch opts.ReadConcernLevel {
		case "majority":
			readConcern = readconcern.Majority()
		case "local":
			readConcern = readconcern.Local()
		case "linearizable":
			readConcern = readconcern.Linearizable()
		case "available":
			readConcern = readconcern.Available()
		default:
			// 默认模式或错误处理
			readConcern = readconcern.Local() // 或其他默认行为
		}

		clientOptions.SetReadConcern(readConcern)

	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(opts.ConnectTimeout))
	defer cancel()
	// 获取一个客户端
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	return &MongoService{opts, client}, nil
}

// Ping 初始化数据库连接的函数.
func (o *MongoService) Ping() error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(o.options.ConnectTimeout)*time.Second)
	defer cancel() // 不要忘记调用 cancel

	var result bson.M
	if err := o.client.Database("admin").RunCommand(ctx, bson.D{{"ping", 1}}).Decode(&result); err != nil {
		//log.Error(err)
		return err
	}

	// 检验 'ping' 命令的响应是否是预期的
	if val, ok := result["ok"]; !ok || val != 1.0 {
		//log.Errorf("%-v", err)
		return fmt.Errorf("MongoDB ping failed, expected 'ok': 1.0 but got result: %v", result)
	}

	//log.Debug("Connected to MongoDB successfully.")

	return nil
}

func (o *MongoService) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(o.options.ConnectTimeout))
	defer cancel()
	//if err := o.client.Disconnect(ctx); err != nil {
	//	log.Warnf("mongo client close, error: %s", err)
	//}
	return o.client.Disconnect(ctx)
}

// GetClient return mongo client
func (o *MongoService) GetClient() *mongo.Client {
	return o.client
}
