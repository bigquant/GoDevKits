package logger

import (
	"go.uber.org/zap/zapcore"
)

type LoggerConfig struct {
	DebugMode          bool                 // default false
	FileWriterLevel    zapcore.LevelEnabler // default "info"
	ConsoleWriterLevel zapcore.LevelEnabler // default "info"
	TimePrecision      string               // default millisecond, second or millisecond
	MaxSize            int                  // default 10, max file size in MB
	MaxBackups         int                  // default 7, max old log files
	MaxAge             int                  // default 15, max days to kee logs
	Compress           bool                 // default false, compress logs rotated
	TextFormat         string               // default console, json or console
	LogFileDir         string               // default "/var/app/data/monitoringrulemanager/logs", log file dir
	ServerLogFileName  string               // default "server.log", server log file name
}

func NewLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		DebugMode:          false,
		TimePrecision:      "millisecond",
		MaxSize:            10,
		MaxBackups:         7,
		MaxAge:             15,
		Compress:           false,
		TextFormat:         "console",
		LogFileDir:         "/var/app/log",
		ServerLogFileName:  "server.log",
		FileWriterLevel:    zapcore.InfoLevel,
		ConsoleWriterLevel: zapcore.InfoLevel,
	}
}
