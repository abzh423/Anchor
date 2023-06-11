package logger

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/fatih/color"
)

var (
	// EnableDebug determines whether Debug() or Debugf() should be enabled.
	EnableDebug bool = strings.ToLower(os.Getenv("DEBUG")) == "true"
	// LastLogTime is used to determine whether or not to create a new log file
	// when processing the queue by comparing the current time to this value.
	LastLogTime time.Time = time.Now()
	// ActiveLogFile is the currently opened log file which is only used for
	// a single day. When logging the next day, management Goroutine will 
	// create a new file.
	ActiveLogFile *os.File = nil
	// Queue is the currently queue of log messages to be processed. This is
	// used to prevent race condition while writing to the standard output, as
	// well as writing to the current log file. Any log method may write to this
	// channel, but only the log management Goroutine may read from it.
	Queue chan string = make(chan string)
	// WhiteColor is the color white.
	WhiteColor *color.Color = color.New(color.FgHiWhite)
	// GreenColor is the color green.
	GreenColor *color.Color = color.New(color.FgHiGreen, color.Bold)
	// YellowColor is the color yellow.
	YellowColor *color.Color = color.New(color.FgHiYellow, color.Bold)
	// RedColor is the color red.
	RedColor *color.Color = color.New(color.FgHiRed, color.Bold)
	// BlueColor is the color blue.
	BlueColor *color.Color = color.New(color.FgHiBlue, color.Bold)
	// AnsiRegExp is a regular expression for finding all ANSI escape codes in
	// a string.
	AnsiRegExp *regexp.Regexp = regexp.MustCompile(`[\\u001b\\u009b][[()#;?]*(?:[0-9]{1,4}(?:;[0-9]{0,4})*)?[0-9A-ORZcf-nqry=><]`)
)

// Initialize is responsible for creating the logs directory and calling the
// function that compresses any old logs if this directory already exists.
func Initialize() (err error) {
	// Creates the logs directory in the current working directory. MkdirAll is
	// used to prevent errors when the directory may already exist.
	if err = os.MkdirAll("logs", os.ModePerm); err != nil {
		return
	}

	if err = CompressOldLogs(); err != nil {
		return
	}

	if ActiveLogFile, err = os.OpenFile(path.Join("logs", LastLogTime.Format("01-02-2006.log")), os.O_CREATE|os.O_RDWR, os.ModePerm); err != nil {
		return
	}

	go StartLogManagementGoroutine()

	return
}

// CompressOldLogs searches the logs directory for any log files that were
// created any date before today, and compresses them to save space since log
// files may grow to a very large size.
func CompressOldLogs() error {
	files, err := os.ReadDir("logs")

	if err != nil {
		return err
	}

	todayFile := fmt.Sprintf("%s.log", time.Now().Format("01-02-2006"))

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".log.gz") || file.Name() == todayFile {
			continue
		}

		f, err := os.Open(path.Join("logs", file.Name()))

		if err != nil {
			return err
		}

		out, err := os.OpenFile(path.Join("logs", file.Name()+".gz"), os.O_CREATE|os.O_APPEND, os.ModePerm)

		if err != nil {
			return err
		}

		w := gzip.NewWriter(out)

		if _, err = io.Copy(w, f); err != nil {
			return err
		}

		if err = w.Close(); err != nil {
			return err
		}

		if err = f.Close(); err != nil {
			return err
		}

		if err = out.Close(); err != nil {
			return err
		}

		if err = os.Remove(path.Join("logs", file.Name())); err != nil {
			return err
		}
	}

	return nil
}

// StartLogManagementGoroutine is the Goroutine that processes all log entries
// and writes them to the current log file.
func StartLogManagementGoroutine() {
	var err error

	for line := range Queue {
		os.Stdout.WriteString(line + "\n")

		now := time.Now()

		if now.Year() != LastLogTime.Year() || now.Month() != LastLogTime.Month() || now.Day() != LastLogTime.Day() {
			if err = ActiveLogFile.Close(); err != nil {
				panic(err)
			}

			if ActiveLogFile, err = os.OpenFile(path.Join("logs", now.Format("01-02-2006.log")), os.O_CREATE|os.O_RDWR, os.ModePerm); err != nil {
				return
			}
		}

		if _, err = ActiveLogFile.WriteString(stripansi.Strip(line) + "\n"); err != nil {
			panic(err)
		}

		LastLogTime = now
	}
}

// GetTimeFormat returns the current time and day of the week used in the start
// of each line of the logs.
func GetTimeFormat() string {
	return time.Now().Format("Mon 15:04:05")
}

// Debug logs new values to the standard output as well as the log file.
func Debug(value ...interface{}) {
	Queue <- fmt.Sprintf(
		"%s %s %s",
		WhiteColor.Sprint(GetTimeFormat()),
		BlueColor.Sprint("[DEBUG]"),
		fmt.Sprint(value...),
	)
}

// Debugf logs new values to the standard output as well as the log file using
// the formatting string provided.
func Debugf(format string, args ...interface{}) {
	Queue <- fmt.Sprintf(
		"%s %s %s",
		WhiteColor.Sprint(GetTimeFormat()),
		BlueColor.Sprint("[DEBUG]"),
		fmt.Sprintf(format, args...),
	)
}

// Info logs new values to the standard output as well as the log file.
func Info(value ...interface{}) {
	Queue <- fmt.Sprintf(
		"%s %s %s",
		WhiteColor.Sprint(GetTimeFormat()),
		GreenColor.Sprint("[INFO]"),
		fmt.Sprint(value...),
	)
}

// Infof logs new values to the standard output as well as the log file using
// the formatting string provided.
func Infof(format string, args ...interface{}) {
	Queue <- fmt.Sprintf(
		"%s %s %s",
		WhiteColor.Sprint(GetTimeFormat()),
		GreenColor.Sprint("[INFO]"),
		fmt.Sprintf(format, args...),
	)
}

// Warn logs new values to the standard output as well as the log file.
func Warn(value ...interface{}) {
	Queue <- fmt.Sprintf(
		"%s %s %s",
		WhiteColor.Sprint(GetTimeFormat()),
		YellowColor.Sprint("[WARN]"),
		fmt.Sprint(value...),
	)
}

// Warnf logs new values to the standard output as well as the log file using
// the formatting string provided.
func Warnf(format string, args ...interface{}) {
	Queue <- fmt.Sprintf(
		"%s %s %s",
		WhiteColor.Sprint(GetTimeFormat()),
		YellowColor.Sprint("[WARN]"),
		fmt.Sprintf(format, args...),
	)
}

// Error logs new values to the standard output as well as the log file.
func Error(value ...interface{}) {
	Queue <- fmt.Sprintf(
		"%s %s %s",
		WhiteColor.Sprint(GetTimeFormat()),
		RedColor.Sprint("[ERROR]"),
		fmt.Sprint(value...),
	)
}

// Errorf logs new values to the standard output as well as the log file using
// the formatting string provided.
func Errorf(format string, args ...interface{}) {
	Queue <- fmt.Sprintf(
		"%s %s %s",
		WhiteColor.Sprint(GetTimeFormat()),
		RedColor.Sprint("[ERROR]"),
		fmt.Sprintf(format, args...),
	)
}

// Stop closes all log files and ends the log management Goroutine.
// TODO rewrite log manager to allow queue to be empty and processed before the
// channel is closed and the Goroutine is killed. sync.WaitGroup may be the fix
func Stop() error {
	close(Queue)

	if ActiveLogFile != nil {
		if err := ActiveLogFile.Close(); err != nil {
			return err
		}
	}

	return nil
}
