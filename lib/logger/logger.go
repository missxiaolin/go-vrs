package logger

import (
	"github.com/sirupsen/logrus"
	"os"
	"time"
)

func New() *logrus.Logger {
	l := &logrus.Logger{
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: time.RFC3339Nano,
		},
		Out: 		os.Stderr,
		Level:		logrus.InfoLevel,
		Hooks:		make(logrus.LevelHooks),
		ExitFunc: 	os.Exit,
	}

	l.AddHook(&hook{})

	return l
}
