package logger

import (
	"io"
	"log/slog"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(cfg Config) {
	programLevel := new(slog.LevelVar)
	programLevel.Set(getLevel(cfg.Level))
	slog.SetDefault(newSlog(cfg, programLevel))
	defaultLogger.Store(&Logger{logger: slog.Default(), levelVar: programLevel, syslog: cfg.Syslog})
}

func NewLogger(cfg Config) *Logger {
	programLevel := new(slog.LevelVar)
	programLevel.Set(getLevel(cfg.Level))
	return &Logger{logger: newSlog(cfg, programLevel), levelVar: programLevel, syslog: cfg.Syslog}
}

func newSlog(cfg Config, lev *slog.LevelVar) *slog.Logger {
	var out io.Writer
	out = os.Stdout
	if cfg.Output != "" {
		switch cfg.Output {
		case "stdout":
			out = os.Stdout
		case "stderr":
			out = os.Stderr
		case "file":
			if name := cfg.OutputFile; name != "" {
				out = &lumberjack.Logger{
					Filename:   cfg.OutputFile,
					MaxSize:    int(cfg.RotationSize), // megabytes
					MaxBackups: int(cfg.RotationCount),
					MaxAge:     cfg.MaxAge, //days
					Compress:   true,       // disabled by default
				}
			}
		}
	}
	//var programLevel = new(slog.LevelVar)
	//programLevel.Set(getLevel(cfg.Level))
	switch cfg.Format {
	case "json":
		return slog.New(slog.NewJSONHandler(out, &slog.HandlerOptions{Level: lev}))
	default:
		return slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{Level: lev}))
	}
}

func getLevel(lev LogLevel) slog.Level {
	switch lev {
	case PanicLevel, FatalLevel:
		fallthrough
	case ErrorLevel:
		return slog.LevelError
	case WarnLevel:
		return slog.LevelWarn
	case InfoLevel:
		return slog.LevelInfo
	case DebugLevel:
		return slog.LevelDebug
	default:
		return slog.LevelInfo
	}
}

// Deprecated: 使用slog
// NewLogger1 创建日志模块
func NewLogger1(c Config) (func(), error) {
	logrus.SetLevel(logrus.Level(c.Level))
	switch c.Format {
	case "json":
		logrus.SetFormatter(new(logrus.JSONFormatter))
	default:
		logrus.SetFormatter(new(logrus.TextFormatter))
	}
	// 设定日志输出
	//var file *os.File
	if c.Output != "" {
		switch c.Output {
		case "stdout":
			logrus.SetOutput(os.Stdout)
		case "stderr":
			logrus.SetOutput(os.Stderr)
		case "file":
			if name := c.OutputFile; name != "" {
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
				logrus.SetOutput(writer)
			}
		}
	}

	return func() {
		//if file != nil {
		//	file.Close()
		//}
	}, nil
}
