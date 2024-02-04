package models

import "time"

type DataSource struct {
	BaseModel
	SourceType     SourceType
	SourcePlatform SourcePlatform
	Name           string
	Host           string
	Port           string
	User           string
	Password       string
}

type MysqlModel struct {
	Url      string
	User     string
	Password string
}

type BaseModel struct {
	RecId     string
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
