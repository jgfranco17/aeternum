package outputs

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func PrintColoredMessage(textColor string, source string, message string, args ...any) {
	var selectedColor color.Attribute
	switch strings.ToLower(textColor) {
	case "green":
		selectedColor = color.FgGreen
	case "yellow":
		selectedColor = color.FgYellow
	case "red":
		selectedColor = color.FgRed
	case "blue":
		selectedColor = color.FgBlue
	case "cyan":
		selectedColor = color.FgCyan
	default:
		selectedColor = color.FgWhite
	}
	colorFunc := color.New(selectedColor).SprintFunc()
	fullMessage := fmt.Sprintf(message, args...)
	fmt.Printf("[%s] %s\n", colorFunc(source), fullMessage)
}

func PrintWarn(text string, args ...any) {
	yellow := color.New(color.FgYellow).SprintFunc()
	message := fmt.Sprintf(text, args...)
	fmt.Printf("[%s] %s\n", yellow("WARNING"), message)
}

func PrintError(text string, args ...any) {
	red := color.New(color.FgRed).SprintFunc()
	message := fmt.Sprintf(text, args...)
	fmt.Printf("[%s] %s\n", red("ERROR"), message)
}
