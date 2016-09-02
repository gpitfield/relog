# relog [![GoDoc](https://godoc.org/github.com/gpitfield/relog?status.svg)](https://godoc.org/github.com/gpitfield/relog) [![GoCover](http://gocover.io/_badge/github.com/gpitfield/relog)](http://gocover.io/github.com/gpitfield/relog)

Relog provides prioritized log messages, as well as logging relays and collectors such that a
single call to relog can log to multiple locations conditional on that message's priority.
It defines a Receiver interface which can handle the Print/Panic/Fatal calls from
standard package log, as well as prioritized messages e.g. Alert, Debug.
Receiver is implemented by 1) the Relay type that routes log messages to its registered
Receivers based on the priority of the log message, and 2) the Collector type which logs
the messages it receives via its embedded log.Logger.

Any io.Writer can be registered as a Collector's Logger.

The log levels, and much of the terminology, stem largely from [rfc3164](https://tools.ietf.org/html/rfc3164).

relog implements all the public methods of the standard log package, enabling drop-in replacement.
