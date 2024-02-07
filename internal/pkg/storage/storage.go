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

var S *Storage

func InitStorage(config config.MySQLConfiguration) *Storage {
	db, err := utils.InitMysql(config)
	if err != nil {
		panic(err)
	}

	// panic
	models.AutoMigrate(db)
	S = &Storage{DB: db}
	return S
}

func (s *Storage) GetDataResource(sourceType models.ResourceType, sourcePlatform models.ResourcePlatform) (res []models.DataResource, err error) {
	err = s.DB.Model(&models.DataResource{}).Where("resource_type = ?", sourceType).Where("resource_platform = ?", sourcePlatform).Scan(&res).Error
	return
}

func (s *Storage) GetDataResourceByID(id string) (res *models.DataResource, err error) {
	err = s.DB.Model(&models.DataResource{}).Where("rec_id = ?", id).First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return
}
