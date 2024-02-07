package worker

import (
	"fmt"
	"strings"

	"gitlab.mvalley.com/wind/rime-utils/internal/pkg/storage"
	"gitlab.mvalley.com/wind/rime-utils/pkg/models"
	"gitlab.mvalley.com/wind/rime-utils/pkg/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type primaryKeyModel struct {
	PrimaryKey string `gorm:"column:Column_name"`
}
type showCreateTable struct {
	Table       string
	CreateTable string `gorm:"column:Create Table"`
}

type MysqlJob struct {
	models.SubTask
	stopCh chan struct{}
}

func NewMysqlJob(task models.SubTask) Job {
	return &MysqlJob{
		stopCh:  make(chan struct{}),
		SubTask: task,
	}
}
func (m *MysqlJob) Stop() {
	m.stopCh <- struct{}{}
}

func (m *MysqlJob) GetID() string {
	return m.RecId
}

func (m *MysqlJob) SaveSyncTask() error {
	return storage.S.UpdateSubTask(m.SubTask)
}

func (m *MysqlJob) Run() {
	// 任务开始
	m.SyncStatus = models.SyncStatusDoing

	go m.SaveSyncTask()

	// 连接同步源
	sourceConfig, err := m.SourceConfig.UnmarshalMysqlConfig()
	if err != nil {
		m.setError(err)
		return
	}
	sourceDB, err := utils.GetMysqlClient(sourceConfig)
	if err != nil {
		m.setError(err)
		return
	}

	// 连接目标源
	targetConfig, err := m.TargetConfig.UnmarshalMysqlConfig()
	if err != nil {
		m.setError(err)
		return
	}

	targetDB, err := utils.GetMysqlClient(targetConfig)
	if err != nil {
		m.setError(err)
		return
	}

	// 如果不存在创建表
	if !targetDB.Migrator().HasTable(m.TargetTable) {
		var showCreateTable showCreateTable
		err := sourceDB.Raw(fmt.Sprintf("SHOW CREATE TABLE %s", m.SourceTable)).First(&showCreateTable).Error
		if err != nil {
			m.setError(err)
			return
		}
		createDDl := strings.Replace(showCreateTable.CreateTable, m.SourceTable, m.TargetTable, 1)
		err = targetDB.Exec(createDDl).Error
		if err != nil {
			m.setError(err)
			return
		}

		// TODO: 删除外健依赖
		// 	SELECT
		// 	TABLE_NAME,
		// 	COLUMN_NAME,
		// 	CONSTRAINT_NAME,
		// 	REFERENCED_TABLE_NAME,
		// 	REFERENCED_COLUMN_NAME
		// FROM
		// 	information_schema.KEY_COLUMN_USAGE
		// WHERE
		// 	TABLE_SCHEMA = 'prod_da_pevc_20211216' AND
		// 	TABLE_NAME = 'deal_sources' AND
		// 	CONSTRAINT_NAME LIKE 'FK%';

		// 		ALTER TABLE table_name
		// DROP FOREIGN KEY constraint_name1,
		// DROP FOREIGN KEY constraint_name2,
		// DROP FOREIGN KEY constraint_name3;
	}

	// 查询总数
	var count int64
	err = sourceDB.Table(m.SourceTable).Count(&count).Error
	if err != nil {
		m.setError(err)
		return
	}
	m.TotalCount = count

	// 查询主键
	var pkm primaryKeyModel
	err = sourceDB.Raw(fmt.Sprintf("SHOW INDEX FROM  %s where Key_name = 'PRIMARY'", m.SourceTable)).First(&pkm).Error
	if err != nil {
		m.setError(err)
		return
	}

	// 开始同步数据
	for {
		select {
		case <-m.stopCh:
			// 暂停
			m.SyncStatus = models.SyncStatusPause
			return
		default:
			sourceData, err := getMysqlTableData(sourceDB, m.Next, m.BatchSize, m.SourceTable, pkm.PrimaryKey)
			if err != nil {
				m.setError(err)
				return
			}
			if len(sourceData) == 0 {
				// 结束
				m.SyncStatus = models.SyncStatusDone
				return
			}

			err = insetMysqlTableData(targetDB, m.TargetTable, pkm.PrimaryKey, sourceData)
			if err != nil {
				m.setError(err)
				return
			}
			m.Next = fmt.Sprint(sourceData[len(sourceData)-1][pkm.PrimaryKey])
			m.Batch++
		}
	}
}

func (m *MysqlJob) setError(err error) {
	// 错误
	m.SyncStatus = models.SyncStatusError
	m.Error = err
}

func getMysqlTableData(db *gorm.DB, next interface{}, limit int64, tableName, primaryKey string) (res []map[string]interface{}, err error) {
	if next != "" {
		err = db.Raw(fmt.Sprintf("select * from %s where %s  > ? order by %s limit ?", tableName, primaryKey, primaryKey), next, limit).Scan(&res).Error
		return
	}
	err = db.Raw(fmt.Sprintf("select * from %s order by %s limit ?", tableName, primaryKey), limit).Scan(&res).Error
	return
}

func insetMysqlTableData(db *gorm.DB, tableName string, primaryKey string, data []map[string]interface{}) error {
	keys := make([]string, 0)
	for key := range data[0] {
		keys = append(keys, key)
	}

	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: primaryKey}},
		DoUpdates: clause.AssignmentColumns(keys),
	}).Table(tableName).Create(&data).Error
}
