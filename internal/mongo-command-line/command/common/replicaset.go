package common

import (
	"DatabaseManage/internal/mongo-command-line/contract"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	localDb              = "local"
	adminDb              = "admin"
	replicaCetCollection = "system.replset"
)

func AddMembersToReplicaSet(ctx context.Context, client *mongo.Client, members []contract.Members) error {
	replSetColl := client.Database(localDb).Collection(replicaCetCollection)

	// 检查是否只有一个文档
	count, err := replSetColl.CountDocuments(ctx, bson.D{})
	if err != nil {
		return fmt.Errorf("error counting replset documents: %w", err)
	}
	if count > 1 {
		return fmt.Errorf("local.system.replset has unexpected contents")
	}

	// 获取现有配置
	var config bson.M
	if err := replSetColl.FindOne(ctx, bson.D{}).Decode(&config); err != nil {
		return fmt.Errorf("unable to retrieve replica set config from local.system.replset: %w", err)
	}

	// 增加版本号
	if version, ok := config["version"].(int32); ok {
		config["version"] = version + 1
	} else {
		return fmt.Errorf("replica set config missing or invalid version")
	}

	// 解析现有成员
	membersList, err := parseExistingMembers(config)
	if err != nil {
		return fmt.Errorf("error parsing existing members: %w", err)
	}

	// 添加新成员
	newMembers, err := generateNewMembers(members, membersList)
	if err != nil {
		return fmt.Errorf("error generating new members: %w", err)
	}
	config["members"] = newMembers

	// 执行 replSetReconfig
	command := bson.D{{"replSetReconfig", config}}
	adminDB := client.Database(adminDb)
	var result bson.M
	if err := adminDB.RunCommand(ctx, command).Decode(&result); err != nil {

		var configStr []byte
		configStr, _ = json.Marshal(config)

		return fmt.Errorf("replSetReconfig command failed: %w, new config: %s", err, configStr)
	}

	return nil
}

// Helper function: Parse existing members
func parseExistingMembers(config bson.M) ([]bson.M, error) {
	members, ok := config["members"].(bson.A)
	if !ok {
		return nil, fmt.Errorf("members field missing or invalid in config")
	}

	parsedMembers := make([]bson.M, 0, len(members))
	for _, member := range members {
		memberMap, ok := member.(bson.M)
		if !ok {
			return nil, fmt.Errorf("invalid member format in config: %v", member)
		}
		parsedMembers = append(parsedMembers, memberMap)
	}
	return parsedMembers, nil
}

// Helper function: Generate new members
func generateNewMembers(newMembers []contract.Members, existingMembers []bson.M) ([]bson.M, error) {
	maxID := findMaxID(existingMembers)

	for _, member := range newMembers {
		bsonData := member.BsonData()
		if bsonData == nil {
			return nil, fmt.Errorf("invalid member: failed to marshal")
		}

		var newMember bson.M
		if err := bson.Unmarshal(bsonData, &newMember); err != nil {
			return nil, fmt.Errorf("failed to unmarshal new member: %v", err)
		}

		maxID++
		newMember["_id"] = maxID
		existingMembers = append(existingMembers, newMember)
	}

	return existingMembers, nil
}

// Helper function: Find max _id
func findMaxID(members []bson.M) int32 {
	maxID := int32(0)
	for _, member := range members {
		if id, ok := member["_id"].(int32); ok && id > maxID {
			maxID = id
		}
	}
	return maxID
}
