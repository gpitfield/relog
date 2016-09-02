/*
	Package relog enables logging relays and collectors, and log message priorities, such that a
	single call to relog can log to multiple locations conditional on that message's priority.
	It defines a Receiver interface which can handle the Print/Panic/Fatal calls from
	standard package log, as well as prioritized messages e.g. Alert, Debug.
	Receiver is implemented by 1) the Relay type that routes log messages to its registered
	Receivers based on the priority of the log message, and 2) the Collector type which logs
	the messages it receives via its embedded log.Logger.
	Any io.Writer can be registered as a Collector's Logger.

	The log levels, and much of the terminology, stem largely from https://tools.ietf.org/html/rfc3164.

	relog implements all the public methods of the standard log package, enabling drop-in replacement.
*/
package relog

import "io"

const (
	LEmerg = iota
	LAlert
	LCritical
	LError
	LWarn
	LNotice
	LInfo
	LDebug
)

// MaskFlags op constants
const (
	NONE = iota
	AND
	OR
	XOR
	ANDNOT
)

var severities = []string{"EMERGENCY", "ALERT", "CRITICAL", "ERROR", "WARNING", "NOTICE", "INFO", "DEBUG"}

// A Receiver, typically either a Relay or a Collector, relays or logs information based on its role.
type Receiver interface {
	Log(severity int, calldepth int, v ...interface{})                 // Log v at given severity level
	Logf(severity int, calldepth int, format string, v ...interface{}) // Logf v at given severity level
	Logln(severity int, calldepth int, v ...interface{})               // Logln v at given severity level
	Output(calldepth int, s string) error                              // Output s using given calldepth
	SetOutput(w io.Writer)                                             // Set the Receiver's Output location
	SetFlags(flag int, maskOp int)                                     // Set the Receiver's output flags using pkg log flags
	SetPrefix(prefix string)                                           // Set the value to prepend to Receiver's log statements
	SetVerbosity(verbosity int)                                        // Set the level at or above which receiver will generate a log statement
}
