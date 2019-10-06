package io

import "time"

// Logger writes log messages to a buffer
type Logger struct {
	Writter              Writter
	Formatter            func(Log) []byte
	ShouldPool           bool
	spool                []Log
	SpoolTriggerSeverity Severity
}

// Severity level for log messages
type Severity int

const severityLow Severity = 0
const severityMedium Severity = 1
const severityHigh Severity = 2
const severityUrgent Severity = 3
const severityCritical Severity = 4

// Log entry
type Log struct {
	Severity  Severity
	Message   string
	Timestamp time.Time
	Context   interface{}
}

func (logger Logger) log(severity Severity, message string, context interface{}) {
	logger.spool = append(logger.spool, Log{
		Severity:  severity,
		Message:   message,
		Context:   context,
		Timestamp: time.Now(),
	})

	if severity >= logger.SpoolTriggerSeverity || logger.ShouldPool == false {
		logger.emptySpool()
	}
}

func (logger Logger) emptySpool() {
	for _, log := range logger.spool {
		logger.Writter.Write(logger.Formatter(log))
	}

	logger.spool = []Log{}
}
