package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bearstech/restartd/listen"
	"github.com/bearstech/restartd/model"
	"os"
	"os/signal"
	"syscall"
)

var GITCOMMIT string

func main() {
	if len(os.Args) > 0 {
		if os.Args[1] == "-V" {
			fmt.Printf("Restartd daemon git:%s\n", GITCOMMIT)
			return
		}
	}
	fldr := os.Getenv("RESTARTD_SOCKET_FOLDER")
	if fldr == "" {
		fldr = "/tmp/"
	}
	log.Info("Socket folder is ", fldr)

	r := listen.New(fldr)
	defer r.Cleanup()

	confFolder := os.Getenv("RESTARTD_CONF")

	if confFolder == "" {
		confFolder = "/etc/restartd/conf.d"
	}
	log.Info("Conf folder is ", confFolder)

	configs := func() {
		confs, err := ReadConfFolder(confFolder)
		if err != nil {
			panic(err)
		}
		if len(confs) == 0 {
			log.Error("No conf found. Add some yml file in " + confFolder)
			//os.Exit(-1)
		}
		for _, conf := range confs {
			err = r.AddUser(conf.User,
				model.NewProtocolHandler(
					&Handler{Services: conf.Services,
						user: conf.User,
					}))
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
