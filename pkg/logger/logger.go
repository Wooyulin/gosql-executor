package logger

import (
	"log"
	"os"
)

type Logger struct {
	logger *log.Logger
}

func NewLogger() *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (l *Logger) Info(msg string) {
	l.logger.Printf("[INFO] %s", msg)
}

func (l *Logger) Error(msg string, err error) {
	l.logger.Printf("[ERROR] %s: %v", msg, err)
}

func (l *Logger) Fatal(msg string, err error) {
	l.logger.Fatalf("[FATAL] %s: %v", msg, err)
}
