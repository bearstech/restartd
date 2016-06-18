package protocol

import (
	"bytes"
	"testing"
)

func TestParseCommand(t *testing.T) {
	cmd, err := ParseCommand("status")
	if err != nil {
		t.Fail()
	}
	if cmd != Message_status {
		t.Fail()
	}
}

func TestParseUnknownCommand(t *testing.T) {
	_, err := ParseCommand("shproutz")
	if err == nil {
		t.Fail()
	}
}

func TestMessage(t *testing.T) {
	service := "toto"
	cmd := Message_status
	m := Message{
		Service: &service,
		Command: &cmd,
	}
	network := new(bytes.Buffer) // Stand-in for the network.
	err := Write(network, &m)
	if err != nil {
		t.Error("Write trouble : ", err)
	}
	s := network.String()
	t.Log("Wire : ", len(s), s)
	var m2 Message
	err = Read(network, &m2)
	if err != nil {
		t.Error("Read trouble : ", err)
	}
}
