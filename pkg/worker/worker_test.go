package worker

import (
	"encoding/json"
	"testing"
	"time"

	"gitlab.mvalley.com/wind/rime-utils/internal/pkg/config"
	"gitlab.mvalley.com/wind/rime-utils/internal/pkg/storage"
	"gitlab.mvalley.com/wind/rime-utils/pkg/models"
)

func TestWorker(t *testing.T) {
	source := config.MySQLConfiguration{
		Host:     "10.220.33.21",
		Port:     "3306",
		User:     "root",
		Password: "root",
		DBName:   "test_da_pevc_v1",
		LogMode:  config.None,
	}
	sourceB, err := json.Marshal(source)
	if err != nil {
		panic(err)
	}

	w := NewWorker(storage.InitStorage(config.MySQLConfiguration{
		Host:     "localhost",
		Port:     "3306",
		User:     "root",
		Password: "root",
		DBName:   "test",
	}), 10)
	w.Run()
	w.Append(&MysqlJob{
		SubTask: models.SubTask{
			BaseModel: models.BaseModel{
				RecId:     "test1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			ParentTaskID: "test",

			// // 数据库配置
			SourceConfig: models.SourceConfig(sourceB),
			SourceTable:  "funds",
			TargetConfig: models.SourceConfig(sourceB),
			TargetTable:  "copy_funds",
			BatchSize:    2000,
		},
	})
	select {}
}
