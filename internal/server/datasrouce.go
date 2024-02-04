package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gitlab.mvalley.com/wind/rime-utils/pkg/models"
)

type DataSourceRequest struct {
	SourcePlatform models.SourcePlatform `json:"source_platform"`
	SourceType     models.SourceType     `json:"source_type"`
}

type DataSourceResponse struct {
	DataSourceList []DataSource `json:"data_source_list"`
}

type DataSource struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	IP   string `json:"ip"`
}

// 获取数据源列表
func (s *Server) GetDataSource(ctx *gin.Context, req DataSourceRequest) (*DataSourceResponse, error) {
	res := make([]DataSource, 0)
	sourceList, err := s.storage.GetDataSource(req.SourceType, req.SourcePlatform)
	if err != nil {
		return nil, err
	}
	for i := range sourceList {
		res = append(res, DataSource{
			ID:   sourceList[i].RecId,
			Name: sourceList[i].Name,
			IP:   sourceList[i].Host,
		})
	}
	return &DataSourceResponse{
		DataSourceList: res,
	}, nil
}

type GetMysqlDatabasesRequest struct {
	DataSourceID string `json:"data_source_id"`
}

type GetMysqlDatabasesResponse struct {
	Databases []string `json:"databases"`
}

// Mysql
func (s *Server) GetMysqlDatabases(ctx *gin.Context, req GetMysqlDatabasesRequest) (*GetMysqlDatabasesResponse, error) {
	ds, err := s.storage.GetDataSourceByID(req.DataSourceID)
	if err != nil {
		return nil, err
	}
	if ds == nil {
		return nil, fmt.Errorf("datasource not fount")
	}

	db, err := s.storage.GetMysqlClient(ds.Host, ds.Port, ds.User, ds.Password, "information_schema")
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
	DataSourceID string `json:"data_source_id"`
	DataBase     string `json:"data_base"`
}

type GetMysqlTablesResponse struct {
	Tables []string `json:"tables"`
}

func (s *Server) GetMysqlTables(ctx *gin.Context, req GetMysqlTablesRequest) (*GetMysqlTablesResponse, error) {
	ds, err := s.storage.GetDataSourceByID(req.DataSourceID)
	if err != nil {
		return nil, err
	}
	if ds == nil {
		return nil, fmt.Errorf("datasource not fount")
	}
	db, err := s.storage.GetMysqlClient(ds.Host, ds.Port, ds.User, ds.Password, req.DataBase)
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
