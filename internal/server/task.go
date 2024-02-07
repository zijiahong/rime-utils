package server

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gitlab.mvalley.com/wind/rime-utils/pkg/models"
)

type GetTasksRequest struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type GetTasksResponse struct {
	Tasks []models.Task
}

// TODO:
func (s *Server) GetTasks(ctx *gin.Context, req GetTasksRequest) (*GetTasksResponse, error) {
	// 处理密码
	tasks, err := s.storage.GetTasks(req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}
	return &GetTasksResponse{Tasks: tasks}, nil
}

type GetSubTasksRequest struct {
	TaskID string `json:"task_id"`
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
}

type GetSubTasksResponse struct {
	SubTasks []models.SubTask
}

func (s *Server) GetSubTasks(ctx *gin.Context, req GetSubTasksRequest) (*GetSubTasksResponse, error) {
	// 处理密码
	subTasks, err := s.storage.GetSubTasksByParentId(req.TaskID)
	if err != nil {
		return nil, err
	}
	return &GetSubTasksResponse{SubTasks: subTasks}, nil
}

type CreateTaskRequest struct {
	SourceResourceID string              `json:"source_resource_id"`
	SourceDataBase   string              `json:"source_data_base"`
	TargetResourceID string              `json:"target_resource_id"`
	TargetDataBase   string              `json:"target_data_base"`
	SyncTables       map[string]string   `json:"sync_tables"`
	SourceType       models.ResourceType `json:"source_type"`
}

type EmptyResponse struct {
}

// 创建主任务同时创建子任务
func (s *Server) CreateTask(ctx *gin.Context, req CreateTaskRequest) (*EmptyResponse, error) {
	// TODO 判断数据库知否存在
	source, err := s.storage.GetDataResourceByID(req.SourceResourceID)
	if err != nil {
		return nil, err
	}

	sourceConfig, err := source.ResourceConfig.FillDataBase(req.SourceDataBase, req.SourceType)
	if err != nil {
		return nil, err
	}

	target, err := s.storage.GetDataResourceByID(req.TargetResourceID)
	if err != nil {
		return nil, err
	}

	targetConfig, err := target.ResourceConfig.FillDataBase(req.TargetDataBase, req.SourceType)
	if err != nil {
		return nil, err
	}

	// TODO: 主站检查

	var subTasks []models.SubTask
	taskID := uuid.New().String()
	task := models.Task{
		BaseModel: models.BaseModel{
			RecId:     taskID,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		SourceConfig: sourceConfig,
		TargetConfig: targetConfig,
		SyncStatus:   models.SyncStatusInit,
		SourceType:   req.SourceType,
	}
	// todo 支持一次性同步全部
	for sourceTable, TargetTable := range req.SyncTables {
		subTasks = append(subTasks, models.SubTask{
			BaseModel: models.BaseModel{
				RecId:     uuid.New().String(),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			ParentTaskID: taskID,
			SourceConfig: sourceConfig,
			TargetConfig: targetConfig,
			SourceTable:  sourceTable,
			TargetTable:  TargetTable,
			SyncStatus:   models.SyncStatusInit,
			SourceType:   req.SourceType,
		})
	}
	if len(subTasks) == 0 {
		return nil, errors.New("请选择要同步的表")
	}
	err = s.storage.SaveTaskWithTX(task, subTasks)
	if err != nil {
		return nil, err
	}

	s.w.AppendJobs(subTasks)

	return &EmptyResponse{}, nil
}

type TaskRequest struct {
	TaskID string `json:"task_id"`
}

// 删除主任务同时删除子任务
func (s *Server) DeleteTask(ctx *gin.Context, req TaskRequest) (*EmptyResponse, error) {
	subTasks, err := s.storage.GetSubTasksByParentId(req.TaskID)
	if err != nil {
		return nil, err
	}
	for i := range subTasks {
		s.w.Stop(subTasks[i].RecId)
	}

	err = s.storage.UpdateTaskStatusWithTX(req.TaskID, models.SyncStatusDelete)
	if err != nil {
		return nil, err
	}
	return &EmptyResponse{}, nil
}

type SubTasksRequest struct {
	TaskIDs []string `json:"task_ids"`
}

// 删除子任务
func (s *Server) DeleteSubTask(ctx *gin.Context, req SubTasksRequest) (*EmptyResponse, error) {
	for i := range req.TaskIDs {
		s.w.Stop(req.TaskIDs[i])
	}

	return &EmptyResponse{}, s.storage.UpdateSubTaskStatusWithTX(req.TaskIDs, models.SyncStatusDelete)
}

// 暂停子任务
func (s *Server) StopSubTasks(ctx *gin.Context, req SubTasksRequest) (*EmptyResponse, error) {
	for i := range req.TaskIDs {
		s.w.Stop(req.TaskIDs[i])
	}
	return &EmptyResponse{}, s.storage.UpdateSubTaskStatusWithTX(req.TaskIDs, models.SyncStatusPause)
}

// 恢复同步任务
func (s *Server) ResumeSubTask(ctx *gin.Context, req SubTasksRequest) (*EmptyResponse, error) {
	subTasks, err := s.storage.GetSubTasksByIds(req.TaskIDs)
	if err != nil {
		return nil, err
	}
	s.w.AppendJobs(subTasks)

	return &EmptyResponse{}, nil
}
