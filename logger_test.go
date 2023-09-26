package logger

import (
	"context"
	"fmt"
	"log/slog"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestMain(m *testing.M) {
	InitLogger(Config{Level: DebugLevel, Output: "stdout", Format: "json", Syslog: Syslog{
		ServiceName: "test",
		ProjectId:   "default",
	}})
	m.Run()
}

func TestLog(t *testing.T) {
	t.Log(IsLevelEnabled(DebugLevel))
	WithField("field", "field1").Debugf("debug,%d", 1)
	WithField("field2", "field22").WithField("field3", "field33").Debugln(1)
	Debugf("Debugf,%d", 1)
	Debugln(1)
	Infof("Infof,%d", 2)
	Infoln(2)
	Warnf("Warnf,%d", 3)
	Warnln(3)
	Errorf("Errorf,%d", 4)
	Errorln(4)
	Printf("Printf,%d", 2)
	Println("Println")
	Fatalf("Fatalf,%d", 5)
	Fatalln(5)
}

func TestNewLogger(t *testing.T) {
	l := NewLogger(Config{Level: DebugLevel, Output: "stdout", Format: "json"})
	l.Debugln(123)
	l.WithField("a", 1).Debugln(1234)
	l.WithField("b", 2).Debugln(12345)
}

func Test_Slog(t *testing.T) {
	slog.Debug("1", slog.Any("a", map[string]interface{}{"b": 11}), slog.Time("times", time.Now().Local()))
	slog.Debug("1", slog.Any("a", map[string]interface{}{"b": 11}), slog.Any("times", time.Now().Local()))
	slog.Info("1", slog.Any("a", map[string]interface{}{"b": 11}), slog.Time("times", time.Now().Local()))
	slog.Warn("1", slog.Any("a", map[string]interface{}{"b": 11}), slog.Time("times", time.Now().Local()))
	slog.Error("1", slog.Any("a", map[string]interface{}{"b": 11}), slog.Time("times", time.Now().Local()))
}

func Test_diff(t *testing.T) {
	logrus.WithField("a", map[string]interface{}{"b": 11}).WithField("url", "url1").Println("test")
	slog.With("a", map[string]interface{}{"b": 11}, "url", "url1").Debug("12")
	slog.With("v", 1).With("v2", 2).Debug("123")
}

func Test_ctx(t *testing.T) {
	ctx := context.Background()
	ctx = NewTagContext(ctx, "__request__")
	ctx = NewTraceIDContext(ctx, "traecid")

	t.Log(FromTagContext(ctx))
	t.Log(FromTraceIDContext(ctx))
}

func Test_syslog(t *testing.T) {
	ctx := context.Background()
	Debugln("a")
	DebugContext(ctx, "11")
	WithContext(ctx).Debugln("12")
	ctx = NewModuleContext(ctx, "m")
	WithContext(ctx).Debugln("13")
	l := NewData(map[string]interface{}{"a": 1})
	l.WithContext(ctx).Debugln(123)
	l.DebugContext(ctx, "def,%d", 1)
	l.InfoContext(ctx, "def,%d", 1)
	l.WarnContext(ctx, "def,%d", 1)
	l.ErrorContext(ctx, "def,%d", 1)
	ctx = NewProjectContext(ctx, "project")
	l.DebugContext(ctx, "def,%d", 1)
	l.InfoContext(ctx, "def,%d", 1)
	l.WarnContext(ctx, "def,%d", 1)
	l.ErrorContext(ctx, "def,%d", 1)
	type test struct {
		Field string `json:"field"`
	}
	l = NewData(test{Field: "field"})
	l.DebugContext(ctx, "def,%d", 1)
	DebugContext(ctx, "def,%d", 1)
	InfoContext(ctx, "def,%d", 1)
	WarnContext(ctx, "def,%d", 1)
	ErrorContext(ctx, "def,%d", 1)
}

func TestNewServiceContext(t *testing.T) {
	SetVersion("1.12")
	ctx := NewServiceContext(context.Background(), "sev")
	ctx = NewStackContext(ctx, fmt.Errorf("stack"))
	ctx = NewUserIDContext(ctx, "admin")
	WithContext(ctx).Println(12)
}
