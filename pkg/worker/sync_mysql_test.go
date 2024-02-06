package worker

import (
	"fmt"
	"strings"
	"testing"

	"gitlab.mvalley.com/wind/rime-utils/internal/pkg/config"
	"gitlab.mvalley.com/wind/rime-utils/pkg/utils"
)

func TestGetCreateTableDDl(test *testing.T) {
	db, err := utils.InitMysql(config.MySQLConfiguration{
		Host:     "localhost",
		Port:     "3306",
		User:     "root",
		Password: "root",
		DBName:   "research_dataservice",
	})
	if err != nil {
		panic(err)
	}
	var showCreateTable showCreateTable
	err = db.Raw("SHOW CREATE TABLE reading_favorites").First(&showCreateTable).Error
	if err != nil {
		panic(err)
	}

	createDDl := strings.Replace(showCreateTable.CreateTable, "reading_favorites", "reading_favorites1", 1)
	err = db.Exec(createDDl).Error
	if err != nil {
		panic(err)
	}

}

func TestGetPrimaryKey(test *testing.T) {
	db, err := utils.InitMysql(config.MySQLConfiguration{
		Host:     "localhost",
		Port:     "3306",
		User:     "root",
		Password: "root",
		DBName:   "test",
	})
	if err != nil {
		panic(err)
	}

	var primaryKeyModel primaryKeyModel
	err = db.Raw(fmt.Sprintf("SHOW INDEX FROM  %s where Key_name = 'PRIMARY'", "users")).First(&primaryKeyModel).Error
	if err != nil {
		panic(err)
	}

	fmt.Println(primaryKeyModel.PrimaryKey)
}

func TestGetAndInsertMysqlData(test *testing.T) {
	db, err := utils.InitMysql(config.MySQLConfiguration{
		Host:     "localhost",
		Port:     "3306",
		User:     "root",
		Password: "root",
		DBName:   "research_dataservice",
	})
	if err != nil {
		panic(err)
	}

	next := ""
	batchSize := 2
	sourceTable := "reading_favorites"
	sourceData, err := getMysqlTableData(db, next, int64(batchSize), sourceTable, "report_id")
	if err != nil {
		panic(err)
	}
	// if len(sourceData) == 0 {
	// 	panic(err)
	// }
	err = insetMysqlTableData(db, "reading_favorites1", "report_id", sourceData)
	if err != nil {
		panic(err)
	}

}

func TestGetColumns(test *testing.T) {
	db, err := utils.InitMysql(config.MySQLConfiguration{
		Host:     "localhost",
		Port:     "3306",
		User:     "root",
		Password: "root",
		DBName:   "research_dataservice",
	})
	if err != nil {
		panic(err)
	}

	var columns []map[string]interface{}
	if err := db.Raw("DESC " + "reading_favorites").Scan(&columns).Error; err != nil {
		fmt.Println("Error getting columns:", err)
		return
	}

	fmt.Println(columns)
}
