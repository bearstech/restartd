package main

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	application "github.com/bearstech/ascetic-rpc/app"
	"github.com/bearstech/ascetic-rpc/server"
	"github.com/bearstech/restartd/restartd"
	"github.com/urfave/cli"
)

var GITCOMMIT string

type RestartServer struct {
	confs      []*restartd.Conf
	servers    *server.ServerUsers
	confFolder string
	prefix     bool
}

func NewRestartServer(prefix bool) (*RestartServer, error) {
	fldr := os.Getenv("RESTARTD_SOCKET_FOLDER")
	if fldr == "" {
		fldr = "/tmp/restartd"
	}

	servers, err := server.NewServerUsers(fldr, "restart.sock").WithGroup("restartd")
	if err != nil {
		return nil, err
	}
	confFolder := os.Getenv("RESTARTD_CONF")
	if confFolder == "" {
		confFolder = "/etc/restartd/conf.d"
	}
	log.Info("Conf folder is ", confFolder)

	err = servers.MakeFolder()
	if err != nil {
		return nil, err
	}

	return &RestartServer{
		servers:    servers,
		confFolder: confFolder,
		prefix:     prefix,
	}, nil
}

func (rs *RestartServer) Config() error {
	confs, err := restartd.ReadConfFolder(rs.confFolder)
	if err != nil {
		return err
	}
	if len(confs) == 0 {
		log.Error("No conf found. Add some yml file in " + rs.confFolder)
		//os.Exit(-1)
	}
	olds := make(map[string]bool)
	for name, _ := range rs.servers.Names {
		olds[name] = true
	}
	for _, conf := range confs {
		delete(olds, conf.User)
		r := &restartd.Restartd{
			PrefixService: rs.prefix,
			User:          conf.User,
			Services:      conf.Services,
		}
		myserver, err := rs.servers.AddUser(conf.User)
		if err != nil {
			return err
		}
		myserver.Register("statusAll", r.StatusAll)
		myserver.Register("status", r.Status)
		myserver.Register("start", r.Start)
		myserver.Register("stop", r.Stop)
		myserver.Register("restart", r.Restart)
		myserver.Register("reload", r.Reload)

		log.Info("Add user ", conf.User)
	}
	for name, _ := range olds {
		rs.servers.RemoveUser(name)
	}
	log.Info("Number of users : ", len(confs))
	return nil
}

func (rs *RestartServer) GetServers() *server.ServerUsers {
	return rs.servers
}

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

		rs, err := NewRestartServer(prefix)
		if err != nil {
			return err
		}

		return application.NewApplication(rs).Start()
	}

	if err := app.Run(os.Args); err != nil {
		log.Error(err)
		os.Exit(-1)
	}
}
