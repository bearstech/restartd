package main

import (
	"../listen"
	"../protocol"
	"encoding/gob"
	"fmt"
	"io"
)

type Handler struct {
}

func (h *Handler) Handle(req io.Reader, resp io.Writer) {
	var msg protocol.Message
	dec := gob.NewDecoder(req)
	err := dec.Decode(&msg)
	fmt.Println(msg)
	var r protocol.Response
	if err != nil {
		r = protocol.Response{1, err.Error()}
	} else {
		r = h.HandleMessage(msg)
	}
	enc := gob.NewEncoder(resp)
	err = enc.Encode(&r)
	if err != nil {
		panic(err)
	}
}

func (h *Handler) HandleMessage(msg protocol.Message) protocol.Response {
	return protocol.Response{0, fmt.Sprintf("%s was sent to %s", msg.Command.Command(), msg.Service)}
}

func main() {
	r := listen.New("/tmp")
	err := r.AddUser("alice", &Handler{})
	if err != nil {
		panic(err)
	}
	r.Listen()

}
