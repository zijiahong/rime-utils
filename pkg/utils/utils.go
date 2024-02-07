package utils

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/elastic/go-elasticsearch/v7"
	cfg "gitlab.mvalley.com/wind/rime-utils/internal/pkg/config"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	mysqlDriver "gorm.io/driver/mysql"
	"gorm.io/gorm"
	gorm_logger "gorm.io/gorm/logger"
)

func InitMysql(config cfg.MySQLConfiguration) (*gorm.DB, error) {
	charset := "utf8"
	loc := "Local"
	if config.Charset != "" {
		charset = config.Charset
	}
	if config.TimeZone != "" {
		loc = config.TimeZone
	}

	url := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset="+charset+"&parseTime=True&loc="+loc+"&multiStatements=True",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	db, err := gorm.Open(mysqlDriver.Open(url), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(2)
	sqlDB.SetMaxOpenConns(20)
	if config.LogMode == cfg.None {
		db.Logger = db.Logger.LogMode(gorm_logger.Silent)
	} else {
		db.Logger = db.Logger.LogMode(gorm_logger.Info)
	}
	if err != nil {
		return nil, err
	}
	return db, nil
}

// InitMongoDB ...
func InitMongoDB(config cfg.MongoDBConfiguration) (*mongo.Database, error) {
	startedCommands := sync.Map{}
	logger := log.New(os.Stdout, "\r\n", log.LstdFlags)
	cmdMonitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			startedCommands.Store(evt.RequestID, evt.Command)
		},
		Succeeded: func(_ context.Context, evt *event.CommandSucceededEvent) {
			command, _ := startedCommands.Load(evt.RequestID)
			logger.Printf("Mongo Command:[%v ms] %v",
				float64(evt.DurationNanos)/float64(time.Millisecond),
				command,
			)
			startedCommands.Delete(evt.RequestID)
		},
		Failed: func(_ context.Context, evt *event.CommandFailedEvent) {
			command, _ := startedCommands.Load(evt.RequestID)
			logger.Printf("Mongo Command:[%v ms] %v Failure: %v \n",
				float64(evt.DurationNanos)/float64(time.Millisecond),
				command,
				evt.Failure,
			)
			startedCommands.Delete(evt.RequestID)
		},
	}
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI(config.Host).SetMaxPoolSize(200)
	if config.Debug {
		clientOptions.SetMonitor(cmdMonitor)
	}

	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}
	mongoDB := client.Database(config.DBName)
	return mongoDB, nil
}

func PerformMongoDBInsert(documents []interface{}, collection *mongo.Collection) error {

	if collection == nil {
		return errors.New("collection can not be nil")
	}
	if len(documents) != 0 {
		_, err := collection.InsertMany(context.TODO(), documents)
		if err != nil {
			return err
		}
	}
	return nil
}

// InitElasticsearch ...
func InitElasticsearch(config cfg.ESConfiguration) (*elasticsearch.Client, error) {
	var cfg = elasticsearch.Config{
		Addresses: config.Host,
		Username:  config.User,
		Password:  config.Password,
		Transport: &http.Transport{
			MaxIdleConnsPerHost:   10,
			ResponseHeaderTimeout: time.Duration(config.ResponseHeaderTimeoutSeconds) * time.Second,
			DialContext:           (&net.Dialer{Timeout: time.Second}).DialContext,
			TLSClientConfig: &tls.Config{
				MinVersion:         tls.VersionTLS11,
				InsecureSkipVerify: true,
			},
		},
	}

	esClient, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	ping, err := esClient.Ping()
	if err != nil {
		panic(err)
	}

	if ping.IsError() || ping.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("%v", ping))
	}
	return esClient, nil
}
