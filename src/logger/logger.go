package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

type Logger struct {
	LogLevel int
}

func NewLogger(level int) *Logger {
	return &Logger{
		LogLevel: level,
	}
}

func (l *Logger) SetLogLevel(level int) {
	l.LogLevel = level
}

func (l Logger) PrintPrefix() {
	color.New(color.FgHiBlack).Print("[")
	color.New(color.Bold).Print(time.Now().Format("15:04:05"))
	color.New(color.FgHiBlack).Print("]")
	fmt.Print(" ")
}

func (l Logger) Debug(args ...interface{}) {
	if l.LogLevel > 0 {
		return
	}

	l.PrintPrefix()

	color.New(color.FgHiBlack).Print("[")
	color.New(color.FgHiBlue, color.Bold).Print("DEBUG")
	color.New(color.FgHiBlack).Print("]")
	fmt.Print(" ")

	fmt.Println(args...)
}

func (l Logger) Debugf(format string, args ...interface{}) {
	if l.LogLevel > 0 {
		return
	}

	l.PrintPrefix()

	color.New(color.FgHiBlack).Print("[")
	color.New(color.FgHiBlue, color.Bold).Print("DEBUG")
	color.New(color.FgHiBlack).Print("]")
	fmt.Print(" ")

	fmt.Printf(format, args...)
}

func (l Logger) Info(args ...interface{}) {
	if l.LogLevel > 1 {
		return
	}

	l.PrintPrefix()

	color.New(color.FgHiBlack).Print("[")
	color.New(color.FgHiGreen, color.Bold).Print("INFO")
	color.New(color.FgHiBlack).Print("]")
	fmt.Print(" ")

	fmt.Println(args...)
}

func (l Logger) Infof(format string, args ...interface{}) {
	if l.LogLevel > 1 {
		return
	}

	l.PrintPrefix()

	color.New(color.FgHiBlack).Print("[")
	color.New(color.FgHiGreen, color.Bold).Print("INFO")
	color.New(color.FgHiBlack).Print("]")
	fmt.Print(" ")

	fmt.Printf(format, args...)
}

func (l Logger) Warn(args ...interface{}) {
	if l.LogLevel > 2 {
		return
	}

	l.PrintPrefix()

	color.New(color.FgHiBlack).Print("[")
	color.New(color.FgHiYellow, color.Bold).Print("WARN")
	color.New(color.FgHiBlack).Print("]")
	fmt.Print(" ")

	fmt.Println(args...)
}

func (l Logger) Warnf(format string, args ...interface{}) {
	if l.LogLevel > 2 {
		return
	}

	l.PrintPrefix()

	color.New(color.FgHiBlack).Print("[")
	color.New(color.FgHiYellow, color.Bold).Print("WARN")
	color.New(color.FgHiBlack).Print("]")
	fmt.Print(" ")

	fmt.Printf(format, args...)
}

func (l Logger) Error(args ...interface{}) {
	l.PrintPrefix()

	color.New(color.FgHiBlack).Print("[")
	color.New(color.FgHiRed, color.Bold).Print("ERROR")
	color.New(color.FgHiBlack).Print("]")
	fmt.Print(" ")

	fmt.Println(args...)
}

func (l Logger) Errorf(format string, args ...interface{}) {
	l.PrintPrefix()

	color.New(color.FgHiBlack).Print("[")
	color.New(color.FgHiRed, color.Bold).Print("ERROR")
	color.New(color.FgHiBlack).Print("]")
	fmt.Print(" ")

	fmt.Printf(format, args...)
}

func (l Logger) Fatal(args ...interface{}) {
	l.PrintPrefix()

	color.New(color.FgHiBlack).Print("[")
	color.New(color.FgHiRed, color.Bold).Print("FATAL")
	color.New(color.FgHiBlack).Print("]")
	fmt.Print(" ")

	fmt.Println(args...)

	os.Exit(1)
}

func (l Logger) Fatalf(format string, args ...interface{}) {
	l.PrintPrefix()

	color.New(color.FgHiBlack).Print("[")
	color.New(color.FgHiRed, color.Bold).Print("FATAL")
	color.New(color.FgHiBlack).Print("]")
	fmt.Print(" ")

	fmt.Printf(format, args...)

	os.Exit(1)
}
