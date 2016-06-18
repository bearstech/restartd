package main

import (
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
	err := protocol.Read(req, &msg)
	fmt.Println(msg)
	var r protocol.Response
	if err != nil {
		log.Error("Error while reading a command: ", err)
		oups := int32(1)
		msg := err.Error()
		r = protocol.Response{
			Code:    &oups,
			Message: &msg,
		}
	} else {
		r = h.HandleMessage(msg)
	}
	err = protocol.Write(resp, &r)
	if err != nil {
		panic(err)
	}
}

func (h *Handler) HandleMessage(msg protocol.Message) protocol.Response {
	ok := int32(0)
	message := fmt.Sprintf("%s was sent to %s", msg.Command.String(), msg.Service)
	return protocol.Response{
		Code:    &ok,
		Message: &message,
	}
}

func main() {
	fldr := os.Getenv("RESTARTD_SOCKET_FOLDER")
	if fldr == "" {
		fldr = "/tmp/"
	}
	log.Info("Socket folder is ", fldr)

	r := listen.New(fldr)
	defer r.Cleanup()

	conf_folder := os.Getenv("RESTARTD_CONF")

	if conf_folder == "" {
		conf_folder = "/etc/restartd/conf.d"
	}
	log.Info("Conf folder is ", conf_folder)

	configs := func() {
		confs, err := ReadConfFolder(conf_folder)
		if err != nil {
			panic(err)
		}
		if len(confs) == 0 {
			log.Error("No conf found. Add some yml file in " + conf_folder)
			//os.Exit(-1)
		}
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
