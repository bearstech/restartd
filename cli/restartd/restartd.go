package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/bearstech/ascetic-rpc/server"
	"github.com/bearstech/restartd/restartd"
	"github.com/urfave/cli"
	"os"
	"os/signal"
	"syscall"
)

var GITCOMMIT string

func main() {

	var prefix bool = true

	app := cli.NewApp()
	app.Version = "git:" + GITCOMMIT

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "Version, V",
			Usage: "Version",
		},
		cli.BoolFlag{
			Name:  "no-prefix, p",
			Usage: "Disable prefix for unit names",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.Bool("V") {
			fmt.Printf("Restartd daemon git:%s\n", GITCOMMIT)
			return nil
		}

		if c.Bool("p") {
			prefix = false
		}

		fldr := os.Getenv("RESTARTD_SOCKET_FOLDER")
		if fldr == "" {
			fldr = "/tmp/restartd"
		}

		err := os.Chmod(fldr, os.FileMode(0755))
		if err != nil {
			return err
		}
		servers := server.NewServerUsers(fldr, "restart.sock")

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
				r := &restartd.Restartd{
					PrefixService: prefix,
					User:          conf.User,
					Services:      conf.Services,
				}
				myserver, err := servers.AddUser(conf.User)
				if err != nil {
					panic(err)
				}
				myserver.Register("statusAll", r.StatusAll)
				myserver.Register("status", r.Status)
				myserver.Register("start", r.Start)
				myserver.Register("stop", r.Stop)
				myserver.Register("restart", r.Restart)
				myserver.Register("reload", r.Reload)

				log.Info("Add user ", conf.User)
			}
			log.Info("Number of users : ", len(confs))
		}
		// initial config
		configs()

		cc := make(chan os.Signal, 1)
		signal.Notify(cc, os.Interrupt, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGTERM)
		go func() {
			for {
				s := <-cc
				log.Info("Signal : ", s)
				switch s {
				case os.Interrupt:
					servers.Stop()
				case syscall.SIGTERM:
					servers.Stop()
				case syscall.SIGHUP:
					configs()
				}
			}
		}()

		// listen and block
		servers.Serve()
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		// FIXME yell to STDERR
		fmt.Println(err)
		os.Exit(-1)
	}
}
