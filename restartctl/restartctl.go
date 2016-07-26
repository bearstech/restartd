package main

import (
	"fmt"
	"github.com/bearstech/restartd/model"
	"github.com/bearstech/restartd/protocol"
	"net"
	"os"
	"strings"
)

var GITCOMMIT, VERSION string

func main() {
	if len(os.Args) == 1 {
		return
	}
	prems := os.Args[1]
	if strings.HasPrefix(prems, "-") {
		switch prems {
		case "-v":
			fmt.Printf("Restartcl %s\n", GITCOMMIT)
		case "-h":
			fmt.Println(`Restartcl

Try something like:

	restartctl toto start
			`)
		}
		return
	}

	service := prems
	command, err := model.ParseCommand(os.Args[2])
	if err != nil {
		panic(err)
	}
	msg := model.Message{
		Service: &service,
		Command: &command,
	}
	socket := os.Getenv("RESTARTCTL_SOCKET")
	if socket == "" {
		socket = "/tmp/restartctl"
	}
	conn, err := net.DialUnix("unix", nil, &net.UnixAddr{socket,
		"unix"})
	if err != nil {
		panic(err)
	}
	err = protocol.Write(conn, &msg)
	if err != nil {
		panic(err)
	}
	var response model.Response
	err = protocol.Read(conn, &response)
	if err != nil {
		panic(err)
	}
	fmt.Println(*response.Code, ":", *response.Message)
	os.Exit(int(*response.Code))
}
