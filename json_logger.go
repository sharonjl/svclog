package svclog

import (
  "bytes"
  "encoding/json"
  "fmt"
  "os"
  "time"
)

type jsonLogger struct {
  kv []interface{}
}

func (jsl jsonLogger) With(kv ...interface{}) Logger {
  N := len(kv)
  if N%2 != 0 {
    N = N - 1
  }
  return jsonLogger{kv: append(kv[:N], jsl.kv...)}
}

func (jsl jsonLogger) print() {
  buf := bytes.Buffer{}
  buf.WriteRune('{')
  for i := 0; i < len(jsl.kv); i += 2 {
    if i != 0 {
      buf.WriteString(",")
    }
    key, _ := json.Marshal(jsl.kv[i])
    buf.Write(key)
    buf.WriteString(":")
    val, _ := json.Marshal(jsl.kv[i+1])
    buf.Write(val)
  }
  buf.WriteRune('}')
  fmt.Println(buf.String())
}

func (jsl jsonLogger) Print(message string, kv ...interface{}) {
  args := append([]interface{}{
    "time", time.Now().UTC().Format(time.RFC3339),
    "message", message,
  }, kv...)
  jsl.With(args...).(jsonLogger).print()
}

func (jsl jsonLogger) Fatal(message string, kv ...interface{}) {
  jsl.Print(message, kv...)
  os.Exit(-1)
}

func NewJSONLogger() Logger {
  return jsonLogger{}
}
