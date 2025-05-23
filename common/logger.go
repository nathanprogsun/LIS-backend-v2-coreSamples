package common

import (
	"fmt"
	"log" // use built-in logger only when zap logger creation has failed
	"time"

	"github.com/hibiken/asynq"
	"go-micro.dev/v4/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var ZapLogger *zap.Logger

type AuditLogEntry struct {
	EventID             string
	ServiceName         string
	ServiceType         string
	EventName           string
	EntityType          string
	EntityID            string
	User                string
	Entrypoint          string
	EntitySnapshot      string
	AttributeName       string
	AttributeValuePrior string
	AttributeValuePost  string
}

// init a zap logger
func InitZapLogger(level string) {
	logLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		log.Fatal(err)
	}
	var encoderConfig = zapcore.EncoderConfig{
		MessageKey: "message",

		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,

		TimeKey:    "time",
		EncodeTime: zapcore.ISO8601TimeEncoder,

		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,
	}
	// Init zap logger config with custom encoder config
	var zapConfig = zap.Config{
		Encoding:      "json",
		OutputPaths:   []string{"stdout"},
		EncoderConfig: encoderConfig,
		Level:         zap.NewAtomicLevelAt(logLevel),
	}
	ZapLogger, _ = zapConfig.Build(zap.AddCallerSkip(1))
}

func LogLatency(process string, start time.Time) {
	fields := []zapcore.Field{}
	latency := time.Now().Sub(start)
	fields = append(fields, zap.Duration("process_latency", latency))
	ZapLogger.Debug(process, fields...)
}

func InfoFields(msg string, fields ...zapcore.Field) {
	ZapLogger.Info(msg, fields...)
}

func ErrorFields(msg string, fields ...zapcore.Field) {
	ZapLogger.Error(msg, fields...)
}

func Infof(template string, args ...interface{}) {
	// in case logger isn't initialized yet
	if ZapLogger == nil {
		logger.Infof(template, args...)
		return
	}
	ZapLogger.Info(fmt.Sprintf(template, args...))
}

func Fatalf(msg string, err error) {
	// in case logger isn't initialized yet
	if ZapLogger == nil {
		logger.Fatalf(msg, err)
	}
	ZapLogger.Fatal(msg, zap.NamedError("fatal error", err))
}

func Debugf(template string, args ...interface{}) {
	// in case logger isn't initialized yet
	if ZapLogger == nil {
		logger.Debugf(template, args...)
		return
	}
	ZapLogger.Debug(fmt.Sprintf(template, args...))
}

func Errorf(msg string, err error) {
	// in case logger isn't initialized yet
	if ZapLogger == nil {
		logger.Error(msg, err)
		return
	}
	ZapLogger.Error(msg, zap.Error(err))
}

func Error(err error) {
	// in case logger isn't initialized yet
	if ZapLogger == nil {
		logger.Error(err)
		return
	}
	ZapLogger.Error("", zap.Error(err))
}

func Fatal(err error) {
	// in case logger isn't initialized yet
	if ZapLogger == nil {
		logger.Fatal(err)
	}
	ZapLogger.Fatal("", zap.NamedError("fatal error", err))
}

func LogTaskInfo(info *asynq.TaskInfo) {
	fields := []zapcore.Field{
		zap.String("task_id", info.ID),
		zap.String("task_queue", info.Queue),
		zap.String("task_type", info.Type),
		zap.Time("completed_at", info.CompletedAt),
		zap.Error(fmt.Errorf(info.LastErr)),
		zap.Int("retries", info.Retried),
		zap.Time("last_failed_at", info.LastFailedAt),
		zap.Bool("orphaned", info.IsOrphaned),
		zap.Time("next_process_at", info.NextProcessAt),
	}
	if Env.LogLevel == "debug" {
		fields = append(fields, zap.String("task_payload", string(info.Payload)))
	}
	ZapLogger.Info("asynq task", fields...)
}

func ErrorLogger(template string, args ...interface{}) {
	ZapLogger.Error(fmt.Sprintf(template, args...))
}

func RecordAuditLog(entry AuditLogEntry) {
	if ZapLogger == nil {
		log.Fatal("ZapLogger is not initialized. Unable to record audit log.")
	}
	fields := []zapcore.Field{
		zap.String("event_id", entry.EventID),
		zap.String("service_name", entry.ServiceName),
		zap.String("service_type", entry.ServiceType),
		zap.String("event_name", entry.EventName),
		zap.String("entity_type", entry.EntityType),
		zap.String("entity_id", entry.EntityID),
		zap.String("user", entry.User),
		zap.String("entrypoint", entry.Entrypoint),
		zap.String("entity_snapshot", entry.EntitySnapshot),
		zap.String("attribute_name", entry.AttributeName),
		zap.String("attribute_value_prior", entry.AttributeValuePrior),
		zap.String("attribute_value_post", entry.AttributeValuePost),
	}
	ZapLogger.Info("Audit Log", fields...)
}

func Info(msg string) {
	// in case logger isn't initialized yet
	if ZapLogger == nil {
		logger.Info(msg)
		return
	}
	ZapLogger.Info(msg)
}

func Warn(msg string) {
	// in case logger isn't initialized yet
	if ZapLogger == nil {
		logger.Warn(msg)
		return
	}
	ZapLogger.Warn(msg)
}
