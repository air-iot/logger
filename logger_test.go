package logger

import (
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"
)

func Test_1(t *testing.T) {
	NewLogger(Log{
		Level:         5,
		Output:        "file",
		OutputFile:    "./log/test.log",
		RotationTime:  1,
		RotationCount: 1,
	})
	for i := 0; i < 100; i++ {
		Debugln(i)
	}

	time.Sleep(time.Minute * 2)
	for i := 0; i < 100; i++ {
		Debugln(i)
	}
	time.Sleep(time.Minute * 2)
	for i := 0; i < 100; i++ {
		Debugln(i)
	}
	println("结束")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-sc
}

func BenchmarkLogger(b *testing.B) {
	NewLogger(Log{
		Level:         5,
		Output:        "file",
		OutputFile:    "./log/test.log",
		RotationTime:  1,
		RotationCount: 1,
	})
	for i := 0; i < b.N; i++ {
		Debugln(i)
	}

}
