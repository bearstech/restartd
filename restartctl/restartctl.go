package main

import (
	"errors"
	"fmt"
	"github.com/bearstech/restartd/model"
	"github.com/bearstech/restartd/protocol"
	"github.com/urfave/cli"
	"net"
	"os"
)

var GITCOMMIT string
var VERSION string

func main() {

	app := cli.NewApp()
	app.Version = "git:" + GITCOMMIT
	app.Usage = "Restartcl is a CLI for Restartd"
	app.HideVersion = true

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "Version, V",
			Usage: "Version",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.Bool("V") {
			fmt.Printf("Restartcl CLI git:%s\n", GITCOMMIT)
			return nil
		}
		if c.NArg() == 0 {
			fmt.Println("Usage: service < option > | --status-all | [ service_name [ command | --full-restart ] ]")
			return nil
		}
		if c.NArg() == 1 {
			return errors.New("You need 2 arguments, a service and an action")
		}

		service := c.Args().Get(0)
		command, err := model.ParseCommand(c.Args().Get(1))
		if err != nil {
			return err
		}

		socket := os.Getenv("RESTARTCTL_SOCKET")
		if socket == "" {
			socket = "/tmp/restartctl.sock"
		}
		conn, err := net.DialUnix("unix", nil, &net.UnixAddr{Name: socket,
			Net: "unix"})
		if err != nil {
			return err
		}

		msg := model.Message{
			Service: &service,
			Command: &command,
		}
		err = protocol.Write(conn, &msg)
		if err != nil {
			return err
		}

		var response model.Response
		err = protocol.Read(conn, &response)
		if err != nil {
			return err
		}

		fmt.Println(*response.Code, ":", *response.Message)
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		//FIXME os.Exit(int(*response.Code))
		os.Exit(-1)
	}
}
