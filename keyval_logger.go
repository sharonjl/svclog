package svclog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type consoleColor string

const (
	ColorNone   consoleColor = "\033[0m"
	ColorRed    consoleColor = "\033[31m"
	ColorGreen  consoleColor = "\033[32m"
	ColorYellow consoleColor = "\033[33m"
	ColorBlue   consoleColor = "\033[34m"
	ColorPurple consoleColor = "\033[35m"
	ColorCyan   consoleColor = "\033[36m"
	ColorGray   consoleColor = "\033[37m"
	ColorWhite  consoleColor = "\033[97m"
)

type kvLogger struct {
	kv       []interface{}
	keyColor consoleColor
}

func (kvl kvLogger) With(kv ...interface{}) Logger {
	N := len(kv)
	if N%2 != 0 {
		N = N - 1
	}

	newKVL := kvLogger{
		kv:       append(kv[:N], kvl.kv...),
		keyColor: kvl.keyColor,
	}
	return newKVL
}

func (kvl kvLogger) print() {
	buf := bytes.Buffer{}
	for i := 0; i < len(kvl.kv); i += 2 {
		if i != 0 {
			buf.WriteString(" ")
		}
		if kvl.keyColor == ColorNone {
			buf.WriteString(fmt.Sprintf("%s", kvl.kv[i]))
		} else {
			buf.WriteString(string(kvl.keyColor))
			buf.WriteString(fmt.Sprintf("%s", kvl.kv[i]))
			buf.WriteString(string(ColorNone))
		}
		buf.WriteString("=")
		val, _ := json.Marshal(kvl.kv[i+1])
		buf.Write(val)
	}

	mu.Lock()
	defer mu.Unlock()
	fmt.Println(buf.String())
}

func (kvl kvLogger) Print(message string, kv ...interface{}) {
	args := append([]interface{}{
		"time", time.Now().UTC().Format(time.RFC3339),
		"message", message,
	}, kv...)
	kvl.With(args...).(kvLogger).print()
}

func (kvl kvLogger) Fatal(message string, kv ...interface{}) {
	kvl.Print(message, kv...)
	os.Exit(-1)
}

func NewKeyvalLogger(keyColor consoleColor) Logger {
	return kvLogger{keyColor: keyColor}
}
