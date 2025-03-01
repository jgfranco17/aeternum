package logging

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
)

var (
	stringToLogLevel map[string]logrus.Level
)

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

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// Create a custom format for the log message
	level := strings.ToUpper(entry.Level.String())
	timestamp := entry.Time.Format(time.TimeOnly)
	colorFunc := color.New(setOutputColorPerLevel(level)).SprintFunc()
	logMessage := fmt.Sprintf("[%s][%s] %s\n", colorFunc(level), timestamp, entry.Message)
	return []byte(logMessage), nil
}

func setOutputColorPerLevel(level string) color.Attribute {
	var selectedColor color.Attribute
	switch level {
	case "DEBUG":
		selectedColor = color.FgCyan
	case "INFO":
		selectedColor = color.FgGreen
	case "WARN", "WARNING":
		selectedColor = color.FgYellow
	case "ERROR", "PANIC", "FATAL":
		selectedColor = color.FgRed
	default:
		selectedColor = color.FgWhite
	}
	return selectedColor
}
