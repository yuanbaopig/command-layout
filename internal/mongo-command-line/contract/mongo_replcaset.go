package contract

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
)

type MongoReplInitConfig struct {
	ID        string    `json:"_id" bson:"_id"`
	Members   []Members `json:"members" bson:"members"`
	Configsvr bool      `json:"configsvr,omitempty" bson:"configsvr,omitempty"`
}

func (m *MongoReplInitConfig) String() string {
	if m == nil {
		return ""
	}

	bytes, err := json.Marshal(m)
	if err != nil {
		return ""
	}

	return string(bytes)
}

func (m *MongoReplInitConfig) Set(s string) error {
	err := json.Unmarshal([]byte(s), m)
	if err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	return nil
}

func (m *MongoReplInitConfig) Type() string {
	return "json"
}

type Members struct {
	ID           int    `json:"_id" bson:"_id"`
	Host         string `json:"host" bson:"host"`
	ArbiterOnly  bool   `json:"arbiterOnly" bson:"arbiterOnly"`
	BuildIndexes bool   `json:"buildIndexes,omitempty" bson:"buildIndexes,omitempty"`
	Hidden       bool   `json:"hidden,omitempty" bson:"hidden,omitempty"`
	Priority     int    `json:"priority" bson:"priority"`
	SlaveDelay   int    `json:"slaveDelay,omitempty" bson:"slaveDelay,omitempty"`
	Votes        int    `json:"votes,omitempty" bson:"votes,omitempty"`
}

func (m *Members) BsonData() []byte {
	bsonData, err := bson.Marshal(m)
	if err != nil {
		return nil

	}
	return bsonData
}
