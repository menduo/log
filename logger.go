package log

import (
	"fmt"
	"io"
	olog "log"
	"os"
	"sync"
	"time"
)

type Logger struct {
	_log         *olog.Logger
	Level        LogLevel
	highlighting bool

	dailyRolling bool
	hourRolling  bool
	sizeRolling  int8
	maxBackup    int8

	Prefix    string
	FileName  string
	logSuffix string
	fd        *os.File

	depth int
	lock  sync.Mutex
}

// Copy A logger, use its level, filename, prefix, depth, etc
func (l *Logger) Copy(resetDepth bool) *Logger {
	nlg := Newlogger(os.Stderr, "")
	nlg.SetLevel(l.Level)
	nlg.SetPrefix(l.Prefix)
	nlg.SetHighlighting(l.highlighting)
	if resetDepth {
		nlg.SetDepth(DeautltDepth)
	}

	if l.FileName != "" {
		err := nlg.SetOutputByName(l.FileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error on SetOutputByName: %s", err.Error())
		}
	}

	if l.dailyRolling {
		nlg.SetRotateByDay()
	}
	if l.hourRolling {
		nlg.SetRotateByHour()
	}

	nlg.SetDepth(l.depth)
	return nlg
}

func (l *Logger) SetHighlighting(highlighting bool) {
	l.highlighting = highlighting
}

func (l *Logger) SetFlags(flags int) {
	l._log.SetFlags(flags)
}

func (l *Logger) SetLevel(level LogLevel) {
	l.Level = level
}

func (l *Logger) SetLevelByString(level string) {
	l.Level = StringToLogLevel(level)
}

func (l *Logger) SetRotateByDay() {
	l.dailyRolling = true
	l.logSuffix = genDayTime(time.Now())
}

func (l *Logger) SetRotateByHour() {
	l.hourRolling = true
	l.logSuffix = genHourTime(time.Now())
}

func (l *Logger) SetDepth(depth int) {
	l.depth = depth
}

func (l *Logger) SetRotateByMaxSize(size, backup int8) {
	l.sizeRolling = size
	l.maxBackup = backup
	l.logSuffix = genHourTime(time.Now())
}

func (l *Logger) rotate() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	var suffix string
	if l.dailyRolling {
		suffix = genDayTime(time.Now())
	} else if l.hourRolling {
		suffix = genHourTime(time.Now())
	} else {
		return nil
	}

	// Notice: if suffix is not equal to l.LogSuffix, then rotate
	if suffix != l.logSuffix {
		err := l.doRotate(suffix)
		if err != nil {
			return err
		}
	}

	return nil
}

// SetPrefix sets the output prefix for the logger.
func (l *Logger) SetPrefix(prefix string) {
	l.Prefix = prefix
	l._log.SetPrefix(prefix)
}

func (l *Logger) doRotate(suffix string) error {
	// Notice: Not check error, is this ok?
	l.fd.Close()

	lastFileName := l.FileName + "." + l.logSuffix
	err := os.Rename(l.FileName, lastFileName)
	if err != nil {
		return err
	}

	err = l.SetOutputByName(l.FileName)
	if err != nil {
		return err
	}

	l.logSuffix = suffix

	return nil
}

func (l *Logger) SetOutput(out io.Writer) {
	l._log = olog.New(out, l._log.Prefix(), l._log.Flags())
}

func (l *Logger) SetOutputByName(path string) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		olog.Fatal(err)
	}

	l.SetOutput(f)

	l.FileName = path
	l.fd = f

	return err
}

func (l *Logger) log(t LogType, v ...interface{}) {
	if l.Level|LogLevel(t) != l.Level {
		return
	}

	err := l.rotate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	v1 := make([]interface{}, len(v)+2)
	logStr, logColor := LogTypeToString(t)
	if l.highlighting {
		v1[0] = "\033" + logColor + "m[" + logStr + "]"
		copy(v1[1:], v)
		v1[len(v)+1] = "\033[0m"
	} else {
		v1[0] = "[" + logStr + "]"
		copy(v1[1:], v)
		v1[len(v)+1] = ""
	}

	s := fmt.Sprintln(v1...)
	l.Output(s)
}

func (l *Logger) logf(t LogType, format string, v ...interface{}) {
	if l.Level|LogLevel(t) != l.Level {
		return
	}

	err := l.rotate()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}

	logStr, logColor := LogTypeToString(t)
	var s string
	if l.highlighting {
		s = "\033" + logColor + "m[" + logStr + "] " + fmt.Sprintf(format, v...) + "\033[0m"
	} else {
		s = "[" + logStr + "] " + fmt.Sprintf(format, v...)
	}
	l.Output(s)
}

func (l *Logger) Output(s string) {
	l._log.Output(l.depth, s)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.log(LOG_FATAL, v...)
	os.Exit(-1)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logf(LOG_FATAL, format, v...)
	os.Exit(-1)
}

func (l *Logger) Error(v ...interface{}) {
	l.log(LOG_ERROR, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.logf(LOG_ERROR, format, v...)
}

func (l *Logger) Warning(v ...interface{}) {
	l.log(LOG_WARNING, v...)
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	l.logf(LOG_WARNING, format, v...)
}

func (l *Logger) Debug(v ...interface{}) {
	l.log(LOG_DEBUG, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logf(LOG_DEBUG, format, v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.log(LOG_INFO, v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.logf(LOG_INFO, format, v...)
}

func (l *Logger) Print(v ...interface{}) {
	l.log(LOG_FATAL, v...)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.logf(LOG_FATAL, format, v...)
}

func (l *Logger) Println(v ...interface{}) {
	l.log(LOG_FATAL, v...)
}

func (l *Logger) Panic(v ...interface{}) {
	l.log(LOG_FATAL, v...)
	s := fmt.Sprint(v...)
	panic(s)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	l.logf(LOG_FATAL, format, v...)
	s := fmt.Sprintf(format, v...)
	panic(s)
}

func (l *Logger) Panicln(v ...interface{}) {
	l.log(LOG_FATAL, v...)
	s := fmt.Sprintln(v...)
	panic(s)
}
