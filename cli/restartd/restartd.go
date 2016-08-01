package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bearstech/restartd/listen"
	"github.com/bearstech/restartd/model"
	"github.com/bearstech/restartd/restartd"
	"github.com/urfave/cli"
	"os"
	"os/signal"
	"syscall"
)

var GITCOMMIT string

func main() {

	app := cli.NewApp()
	app.Version = "git:" + GITCOMMIT

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "Version, V",
			Usage: "Version",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.Bool("V") {
			fmt.Printf("Restartd daemon git:%s\n", GITCOMMIT)
			return nil
		}

		fldr := os.Getenv("RESTARTD_SOCKET_FOLDER")
		if fldr == "" {
			fldr = "/tmp/restartd"
		}

		_, err := os.Stat(fldr)
		if err != nil && os.IsExist(err) {
			panic(err)
		}

		if os.IsNotExist(err) {
			err = os.Mkdir(fldr, 0644)
			if err != nil {
				panic(err)
			}
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
			confs, err := restartd.ReadConfFolder(confFolder)
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
						&restartd.Handler{Services: conf.Services,
							User: conf.User,
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

		cc := make(chan os.Signal, 1)
		signal.Notify(cc, os.Interrupt, syscall.SIGHUP, syscall.SIGUSR1)
		go func() {
			for {
				s := <-cc
				log.Info("Signal : ", s)
				switch s {
				case os.Interrupt:
					r.Stop()
				case syscall.SIGTERM:
					r.Stop()
				case syscall.SIGHUP:
					configs()
				}
			}
		}()

		// listen and block
		r.Listen()
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		os.Exit(-1)
	}
}
