// high level log wrapper, so it can output different log based on level
package log

import (
	"io"
	"log"
	"os"
)

const (
	Ldate         = log.Ldate
	Llongfile     = log.Llongfile
	Lmicroseconds = log.Lmicroseconds
	Lshortfile    = log.Lshortfile
	LstdFlags     = log.LstdFlags
	Ltime         = log.Ltime
)

type (
	LogLevel int
	LogType int
)

const (
	LOG_FATAL   = LogType(0x1)
	LOG_ERROR   = LogType(0x2)
	LOG_WARNING = LogType(0x4)
	LOG_INFO    = LogType(0x8)
	LOG_DEBUG   = LogType(0x10)
)

const (
	LOG_LEVEL_NONE  = LogLevel(0x0)
	LOG_LEVEL_FATAL = LOG_LEVEL_NONE | LogLevel(LOG_FATAL)
	LOG_LEVEL_ERROR = LOG_LEVEL_FATAL | LogLevel(LOG_ERROR)
	LOG_LEVEL_WARN  = LOG_LEVEL_ERROR | LogLevel(LOG_WARNING)
	LOG_LEVEL_INFO  = LOG_LEVEL_WARN | LogLevel(LOG_INFO)
	LOG_LEVEL_DEBUG = LOG_LEVEL_INFO | LogLevel(LOG_DEBUG)
	LOG_LEVEL_ALL   = LOG_LEVEL_DEBUG
)

const FORMAT_TIME_DAY string = "2006-01-02"
const FORMAT_TIME_HOUR string = "2006-01-02-15"

var _log *Logger = New()

func init() {
	_log.SetDepth(5)
}

func SetLevel(level LogLevel) {
	_log.SetLevel(level)
}
func GetLogLevel() LogLevel {
	return _log.Level
}

func SetOutput(out io.Writer) {
	_log.SetOutput(out)
}

func SetOutputByName(path string) error {
	return _log.SetOutputByName(path)
}

func SetPrefix(prefix string) {
	_log.SetPrefix(prefix)
}

func SetFlags(flags int) {
	_log._log.SetFlags(flags)
}

func Info(v ...interface{}) {
	_log.Info(v...)
}

func Infof(format string, v ...interface{}) {
	_log.Infof(format, v...)
}

func Debug(v ...interface{}) {
	_log.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	_log.Debugf(format, v...)
}

func Warn(v ...interface{}) {
	_log.Warning(v...)
}

func Warnf(format string, v ...interface{}) {
	_log.Warningf(format, v...)
}

func Warning(v ...interface{}) {
	_log.Warning(v...)
}

func Warningf(format string, v ...interface{}) {
	_log.Warningf(format, v...)
}

func Error(v ...interface{}) {
	_log.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	_log.Errorf(format, v...)
}

func Fatal(v ...interface{}) {
	_log.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	_log.Fatalf(format, v...)
}

func Panic(v ...interface{}) {
	_log.Panic(v...)
}

func Panicf(format string, v ...interface{}) {
	_log.Panicf(format, v...)
}

func Panicln(v ...interface{}) {
	_log.Panicln(v...)
}

func Print(v ...interface{}) {
	_log.Info(v...)
}

func Printf(format string, v ...interface{}) {
	_log.Printf(format, v...)
}

func Println(v ...interface{}) {
	_log.Println(v...)
}

func SetLevelByString(level string) {
	_log.SetLevelByString(level)
}

func EnableHighlighting() {
	_log.SetHighlighting(true)
}

func DisableHighlighting() {
	_log.SetHighlighting(false)
}

func SetRotateByDay() {
	_log.SetRotateByDay()
}

func SetRotateByHour() {
	_log.SetRotateByHour()
}

func SetRotateByMaxSize(size, backup int8) {
	_log.SetRotateByMaxSize(size, backup)
}

func New() *Logger {
	l := Newlogger(os.Stderr, "")
	return l
}

func Newlogger(w io.Writer, prefix string) *Logger {
	return &Logger{
		_log:   log.New(w, prefix, Ldate|Ltime|Lshortfile),
		Level:  LOG_LEVEL_ALL,
		depth:  4,
		Prefix: prefix,
	}
}

// Copy A logger
func Copy(l *Logger) *Logger {
	return _log.Copy()
}
