package sqls

import (
	"fmt"
	"github.com/freddyfeng-fy/mucy-core/core"
	"github.com/freddyfeng-fy/mucy-core/logs"
	"go.uber.org/zap"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Driver              string `mapstructure:"driver" json:"driver" yaml:"driver"`
	Host                string `mapstructure:"host" json:"host" yaml:"host"`
	Port                int    `mapstructure:"port" json:"port" yaml:"port"`
	Database            string `mapstructure:"database" json:"database" yaml:"database"`
	UserName            string `mapstructure:"username" json:"username" yaml:"username"`
	Password            string `mapstructure:"password" json:"password" yaml:"password"`
	Charset             string `mapstructure:"charset" json:"charset" yaml:"charset"`
	MaxIdleConns        int    `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns        int    `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"max_open_conns"`
	LogMode             string `mapstructure:"log_mode" json:"log_mode" yaml:"log_mode"`
	EnableFileLogWriter bool   `mapstructure:"enable_file_log_writer" json:"enable_file_log_writer" yaml:"enable_file_log_writer"`
	LogFilename         string `mapstructure:"log_filename" json:"log_filename" yaml:"log_filename"`
}

func InitializeDB(config *Config, logConfig *logs.Config, initTable ...interface{}) {
	// 根据驱动配置进行初始化
	switch config.Driver {
	case "mysql":
		core.App.DB = initMySqlGorm(config, logConfig, initTable...)
	case "postgres":
		core.App.DB = initPostgresGorm(config, logConfig, initTable...)
	default:
		core.App.DB = initMySqlGorm(config, logConfig, initTable...)
	}
}

// 初始化 mysql gorm.DB
func initMySqlGorm(config *Config, logConfig *logs.Config, initTable ...interface{}) *gorm.DB {
	gormConfig := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,                             // 禁用自动创建外键约束
		Logger:                                   getGormLogger(config, logConfig), // 使用自定义 Logger
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
		},
	}
	if config.Database == "" {
		return nil
	}
	dsn := config.UserName + ":" + config.Password + "@tcp(" + config.Host + ":" + strconv.Itoa(config.Port) + ")/" +
		config.Database + "?charset=" + config.Charset + "&parseTime=True&loc=Local"

	mysqlConfig := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         191,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   false, // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig); err != nil {
		core.App.Log.Error("mysql connect failed, err:", zap.Any("err", err))
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
		initTables(db, initTable...)
		return db
	}
}

func initPostgresGorm(config *Config, logConfig *logs.Config, initTable ...interface{}) *gorm.DB {
	gormConfig := &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,                             // 禁用自动创建外键约束
		Logger:                                   getGormLogger(config, logConfig), // 使用自定义 Logger
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "t_",
			SingularTable: true,
		},
	}

	if config.Database == "" {
		return nil
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai",
		config.Host, config.UserName, config.Password, config.Database, config.Port)
	if db, err := gorm.Open(postgres.Open(dsn), gormConfig); err != nil {
		core.App.Log.Error("mysql connect failed, err:", zap.Any("err", err))
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(config.MaxIdleConns)
		sqlDB.SetMaxOpenConns(config.MaxOpenConns)
		initTables(db, initTable...)
		return db
	}
}

func initTables(db *gorm.DB, models ...interface{}) {
	err := db.AutoMigrate(models...)
	if err != nil {
		core.App.Log.Error("migrate table failed", zap.Any("err", err))
		os.Exit(0)
	}
}

// 自定义 gorm Writer
func getGormLogWriter(config *Config, logConfig *logs.Config) logger.Writer {
	var writer io.Writer

	// 是否启用日志文件
	if config.EnableFileLogWriter {
		// 自定义 Writer
		writer = &lumberjack.Logger{
			Filename:   logConfig.RootDir + "/" + config.LogFilename,
			MaxSize:    logConfig.MaxSize,
			MaxBackups: logConfig.MaxBackups,
			MaxAge:     logConfig.MaxAge,
			Compress:   logConfig.Compress,
		}
	} else {
		// 默认 Writer
		writer = os.Stdout
	}
	return log.New(writer, "\r\n", log.LstdFlags)
}

func getGormLogger(config *Config, logConfig *logs.Config) logger.Interface {
	var logMode logger.LogLevel

	switch config.LogMode {
	case "silent":
		logMode = logger.Silent
	case "error":
		logMode = logger.Error
	case "warn":
		logMode = logger.Warn
	case "info":
		logMode = logger.Info
	default:
		logMode = logger.Info
	}

	return logger.New(getGormLogWriter(config, logConfig), logger.Config{
		SlowThreshold:             200 * time.Millisecond,      // 慢 SQL 阈值
		LogLevel:                  logMode,                     // 日志级别
		IgnoreRecordNotFoundError: false,                       // 忽略ErrRecordNotFound（记录未找到）错误
		Colorful:                  !config.EnableFileLogWriter, // 禁用彩色打印
	})
}
