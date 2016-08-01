package model

import (
	"bytes"
	"github.com/bearstech/restartd/protocol"
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
	err := protocol.Write(network, &m)
	if err != nil {
		t.Error("Write trouble : ", err)
	}
	s := network.String()
	t.Log("Wire : ", len(s), s)
	var m2 Message
	err = protocol.Read(network, &m2)
	if err != nil {
		t.Error("Read trouble : ", err)
	}
}

func TestResponse(t *testing.T) {
	success := Response_success
	r := Response{
		Code: &success,
	}
	network := new(bytes.Buffer) // Stand-in for the network.
	err := protocol.Write(network, &r)
	if err != nil {
		t.Error("Write trouble : ", err)
	}
	s := network.String()
	t.Log("Wire : ", len(s), s)
	var r2 Response
	err = protocol.Read(network, &r2)
	if err != nil {
		t.Error("Read trouble : ", err)
	}
}

func TestFatResponse(t *testing.T) {
	started := Response_Statuz_started
	service := "web"
	success := Response_success
	r := Response{
		Code: &success,
		Status: []*Response_Statuz{
			&Response_Statuz{
				Service: &service,
				Code:    &started,
			},
		},
	}
	network := new(bytes.Buffer) // Stand-in for the network.
	err := protocol.Write(network, &r)
	if err != nil {
		t.Error("Write trouble : ", err)
	}
	s := network.String()
	t.Log("Wire : ", len(s), s)
	var r2 Response
	err = protocol.Read(network, &r2)
	if err != nil {
		t.Error("Read trouble : ", err)
	}
}
