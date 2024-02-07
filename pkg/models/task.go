package models

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	BaseModel
	SourceConfig ResourceConfig
	TargetConfig ResourceConfig
	StartTime    time.Time
	EndTime      time.Time
	LastedLog    string
	SyncStatus   SyncStatus
	SourceType   ResourceType
}

func (t *Task) TableName() string {
	return "tasks"
}

type SubTask struct {
	BaseModel
	ParentTaskID string // 主任务

	// 数据库配置
	SourceConfig ResourceConfig `gorm:"type:text"`
	SourceTable  string         // 源表
	TargetConfig ResourceConfig `gorm:"type:text"`
	TargetTable  string         // 目标表
	SourceType   ResourceType

	// 同步程序参数
	Next       string     // 同步游标
	StartTime  time.Time  // 开始时间
	EndTime    time.Time  // 结束时间
	Error      error      `gorm:"type:text"`        // 最后日志
	SyncStatus SyncStatus `gorm:"type:varchar(36)"` // 同步状态
	TotalCount int64      // 总数
	Batch      int64      // 批次
	BatchSize  int64      // 批次数
}

func (t *SubTask) TableName() string {
	return "sub_tasks"
}

type SyncStatus string

const (
	SyncStatusInit   SyncStatus = "init"
	SyncStatusDoing  SyncStatus = "doing"
	SyncStatusPause  SyncStatus = "pause"
	SyncStatusDone   SyncStatus = "done"
	SyncStatusWaring SyncStatus = "waring"
	SyncStatusError  SyncStatus = "error"
	SyncStatusDelete SyncStatus = "delete"
)

func autoMigrateTask(db *gorm.DB) {
	if !db.Migrator().HasTable(&Task{}) {
		err := db.AutoMigrate(
			&Task{},
		)
		if err != nil {
			panic(err)
		}
	}

	if !db.Migrator().HasTable(&SubTask{}) {
		err := db.AutoMigrate(
			&SubTask{},
		)
		if err != nil {
			panic(err)
		}
	}
}
