package main

import (
	"encoding/gob"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bearstech/restartd/listen"
	"github.com/bearstech/restartd/protocol"
	"io"
	"os"
)

type Handler struct {
	services []string
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
	fldr := os.Getenv("RESTARTD_SOCKET_FOLDER")
	if fldr == "" {
		fldr = "/tmp/"
	}
	conf_folder := os.Getenv("RESTARTD_CONF")
	if conf_folder == "" {
		conf_folder = "/etc/restartd/conf.d"
	}
	confs, err := ReadConfFolder(conf_folder)
	if err != nil {
		panic(err)
	}
	r := listen.New(fldr)
	for _, conf := range confs {
		err = r.AddUser(conf.User, &Handler{conf.Services})
		if err != nil {
			panic(err)
		}
		log.Info("Add user ", conf.User)
	}
	log.Info("Number of users : ", len(confs))
	r.Listen()
}
