package core

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"
)

var stringToLogLevel map[string]logrus.Level

func init() {

	stringToLogLevel = map[string]logrus.Level{
		"DEBUG": logrus.DebugLevel,
		"INFO":  logrus.InfoLevel,
		"WARN":  logrus.WarnLevel,
		"ERROR": logrus.ErrorLevel,
		"PANIC": logrus.PanicLevel,
		"FATAL": logrus.FatalLevel,
		"TRACE": logrus.TraceLevel,
	}
}

// Returns an instance of the logger adding the fields found in the context.
func FromContext(ctx context.Context) *logrus.Entry {
	entry := logrus.WithFields(logrus.Fields{})
	if ctx == nil {
		return entry
	}

	fields := []string{
		RequestId,
		Version,
		Environment,
	}

	for _, field := range fields {
		value := ctx.Value(field)
		if value != nil {
			entry = entry.WithField(string(field), value)
		}
	}

	return entry
}

func SetLevel(level string) {
	loglevel, ok := stringToLogLevel[strings.ToUpper(level)]

	if ok {
		logrus.SetLevel(loglevel)
	}

}
