package relog

import (
	"fmt"
	"io"
	"log"
	"os"
)

// Relay forwards log messages to its receivers based on its verbosity value.
// Relay implements the Receiver interface.
// To prevent a single call to Fatal[f|ln] or Panic[f|ln] from spawning multiple calls to os.Exit()
// or panic() (in the case of the mirrored calls from standard package log), the Relay
// forwards those messages to its Receivers as prioritized messages, handling os.Exit() and panic() itself, with
// Panic[f|ln] and Fatal[f|ln] forwarded to the Receivers' Emerg[f|ln] function, and Print[f|ln]
// forwarded to the Receivers' Notice[f|ln].
type Relay struct {
	receivers []Receiver
	prefix    string
	flag      int
	verbosity int
	calldepth int
}

// TODO: initialize this to point to sys.log
var std Relay = Relay{
	verbosity: LDebug,
	calldepth: 3,
	receivers: []Receiver{NewCollector(os.Stderr, LDebug, "", log.Lshortfile|log.LstdFlags)},
}

// New creates a new Relay with no receivers.
func New(verbosity int, prefix string, flag int) *Relay {
	return &Relay{
		verbosity: verbosity,
		prefix:    prefix,
		flag:      flag,
		calldepth: 2,
	}
}

// New creates a new Relay with a collector to StdErr
func NewStdLog(verbosity int, prefix string, flag int) *Relay {
	return &Relay{
		verbosity: verbosity,
		prefix:    prefix,
		flag:      flag,
		calldepth: 2,
		receivers: []Receiver{NewCollector(os.Stderr, verbosity, "", flag)},
	}
}

// AddWriter creates a Collector and adds it to the Relay's receivers
func (r *Relay) AddWriter(w io.Writer, verbosity int, prefix string, flag int) {
	r.receivers = append(r.receivers, NewCollector(w, verbosity, prefix, flag))
}

// AddReceiver adds a Receiver to the Relay.
func (r *Relay) AddReceiver(rcvr Receiver) {
	r.receivers = append(r.receivers, rcvr)
}

// SetFlags sets the Relay's flag via a masking operation, and calls SetFlags for its Receivers with its own flags as the mask.
func SetFlags(flag int) { std.SetFlags(flag, NONE) }
func (r *Relay) SetFlags(flag int, maskOp int) {
	switch maskOp {
	case NONE:
		r.flag = flag
	case AND:
		r.flag = r.flag & flag
	case OR:
		r.flag = r.flag | flag
	case XOR:
		r.flag = r.flag ^ flag
	case ANDNOT:
		r.flag = r.flag &^ flag
	}
	for i, _ := range r.receivers {
		r.receivers[i].SetFlags(r.flag, maskOp)
	}
}

// Flags returns the output flags for the Relay
func Flags() int            { return std.Flags() }
func (r *Relay) Flags() int { return r.flag }

// SetPrefix sets the Relay's prefix which is prepended to log statements.
func SetPrefix(prefix string) { std.SetPrefix(prefix) }
func (r *Relay) SetPrefix(prefix string) {
	r.prefix = prefix
}

// Prefix returns the log prefix for the Relay
func Prefix() string            { return std.Prefix() }
func (r *Relay) Prefix() string { return r.prefix }

// SetOutput sets the standard Relay's receiver output.
func SetOutput(w io.Writer) {
	std.receivers[0].SetOutput(w)
}

// SetOutput is a null function for interface compatibility; set Collector outputs directly instead.
func (r *Relay) SetOutput(w io.Writer) {}

// Output writes the output for a logging event. Only provided for compatibility with standard log package.
func Output(calldepth int, s string) { std.Output(calldepth, s) }
func (r *Relay) Output(calldepth int, s string) {
	for i, _ := range r.receivers {
		r.receivers[i].Output(calldepth, s)
	}
}

// SetVerbosity sets the Relay's verbosity.
func SetVerbosity(verbosity int) { std.SetVerbosity(verbosity) }
func (r *Relay) SetVerbosity(verbosity int) {
	r.verbosity = verbosity
}

// Log forwards messages to the each receiver's Log function.
func (r *Relay) Log(severity int, calldepth int, v ...interface{}) {
	if r.verbosity < severity {
		return
	}
	v = append([]interface{}{r.prefix}, v...)
	calldepth++ // increment for this frame
	for i, _ := range r.receivers {
		r.receivers[i].Log(severity, calldepth, v...)
	}
}

