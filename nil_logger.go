package svclog

import "os"

type nilLogger struct {
}

func (n nilLogger) With(kv ...interface{}) Logger {
  return nilLogger{}
}

func (n nilLogger) Print(message string, kv ...interface{}) {
  // do nothing
}

func (n nilLogger) Fatal(message string, kv ...interface{}) {
  os.Exit(-1)
}

func NewNilLogger() Logger {
  return nilLogger{}
}
