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

type SourceType string

const (
	SourceTypeElasticSearch SourceType = "elastic_search"
	SourceTypeMongo         SourceType = "mongo"
	SourceTypeMySQL         SourceType = "mysql"
)

type SourcePlatform string

const (
	SourcePlatformProd SourcePlatform = "prod" // 正式
	SourcePlatformDev  SourcePlatform = "dev"  // 开发
	SourcePlatformTest SourcePlatform = "test" // 测试
)

type SourceConfig string

func (s SourceConfig) UnmarshalMysqlConfig() (res config.MySQLConfiguration, err error) {
	err = json.Unmarshal([]byte(s), &res)
	return
}

// TODO
func (s SourceConfig) UnmarshalESConfig() (res config.MySQLConfiguration, err error) {
	err = json.Unmarshal([]byte(s), &res)
	return
}

// TODO
func (s SourceConfig) UnmarshalMongoConfig() (res config.MongoDBConfiguration, err error) {
	err = json.Unmarshal([]byte(s), &res)
	return
}

func AutoMigrate(db *gorm.DB) {
	autoMigrateTask(db)
	autoMigrateSource(db)
}
