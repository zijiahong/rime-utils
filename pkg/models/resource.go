package models

import (
	"encoding/json"
	"time"

	"gitlab.mvalley.com/wind/rime-utils/internal/pkg/config"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DataResource struct {
	BaseModel
	ResourceType     ResourceType
	ResourcePlatform ResourcePlatform
	Name             string
	ResourceConfig   ResourceConfig `gorm:"type:text"`
}

func (t *DataResource) TableName() string {
	return "data_resources"
}

func autoMigrateResource(db *gorm.DB) {
	db.Model(&DataResource{}).Name()
	if !db.Migrator().HasTable(&DataResource{}) {
		err := db.AutoMigrate(&DataResource{})
		if err != nil {
			panic(err)
		}
	}

}

func initDataResourceList(db *gorm.DB) {
	list := make([]DataResource, 0)

	// mysql
	mysql21 := config.MySQLConfiguration{Host: "192.168.88.201", Port: "21306", User: "db_viewer", Password: "JHGCh1bmrGZBpohvDEPY", LogMode: config.None}
	b21, _ := json.Marshal(mysql21)
	list = append(list, DataResource{BaseModel: BaseModel{RecId: "1", CreatedAt: time.Now(), UpdatedAt: time.Now()}, ResourceType: SourceTypeMySQL, ResourcePlatform: SourcePlatformProd, Name: "正式站-10.20.70.21:3306", ResourceConfig: ResourceConfig(b21)})

	mysql24 := config.MySQLConfiguration{Host: "192.168.88.201", Port: "24306", User: "db_viewer", Password: "mHiP0M0b6J09riCLGimK", LogMode: config.None}
	b24, _ := json.Marshal(mysql24)
	list = append(list, DataResource{BaseModel: BaseModel{RecId: "2", CreatedAt: time.Now(), UpdatedAt: time.Now()}, ResourceType: SourceTypeMySQL, ResourcePlatform: SourcePlatformProd, Name: "正式站-10.20.70.24:3306", ResourceConfig: ResourceConfig(b24)})

	mysql33 := config.MySQLConfiguration{Host: "10.220.33.21", Port: "3306", User: "root", Password: "root", LogMode: config.None}
	b33, _ := json.Marshal(mysql33)
	list = append(list, DataResource{BaseModel: BaseModel{RecId: "3", CreatedAt: time.Now(), UpdatedAt: time.Now()}, ResourceType: SourceTypeMySQL, ResourcePlatform: SourcePlatformDev, Name: "开发站-10.220.33.21:3306", ResourceConfig: ResourceConfig(b33)})

	db.Model(&DataResource{}).Clauses(clause.Insert{Modifier: "IGNORE"}).CreateInBatches(list, 100)
}
