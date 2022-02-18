package log

import (
	"time"

	"github.com/fatih/color"
)

const (
	LogLevelDebug int = iota
	LogLevelInfo
	LogLevelWarning
	LogLevelError
)

var (
	darkGray        *color.Color = color.New(color.FgHiBlack)
	blueBold        *color.Color = color.New(color.FgHiBlue, color.Bold)
	greenBold       *color.Color = color.New(color.FgHiGreen, color.Bold)
	yellowBold      *color.Color = color.New(color.FgHiYellow, color.Bold)
	redBold         *color.Color = color.New(color.FgHiRed, color.Bold)
	reset           *color.Color = color.New(color.Reset)
	currentLogLevel int          = LogLevelDebug
)

func SetLogLevel(level int) {
	currentLogLevel = level
}

func timeFormat() string {
	return time.Now().Format("Mon 15:04:05")
}

func Debug(value interface{}) {
	if currentLogLevel > LogLevelDebug {
		return
	}

	darkGray.Print(timeFormat())
	reset.Print(" ")
	blueBold.Print("[DEBUG]")
	reset.Print(" ")
	reset.Println(value)
}

func Debugf(format string, args ...interface{}) {
	if currentLogLevel > LogLevelDebug {
		return
	}

	darkGray.Print(timeFormat())
	reset.Print(" ")
	blueBold.Print("[DEBUG]")
	reset.Print(" ")
	reset.Printf(format, args...)
}

func Info(value interface{}) {
	if currentLogLevel > LogLevelInfo {
		return
	}

	darkGray.Print(timeFormat())
	reset.Print(" ")
	greenBold.Print("[INFO]")
	reset.Print(" ")
	reset.Println(value)
}

func Infof(format string, args ...interface{}) {
	if currentLogLevel > LogLevelInfo {
		return
	}

	darkGray.Print(timeFormat())
	reset.Print(" ")
	greenBold.Print("[INFO]")
	reset.Print(" ")
	reset.Printf(format, args...)
}

func Warn(value interface{}) {
	if currentLogLevel > LogLevelWarning {
		return
	}

	darkGray.Print(timeFormat())
	reset.Print(" ")
	yellowBold.Print("[WARN]")
	reset.Print(" ")
	reset.Println(value)
}

func Warnf(format string, args ...interface{}) {
	if currentLogLevel > LogLevelInfo {
		return
	}

	darkGray.Print(timeFormat())
	reset.Print(" ")
	yellowBold.Print("[WARN]")
	reset.Print(" ")
	reset.Printf(format, args...)
}

func Error(value interface{}) {
	if currentLogLevel > LogLevelError {
		return
	}

	darkGray.Print(timeFormat())
	reset.Print(" ")
	redBold.Print("[ERROR]")
	reset.Print(" ")
	reset.Println(value)
}

func Errorf(format string, args ...interface{}) {
	if currentLogLevel > LogLevelError {
		return
	}

	darkGray.Print(timeFormat())
	reset.Print(" ")
	redBold.Print("[ERROR]")
	reset.Print(" ")
	reset.Printf(format, args...)
}
