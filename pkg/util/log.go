package util

import (
	"fmt"
	"os"
	"strings"
)

type logLevel int

const (
	levelNone  logLevel = iota // suppress all
	levelWarn                  // LOG_LEVEL=warn
	levelInfo                  // LOG_LEVEL=info
	levelDebug                 // LOG_LEVEL=debug
	levelTrace                 // LOG_LEVEL=trace (most verbose)
)

var currentLevel = parseLogLevel(os.Getenv("LOG_LEVEL"))

func parseLogLevel(s string) logLevel {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "trace":
		return levelTrace
	case "debug":
		return levelDebug
	case "info":
		return levelInfo
	case "warn":
		return levelWarn
	default:
		return levelNone
	}
}

func logf(level logLevel, prefix, format string, args ...interface{}) {
	if currentLevel >= level {
		fmt.Fprintf(os.Stderr, prefix+format, args...)
	}
}

// Warnf prints to stderr when LOG_LEVEL=warn or higher.
func Warnf(format string, args ...interface{}) {
	logf(levelWarn, "[warn] ", format, args...)
}

// Infof prints to stderr when LOG_LEVEL=info or higher.
func Infof(format string, args ...interface{}) {
	logf(levelInfo, "[info] ", format, args...)
}

// Debugf prints to stderr when LOG_LEVEL=debug or higher.
func Debugf(format string, args ...interface{}) {
	logf(levelDebug, "[debug] ", format, args...)
}

// Tracef prints to stderr when LOG_LEVEL=trace.
func Tracef(format string, args ...interface{}) {
	logf(levelTrace, "[trace] ", format, args...)
}
