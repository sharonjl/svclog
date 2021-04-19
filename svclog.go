package svclog

import "sync"

type Logger interface {
	With(kv ...interface{}) Logger
	Print(message string, kv ...interface{})
	Fatal(message string, kv ...interface{})
}

var (
	mu            = sync.Mutex{}
	defaultLogger = NewKeyvalLogger(ColorNone)
)

func With(kv ...interface{}) Logger {
	mu.Lock()
	defer mu.Unlock()
	return defaultLogger.With(kv...)
}

func Print(message string, kv ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	defaultLogger.Print(message, kv...)
}

func SetLogger(logger Logger) {
	mu.Lock()
	defer mu.Unlock()
	defaultLogger = logger
}
