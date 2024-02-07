package server

import (
	"github.com/gin-gonic/gin"
	"gitlab.mvalley.com/wind/rime-utils/pkg/worker"
)

type WorkerInfoResponse struct {
	MaxQueueCount      int64 `json:"max_count"`
	MaxWorkerCount     int64 `json:"max_worker_count"`
	CurrentWorkerCount int64 `json:"current_worker_count"`
	RemainWorkerCount  int64 `json:"remain_worker_count"`
}

// 获取同步任务队列信息
func (s *Server) GetWorkerInfo(ctx *gin.Context, req EmptyResponse) (*WorkerInfoResponse, error) {

	current := int64(len(s.w.RunCount))
	return &WorkerInfoResponse{
		MaxQueueCount:      worker.MaxQueueCount,
		MaxWorkerCount:     worker.MaxWorkerCount,
		CurrentWorkerCount: current,
		RemainWorkerCount:  worker.MaxWorkerCount - current,
	}, nil
}