// Logf forwards messages to the each receiver's Logf function.
func (r *Relay) Logf(severity int, calldepth int, format string, v ...interface{}) {
	if r.verbosity < severity {
		return
	}
	if r.prefix != "" {
		format = "%s " + format
		v = append([]interface{}{r.prefix}, v...)
	}
	calldepth++ // increment for this frame
	for i, _ := range r.receivers {
		r.receivers[i].Logf(severity, calldepth, format, v...)
	}
}

// Logln forwards messages to the each receiver's Logln function.
func (r *Relay) Logln(severity int, calldepth int, v ...interface{}) {
	if r.verbosity < severity {
		return
	}
	if r.prefix != "" {
		v = append([]interface{}{r.prefix}, v...)
	}
	calldepth++ // increment for this frame
	for i, _ := range r.receivers {
		r.receivers[i].Logln(severity, calldepth, v...)
	}
}

// Fatal is equivalent to a call to r.Emerg followed by a call to os.Exit(1).
func Fatal(v ...interface{}) { std.Fatal(v...) }
func (r *Relay) Fatal(v ...interface{}) {
	r.Log(LEmerg, r.calldepth, v...)
	os.Exit(1)
}

// Fatalf is equivalent to a call to r.Emergf followed by a call to os.Exit(1).
func Fatalf(format string, v ...interface{}) { std.Fatalf(format, v...) }
func (r *Relay) Fatalf(format string, v ...interface{}) {
	r.Logf(LEmerg, r.calldepth, format, v...)
	os.Exit(1)
}

// Fatalln is equivalent to a call to r.Emergln followed by a call to os.Exit(1).
func Fatalln(v ...interface{}) { std.Fatalln(v...) }
func (r *Relay) Fatalln(v ...interface{}) {
	r.Logln(LEmerg, r.calldepth, v...)
	os.Exit(1)
}

// Panic is equivalent to a call to r.Emerg followed by a call to panic().
func Panic(v ...interface{}) { std.Panic(v...) }
func (r *Relay) Panic(v ...interface{}) {
	r.Log(LEmerg, r.calldepth, v...)
	v = append([]interface{}{r.prefix}, v...)
	panic(fmt.Sprint(v...))
}

// Panicf is equivalent to a call to r.Logf at severity Emerg followed by a call to panic().
func Panicf(format string, v ...interface{}) { std.Panicf(format, v...) }
func (r *Relay) Panicf(format string, v ...interface{}) {
	r.Logf(LEmerg, r.calldepth, format, v...)
	msg := fmt.Sprintf(format, v...)
	if r.prefix != "" {
		panic(fmt.Sprintf("%s %s", r.prefix, msg))
	} else {
		panic(msg)
	}
}

// Panicln is equivalent to a call to r.Emergln followed by a call to panic().
func Panicln(v ...interface{}) { std.Panicln(v...) }
func (r *Relay) Panicln(v ...interface{}) {
	r.Logln(LEmerg, r.calldepth, v...)
	v = append([]interface{}{r.prefix}, v...)
	panic(fmt.Sprintln(v...))
}

// Print is equivalent to a call to r.Log at severity Notice.
func Print(v ...interface{}) { std.Print(v...) }
func (r *Relay) Print(v ...interface{}) {
	r.Log(LNotice, r.calldepth, v...)
}

// Printf is equivalent to a call to r.Logf at severity Notice.
func Printf(format string, v ...interface{})            { std.Printf(format, v...) }
func (r *Relay) Printf(format string, v ...interface{}) { r.Logf(LNotice, r.calldepth, format, v...) }

// Println is equivalent to a call to r.Logln at severity Notice.
func Println(v ...interface{})            { std.Println(v...) }
func (r *Relay) Println(v ...interface{}) { r.Logln(LNotice, r.calldepth, v...) }

// Emerg calls Log with severity Emerg.
func Emerg(v ...interface{})            { std.Emerg(v...) }
func (r *Relay) Emerg(v ...interface{}) { r.Log(LEmerg, r.calldepth, v...) }

// Emergf calls Logf with severity Emerg.
func Emergf(format string, v ...interface{})            { std.Emergf(format, v...) }
func (r *Relay) Emergf(format string, v ...interface{}) { r.Logf(LEmerg, r.calldepth, format, v...) }

// Emergln calls Logln with severity Emerg.
func Emergln(v ...interface{})            { std.Emergln(v...) }
func (r *Relay) Emergln(v ...interface{}) { r.Logln(LEmerg, r.calldepth, v...) }

// Alert calls Log with severity Alert.
func Alert(v ...interface{})            { std.Alert(v...) }
func (r *Relay) Alert(v ...interface{}) { r.Log(LAlert, r.calldepth, v...) }

