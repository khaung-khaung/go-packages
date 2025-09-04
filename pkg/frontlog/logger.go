package frontlog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

type LogConfig struct {
	Enabled         bool
	Name            string
	TimestampInName bool
}

// InitLogger initializes the global logger with custom encoding
func InitLogger(fileConfig LogConfig, level string) error {
	logLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		return err
	}

	// Create custom encoder config
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "severity",
		MessageKey:     "message",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  "stacktrace",
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
	}

	encoder := zapcore.NewJSONEncoder(encoderConfig)

	// Create output writers
	writers := []zapcore.WriteSyncer{zapcore.AddSync(os.Stdout)}

	if fileConfig.Enabled {
		filePath, err := renderLogPath(fileConfig.Name, fileConfig.TimestampInName)
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return err
		}

		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		writers = append(writers, zapcore.AddSync(file))
	}

	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(writers...),
		logLevel,
	)

	// Create logger with caller info
	Logger = zap.New(core, zap.AddCaller())
	return nil
}

func FormatJSON(message string, v interface{}) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Printf("Error formatting JSON: %v", err)

		return fmt.Sprintf("Error formatting JSON: %v", err)
	}

	jsonString := string(b)
	log.Printf(message+" Successfully formatted JSON:%s\n", jsonString)
	return jsonString
}

// LogObjectWithEmptyField logs any object under an empty field name
func LogObjectWithEmptyField(msg string, obj interface{}) {
	if Logger == nil {
		return
	}

	// Marshal the object to JSON
	_, err := json.Marshal(obj)
	if err != nil {
		Logger.Error("Failed to marshal object", zap.Error(err))
		return
	}

	// Create a raw JSON field with empty name
	emptyField := zap.Field{
		Key:       "",
		Type:      zapcore.ReflectType,
		Interface: obj,
	}

	// Log with the empty field
	Logger.Info(msg, emptyField)
}

func renderLogPath(tmpl string, includeTimestamp bool) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	data := struct {
		WorkingDir       string
		IncludeTimestamp bool
		Timestamp        time.Time
	}{
		WorkingDir:       wd,
		IncludeTimestamp: includeTimestamp,
		Timestamp:        time.Now(),
	}

	t := template.Must(template.New("path").Parse(tmpl))
	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// GetLogger returns the global logger instance
func GetLogger() *zap.Logger {
	return Logger
}
