package worker

import (
	"context"
	"fmt"

	"gitlab.mvalley.com/wind/rime-utils/internal/pkg/storage"
	"gitlab.mvalley.com/wind/rime-utils/pkg/models"
	"gitlab.mvalley.com/wind/rime-utils/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoJob struct {
	models.SubTask
	stopCh chan struct{}
}

func NewMongoJob(id string, task models.SubTask) Job {
	return &MongoJob{
		stopCh:  make(chan struct{}),
		SubTask: task,
	}
}
func (m *MongoJob) Stop() {
	m.stopCh <- struct{}{}
}

func (m *MongoJob) GetID() string {
	return m.RecId
}

func (m *MongoJob) SaveSyncTask() error {
	return storage.S.UpdateSubTask(m.SubTask)
}

func (m *MongoJob) Run() {
	// 任务开始
	m.SyncStatus = models.SyncStatusDoing
	go m.SaveSyncTask()

	// 连接同步源
	sourceConfig, err := m.SourceConfig.UnmarshalMongoConfig()
	if err != nil {
		m.setError(err)
		return
	}
	sourceDB, err := utils.GetMongoClient(sourceConfig)
	if err != nil {
		m.setError(err)
		return
	}
	sourceCollection := sourceDB.Collection(m.SourceTable)

	// 连接目标源
	targetConfig, err := m.TargetConfig.UnmarshalMongoConfig()
	if err != nil {
		m.setError(err)
		return
	}

	targetDB, err := utils.GetMongoClient(targetConfig)
	if err != nil {
		m.setError(err)
		return
	}

	targetCollection := sourceDB.Collection(m.TargetTable)

	// 如果不存在创建表
	ex, err := containCollection(targetDB, m.TargetTable)
	if !ex {
		// 获取源集合的索引信息
		err := targetDB.CreateCollection(context.Background(), m.TargetTable)
		if err != nil {
			m.setError(err)
			return
		}
		indexes, err := sourceCollection.Indexes().List(context.Background())
		if err != nil {
			m.setError(err)
			return
		}
		// 在目标集合中创建相同的索引
		for indexes.Next(context.Background()) {
			var indexDescription mongo.IndexModel
			if err := indexes.Decode(&indexDescription); err != nil {
				m.setError(err)
				return
			}
			// 创建索引
			_, err := targetCollection.Indexes().CreateOne(context.Background(), indexDescription)
			if err != nil {
				m.setError(err)
				return
			}
		}
	}

	// 查询总数
	var count int64
	count, err = sourceDB.Collection(m.SourceTable).CountDocuments(context.Background(), nil, nil, nil)
	if err != nil {
		m.setError(err)
		return
	}
	m.TotalCount = count

	// 查询主键
	pkm := "_id"

	// 开始同步数据
	for {
		select {
		case <-m.stopCh:
			// 暂停
			m.SyncStatus = models.SyncStatusPause
			return
		default:
			sourceData, err := getMongoData(sourceCollection, m.Next, m.BatchSize, m.SourceTable, pkm)
			if err != nil {
				m.setError(err)
				return
			}
			if len(sourceData) == 0 {
				// 结束
				m.SyncStatus = models.SyncStatusDone
				return
			}

			err = insetMongoData(targetCollection, sourceData)
			if err != nil {
				m.setError(err)
				return
			}
			m.Next = fmt.Sprint(sourceData[len(sourceData)-1][pkm])
			m.Batch++
		}
	}
}

func (m *MongoJob) setError(err error) {
	// 错误
	m.SyncStatus = models.SyncStatusError
	m.Error = err
}

func containCollection(db *mongo.Database, collection string) (bool, error) {
	s, err := db.ListCollectionNames(context.Background(), nil)
	if err != nil {
		return false, err
	}
	for i := range s {
		if s[i] == collection {
			return true, nil
		}
	}
	return false, nil
}

func getMongoData(collection *mongo.Collection, next string, limit int64, tableName, primaryKey string) (res []map[string]interface{}, err error) {
	var cur *mongo.Cursor
	if next != "" {
		filters := bson.M{
			"_id": bson.M{
				"$lt": next,
			},
		}
		cur, err = collection.Find(context.Background(), filters, options.Find().SetLimit(limit).SetSort(bson.D{{"_id", 1}}))
		if err != nil {
			return nil, err
		}

	}

	cur, err = collection.Find(context.Background(), nil, options.Find().SetLimit(limit).SetSort(bson.D{{"_id", 1}}))
	if err != nil {
		return
	}
	err = cur.All(context.Background(), &res)
	return
}

func insetMongoData(collection *mongo.Collection, res []map[string]interface{}) error {
	_, err := collection.InsertMany(context.Background(), toInterfaceArray(res), options.InsertMany().SetOrdered(false))
	return err
}

func toInterfaceArray(input []map[string]interface{}) []interface{} {
	output := make([]interface{}, 0)
	for i := range input {
		output = append(output, input[i])
	}
	return output
}
