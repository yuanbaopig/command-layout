package contract

import (
	"encoding/json"
	"fmt"
	"strings"
)

type MongoShardingConfig struct {
	Sharding []MongoShardNode `json:"sharding"`
}

type MongoShardNode struct {
	ReplicaSet string           `json:"name"`
	ShardNode  []MongoShardHost `json:"addShard"`
	MaxSize    int              `json:"maxSize,omitempty"`
}

type MongoShardHost struct {
	Host string
	Port int
}

func (m *MongoShardingConfig) String() string {
	if m == nil {
		return ""
	}
	bytes, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func (m *MongoShardingConfig) Set(s string) error {
	err := json.Unmarshal([]byte(s), m)
	if err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return nil
}

func (m *MongoShardingConfig) Type() string {
	return "json"
}

// ConvertToShardString 将 MongoShardingConfig 转换为指定字符串格式
func ConvertToShardString(shardNode MongoShardNode) string {
	var hosts []string
	for _, host := range shardNode.ShardNode {
		hosts = append(hosts, fmt.Sprintf("%s:%d", host.Host, host.Port))
	}

	// 拼接 ReplicaSet 和 Host 列表
	shardString := fmt.Sprintf("%s/%s", shardNode.ReplicaSet, strings.Join(hosts, ","))

	return shardString
}
