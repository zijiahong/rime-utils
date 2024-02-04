package models

import "time"

type Task struct {
	BaseModel
	SourceName     string
	TargetName     string
	SourceDatabase string
	TargetDatabase string
	StartTime      time.Time
	EndTime        time.Time
	LastedLog      string
	SyncStatus     SyncStatus
}

type SubTask struct {
	BaseModel
	ParentTaskID string     // 主任务
	SourceTable  string     // 源表
	TargetTable  string     // 目标表
	Next         string     // 同步游标
	StartTime    time.Time  // 开始时间
	EndTime      time.Time  // 结束时间
	LastedLog    string     // 最后日志
	SyncStatus   SyncStatus // 同步状态
	TotalCount   int64      // 同步总数
	SyncCount    int64      // 已经同步数
}

type SyncStatus string

const (
	SyncStatusInit   SyncStatus = "init"
	SyncStatusSync   SyncStatus = "sync"
	SyncStatusPause  SyncStatus = "pause"
	SyncStatusEnd    SyncStatus = "end"
	SyncStatusWaring SyncStatus = "waring"
	SyncStatusError  SyncStatus = "error"
)
