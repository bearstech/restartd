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

func ask(service *string, command *model.Message_Commands) (response *model.Response, err error) {

	socket := os.Getenv("RESTARTCTL_SOCKET")
	if socket == "" {
		socket = "/tmp/restartctl/restartctl.sock"
	}
	conn, err := net.DialUnix("unix", nil, &net.UnixAddr{Name: socket,
		Net: "unix"})
	if err != nil {
		return nil, err
	}

	msg := model.Message{
		Service: service,
		Command: command,
	}
	err = protocol.Write(conn, &msg)
	if err != nil {
		return nil, err
	}

	response = &model.Response{}
	err = protocol.Read(conn, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

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
		cli.BoolFlag{
			Name:  "status-all",
			Usage: "All status",
		},
	}

	app.Action = func(c *cli.Context) error {
		if c.Bool("V") {
			fmt.Printf("Restartcl CLI git:%s\n", GITCOMMIT)
			return nil
		}
		if c.Bool("status-all") {
			service := "--all"
			command := model.Message_status
			response, err := ask(&service, &command)
			if err != nil {
				return err
			}
			//FIXME display blabla
			fmt.Println(response)
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
		response, err := ask(&service, &command)
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
