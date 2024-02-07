package storage

import "gitlab.mvalley.com/wind/rime-utils/pkg/models"

func (s *Storage) UpdateSubTask(m models.SubTask) error {
	return s.DB.Updates(m).Error
}

func (s *Storage) GetTasks(limit, offset int) (res []models.Task, err error) {
	err = s.DB.Model(&models.Task{}).Limit(limit).Offset(offset).Scan(&res).Error
	return
}

func (s *Storage) GetSubTasksByParentId(taskID string) (res []models.SubTask, err error) {
	err = s.DB.Model(&models.SubTask{}).Where("parent_task_id = ?", taskID).Scan(&res).Error
	return
}

func (s *Storage) GetSubTasksByIds(taskIDs []string) (res []models.SubTask, err error) {
	err = s.DB.Model(&models.SubTask{}).Where("rec_id in ?", taskIDs).Scan(&res).Error
	return
}

func (s *Storage) SaveTaskWithTX(task models.Task, subTasks []models.SubTask) error {
	d := s.DB.Begin()

	err := d.Create(&task).Error
	if err != nil {
		d.Rollback()
		return err
	}
	err = d.CreateInBatches(subTasks, 100).Error
	if err != nil {
		d.Rollback()
		return err
	}
	d.Commit()
	return nil
}

func (s *Storage) UpdateTaskStatusWithTX(taskID string, syncStatus models.SyncStatus) error {
	d := s.DB.Begin()

	err := d.Save(&models.Task{BaseModel: models.BaseModel{RecId: taskID}, SyncStatus: syncStatus}).Error
	if err != nil {
		d.Rollback()
		return err
	}
	err = d.Save(&models.SubTask{ParentTaskID: taskID, SyncStatus: syncStatus}).Error
	if err != nil {
		d.Rollback()
		return err
	}
	d.Commit()
	return nil
}

func (s *Storage) UpdateSubTaskStatusWithTX(taskIDs []string, syncStatus models.SyncStatus) error {
	return s.DB.Model(&models.SubTask{}).Where("rec_id in ?", taskIDs).Update("sync_status", syncStatus).Error
}
