package logger

import (
	"encoding/json"
	"log"
	"os"
)

type LogLevel string

const (
	InfoLevel  LogLevel = "INFO"
	WarnLevel  LogLevel = "WARN"
	ErrorLevel LogLevel = "ERROR"
	FatalLevel LogLevel = "FATAL"
)

type LogEntry struct {
	Level   LogLevel `json:"level"`
	Message string   `json:"message"`
	Fields  any      `json:"fields,omitempty"`
}

var stdLogger = log.New(os.Stdout, "", 0)

func Log(level LogLevel, message string, fields any) {
	entry := LogEntry{
		Level:   level,
		Message: message,
		Fields:  fields,
	}
	b, err := json.Marshal(entry)
	if err != nil {
		stdLogger.Printf(`{"level":"ERROR","message":"Failed to marshal log entry","fields":{"error":"%v"}}`, err)
		return
	}
	stdLogger.Println(string(b))
}

func Info(message string, fields any)  { Log(InfoLevel, message, fields) }
func Warn(message string, fields any)  { Log(WarnLevel, message, fields) }
func Error(message string, fields any) { Log(ErrorLevel, message, fields) }
func Fatal(message string, fields any) { Log(FatalLevel, message, fields); os.Exit(1) }
