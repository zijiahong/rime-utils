package worker

import (
	"context"
	"fmt"
	"testing"

	"gitlab.mvalley.com/wind/rime-utils/internal/pkg/config"
	"gitlab.mvalley.com/wind/rime-utils/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMongoContainCollection(t *testing.T) {
	sourceDB, err := utils.GetMongoClient(config.MongoDBConfiguration{
		Host:   "mongodb://root:root@10.220.33.21:27517/realm?authSource=admin",
		DBName: "realm",
	})
	if err != nil {
		panic(err)
	}
	// tableName := ""
	// sourceCollection := sourceDB.Collection(tableName)

	// 如果不存在创建表
	ex, err := containCollection(sourceDB, "next_institution_for_wind_240201")
	if err != nil {
		panic(err)
	}
	fmt.Println(ex)
}

func TestMongoIndexes(t *testing.T) {
	sourceDB, err := utils.GetMongoClient(config.MongoDBConfiguration{
		Host:   "mongodb://root:root@10.220.33.21:27517/realm?authSource=admin",
		DBName: "realm",
	})
	if err != nil {
		panic(err)
	}

	// 如果不存在创建表
	ex, err := containCollection(sourceDB, "next_institution_for_wind_240201_copy")
	if err != nil {
		panic(err)
	}

	if !ex {
		// 获取源集合的索引信息
		err := sourceDB.CreateCollection(context.Background(), "next_institution_for_wind_240201_copy")
		if err != nil {
			panic(err)
		}
		indexes, err := sourceDB.Collection("next_institution_for_wind_240201").Indexes().List(context.Background())
		if err != nil {
			panic(err)
		}
		// 在目标集合中创建相同的索引
		for indexes.Next(context.Background()) {
			var indexDescription MongoIndexModel
			if err := indexes.Decode(&indexDescription); err != nil {
				panic(err)
			}
			if indexDescription.Key["_id"] == 1 {
				continue
			}
			// 创建索引
			fmt.Println(indexDescription)
			_, err := sourceDB.Collection("next_institution_for_wind_240201_copy").Indexes().CreateOne(context.Background(), mongo.IndexModel{
				Keys:    indexDescription.Key,
				Options: options.Index().SetUnique(indexDescription.Unique).SetName(indexDescription.Name),
			})
			if err != nil {
				panic(err)
			}
		}
	}
}

func TestGetCount(t *testing.T) {
	sourceDB, err := utils.GetMongoClient(config.MongoDBConfiguration{
		Host:   "mongodb://root:root@10.220.33.21:27517/realm?authSource=admin",
		DBName: "realm",
	})
	if err != nil {
		panic(err)
	}
	var count int64
	count, err = sourceDB.Collection("next_institution_for_wind_240103").CountDocuments(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}
	fmt.Println(count)
}

func TestSyncMongoData(t *testing.T) {
	sourceDB, err := utils.GetMongoClient(config.MongoDBConfiguration{
		Host:   "mongodb://root:root@10.220.33.21:27517/realm?authSource=admin",
		DBName: "realm",
	})
	if err != nil {
		panic(err)
	}
	a := sourceDB.Collection("next_institution_for_wind_240201")
	b := sourceDB.Collection("next_institution_for_wind_240201_copy")
	pkm := "_id"
	var next string
	for {
		sourceData, err := getMongoData(a, next, 100, pkm)
		if err != nil {
			panic(err)
		}
		if len(sourceData) == 0 {
			// 结束
			return
		}
		err = insetMongoData(b, sourceData)
		if err != nil {
			panic(err)
		}
		fmt.Println(sourceData[len(sourceData)-1][pkm].(primitive.ObjectID).Hex())
		next = sourceData[len(sourceData)-1][pkm].(primitive.ObjectID).Hex()
	}
}
