package model

import (
	"testing"
)

type HandlerTest struct {
}

func (h *HandlerTest) Handle(m Message) (r Response) {
	ok := Response_success
	r = Response{
		Code: &ok,
	}
	return r
}

func TestHandler(t *testing.T) {
	h := HandlerTest{}
	service := "toto"
	cmd := Message_status
	m := Message{
		Service: &service,
		Command: &cmd,
	}
	r := h.Handle(m)
	if *r.Code != Response_success {
		t.Fail()
	}
}
