package utils

import (
	"sync"

	"gitlab.mvalley.com/wind/rime-utils/internal/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

var DBStore *dBStore

type dBStore struct {
	MySqlDatabases map[string]*gorm.DB
	MongoDatabases map[string]*mongo.Database
	mu             sync.Mutex
}

func InitDBStore() {
	DBStore = &dBStore{
		MySqlDatabases: make(map[string]*gorm.DB),
	}
}

func GetMysqlClient(config config.MySQLConfiguration) (*gorm.DB, error) {
	if DBStore == nil {
		InitDBStore()
	}
	DBStore.mu.Lock()
	defer DBStore.mu.Unlock()
	key := config.Host + config.Port + config.DBName
	db, ex := DBStore.MySqlDatabases[key]
	if !ex {
		newDB, err := InitMysql(config)
		if err != nil {
			return nil, err
		}
		DBStore.MySqlDatabases[key] = newDB
		db = newDB
	}
	return db, nil
}

func GetMongoClient(config config.MongoDBConfiguration) (*mongo.Database, error) {
	if DBStore == nil {
		InitDBStore()
	}
	DBStore.mu.Lock()
	defer DBStore.mu.Unlock()
	key := config.Host + config.DBName
	db, ex := DBStore.MongoDatabases[key]
	if !ex {
		newDB, err := InitMongoDB(config)
		if err != nil {
			return nil, err
		}
		DBStore.MongoDatabases[key] = newDB
		db = newDB
	}
	return db, nil
}
