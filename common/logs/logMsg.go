package logs

import (
	"github.com/fatih/color"
	"log"
)

func LogMessage(level, message string) {
	coloredLevel := colorizeLogLevel(level)
	log.Printf("%s %s\n", coloredLevel, message)
}

func colorizeLogLevel(level string) string {
	switch level {
	case "INFO:":
		return color.BlueString(level)
	case "SUCCESS:":
		return color.GreenString(level)
	case "ERROR:":
		return color.RedString(level)
	default:
		return level
	}
}
