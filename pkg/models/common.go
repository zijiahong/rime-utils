package models

import (
	"encoding/json"
	"time"

	"gitlab.mvalley.com/wind/rime-utils/internal/pkg/config"
	"gorm.io/gorm"
)

type BaseModel struct {
	RecId     string `gorm:"type:varchar(36);primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ResourceType string

const (
	SourceTypeElasticSearch ResourceType = "elastic_search"
	SourceTypeMongo         ResourceType = "mongo"
	SourceTypeMySQL         ResourceType = "mysql"
)

type ResourcePlatform string

const (
	SourcePlatformProd ResourcePlatform = "prod" // 正式
	SourcePlatformDev  ResourcePlatform = "dev"  // 开发
	SourcePlatformTest ResourcePlatform = "test" // 测试
)

type ResourceConfig string

func (s ResourceConfig) UnmarshalMysqlConfig() (res config.MySQLConfiguration, err error) {
	err = json.Unmarshal([]byte(s), &res)
	return
}

func (s ResourceConfig) UnmarshalElasticSearchConfig() (res config.ESConfiguration, err error) {
	err = json.Unmarshal([]byte(s), &res)
	return
}

func (s ResourceConfig) UnmarshalMongoConfig() (res config.MongoDBConfiguration, err error) {
	err = json.Unmarshal([]byte(s), &res)
	return
}

func (s ResourceConfig) FillDataBase(db string, t ResourceType) (ResourceConfig, error) {
	switch t {
	case SourceTypeMongo:
		a, err := s.UnmarshalMongoConfig()
		if err != nil {
			return "", err
		}
		a.DBName = db
		b, err := json.Marshal(a)
		return ResourceConfig(b), err
	case SourceTypeMySQL:
		a, err := s.UnmarshalMysqlConfig()
		if err != nil {
			return "", err
		}
		a.DBName = db
		b, err := json.Marshal(a)
		return ResourceConfig(b), err
	}
	return s, nil
}

func AutoMigrate(db *gorm.DB) {
	autoMigrateTask(db)
	autoMigrateResource(db)

	initDataResourceList(db)
}
