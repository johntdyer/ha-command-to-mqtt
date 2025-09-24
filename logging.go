package main

import (
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

// InitLogger initializes the global logger
func InitLogger() {
	logger = logrus.New()
}

// SetLogLevel sets the log level and format based on configuration
func SetLogLevel(logLevel, logFormat string) {
	// Set log format
	switch strings.ToLower(logFormat) {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	case "logfmt":
		// Use TextFormatter with no colors - this produces logfmt-compatible output
		logger.SetFormatter(&logrus.TextFormatter{
			DisableColors:   true,
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: time.RFC3339,
		})
		logger.Warnf("Unknown log format '%s', defaulting to text", logFormat)
	}

	// Set log level
	switch strings.ToLower(logLevel) {
	case "panic":
		logger.SetLevel(logrus.PanicLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "warn", "warning":
		logger.SetLevel(logrus.WarnLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "trace":
		logger.SetLevel(logrus.TraceLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
		logger.Warnf("Unknown log level '%s', defaulting to info", logLevel)
	}

	// Add caller info for debug and trace levels
	if logger.Level >= logrus.DebugLevel {
		logger.SetReportCaller(true)
	}

	logger.Debugf("Log level set to %s, format set to %s", strings.ToUpper(logLevel), strings.ToUpper(logFormat))
}