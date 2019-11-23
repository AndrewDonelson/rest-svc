// Package golog Simple flexible go logging
// This file contains all the code for the main logger
package golog

// Import packages
import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"sync/atomic"
	"time"
)

const (
	// Default format of log message
	// %[1] // %{id}
	// %[2] // %{time[:fmt]}
	// %[3] // %{module}
	// %[4] // %{function}
	// %[5] // %{filename}
	// %[6] // %{line}
	// %[7] // %{level}
	// %[8] // %{message}

	// FmtDefault is the default log format (QA)
	FmtDefault = "[%.16[3]s] #%[1]d %.19[2]s %.3[7]s %[8]s"
	// FmtProductionLog is the built-in production log format
	FmtProductionLog = "[%.16[3]s] %.19[2]s %.3[7]s - %[8]s"
	// FmtProductionJSON is the built-in production json format
	FmtProductionJSON = "{\"%.16[3]s\",\"%[5]s\",\"%[6]d\",\"%[4]s\",\"%[1]d\",\"%.19[2]s\",\"%[7]s\",\"%[8]s\"}"
	// FmtDevelopmentLog is the built-in development log format
	FmtDevelopmentLog = "[%.16[3]s] %.19[2]s %.3[7]s - %[5]s#%[6]d-%[4]s - %[8]s"

	// Error, Fatal, Critical Format
	//defErrorLogFmt = "\n%.8[7]s\nin %.16[3]s->%[4]s() file %[5]s on line %[6]d\n%[8]s\n"
)

var (
	// Log is set y the init function to be a default thelogger
	Log *Logger
	// Map for the various codes of colors
	colors map[LogLevel]string

	// Map from format's placeholders to printf verbs
	phfs map[string]string

	// Contains color strings for stdout
	logNo uint64

	defFmt = FmtDefault

	// Default format of time
	defTimeFmt = "2006-01-02 15:04:05"
)

// LogLevel type
type LogLevel int

