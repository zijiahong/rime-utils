package storage

import (
	"errors"

	"gitlab.mvalley.com/wind/rime-utils/internal/pkg/config"
	"gitlab.mvalley.com/wind/rime-utils/pkg/models"
	"gitlab.mvalley.com/wind/rime-utils/pkg/utils"
	"gorm.io/gorm"
)

type Storage struct {
	DB *gorm.DB
}

var S Storage

func InitStorage(config config.MySQLConfiguration) *Storage {
	db, err := utils.InitMysql(config)
	if err != nil {
		panic(err)
	}

	// panic
	models.AutoMigrate(db)

	return &Storage{DB: db}
}

func (s *Storage) GetDataSource(sourceType models.SourceType, sourcePlatform models.SourcePlatform) (res []models.DataSource, err error) {
	err = s.DB.Model(&models.DataSource{}).Where("source_type = ?", sourceType).Where("source_platform = ?", sourcePlatform).Scan(&res).Error
	return
}

func (s *Storage) GetDataSourceByID(id string) (res *models.DataSource, err error) {
	err = s.DB.Model(&models.DataSource{}).Where("rec = ?", id).First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return
}
