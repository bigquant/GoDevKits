package logger

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func CreateZapFactory(loggerConfig *LoggerConfig, entry func(zapcore.Entry) error) (*zap.Logger, error) {

	if loggerConfig == nil {
		log.Fatal("config is nil, init config before create logger")
	}

	if loggerConfig.DebugMode {
		if logger, err := zap.NewDevelopment(zap.Hooks(entry)); err == nil {
			return logger, nil
		} else {
			log.Fatal("crate zap logger object failed, reason: " + err.Error())
			return nil, err
		}
	}

	// if run in production
	encoderConfig := zap.NewProductionEncoderConfig()

	var recordTimeFormat string
	switch loggerConfig.TimePrecision {
	case "millisecond":
		recordTimeFormat = "2006-01-02 15:04:05.000"
	case "second":
		fallthrough
	default:
		recordTimeFormat = "2006-01-02 15:04:05"
	}
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(recordTimeFormat))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// adpter to es log formate
	encoderConfig.TimeKey = "created_at"

	var encoder zapcore.Encoder
	switch loggerConfig.TextFormat {
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig) // json formate
	case "console":
		fallthrough // normal formate
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig) // normal formate
	}

	// logs writer
	fileName := path.Join(
		loggerConfig.LogFileDir,
		loggerConfig.ServerLogFileName,
	)
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,                // log file path
		MaxSize:    loggerConfig.MaxSize,    // max log file size in MB
		MaxBackups: loggerConfig.MaxBackups, // max log file count
		MaxAge:     loggerConfig.MaxAge,     // max log file age in days
		Compress:   loggerConfig.Compress,   // compress log file
	}
	writer := zapcore.AddSync(lumberJackLogger)
	// write log to file. level >= zap.InfoLevel
	var allZapCore []zapcore.Core
	fileZapCore := zapcore.NewCore(
		encoder,
		writer,
		loggerConfig.FileWriterLevel,
	)
	allZapCore = append(allZapCore, fileZapCore)
	// write log to console, level >= zap.DebugLevel
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleZapCore := zapcore.NewCore(
		consoleEncoder,
		consoleDebugging,
		loggerConfig.ConsoleWriterLevel,
	)
	allZapCore = append(allZapCore, consoleZapCore)
	core := zapcore.NewTee(allZapCore...)
	return zap.New(
		core,
		zap.AddCaller(),
		zap.Hooks(entry),
		zap.AddStacktrace(zap.WarnLevel),
	), nil
}
