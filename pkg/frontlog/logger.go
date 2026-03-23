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

// prettyWriteSyncer wraps a WriteSyncer to pretty-print JSON output
type prettyWriteSyncer struct {
	zapcore.WriteSyncer
}

func (p *prettyWriteSyncer) Write(b []byte) (n int, err error) {
	// Try to parse as JSON
	var jsonData map[string]interface{}
	if err := json.Unmarshal(b, &jsonData); err != nil {
		// If not valid JSON, write as-is
		return p.WriteSyncer.Write(b)
	}

	// Pretty-print the JSON with indentation
	prettyJSON, err := json.MarshalIndent(jsonData, "", "   ")
	if err != nil {
		// If formatting fails, write original
		return p.WriteSyncer.Write(b)
	}

	// Add newline at the end
	prettyJSON = append(prettyJSON, '\n')
	return p.WriteSyncer.Write(prettyJSON)
}

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

	// Create output writers with pretty-printing
	writers := []zapcore.WriteSyncer{&prettyWriteSyncer{zapcore.AddSync(os.Stdout)}}

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
		writers = append(writers, &prettyWriteSyncer{zapcore.AddSync(file)})
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

// LogNested logs a message with a nested object as a proper JSON field
func LogNested(level, msg, fieldName string, obj interface{}) {
	if Logger == nil {
		return
	}

	field := zap.Any(fieldName, obj)

	switch level {
	case "debug":
		Logger.Debug(msg, field)
	case "info":
		Logger.Info(msg, field)
	case "warn":
		Logger.Warn(msg, field)
	case "error":
		Logger.Error(msg, field)
	default:
		Logger.Info(msg, field)
	}
}

// LogNestedFields logs a message with multiple nested objects as JSON fields
func LogNestedFields(level, msg string, fields map[string]interface{}) {
	if Logger == nil {
		return
	}

	zapFields := make([]zap.Field, 0, len(fields))
	for key, value := range fields {
		zapFields = append(zapFields, zap.Any(key, value))
	}

	switch level {
	case "debug":
		Logger.Debug(msg, zapFields...)
	case "info":
		Logger.Info(msg, zapFields...)
	case "warn":
		Logger.Warn(msg, zapFields...)
	case "error":
		Logger.Error(msg, zapFields...)
	default:
		Logger.Info(msg, zapFields...)
	}
}

// LogNestedJSON logs a message with nested JSON string
func LogNestedJSON(level, msg, fieldName, jsonStr string) {
	if Logger == nil {
		return
	}

	var obj interface{}
	if err := json.Unmarshal([]byte(jsonStr), &obj); err != nil {
		Logger.Error("Failed to parse JSON string", zap.Error(err))
		return
	}

	LogNested(level, msg, fieldName, obj)
}

// InfoNested logs an info message with a nested object
func InfoNested(msg, fieldName string, obj interface{}) {
	LogNested("info", msg, fieldName, obj)
}

// DebugNested logs a debug message with a nested object
func DebugNested(msg, fieldName string, obj interface{}) {
	LogNested("debug", msg, fieldName, obj)
}

// WarnNested logs a warning message with a nested object
func WarnNested(msg, fieldName string, obj interface{}) {
	LogNested("warn", msg, fieldName, obj)
}

// ErrorNested logs an error message with a nested object
func ErrorNested(msg, fieldName string, obj interface{}) {
	LogNested("error", msg, fieldName, obj)
}

// LogStructured logs a message with structured nested data
func LogStructured(level, msg string, data interface{}) {
	if Logger == nil {
		return
	}

	// Marshal to JSON and back to ensure proper nested structure
	jsonData, err := json.Marshal(data)
	if err != nil {
		Logger.Error("Failed to marshal structured data", zap.Error(err))
		return
	}

	var structuredData map[string]interface{}
	if err := json.Unmarshal(jsonData, &structuredData); err != nil {
		Logger.Error("Failed to unmarshal structured data", zap.Error(err))
		return
	}

	fields := make([]zap.Field, 0, len(structuredData))
	for key, value := range structuredData {
		fields = append(fields, zap.Any(key, value))
	}

	switch level {
	case "debug":
		Logger.Debug(msg, fields...)
	case "info":
		Logger.Info(msg, fields...)
	case "warn":
		Logger.Warn(msg, fields...)
	case "error":
		Logger.Error(msg, fields...)
	default:
		Logger.Info(msg, fields...)
	}
}
