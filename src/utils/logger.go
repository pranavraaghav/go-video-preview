package utils

import (
	"go.uber.org/zap"
	"strings"
)

// StandardLogger enforces specific log message formats.
type StandardLogger struct {
	*zap.SugaredLogger
}

func NewLogger(config *Config) *StandardLogger {
	var cfg zap.Config

	cfg = zap.NewDevelopmentConfig()
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	return &StandardLogger{SugaredLogger: logger.Sugar()}
}

func (l *StandardLogger) Printf(format string, v ...[]any) {
	if strings.Contains(format, "failed") {
		l.Errorf(format, v)
	} else {
		l.Infof(format, v)
	}
}
