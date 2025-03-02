package log

import (
	"io"
	"log"
	"os"
	"sync"
)

// 两个logger，分别用于记录错误信息和普通信息。
var (
	errorLog = log.New(os.Stdout, "\033[31m[error]\033[0m ", log.LstdFlags|log.Lshortfile)
	infoLog  = log.New(os.Stdout, "\033[34m[info ]\033[0m ", log.LstdFlags|log.Lshortfile)
	loggers  = []*log.Logger{errorLog, infoLog}
	mu       sync.Mutex
)

var (
	Error  = errorLog.Println
	Errorf = errorLog.Printf
	Info   = infoLog.Println
	Infof  = infoLog.Printf
)

const (
	InfoLevel = iota
	ErrorLevel
	Disabled
)

// SetLevel 控制日志的输出级别
func SetLevel(level int) {
	mu.Lock()
	defer mu.Unlock()

	for _, logger := range loggers {
		logger.SetOutput(os.Stdout)
	}

	// 传入Disable时，丢弃错误日志
	if ErrorLevel < level {
		errorLog.SetOutput(io.Discard)
	}

	// 传入 Error或 Disabled时，丢弃普通日志
	if InfoLevel < level {
		infoLog.SetOutput(io.Discard)
	}
}
