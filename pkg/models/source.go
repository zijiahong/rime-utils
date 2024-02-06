package models

import "gorm.io/gorm"

type DataSource struct {
	BaseModel
	SourceType     SourceType
	SourcePlatform SourcePlatform
	Name           string
	SourceConfig   string
}

func (t *DataSource) TableName() string {
	return "data_sources"
}

func autoMigrateSource(db *gorm.DB) {
	db.Model(&DataSource{}).Name()
	if !db.Migrator().HasTable(&DataSource{}) {
		err := db.AutoMigrate(&DataSource{})
		if err != nil {
			panic(err)
		}
	}

}
