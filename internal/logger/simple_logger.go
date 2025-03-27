package logger

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type SimpleLogger struct {
	file   *os.File
	mu     sync.Mutex
	prefix string
}

func NewSimpleLogger(filePath string) (*SimpleLogger, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}

	return &SimpleLogger{
		file:   file,
		prefix: "[%s][%s] - %s\n", // [LEVEL][DATE] - MESSAGE
	}, nil
}

func (sl *SimpleLogger) Close() error {
	if sl.file != nil {
		return sl.file.Close()
	}
	return nil
}

func (sl *SimpleLogger) log(level, message string) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	now := time.Now().Format("2006-01-02")
	logLine := fmt.Sprintf(sl.prefix, level, now, message)

	if _, err := sl.file.WriteString(logLine); err != nil {
		fmt.Printf("failed to write log: %v\n", err)
		return
	}
	fmt.Println(logLine)
}

func (sl *SimpleLogger) Error(message string) {
	sl.log("ERR", message)
}

func (sl *SimpleLogger) Info(message string) {
	sl.log("INFO", message)
}
