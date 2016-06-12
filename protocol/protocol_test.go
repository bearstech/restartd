package protocol

import "testing"

func TestParseCommand(t *testing.T) {
	cmd, err := ParseCommand("status")
	if err != nil {
		t.Fail()
	}
	if cmd != Status {
		t.Fail()
	}
}

func TestParseUnknownCommand(t *testing.T) {
	cmd, err := ParseCommand("shproutz")
	if err == nil {
		t.Fail()
	}
	if cmd != Unknown {
		t.Fail()
	}
}
