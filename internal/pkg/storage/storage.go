package storage

import (
	"errors"
	"sync"

	"gitlab.mvalley.com/wind/rime-utils/pkg/models"
	"gitlab.mvalley.com/wind/rime-utils/pkg/utils"
	"gorm.io/gorm"
)

type Storage struct {
	DB             *gorm.DB
	MySqlDatabases map[string]*gorm.DB
	mu             sync.Mutex
}

func InitStorage() *Storage {
	return nil
}

func (s *Storage) GetMysqlClient(host, port, user, password, dbName string) (*gorm.DB, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := host + port + dbName
	db, ex := s.MySqlDatabases[key]
	if !ex {
		newDB, err := utils.InitMysql(host, port, user, password, dbName)
		if err != nil {
			return nil, err
		}
		s.MySqlDatabases[key] = newDB
		db = newDB
	}
	return db, nil
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
