package logging

import (
	"strings"
	"testing"
	"time"

	"github.com/fatih/color"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSetLoggingLevelPerColor(t *testing.T) {
	expectedColorPerLevel := map[string]color.Attribute{
		// Valid cases
		"DEBUG": color.FgCyan,
		"INFO":  color.FgGreen,
		"WARN":  color.FgYellow,
		"ERROR": color.FgRed,
		// Sample invalid case hits default
		"INVALID": color.FgWhite,
	}
	for level, color := range expectedColorPerLevel {
		assert.Equal(t, color, setOutputColorPerLevel(level))
	}
}

func TestCustomFormatterFormat(t *testing.T) {
	// Create a CustomFormatter instance
	formatter := &CustomFormatter{}

	// Create a logrus Entry
	entry := &logrus.Entry{
		Logger:  logrus.New(),
		Data:    logrus.Fields{},
		Time:    time.Now(),
		Level:   logrus.InfoLevel,
		Message: "This is a test log message",
	}

	// Format the entry
	output, err := formatter.Format(entry)
	assert.NoError(t, err)
	outputStr := string(output)
	assert.Contains(t, outputStr, "INFO")
	assert.Contains(t, outputStr, entry.Message)

	expectedTimestamp := entry.Time.Format(time.TimeOnly)
	assert.Contains(t, outputStr, expectedTimestamp)
	colorFunc := color.New(color.FgGreen).SprintFunc()
	assert.Contains(t, outputStr, colorFunc("INFO"))

	// Ensure the string ends with a newline
	assert.True(t, strings.HasSuffix(outputStr, "\n"), "Log message should end with a newline")
}
