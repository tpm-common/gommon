package logger

import (
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
)

type Logger struct {
	logger logr.Logger
}

func New(service string) Logger {
	zerolog.TimestampFieldName = "timestamp"
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.LevelFieldName = ""

	logsink := zerolog.New(os.Stderr).With().
		Str("service", service).Timestamp().
		CallerWithSkipFrameCount(zerolog.CallerSkipFrameCount + 1).
		Logger()

	logger := zerologr.New(&logsink)

	return Logger{logger: logger}
}

func (l Logger) WithName(name string) Logger {
	return Logger{logger: l.logger.WithName(name)}
}

func (l Logger) WithFields(keysAndValues ...any) Logger {
	return Logger{logger: l.logger.WithValues(keysAndValues...)}
}

func (l Logger) WithLevel(level int) Logger {
	return Logger{logger: l.logger.V(level)}
}

func (l Logger) Info(msg string, keysAndValues ...any) {
	l.logger.Info(msg, keysAndValues...)
}

func (l Logger) Error(err error, msg string, keysAndValues ...any) {
	l.logger.Error(err, msg, keysAndValues...)
}
