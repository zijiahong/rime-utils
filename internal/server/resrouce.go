package server

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.mvalley.com/wind/rime-utils/internal/pkg/config"
	"gitlab.mvalley.com/wind/rime-utils/pkg/models"
	"gitlab.mvalley.com/wind/rime-utils/pkg/utils"
)

type ResourceListRequest struct {
	ResourcePlatform models.ResourcePlatform `json:"resource_platform"`
	ResourceType     models.ResourceType     `json:"resource_type"`
}

type ResourceListResponse struct {
	ResourceList []DataResource `json:"resource_list"`
}

type DataResource struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Source string `json:"source"`
}

// 获取数据源列表
func (s *Server) GetResourceList(ctx *gin.Context, req ResourceListRequest) (*ResourceListResponse, error) {
	res := make([]DataResource, 0)
	resourceList, err := s.storage.GetDataResource(req.ResourceType, req.ResourcePlatform)
	if err != nil {
		return nil, err
	}
	for i := range resourceList {
		res = append(res, DataResource{
			ID:     resourceList[i].RecId,
			Name:   resourceList[i].Name,
			Source: string(resourceList[i].ResourceConfig),
		})
	}
	return &ResourceListResponse{
		ResourceList: res,
	}, nil
}

type GetMysqlDatabasesRequest struct {
	ResourceID string `json:"resource_id"`
}

type GetMysqlDatabasesResponse struct {
	Databases []string `json:"databases"`
}

// Mysql
func (s *Server) GetMysqlDatabases(ctx *gin.Context, req GetMysqlDatabasesRequest) (*GetMysqlDatabasesResponse, error) {
	ds, err := s.storage.GetDataResourceByID(req.ResourceID)
	if err != nil {
		return nil, err
	}
	if ds == nil {
		return nil, fmt.Errorf("datasource not fount")
	}

	var mc config.MySQLConfiguration
	err = json.Unmarshal([]byte(ds.ResourceConfig), &mc)
	if err != nil {
		return nil, err
	}

	mc.DBName = "information_schema"
	db, err := utils.GetMysqlClient(mc)
	if err != nil {
		return nil, err
	}

	// 执行 SQL 查询语句获取所有数据库
	var databases []string
	err = db.Raw("SHOW DATABASES").Scan(&databases).Error
	if err != nil {
		return nil, err
	}
	return &GetMysqlDatabasesResponse{
		Databases: databases,
	}, nil
}

type GetMysqlTablesRequest struct {
	ResourceID string `json:"resource_id"`
	DataBase   string `json:"data_base"`
}

type GetMysqlTablesResponse struct {
	Tables []string `json:"tables"`
}

func (s *Server) GetMysqlTables(ctx *gin.Context, req GetMysqlTablesRequest) (*GetMysqlTablesResponse, error) {
	ds, err := s.storage.GetDataResourceByID(req.ResourceID)
	if err != nil {
		return nil, err
	}
	if ds == nil {
		return nil, fmt.Errorf("datasource not fount")
	}

	var mc config.MySQLConfiguration
	err = json.Unmarshal([]byte(ds.ResourceConfig), &mc)
	if err != nil {
		return nil, err
	}

	mc.DBName = req.DataBase

	db, err := utils.GetMysqlClient(mc)
	if err != nil {
		return nil, err
	}

	// 执行 SQL 查询语句获取所有表
	var tables []struct {
		TableName string `gorm:"column:TABLE_NAME"`
	}
	err = db.Raw("SELECT TABLE_NAME FROM information_schema.tables WHERE TABLE_SCHEMA = ?", req.DataBase).Scan(&tables).Error
	if err != nil {
		return nil, err
	}

	res := &GetMysqlTablesResponse{}
	for i := range tables {
		res.Tables = append(res.Tables, tables[i].TableName)
	}
	return res, nil
}