// Alertf calls Logf with severity Alert.
func Alertf(format string, v ...interface{})            { std.Alertf(format, v...) }
func (r *Relay) Alertf(format string, v ...interface{}) { r.Logf(LAlert, r.calldepth, format, v...) }

// Alertln calls Logln with severity Alert.
func Alertln(v ...interface{})            { std.Alertln(v...) }
func (r *Relay) Alertln(v ...interface{}) { r.Logln(LAlert, r.calldepth, v...) }

// Critical calls Log with severity Critical.
func Critical(v ...interface{})            { std.Critical(v...) }
func (r *Relay) Critical(v ...interface{}) { r.Log(LCritical, r.calldepth, v...) }

// Criticalf calls Logf with severity Critical.
func Criticalf(format string, v ...interface{}) { std.Criticalf(format, v...) }
func (r *Relay) Criticalf(format string, v ...interface{}) {
	r.Logf(LCritical, r.calldepth, format, v...)
}

// Criticalln calls Logln with severity Critical.
func Criticalln(v ...interface{})            { std.Criticalln(v...) }
func (r *Relay) Criticalln(v ...interface{}) { r.Logln(LCritical, r.calldepth, v...) }

// Error calls Log with severity Error.
func Error(v ...interface{})            { std.Error(v...) }
func (r *Relay) Error(v ...interface{}) { r.Log(LError, r.calldepth, v...) }

// Errorf calls Logf with severity Error.
func Errorf(format string, v ...interface{})            { std.Errorf(format, v...) }
func (r *Relay) Errorf(format string, v ...interface{}) { r.Logf(LError, r.calldepth, format, v...) }

// Errorln calls Logln with severity Error.
func Errorln(v ...interface{})            { std.Errorln(v...) }
func (r *Relay) Errorln(v ...interface{}) { r.Logln(LError, r.calldepth, v...) }

// Warn calls Log with severity Warn.
func Warn(v ...interface{})            { std.Warn(v...) }
func (r *Relay) Warn(v ...interface{}) { r.Log(LWarn, r.calldepth, v...) }

// Warnf calls Logf with severity Warn.
func Warnf(format string, v ...interface{})            { std.Warnf(format, v...) }
func (r *Relay) Warnf(format string, v ...interface{}) { r.Logf(LWarn, r.calldepth, format, v...) }

// Warnln calls Logln with severity Warn.
func Warnln(v ...interface{})            { std.Warnln(v...) }
func (r *Relay) Warnln(v ...interface{}) { r.Logln(LWarn, r.calldepth, v...) }

// Notice calls Log with severity Notice.
func Notice(v ...interface{})            { std.Notice(v...) }
func (r *Relay) Notice(v ...interface{}) { r.Log(LNotice, r.calldepth, v...) }

// Noticef calls Logf with severity Notice.
func Noticef(format string, v ...interface{})            { std.Noticef(format, v...) }
func (r *Relay) Noticef(format string, v ...interface{}) { r.Logf(LNotice, r.calldepth, format, v...) }

// Noticeln calls Logln with severity Notice.
func Noticeln(v ...interface{})            { std.Noticeln(v...) }
func (r *Relay) Noticeln(v ...interface{}) { r.Logln(LNotice, r.calldepth, v...) }

// Info calls Log with severity Info.
func Info(v ...interface{})            { std.Info(v...) }
func (r *Relay) Info(v ...interface{}) { r.Log(LInfo, r.calldepth, v...) }

// Infof calls Logf with severity Info.
func Infof(format string, v ...interface{})            { std.Infof(format, v...) }
func (r *Relay) Infof(format string, v ...interface{}) { r.Logf(LInfo, r.calldepth, format, v...) }

// Infoln calls Logln with severity Info.
func Infoln(v ...interface{})            { std.Infoln(v...) }
func (r *Relay) Infoln(v ...interface{}) { r.Logln(LInfo, r.calldepth, v...) }

// Debug calls Log with severity Debug.
func Debug(v ...interface{})            { std.Debug(v...) }
func (r *Relay) Debug(v ...interface{}) { r.Log(LDebug, r.calldepth, v...) }

// Debugf calls Logf with severity Debug.
func Debugf(format string, v ...interface{})            { std.Debugf(format, v...) }
func (r *Relay) Debugf(format string, v ...interface{}) { r.Logf(LDebug, r.calldepth, format, v...) }

// Debugln calls Logln with severity Debug.
func Debugln(v ...interface{})            { std.Debugln(v...) }
func (r *Relay) Debugln(v ...interface{}) { r.Logln(LDebug, r.calldepth, v...) }
