package main

import (
	"encoding/gob"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bearstech/restartd/listen"
	"github.com/bearstech/restartd/protocol"
	"io"
	"os"
	"os/signal"
	"syscall"
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
	log.Info("Socket folder is ", fldr)
	r := listen.New(fldr)
	defer r.Cleanup()
	configs := func() {
		confs, err := ReadConfFolder(conf_folder)
		if err != nil {
			panic(err)
		}
		if len(confs) == 0 {
			log.Error("No conf found. Add some yml file in " + conf_folder)
			//os.Exit(-1)
		}
		log.Info("Conf folder is ", conf_folder)
		for _, conf := range confs {
			err = r.AddUser(conf.User, &Handler{conf.Services})
			if err != nil {
				panic(err)
			}
			log.Info("Add user ", conf.User)
		}
		log.Info("Number of users : ", len(confs))
	}
	// initial config
	configs()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGUSR1)
	go func() {
		for {
			s := <-c
			log.Info("Signal : ", s)
			switch s {
			case os.Interrupt:
				r.Stop()
			case syscall.SIGHUP:
				configs()
			}
		}
	}()

	// listen and block
	r.Listen()
}
