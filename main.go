package main

import (
	"fmt"
	"net"
	"os"
)

type Restartd struct {
	socketHome string
	sockets    map[string]*net.UnixListener
}

func New(socketHome string) *Restartd {
	r := Restartd{
		socketHome,
		make(map[string]*net.UnixListener),
	}
	return &r
}

func (r *Restartd) AddUser(user string) error {
	l, err := net.ListenUnix("unix", &net.UnixAddr{r.socketHome + "/" + user, "unix"})
	if err != nil {
		return err
	}
	r.sockets[user] = l
	go r.listen(user)
	return nil
}

func (r *Restartd) listen(user string) {
	for {
		conn, err := r.sockets[user].AcceptUnix()
		if err != nil {
			panic(err)
		}
		var buff [1024]byte
		n, err := conn.Read(buff[:])
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s: %s\n", user, string(buff[:n]))
	}
}

func (r *Restartd) RemoveUser(user string) {
	delete(r.sockets, user)
	os.Remove(r.socketHome + "/" + user)
}

func main() {
	r := New("/tmp")
	err := r.AddUser("pim")
	if err != nil {
		panic(err)
	}
	err = r.AddUser("pam")
	if err != nil {
		panic(err)
	}
	err = r.AddUser("poum")
	if err != nil {
		panic(err)
	}
}
