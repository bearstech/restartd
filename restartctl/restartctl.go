package main

import (
	"fmt"
	"github.com/bearstech/restartd/protocol"
	"net"
	"os"
)

func main() {
	service := os.Args[1]
	command, err := protocol.ParseCommand(os.Args[2])
	if err != nil {
		panic(err)
	}
	msg := protocol.Message{
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
	var response protocol.Response
	err = protocol.Read(conn, &response)
	if err != nil {
		panic(err)
	}
	fmt.Println(response)
	os.Exit(int(*response.Code))
}
