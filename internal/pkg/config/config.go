package config

type config struct {
	Port        string
	MySQLConfig MySQLConfiguration
}
type MySQLConfiguration struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	LogMode  MySQLLogMode
	Charset  string
	// DSN的loc参数(默认Local)，要与数据库中time_zone变量指定的时区信息一致
	TimeZone string
}

type MongoDBConfiguration struct {
	Host   string
	DBName string
	Debug  bool
}

type ESConfiguration struct {
	Host                         []string
	User                         string
	Password                     string
	ResponseHeaderTimeoutSeconds int
}

var CONFIG config

// MySQLLogMode ...
type MySQLLogMode string

// Console 使用 gorm 的 logger，打印漂亮的sql到控制台
// SlowQuery 使用自定义 logger.Logger,记录慢查询sql到日志
// None 关闭 log 功能
const (
	Console   MySQLLogMode = "console"
	SlowQuery MySQLLogMode = "slow_query"
	None      MySQLLogMode = "none"
)
