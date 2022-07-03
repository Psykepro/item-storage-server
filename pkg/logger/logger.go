package logger

import (
	"fmt"
	"os"

	"github.com/Psykepro/item-storage-server/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger
type logger struct {
	cfg         *config.Logging
	sugarLogger *zap.SugaredLogger
}

// App Logger constructor
func NewLogger(cfg *config.Logging) *logger {
	return &logger{cfg: cfg}
}

// For mapping config logger to app logger levels
var loggerLevelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

// InitFileLogger initializing the logger for file with the parameters from the config
func (l *logger) InitFileLogger() {
	var encoderCfg zapcore.EncoderConfig
	if l.cfg.File.Mode == "Development" {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}
	var encoder zapcore.Encoder
	encoderCfg.MessageKey = "msg"
	encoderCfg.LevelKey = ""
	encoderCfg.CallerKey = ""
	encoderCfg.TimeKey = ""
	encoderCfg.NameKey = ""

	var logFile *os.File
	switch l.cfg.File.Encoding {
	case "console":
		logFile, _ = os.OpenFile("responses.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
		break
	case "json":
		logFile, _ = os.OpenFile("responses.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
		break
	default:
		panic(fmt.Sprintf("not supported logger encoding: '%s'. supported ones are: 'console' and 'json'"))
	}

	logWriter := zapcore.AddSync(logFile)
	l.initLogger(encoder, logWriter, l.getLoggerLevelFromCfg(l.cfg.File.Level))
}

// InitStdOutLogger initializing the logger for stdout with the parameters from the config
func (l *logger) InitStdOutLogger() {
	var (
		encoderCfg zapcore.EncoderConfig
		encoder    zapcore.Encoder
	)

	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.LevelKey = "LEVEL"
	encoderCfg.CallerKey = "CALLER"
	encoderCfg.TimeKey = "TIME"
	encoderCfg.NameKey = "NAME"
	encoderCfg.MessageKey = "MSG"

	switch l.cfg.StdOut.Encoding {
	case "console":
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
		break
	case "json":
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
		break
	default:
		panic(fmt.Sprintf("not supported logger encoding: '%s'. supported ones are: 'console' and 'json'"))
	}

	l.initLogger(encoder, zapcore.AddSync(os.Stdout), l.getLoggerLevelFromCfg(l.cfg.StdOut.Level))
}

// initLogger initializing the logger with the parameters from the config
func (l *logger) initLogger(encoder zapcore.Encoder, writeSyncer zapcore.WriteSyncer, logLevel zapcore.Level) {

	core := zapcore.NewCore(encoder, writeSyncer, zap.NewAtomicLevelAt(logLevel))

	logg := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	l.sugarLogger = logg.Sugar()
	if err := l.sugarLogger.Sync(); err != nil {
		l.sugarLogger.Error(err)
	}
}

// Utility methods

func (l *logger) getLoggerLevelFromCfg(levelFromCfg string) zapcore.Level {
	level, exist := loggerLevelMap[levelFromCfg]
	if !exist {
		return zapcore.DebugLevel
	}

	return level
}

// Logging methods

func (l *logger) Debug(args ...any) {
	l.sugarLogger.Debug(args...)
}

func (l *logger) Debugf(template string, args ...any) {
	l.sugarLogger.Debugf(template, args...)
}

func (l *logger) Info(args ...any) {
	l.sugarLogger.Info(args...)
}

func (l *logger) Infof(template string, args ...any) {
	l.sugarLogger.Infof(template, args...)
}

func (l *logger) Warn(args ...any) {
	l.sugarLogger.Warn(args...)
}

func (l *logger) Warnf(template string, args ...any) {
	l.sugarLogger.Warnf(template, args...)
}

func (l *logger) Error(args ...any) {
	l.sugarLogger.Error(args...)
}

func (l *logger) Errorf(template string, args ...any) {
	l.sugarLogger.Errorf(template, args...)
}

func (l *logger) DPanic(args ...any) {
	l.sugarLogger.DPanic(args...)
}

func (l *logger) DPanicf(template string, args ...any) {
	l.sugarLogger.DPanicf(template, args...)
}

func (l *logger) Panic(args ...any) {
	l.sugarLogger.Panic(args...)
}

func (l *logger) Panicf(template string, args ...any) {
	l.sugarLogger.Panicf(template, args...)
}

func (l *logger) Fatal(args ...any) {
	l.sugarLogger.Fatal(args...)
}

func (l *logger) Fatalf(template string, args ...any) {
	l.sugarLogger.Fatalf(template, args...)
}
