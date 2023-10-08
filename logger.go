package logger

import (
	"context"
	"fmt"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"os"
	"sync/atomic"
)

var defaultLogger atomic.Value

func init() {
	defaultLogger.Store(&Logger{logger: slog.Default()})
}

// Default returns the default Logger.
func Default() *Logger { return defaultLogger.Load().(*Logger) }

// Logger slog官方库
type Logger struct {
	logger *slog.Logger
	syslog Syslog
}

func (l *Logger) IsLevelEnabled(lev LogLevel) bool {
	return l.logger.Enabled(context.Background(), getLevel(lev))
}

func (l *Logger) Debugf(format string, args ...any) {
	l.logger.Debug(fmt.Sprintf(format, args...))
}

func (l *Logger) Debugln(args ...any) {
	msg := fmt.Sprintln(args...)
	l.logger.Debug(msg[:len(msg)-1])
}

func (l *Logger) Infof(format string, args ...any) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Infoln(args ...any) {
	msg := fmt.Sprintln(args...)
	l.logger.Info(msg[:len(msg)-1])
}

func (l *Logger) Warnf(format string, args ...any) {
	l.logger.Warn(fmt.Sprintf(format, args...))
}

func (l *Logger) Warnln(args ...any) {
	msg := fmt.Sprintln(args...)
	l.logger.Warn(msg[:len(msg)-1])
}

func (l *Logger) Errorf(format string, args ...any) {
	l.logger.Error(fmt.Sprintf(format, args...))
}

func (l *Logger) Errorln(args ...any) {
	msg := fmt.Sprintln(args...)
	l.logger.Error(msg[:len(msg)-1])
}

func (l *Logger) Fatalf(format string, args ...any) {
	l.logger.Error(fmt.Sprintf(format, args...))
	os.Exit(1)
}

func (l *Logger) Fatalln(args ...any) {
	msg := fmt.Sprintln(args...)
	l.logger.Error(msg[:len(msg)-1])
	os.Exit(1)
}

func (l *Logger) Printf(format string, args ...any) {
	l.logger.Info(fmt.Sprintf(format, args...))
}

func (l *Logger) Println(args ...any) {
	msg := fmt.Sprintln(args...)
	l.logger.Info(msg[:len(msg)-1])
}

func (l *Logger) DebugContext(ctx context.Context, format string, args ...any) {
	l.WithField(getFields(ctx)...).Debugf(format, args...)
}

func (l *Logger) InfoContext(ctx context.Context, format string, args ...any) {
	l.WithField(getFields(ctx)...).Infof(format, args...)
}

func (l *Logger) WarnContext(ctx context.Context, format string, args ...any) {
	l.WithField(getFields(ctx)...).Warnf(format, args...)
}

func (l *Logger) ErrorContext(ctx context.Context, format string, args ...any) {
	l.WithField(getFields(ctx)...).Errorf(format, args...)
}

func (l *Logger) WithField(args ...any) *Logger {
	c := *l
	c.logger = l.logger.With(args...)
	return &c
}

func (l *Logger) NewData(data any) *Logger {
	return l.WithField(LogDataKey, data)
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	return l.WithField(getFields(ctx)...)
}

func IsLevelEnabled(level LogLevel) bool {
	return Default().IsLevelEnabled(level)
}

func Debugf(format string, args ...any) {
	Default().Debugf(format, args...)
}

func Debugln(args ...any) {
	Default().Debugln(args...)
}

func Infof(format string, args ...any) {
	Default().Infof(format, args...)
}

func Infoln(args ...any) {
	Default().Infoln(args...)
}

func Warnf(format string, args ...any) {
	Default().Warnf(format, args...)
}

func Warnln(args ...any) {
	Default().Warnln(args...)
}

func Errorf(format string, args ...any) {
	Default().Errorf(format, args...)
}

func Errorln(args ...any) {
	Default().Errorln(args...)
}

func Printf(format string, args ...any) {
	Default().Printf(format, args...)
}

func Println(args ...any) {
	Default().Println(args...)
}

func Fatalf(format string, args ...any) {
	Default().Fatalf(format, args...)
}

func Fatalln(args ...any) {
	Default().Fatalln(args...)
}

func WithField(args ...any) *Logger {
	return Default().WithField(args...)
}

func DebugContext(ctx context.Context, format string, args ...any) {
	Default().DebugContext(ctx, format, args...)
}

func InfoContext(ctx context.Context, format string, args ...any) {
	Default().InfoContext(ctx, format, args...)
}

func WarnContext(ctx context.Context, format string, args ...any) {
	Default().WarnContext(ctx, format, args...)
}

func ErrorContext(ctx context.Context, format string, args ...any) {
	Default().ErrorContext(ctx, format, args...)
}

func NewData(data any) *Logger {
	return WithField(LogDataKey, data)
}

func NewDataContext(ctx context.Context, data any) *Logger {
	fields := getFields(ctx)
	fields = append(fields, LogDataKey, data)
	return WithField(fields...)
}

//func withSyslogContext(ctx context.Context, data any) *Logger {
//	spanCtx := trace.SpanContextFromContext(ctx)
//	var traceID, spanID string
//	if spanCtx.HasTraceID() {
//		traceID = spanCtx.TraceID().String()
//	}
//	if spanCtx.HasSpanID() {
//		spanID = spanCtx.SpanID().String()
//	}
//	return WithField(logTypeKey, logTypeValue, logIdKey, primitive.NewObjectID().Hex(), ServiceKey, FromServiceContext(ctx),
//		ProjectIdKey, FromProjectContext(ctx), ModuleKey, FromModuleContext(ctx), LogDataKey, data, traceIdKey, traceID, spanIdKey, spanID)
//}

