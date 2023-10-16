package logs

import (
	"github.com/freddyfeng-fy/mucy-core/core"
	"github.com/freddyfeng-fy/mucy-core/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

type Config struct {
	Level      string `mapstructure:"level" json:"level" yaml:"level"`
	RootDir    string `mapstructure:"root_dir" json:"root_dir" yaml:"root_dir"`
	Filename   string `mapstructure:"filename" json:"filename" yaml:"filename"`
	Format     string `mapstructure:"format" json:"format" yaml:"format"`
	ShowLine   bool   `mapstructure:"show_line" json:"show_line" yaml:"show_line"`
	MaxBackups int    `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"`
	MaxSize    int    `mapstructure:"max_size" json:"max_size" yaml:"max_size"` // MB
	MaxAge     int    `mapstructure:"max_age" json:"max_age" yaml:"max_age"`    // day
	Compress   bool   `mapstructure:"compress" json:"compress" yaml:"compress"`
}

var (
	level   zapcore.Level // zap 日志等级
	options []zap.Option  // zap 配置项
)

func InitializeLog(config *Config) {
	// 创建根目录
	createRootDir(config)

	// 设置日志等级
	setLogLevel(config)

	if config.ShowLine {
		options = append(options, zap.AddCaller())
	}

	// 初始化 zap
	coreLog := zapcore.NewTee(getZapCore(config), getZapConsoleCore(config))
	core.App.Log = zap.New(coreLog, options...)
	core.App.Log.Info("log init success!")
}

func createRootDir(config *Config) {
	if ok, _ := utils.PathExists(config.RootDir); !ok {
		_ = os.Mkdir(config.RootDir, os.ModePerm)
	}
}

func setLogLevel(config *Config) {
	switch config.Level {
	case "debug":
		level = zap.DebugLevel
		options = append(options, zap.AddStacktrace(level))
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
		options = append(options, zap.AddStacktrace(level))
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
}

// 扩展 Zap
func getZapCore(config *Config) zapcore.Core {
	var encoder zapcore.Encoder

	// 调整编码器默认配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("[" + "2006-01-02 15:04:05.000" + "]"))
	}
	encoderConfig.EncodeLevel = func(l zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(l.String())
	}

	// 设置编码器
	if config.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	return zapcore.NewCore(encoder, getLogWriter(config), level)
}

// 扩展 Zap
func getZapConsoleCore(config *Config) zapcore.Core {
	var encoder zapcore.Encoder

	// 调整编码器默认配置
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("[" + "2006-01-02 15:04:05.000" + "]"))
	}
	encoderConfig.EncodeLevel = func(l zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(l.String())
	}

	// 设置编码器
	if config.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	return zapcore.NewCore(encoder, zapcore.Lock(os.Stdout), level)
}

// 使用 lumberjack 作为日志写入器
func getLogWriter(config *Config) zapcore.WriteSyncer {
	file := &lumberjack.Logger{
		Filename:   config.RootDir + "/" + config.Filename,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}

	return zapcore.AddSync(file)
}
