package model

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bearstech/restartd/protocol"
	"io"
)

type Handler interface {
	Handle(m Message) (resp Response)
}

type ProtocolHandler struct {
	handler Handler
}

func NewProtocolHandler(handler Handler) *ProtocolHandler {
	return &ProtocolHandler{
		handler: handler,
	}
}

func (h *ProtocolHandler) Handle(req io.Reader, resp io.Writer) {
	var msg Message
	err := protocol.Read(req, &msg)
	fmt.Println(msg)
	var r Response
	if err != nil {
		log.Error("Error while reading a command: ", err)
		oups := Response_err_reading
		msg := err.Error()
		r = Response{
			Code:    &oups,
			Message: &msg,
		}
	} else {
		r = h.handler.Handle(msg)
	}
	err = protocol.Write(resp, &r)
	if err != nil {
		panic(err)
	}
}
