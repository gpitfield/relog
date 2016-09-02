package relog

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

const collectorTestMessage = "collector\n"

var CollectorTests = []ReceiverTest{
	{LEmerg, LEmerg, 0, "", "[EMERGENCY] " + collectorTestMessage, 22},
	{LEmerg, LError, 0, "", "[EMERGENCY] " + collectorTestMessage, 22},
	{LDebug, LError, 0, "", "", 0},
	{LError, LError, 0, "", "[ERROR] " + collectorTestMessage, 18},
	{LError, LError, 0, "Prefix", "Prefix[ERROR] " + collectorTestMessage, 24},
	{LError, LError, log.Lshortfile, "", "collector_test.go", 40},
}

func TestCollector(t *testing.T) {
	var output bytes.Buffer
	collector := NewCollector(&output, LError, "", 0)
	for _, test := range CollectorTests {
		collector.SetVerbosity(test.verbosity)
		collector.SetPrefix(test.prefix)
		collector.SetFlags(test.flag, NONE)

		output.Reset()
		collector.Log(test.severity, 1, collectorTestMessage) // Test Log()
		result := output.String()
		if !strings.Contains(result, test.match) {
			t.Errorf("Collector Log messages didn't match\nEXP: %s^\nGOT: %s^", test.match, result)
		} else if len(result) != test.len {
			t.Errorf("Collector Log messages wrong length EXP: %d GOT: %d\n%s", test.len, len(result), result)
		}

		output.Reset()
		collector.Logf(test.severity, 1, "%s", collectorTestMessage) // Test Logf()
		result = output.String()
		if !strings.Contains(result, test.match) {
			t.Errorf("Collector Logf messages didn't match\nEXP: %s^\nGOT: %s^", test.match, result)
		} else if len(result) != test.len {
			t.Errorf("Collector Logf messages wrong length EXP: %d GOT: %d\n%s", test.len, len(result), result)
		}

		output.Reset()
		collector.Logln(test.severity, 1, collectorTestMessage) // Test Logln()
		result = output.String()
		if len(result) > 0 {
			result = result[0 : len(result)-1]
		}
		if !strings.Contains(result, test.match) {
			t.Errorf("Collector Logln messages didn't match\nEXP: %s^\nGOT: %s^", test.match, result)
		} else if len(result) != test.len {
			t.Errorf("Collector Logln messages wrong length EXP: %d GOT: %d\n%s", test.len, len(result), result)
		}
	}
}
