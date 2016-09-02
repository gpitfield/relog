package relog

import (
	"bytes"
	"strings"
	"testing"
)

const relayTestMessage = "relay\n"

type ReceiverTest struct {
	severity  int
	verbosity int
	flag      int
	prefix    string
	match     string
	len       int
}

var PrintTests = []ReceiverTest{
	{LEmerg, LNotice, 0, "", "[NOTICE] " + relayTestMessage, 15},
	{LEmerg, LAlert, 0, "", "", 0},
}

func TestRelayPrint(t *testing.T) {
	var output bytes.Buffer
	SetOutput(&output)
	for _, test := range PrintTests {
		SetVerbosity(test.verbosity)
		SetPrefix(test.prefix)
		SetFlags(test.flag)

		output.Reset()
		Print(relayTestMessage) // Test Print()
		result := output.String()
		if !strings.Contains(result, test.match) {
			t.Errorf("Relay Print messages didn't match\nEXP: %s^\nGOT: %s^", test.match, result)
		} else if len(result) != test.len {
			t.Errorf("Relay Print messages wrong length EXP: %d GOT: %d\n%s", test.len, len(result), result)
		}

		output.Reset()
		Printf("%s", relayTestMessage) // Test Printf()
		result = output.String()
		if !strings.Contains(result, test.match) {
			t.Errorf("Relay Println messages didn't match\nEXP: %s^\nGOT: %s^", test.match, result)
		} else if len(result) != test.len {
			t.Errorf("Relay Println messages wrong length EXP: %d GOT: %d\n%s", test.len, len(result), result)
		}

		output.Reset()
		Println(relayTestMessage) // Test Println()
		result = output.String()
		if len(result) > 0 {
			result = result[0 : len(result)-1]
		}
		if !strings.Contains(result, test.match) {
			t.Errorf("Relay Printf messages didn't match\nEXP: %s^\nGOT: %s^", test.match, result)
		} else if len(result) != test.len {
			t.Errorf("Relay Printf messages wrong length EXP: %d GOT: %d\n%s", test.len, len(result), result)
		}
	}
}