//func SyslogDataDebug(ctx context.Context, data any, format string, args ...any) {
//	withSyslogContext(ctx, data).Debugf(format, args...)
//}
//
//func SyslogDataInfo(ctx context.Context, data any, format string, args ...any) {
//	withSyslogContext(ctx, data).Infof(format, args...)
//}
//
//func SyslogDataWarn(ctx context.Context, data any, format string, args ...any) {
//	withSyslogContext(ctx, data).Warnf(format, args...)
//}
//
//func SyslogDataError(ctx context.Context, data any, format string, args ...any) {
//	withSyslogContext(ctx, data).Errorf(format, args...)
//}
//
//func SyslogDebug(ctx context.Context, format string, args ...any) {
//	SyslogDataDebug(ctx, "", format, args...)
//}
//
//func SyslogInfo(ctx context.Context, format string, args ...any) {
//	SyslogDataInfo(ctx, "", format, args...)
//}
//
//func SyslogWarn(ctx context.Context, format string, args ...any) {
//	SyslogDataWarn(ctx, "", format, args...)
//}
//
//func SyslogError(ctx context.Context, format string, args ...any) {
//	SyslogDataError(ctx, "", format, args...)
//}

const (
	TraceIDKey = "trace_id"
	UserIDKey  = "user_id"
	TagKey     = "tag"
	//VersionKey   = "version"
	StackKey     = "stack"
	ServiceKey   = "service"
	ModuleKey    = "module"
	ProjectIdKey = "projectId"
	LogDataKey   = "data"
	logIdKey     = "id"
	logTypeKey   = "logType"
	logTypeValue = "__syslog__"
	traceIdKey   = "traceId"
	spanIdKey    = "spanId"
)

var (
	version string
)

type (
	traceIDKey struct{}
	userIDKey  struct{}
	tagKey     struct{}
	stackKey   struct{}
	serviceKey struct{}
	moduleKey  struct{}
	projectKey struct{}
)

// SetVersion 设定版本
func SetVersion(v string) {
	version = v
}

func NewServiceContext(ctx context.Context, service string) context.Context {
	return context.WithValue(ctx, serviceKey{}, service)
}

func FromServiceContext(ctx context.Context) string {
	v := ctx.Value(serviceKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return Default().syslog.ServiceName
}

func NewModuleContext(ctx context.Context, module string) context.Context {
	return context.WithValue(ctx, moduleKey{}, module)
}

func FromModuleContext(ctx context.Context) string {
	v := ctx.Value(moduleKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func NewProjectContext(ctx context.Context, projectId string) context.Context {
	return context.WithValue(ctx, projectKey{}, projectId)
}

func FromProjectContext(ctx context.Context) string {
	v := ctx.Value(projectKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return Default().syslog.ProjectId
}

// NewTraceIDContext 创建跟踪ID上下文
func NewTraceIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

// FromTraceIDContext 从上下文中获取跟踪ID
func FromTraceIDContext(ctx context.Context) string {
	v := ctx.Value(traceIDKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// NewUserIDContext 创建用户ID上下文
func NewUserIDContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

// FromUserIDContext 从上下文中获取用户ID
func FromUserIDContext(ctx context.Context) string {
	v := ctx.Value(userIDKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// NewTagContext 创建Tag上下文
func NewTagContext(ctx context.Context, tag string) context.Context {
	return context.WithValue(ctx, tagKey{}, tag)
}

// FromTagContext 从上下文中获取Tag
func FromTagContext(ctx context.Context) string {
	v := ctx.Value(tagKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// NewStackContext 创建Stack上下文
func NewStackContext(ctx context.Context, stack error) context.Context {
	return context.WithValue(ctx, stackKey{}, stack)
}

// FromStackContext 从上下文中获取Stack
func FromStackContext(ctx context.Context) error {
	v := ctx.Value(stackKey{})
	if v != nil {
		if s, ok := v.(error); ok {
			return s
		}
	}
	return nil
}

// WithContext Use context create entry
func WithContext(ctx context.Context) *Logger {
	return Default().WithContext(ctx)
}

func getFields(ctx context.Context) []any {
	if ctx == nil {
		ctx = context.Background()
	}
	fields := make([]any, 0)
	if v := FromTraceIDContext(ctx); v != "" {
		fields = append(fields, TraceIDKey, v)
	}
	if v := FromUserIDContext(ctx); v != "" {
		fields = append(fields, UserIDKey, v)
	}
	if v := FromTagContext(ctx); v != "" {
		fields = append(fields, TagKey, v)
	}
	if v := FromStackContext(ctx); v != nil {
		fields = append(fields, StackKey, fmt.Sprintf("%+v", v))
	}
	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		fields = append(fields, traceIdKey, spanCtx.TraceID().String())
	}
	if spanCtx.HasSpanID() {
		fields = append(fields, spanIdKey, spanCtx.SpanID().String())
	}
	if v := FromServiceContext(ctx); v != "" {
		fields = append(fields, logTypeKey, logTypeValue, ServiceKey, v)
	}
	if v := FromModuleContext(ctx); v != "" {
		fields = append(fields, ModuleKey, v)
	}
	if v := FromProjectContext(ctx); v != "" {
		fields = append(fields, ProjectIdKey, v)
	}
	return fields
}
