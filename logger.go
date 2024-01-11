package logger

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync/atomic"

	"go.opentelemetry.io/otel/trace"
)

var defaultLogger atomic.Value

func init() {
	programLevel := new(slog.LevelVar)
	programLevel.Set(slog.LevelInfo)
	defaultLogger.Store(&Logger{logger: slog.Default(), levelVar: programLevel})
}

type Focus int

const (
	FocusNotice Focus = 1
)

// Default returns the default Logger.
func Default() *Logger { return defaultLogger.Load().(*Logger) }

// Logger slog官方库
type Logger struct {
	logger   *slog.Logger
	syslog   Syslog
	levelVar *slog.LevelVar
}

func (l *Logger) SetLevel(logLevel LogLevel) {
	l.levelVar.Set(getLevel(logLevel))
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

func SetLevel(logLevel LogLevel) {
	Default().SetLevel(logLevel)
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
	ExtraKey     = "key"
	GroupKey     = "group"
	SuggestKey   = "suggest"
	FocusKey     = "focus"
	DeviceKey    = "device"
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
	extraKey   struct{}
	groupKey   struct{}
	suggestKey struct{}
	focusKey   struct{}
	deviceKey  struct{}
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

func NewExtraKeyContext(ctx context.Context, val string) context.Context {
	return context.WithValue(ctx, extraKey{}, val)
}

func FromExtraKeyContext(ctx context.Context) string {
	v := ctx.Value(extraKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
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

func NewGroupContext(ctx context.Context, group string) context.Context {
	return context.WithValue(ctx, groupKey{}, group)
}

func FromGroupContext(ctx context.Context) string {
	v := ctx.Value(groupKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func NewSuggestContext(ctx context.Context, suggest string) context.Context {
	return context.WithValue(ctx, suggestKey{}, suggest)
}

func FromSuggestContext(ctx context.Context) string {
	v := ctx.Value(suggestKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func NewFocusContext(ctx context.Context, f Focus) context.Context {
	return context.WithValue(ctx, focusKey{}, f)
}

func FromFocusContext(ctx context.Context) Focus {
	v := ctx.Value(focusKey{})
	if v != nil {
		if s, ok := v.(Focus); ok {
			return s
		}
	}
	return 0
}

func NewDeviceContext(ctx context.Context, f string) context.Context {
	return context.WithValue(ctx, deviceKey{}, f)
}

func FromDeviceContext(ctx context.Context) string {
	v := ctx.Value(deviceKey{})
	if v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func NewErrorContext(ctx context.Context, err error) context.Context {
	if v, ok := err.(*LogError); ok {
		ctx = NewSuggestContext(ctx, v.Suggest)
		ctx = NewFocusContext(ctx, v.Focus)
		return ctx
	}
	return context.WithoutCancel(ctx)
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

func NewPGTContext(ctx context.Context, projectId, groupId, tableId string) context.Context {
	return NewGroupContext(NewExtraKeyContext(NewProjectContext(ctx, projectId), tableId), groupId)
}

func NewPGTMContext(ctx context.Context, projectId, groupId, tableId, module string) context.Context {
	return NewModuleContext(NewGroupContext(NewExtraKeyContext(NewProjectContext(ctx, projectId), tableId), groupId), module)
}

func NewPGTDContext(ctx context.Context, projectId, groupId, tableId, deviceId string) context.Context {
	return NewDeviceContext(NewGroupContext(NewExtraKeyContext(NewProjectContext(ctx, projectId), tableId), groupId), deviceId)
}

func NewPGTDMContext(ctx context.Context, projectId, groupId, tableId, deviceId, module string) context.Context {
	return NewModuleContext(NewDeviceContext(NewGroupContext(NewExtraKeyContext(NewProjectContext(ctx, projectId), tableId), groupId), deviceId), module)
}

func NewPTContext(ctx context.Context, projectId, tableId string) context.Context {
	return NewProjectContext(NewExtraKeyContext(ctx, tableId), projectId)
}

func NewPTMContext(ctx context.Context, projectId, tableId, module string) context.Context {
	return NewModuleContext(NewProjectContext(NewExtraKeyContext(ctx, tableId), projectId), module)
}

func NewGTContext(ctx context.Context, groupId, tableId string) context.Context {
	return NewGroupContext(NewExtraKeyContext(ctx, tableId), groupId)
}

func NewGTMContext(ctx context.Context, groupId, tableId, module string) context.Context {
	return NewModuleContext(NewGroupContext(NewExtraKeyContext(ctx, tableId), groupId), module)
}

func NewGTDContext(ctx context.Context, groupId, tableId, deviceId string) context.Context {
	return NewDeviceContext(NewGroupContext(NewExtraKeyContext(ctx, tableId), groupId), deviceId)
}

func NewGTDMContext(ctx context.Context, groupId, tableId, deviceId, module string) context.Context {
	return NewModuleContext(NewDeviceContext(NewGroupContext(NewExtraKeyContext(ctx, tableId), groupId), deviceId), module)
}

func NewTDContext(ctx context.Context, tableId, deviceId string) context.Context {
	return NewDeviceContext(NewExtraKeyContext(ctx, tableId), deviceId)
}

func NewTDMContext(ctx context.Context, tableId, deviceId, module string) context.Context {
	return NewModuleContext(NewDeviceContext(NewExtraKeyContext(ctx, tableId), deviceId), module)
}

func NewPMContext(ctx context.Context, projectId, module string) context.Context {
	return NewModuleContext(NewProjectContext(ctx, projectId), module)
}

func NewTMContext(ctx context.Context, tableId, module string) context.Context {
	return NewModuleContext(NewExtraKeyContext(ctx, tableId), module)
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
	if v := FromExtraKeyContext(ctx); v != "" {
		fields = append(fields, ExtraKey, v)
	}
	if v := FromGroupContext(ctx); v != "" {
		fields = append(fields, GroupKey, v)
	}
	if v := FromSuggestContext(ctx); v != "" {
		fields = append(fields, SuggestKey, v)
	}
	if v := FromFocusContext(ctx); v > 0 {
		fields = append(fields, FocusKey, v)
	}
	if v := FromDeviceContext(ctx); v != "" {
		fields = append(fields, DeviceKey, v)
	}
	return fields
}
