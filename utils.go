package log

import (
	"strings"
	"time"
)

var StringToLevelMap = map[string]LogLevel{
	"fatal":   LOG_LEVEL_FATAL,
	"error":   LOG_LEVEL_ERROR,
	"warn":    LOG_LEVEL_WARN,
	"warning": LOG_LEVEL_WARN,
	"debug":   LOG_LEVEL_DEBUG,
	"info":    LOG_LEVEL_INFO,
	"all":     LOG_LEVEL_ALL,
	"":        LOG_LEVEL_ALL,
}

func StringToLogLevel(level string) LogLevel {
	level = strings.ToLower(level)
	return StringToLevelMap[level]
}

func LogTypeToString(t LogType) (string, string) {
	switch t {
	case LOG_FATAL:
		return "F", "[0;31"
	case LOG_ERROR:
		return "E", "[0;31"
	case LOG_WARNING:
		return "W", "[0;33"
	case LOG_DEBUG:
		return "D", "[0;36"
	case LOG_INFO:
		return "I", "[0;37"
	}
	// unknow
	return "?", "[0;37"
}

func genDayTime(t time.Time) string {
	return t.Format(FORMAT_TIME_DAY)
}

func genHourTime(t time.Time) string {
	return t.Format(FORMAT_TIME_HOUR)
}
