package logger

type LogLevel uint32

const (
	PanicLevel LogLevel = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
)

// Config 日志配置参数
type Config struct {
	Level         LogLevel
	Format        string
	Output        string
	OutputFile    string
	MaxAge        int   // 日志最大保存分钟
	RotationTime  int   // 日志分割分钟
	RotationSize  int64 // 日志分割文件大小
	RotationCount uint  // 日志保存个数
	Syslog        Syslog
}

type Syslog struct {
	ServiceName string
	ProjectId   string
}
