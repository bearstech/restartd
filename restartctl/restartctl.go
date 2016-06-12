package main

import (
	"../protocol"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

func main() {
	service := os.Args[1]
	command, err := protocol.ParseCommand(os.Args[2])
	if err != nil {
		panic(err)
	}
	msg := protocol.Message{service, command}
	conn, err := net.DialUnix("unix", nil, &net.UnixAddr{"/tmp/alice", "unix"})
	if err != nil {
		panic(err)
	}
	enc := gob.NewEncoder(conn)
	enc.Encode(&msg)
	dec := gob.NewDecoder(conn)
	var response protocol.Response
	dec.Decode(&response)
	fmt.Println(response)
}
