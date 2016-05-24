package main

import (
	"fmt"
	"net"
	"os"
)

type Handler interface {
	Handle(req []byte) []byte
}

type Restartd struct {
	socketHome string
	sockets    map[string]*net.UnixListener
	bus        chan bool
	handler    Handler
}

func New(socketHome string, handler Handler) *Restartd {
	r := Restartd{
		socketHome,
		make(map[string]*net.UnixListener),
		make(chan bool),
		handler,
	}
	return &r
}

func (r *Restartd) AddUser(user string) error {
	os.Remove(r.socketHome + "/" + user)
	l, err := net.ListenUnix("unix", &net.UnixAddr{r.socketHome + "/" + user, "unix"})
	if err != nil {
		return err
	}
	r.sockets[user] = l
	go r.listen(user)
	return nil
}

func (r *Restartd) RemoveUser(user string) {
	delete(r.sockets, user)
	os.Remove(r.socketHome + "/" + user)
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

func (r *Restartd) Listen() {
	defer r.Cleanup()
	<-r.bus
}

func (r *Restartd) Cleanup() {
	for user, _ := range r.sockets {
		r.RemoveUser(user)
	}
	fmt.Printf("bye")
}

type Echo struct {
}

func (e Echo) Handle(req []byte) []byte {
	return req
}

func main() {
	r := New("/tmp", Echo{})
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
	r.Listen()
}
