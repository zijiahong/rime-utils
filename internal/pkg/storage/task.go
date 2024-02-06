package storage

import "gitlab.mvalley.com/wind/rime-utils/pkg/models"

func (s *Storage) UpdateSubTask(m models.SubTask) error {
	return s.DB.Updates(m).Error
}
