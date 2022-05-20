package logger

import "go.uber.org/zap/zapcore"

func ZapLogHandler(entry zapcore.Entry) error {
	go func(paramEntry zapcore.Entry) {
		// handle log here
	}(entry)
	return nil
}