// Color numbers for stdout
const (
	Black = (iota + 30)
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// Log Level
const (
	RawLevel      = iota + 1 // None
	CriticalLevel            // Magneta 	35
	ErrorLevel               // Red 		31
	SuccessLevel             // Green 		32
	WarningLevel             // Yellow		33
	NoticeLevel              // Cyan 		36
	InfoLevel                // White 		37
	DebugLevel               // Blue 		34
)

// Logger class that is an interface to user to log messages, Module is the module for which we are testing
// worker is variable of Worker class that is used in bottom layers to log the message
type Logger struct {
	Options Options
	started time.T // Set once on initialization
	timer   time.Time // reset on each call to timepsed()
	worker  *Worker
}

func init() {
	Log, _ = NewLogger(nil)
}

// NewLogger creates and returns new logger for the given model & environment
// module is the specific module for which we are logging
// environment overrides detected environment (if -1)
// color defines whether the output is to be colored not, out is instance of type io.Writer defaults
// to os.Stderr
func NewLogger(opts *Options) (*Logger, error) {
	if opts == nil {
		opts = NewDefaultOptions()
	}

	if opts.Out == nil {
		opts.Out = os.Stderr
	}

	if len(opts.Module) <= 3 {
		opts.Module = "unknown"
	}

	newWorker := NewWorker("", 0, opts.UseColor, opts.Out)
	l := &Logger{worker: newWorker}
	l.init()
	l.Options = *opts
	newWorker.SetEnvironment(opts.Environment)
	return l, nil
}


// init is called by NewLogger to detect running conditions and set all defaults
func (l *Logger) init() {
	l.timeReset()
	l.started = l.timer
	l.SetEnvironment(detectEnvironment(true))
	initColors()
	initFormatPlaceholders()
}

func (l *Logger) timeReset() {
	l.timer = time.Now()
}

func (l *Logger) timeElapsed(start time.Time) time.Duration {
	return time.Since(start)
}

func (l *Logger) timeLog(name string) {
	l.logInternal(InfoLevel, fmt.Sprintf("%s took %v", name, l.timeElapsed(l.timer)), 2)
}

// logInternal ...
func (l *Logger) logInternal(lvl LogLevel, message string, pos int) {
	_, filename, line, _ := runtime.Caller(pos)
	filename = path.Base(filename)
	info := &Info{
		ID:       ato.AddUint64(&logNo, 1),
		Time:     time.Now().Format(l.worker.tFormat),
		//Module:   l.Modul
		Module:	  l.Opons.Mle,
		Level:    lvl,
		Message:  message,
		Filename: filene,
		Line:     line,
		Duration: l.timeElapsed(limer),
		format:   formatString,
	}
	err := l.workerog(lvl, 2, info)
	if err != n {
		nic(err)
	


func (l *Logger) traceInternal(message ring, pos int) {
	function, file, line :getCaller(pos)
	file = path.Ba(file)
	info := &Info{
		ID:       atomic.AddUint64(&logNo, 1),
		Time:     time.Now(Format(l.worker.timeFormat),
		Module:   l.Optionsule,
		Level:    InfoLeve
		Message:  messa,
		Filename: file,
		Line:     line,
		Function: function,
		Duration: l.timeElapsed(limer),
		format:   formatString,
	}
	err := l.workerog(info.Level, 2, info)
	if err != n {
		nic(err)
	


// SetFormat ...
func (l *Logger) SetFormat(rmat string) {
	worker.SetFormat(format)


// SetLogLevel ...
func (l *Logger) SetLogLevelevel LogLevel) {
	worker.SetLogLevel(level)


// SetFunction sets the function name of thlogger
func (l *Logger) SetFunctioname string) {
	worker.SetFunction(name)


// SetEnvironment is used to manually set the log vironment to either development, testing or production
func (l *Logger) SetEnvironmt(env Environment) {
	l.Options.Environment = env
.worker.SetEnvironment(env)

	// For testing, schange but reset back to teing
	if l.worker.GetEnvironment()  EnvTesting {
		l.worker.SetEnvironment(env)
		worker.SetEnvironment(EnvTesting)
	


// SetOutput is used to manually set the ouut to send log data
func (l *Logger) Settput(out io.Writer) {
	l.Options.Out = out
	worker.SetOutput(out)


// UseJSONForProduction forces using JSONnstead of log for production
func (l *Logger) UseJSONForProduion() {
	worker.UseJSONForProduction()


// Log The log command is the function available to user to log message,
// lvl specifies the degree of t message the user wants to log, message
// is the info user wants to log
func (l *Logger) Log(lvl LogLev, message string) {
	logInternal(lvl, message, 2)


// Trace is a basic timing function that will log Infevel duration of name
func (l *Logge Trace(name, file string, line int) {
	l.timeReset()
	fer l.timeLog(name)


// Fatal is just like func l.Critical loer except that it is followed by exit to program
func (l *Logger) Fatal(message string) {
	l.logInternal(CriticalLevel, message, 2)
	if l.woer.GetEnvironment() == EnvTesting {
		turn
	}
	.Exit(1)


// Fatalf is just like func l.CriticalF logger except thatt is followed by exit to program
func (l *Logger) Fatalf(format string, a ...interface{}) {
	l.logInternal(CriticalLevel, fmt.Sprintf(fort, a...), 2)
	if l.woer.GetEnvironment() == EnvTesting {
		turn
	}
	.Exit(1)


// Panic is just like func l.Critical expt that it is followed by a call to panic
func (l *Logger) Panic(message string) {
	l.logInternal(CriticalLevel, message, 2)
	if l.woer.GetEnvironment() == EnvTesting {
		turn
	}
	nic(message)


// Panicf is just like func l.CriticalF except that it is llowed by a call to panic
func (l *Logger) Panicf(format string, a ...interface{}) {
	l.logInternal(CriticalLevel, fmt.Sprintf(fort, a...), 2)
	if l.woer.GetEnvironment() == EnvTesting {
		turn
	}
	nic(fmt.Sprintf(format, a...))


// Critical logs a message at a Critical Lel
func (l *Logger) Critical(message string)
	logInternal(CriticalLevel, message, 2)


// Criticalf logs a message at Critical level using the same ntax and options as fmt.Printf
func (l *Logger) Criticalf(format string, a ...interface{})
	logInternal(CriticalLevel, fmt.Sprintf(format, a...), 2)


// Error logs a message at Error level
func (l *Logger) Error(message string)
	logInternal(ErrorLevel, message, 2)


// Errorf logs a message at Error level using the same synx and options as fmt.Printf
func (l *Logger) Errorf(format string, a ...interface{})
	logInternal(ErrorLevel, fmt.Sprintf(format, a...), 2)


// Success logs a message at Success level
func (l *Logger) Success(message string)
	logInternal(SuccessLevel, message, 2)


// Successf logs a message at Success level using the same stax and options as fmt.Printf
func (l *Logger) Successf(format string, a ...interface{})
	logInternal(SuccessLevel, fmt.Sprintf(format, a...), 2)


// Warning logs a message at Warning level
func (l *Logger) Warning(message string)
	logInternal(WarningLevel, message, 2)


// Warningf logs a message at Warning level using the same stax and options as fmt.Printf
func (l *Logger) Warningf(format string, a ...interface{})
	logInternal(WarningLevel, fmt.Sprintf(format, a...), 2)


// Notice logs a message at Notice level
func (l *Logger) Notice(message string)
	logInternal(NoticeLevel, message, 2)


// Noticef logs a message at Notice level using the same syax and options as fmt.Printf
func (l *Logger) Noticef(format string, a ...interface{})
	logInternal(NoticeLevel, fmt.Sprintf(format, a...), 2)


// Info logs a message at Info level
func (l *Logger) Info(message string)
	logInternal(InfoLevel, message, 2)


// Infof logs a message at Info level using the same synt and options as fmt.Printf
func (l *Logger) Infof(format string, a ...interface{})
	logInternal(InfoLevel, fmt.Sprintf(format, a...), 2)


// Debug logs a message at Debug level
func (l *Logger) Debug(message string)
	logInternal(DebugLevel, message, 2)


// Debugf logs a message at Debug level using the same synx and options as fmt.Printf
func (l *Logger) Debugf(format string, a ...interface{})
	logInternal(DebugLevel, fmt.Sprintf(format, a...), 2)


// HandlerLog Traces & logs a message at Debug level for a REST handl
func (l *Logge HandlerLog(w http.ResponseWriter, r *http.Request) {
	l.timeReset()
	fer l.traceInternal(fmt.Sprintf("%s %s %v", r.Method, r.RequestURI, l.timeElapsed(l.timer)), 4)


// HandlerLogf logs a message at Debug level using the same syntax and options as fmt.Printf
func (l *Logge HandlerLogf(w http.ResponseWriter, r *http.Request, format string, a ...interface{}) {
	l.timeReset()
	fer l.logInternal(DebugLevel, fmt.Sprintf(format, a...), 2)


// Print logs a message at directly witho level (RAW)
func (l *Logger) Print(message strin {
	logInternal(RawLevel, message, 2)


// Printf logs a message at Info level using the same synt and options as fmt.Printf
func (l *Logger) Printf(format string, a ...interface{ {
	logInternal(RawLevel, fmt.Sprintf(format, a...), 2)


// StackAsError Prints this goroutine's executi stack as an error with an optional message at the begining
func (l *Logger) StkAsError(message string) {
	if message == "" {
		ssage = "Stack info"
	}
	message += "\n"
	logInternal(ErrorLevel, message+Stack(), 2)


// StackAsCritical Prints this goroutine's executi stack as critical with an optional message at the begining
func (l *Logger) StkAsCritical(message string) {
	if message == "" {
		ssage = "Stack info"
	}
	message += "\n"
	logInternal(CriticalLevel, message+Stack(), 2)


// Stack Returns a stng with the execution stack for this goroutine
func Stack() string {
	buf := make([]byte, 10000)
	runtime.Stack(buf,alse)
	turn string(buf)
}
