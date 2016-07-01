package main

import (
	"fmt"
	"github.com/bearstech/restartd/model"
	"github.com/bearstech/restartd/protocol"
	"net"
	"os"
)

func main() {
	service := os.Args[1]
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
