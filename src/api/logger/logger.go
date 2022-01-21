package logger

import (
	"fmt"
	"sync"
	"time"

	"github.com/fatih/color"
)

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
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
	whiteBold       *color.Color = color.New(color.FgHiWhite, color.Bold)
	reset           *color.Color = color.New(color.Reset)
	currentLogLevel LogLevel     = LogLevelDebug
	mutex           *sync.Mutex  = &sync.Mutex{}
	OnMessage       *chan string = nil
)

func init() {
	onMessage := make(chan string, 1000)
	OnMessage = &onMessage
}

func SetLogLevel(level LogLevel) {
	currentLogLevel = level
}

func timeFormat() string {
	return time.Now().Format("15:04:05")
}

func Debug(component, value interface{}) {
	if currentLogLevel > LogLevelDebug {
		return
	}

	mutex.Lock()

	// Time
	darkGray.Print("[")
	blueBold.Print(timeFormat())
	reset.Print("")
	darkGray.Print("]")

	reset.Print(" ")

	// Log level
	darkGray.Print("[")
	blueBold.Print("DEBUG")
	reset.Print("")
	darkGray.Print("]")

	reset.Print(" ")

	// Component
	darkGray.Print("[")
	whiteBold.Print(component)
	reset.Print("")
	darkGray.Print("]")

	reset.Print(" ")
	reset.Println(value)

	mutex.Unlock()

	*OnMessage <- fmt.Sprintf("[%s] [DEBUG] [%s] %v\n", timeFormat(), component, value)
}

func Debugf(component, format string, args ...interface{}) {
	if currentLogLevel > LogLevelDebug {
		return
	}

	mutex.Lock()

	// Time
	darkGray.Print("[")
	blueBold.Print(timeFormat())
	reset.Print("")
	darkGray.Print("]")

	reset.Print(" ")

	// Log level
	darkGray.Print("[")
	blueBold.Print("DEBUG")
	reset.Print("")
	darkGray.Print("]")

	reset.Print(" ")

	// Component
	darkGray.Print("[")
	whiteBold.Print(component)
	reset.Print("")
	darkGray.Print("]")

	reset.Print(" ")
	reset.Printf(format, args...)

	mutex.Unlock()

	*OnMessage <- fmt.Sprintf("[%s] [DEBUG] [%s] %s", timeFormat(), component, fmt.Sprintf(format, args...))
}

func Info(component, value interface{}) {
	if currentLogLevel > LogLevelInfo {
		return
	}

	mutex.Lock()

	// Time
	darkGray.Print("[")
	blueBold.Print(timeFormat())
	reset.Print("")
	darkGray.Print("]")

	reset.Print(" ")

	// Log level
	darkGray.Print("[")
	greenBold.Print("INFO")
	reset.Print("")
	darkGray.Print("]")

	reset.Print(" ")

	// Component
	darkGray.Print("[")
	whiteBold.Print(component)
	reset.Print("")
	darkGray.Print("]")

	reset.Print(" ")
	reset.Println(value)

	mutex.Unlock()

	*OnMessage <- fmt.Sprintf("[%s] [INFO] [%s] %v\n", timeFormat(), component, value)
}

func Infof(component, format string, args ...interface{}) {
	if currentLogLevel > LogLevelInfo {
		return
	}

	mutex.Lock()

	// Time
	darkGray.Print("[")
	blueBold.Print(timeFormat())
	reset.Print("")
	darkGray.Print("]")

	reset.Print(" ")

	// Log level
	darkGray.Print("[")
	greenBold.Print("INFO")
	reset.Print("")
	darkGray.Print("]")

	reset.Print(" ")

	// Component
	darkGray.Print("[")
	whiteBold.Print(component)
	reset.Print("")
	darkGray.Print("]")

	reset.Print(" ")
	reset.Printf(format, args...)

	mutex.Unlock()

	*OnMessage <- fmt.Sprintf("[%s] [INFO] [%s] %s", timeFormat(), component, fmt.Sprintf(format, args...))
}

func Warn(component, value interface{}) {
	if currentLogLevel > LogLevelWarning {
		return
	}

	mutex.Lock()

	// Time
	darkGray.Print("[")
	blueBold.Print(timeFormat())
	reset.Print("")
	darkGray.Print("]")

	reset.Print(" ")

	// Log level
	darkGray.Print("[")
	yellowBold.Print("WARNING")
	reset.Print("")
	darkGray.Print("]")

	reset.Print(" ")

	// Component
	darkGray.Print("[")
	whiteBold.Print(component)
	reset.Print("")
	darkGray.Print("]")

	reset.Print(" ")
	reset.Println(value)

	mutex.Unlock()

	*OnMessage <- fmt.Sprintf("[%s] [WARNING] [%s] %v\n", timeFormat(), component, value)
}

func Warnf(component, format string, args ...interface{}) {
	if currentLogLevel > LogLevelInfo {
		return
	}

	mutex.Lock()

	// Time
	darkGray.Print("[")
	blueBold.Print(timeFormat())
	reset.Print("")
	darkGray.Print("]")

	reset.Print(" ")

	// Log level
	darkGray.Print("[")
	yellowBold.Print("WARNING")
	reset.Print("")
	darkGray.Print("]")

	reset.Print(" ")

	// Component
	darkGray.Print("[")
	whiteBold.Print(component)
	reset.Print("")
	darkGray.Print("]")

	reset.Print(" ")
	reset.Printf(format, args...)

	mutex.Unlock()

	*OnMessage <- fmt.Sprintf("[%s] [WARNING] [%s] %s", timeFormat(), component, fmt.Sprintf(format, args...))
}

func Error(component, value interface{}) {
	if currentLogLevel > LogLevelError {
		return
	}

	mutex.Lock()

	// Time
	darkGray.Print("[")
	blueBold.Print(timeFormat())
	reset.Print("")
	darkGray.Print("]")

	reset.Print(" ")

	// Log level
	darkGray.Print("[")
	redBold.Print("ERROR")
	reset.Print("")
	darkGray.Print("]")

	reset.Print(" ")

	// Component
	darkGray.Print("[")
	whiteBold.Print(component)
	reset.Print("")
	darkGray.Print("]")

	reset.Print(" ")
	reset.Println(value)

	mutex.Unlock()

	*OnMessage <- fmt.Sprintf("[%s] [ERROR] [%s] %v\n", timeFormat(), component, value)
}

func Errorf(component, format string, args ...interface{}) {
	if currentLogLevel > LogLevelError {
		return
	}

	mutex.Lock()

	// Time
	darkGray.Print("[")
	blueBold.Print(timeFormat())
	reset.Print("")
	darkGray.Print("]")

	reset.Print(" ")

	// Log level
	darkGray.Print("[")
	redBold.Print("ERROR")
	reset.Print("")
	darkGray.Print("]")

	reset.Print(" ")

	// Component
	darkGray.Print("[")
	whiteBold.Print(component)
	reset.Print("")
	darkGray.Print("]")

	reset.Print(" ")
	reset.Printf(format, args...)

	mutex.Unlock()

	*OnMessage <- fmt.Sprintf("[%s] [ERROR] [%s] %s", timeFormat(), component, fmt.Sprintf(format, args...))
}
