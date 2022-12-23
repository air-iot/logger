package logger

import (
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
)

// LogHook 日志钩子
type LogHook string

// IsGorm 是否是gorm钩子
func (h LogHook) IsGorm() bool {
	return h == "gormdb"
}

// IsMongo 是否是mongo钩子
func (h LogHook) IsMongo() bool {
	return h == "mongo"
}

// Log 日志配置参数
type Log struct {
	Level         int
	Format        string
	Output        string
	OutputFile    string
	MaxAge        int   // 日志最大保存分钟
	RotationTime  int   // 日志分割分钟
	RotationSize  int64 // 日志分割文件大小
	RotationCount uint  // 日志保存个数
	EnableHook    bool
	HookLevels    []string
	Hook          LogHook
	HookMaxThread int
	HookMaxBuffer int
}

// LogGormHook 日志gorm钩子配置
type LogGormHook struct {
	DBType       string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
	Table        string
}

// NewLogger 创建日志模块
func NewLogger(c Log) (func(), error) {
	SetLevel(c.Level)
	SetFormatter(c.Format)

	// 设定日志输出
	//var file *os.File
	if c.Output != "" {
		switch c.Output {
		case "stdout":
			SetOutput(os.Stdout)
		case "stderr":
			SetOutput(os.Stderr)
		case "file":
			if name := c.OutputFile; name != "" {
				//_ = os.MkdirAll(filepath.Dir(name), 0777)
				//
				//f, err := os.OpenFile(name, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
				//if err != nil {
				//	return nil, err
				//}
				//file = f
				writer, err := rotatelogs.New(
					name+".%Y%m%d%H%M",
					rotatelogs.WithLinkName(name), // 生成软链，指向最新日志文件
					rotatelogs.WithRotationSize(c.RotationSize),
					rotatelogs.WithRotationCount(c.RotationCount),
					rotatelogs.WithMaxAge(time.Duration(c.MaxAge)*time.Minute),             // 文件最大保存时间
					rotatelogs.WithRotationTime(time.Duration(c.RotationTime)*time.Minute), // 日志切割时间间隔
				)
				if err != nil {
					return nil, err
				}
				SetOutput(writer)
			}
		}
	}

	return func() {
		//if file != nil {
		//	file.Close()
		//}
	}, nil
}
