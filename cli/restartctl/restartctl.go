package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/bearstech/ascetic-rpc/client"
	"github.com/bearstech/restartd/restartd"
	"github.com/urfave/cli"
)

var GITCOMMIT string
var VERSION string

func main() {

	app := cli.NewApp()
	app.Version = "git:" + GITCOMMIT
	app.Usage = fmt.Sprintf("%s is a CLI for Restartd", app.Name)
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
		socket := os.Getenv("RESTARTCTL_SOCKET")
		if socket == "" {
			socket = "/tmp/restartctl/restart.sock"
		}
		cl, err := client.NewClientUnix(socket)
		if err != nil {
			return err
		}
		if c.Bool("status-all") {
			var status restartd.Status
			err = cl.Do("statusAll", nil, &status)
			if err != nil {
				return err
			}
			//FIXME display blabla
			fmt.Println(status)
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
		function := c.Args().Get(1)
		if function == "status" {
			var states restartd.Status
			err = cl.Do(function, &restartd.Service{Name: service}, &states)
			if err != nil {
				return err
			}
			for _, state := range states.Status {
				fmt.Printf("%s %s\n", state.Name, restartd.Status_States_name[int32(state.State)])
			}
			return nil
		}
		err = cl.Do(function, &restartd.Service{Name: service}, nil)
		if err != nil {
			return err
		}
		fmt.Println("Done")
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		//FIXME os.Exit(int(*response.Code))
		os.Exit(-1)
	}
}
