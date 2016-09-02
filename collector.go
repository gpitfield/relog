package relog

import (
	"fmt"
	"io"
	"log"
)

// Collector is a wrapper on the golang log package type Logger, with the addition
// of a verbosity parameter to control what level of log messages should be sent to its Logger.
// Collector implements the Receiver interface.
// The Collector's logger's flag are derived as a bitwise OR of the Relay and Collector's flag values.
// The Collector's logger's prefix is set when a Collector is registered with a Relay, and cannot
// be changed.
type Collector struct {
	logger    *log.Logger
	verbosity int
	flag      int // stored at the Collector level to allow masking modifications
}

// NewCollector creates a new Collector using the provided io.Writer and settings.
func NewCollector(w io.Writer, verbosity int, prefix string, flag int) *Collector {
	return &Collector{
		verbosity: verbosity,
		flag:      flag,
		logger:    log.New(w, prefix, flag),
	}
}

// SetFlags sets the Collector's flag via a masking operation, and sets the logger's flag to the resultant value.
func (c *Collector) SetFlags(flag int, maskOp int) {
	switch maskOp {
	case NONE:
		c.flag = flag
	case AND:
		c.flag = c.flag & flag
	case OR:
		c.flag = c.flag | flag
	case XOR:
		c.flag = c.flag ^ flag
	case ANDNOT:
		c.flag = c.flag &^ flag
	}
	c.logger.SetFlags(c.flag)
}

// SetPrefix sets the Collectors's logger's prefix.
func (c *Collector) SetPrefix(prefix string) {
	c.logger.SetPrefix(prefix)
}

// Prefix returns the Collectors's logger's prefix.
func (c *Collector) Prefix() string {
	return c.logger.Prefix()
}

func (c *Collector) SetOutput(w io.Writer) {
	c.logger.SetOutput(w)
}

// Output prepends the severity label and calls the Collector's logger.
func (c *Collector) Output(calldepth int, s string) error {
	return c.logger.Output(calldepth+1, s)
}

// SetVerbosity sets the Collector's verbosity. Messages of lower priority than the verbosity are not logged.
func (c *Collector) SetVerbosity(verbosity int) {
	c.verbosity = verbosity
}

// Verbosity returns the Collector's verbosity.
func (c *Collector) Verbosity() int {
	return c.verbosity
}

// Log generates the log string and calls Output
func (c *Collector) Log(severity int, calldepth int, v ...interface{}) {
	if c.verbosity >= severity {
		c.Output(calldepth+1, "["+severities[severity]+"] "+fmt.Sprint(v...))
	}
}

// Logf generates the log string and calls Output.
func (c *Collector) Logf(severity int, calldepth int, format string, v ...interface{}) {
	if c.verbosity >= severity {
		c.Output(calldepth+1, "["+severities[severity]+"] "+fmt.Sprintf(format, v...))
	}
}

// Logln generates the log string and calls Output.
func (c *Collector) Logln(severity int, calldepth int, v ...interface{}) {
	if c.verbosity >= severity {
		c.Output(calldepth+1, "["+severities[severity]+"] "+fmt.Sprintln(v...))
	}
}
