package listen

import (
	"fmt"
	"net"
	"os"
	"os/signal"
)

type Handler interface {
	Handle(req []byte) []byte
}

type channel struct {
	name    string
	socket  *net.UnixListener
	handler Handler
}

type Restartd struct {
	socketHome string
	sockets    map[string]channel
	bus        chan bool
}

func New(socketHome string) *Restartd {
	r := Restartd{
		socketHome,
		make(map[string]channel),
		make(chan bool),
	}
	return &r
}

func (r *Restartd) AddUser(user string, handler Handler) error {
	os.Remove(r.socketHome + "/" + user)
	l, err := net.ListenUnix("unix", &net.UnixAddr{r.socketHome + "/" + user, "unix"})
	if err != nil {
		return err
	}
	r.sockets[user] = channel{
		user,
		l,
		handler,
	}
	go r.listen(user)
	return nil
}

func (r *Restartd) RemoveUser(user string) {
	delete(r.sockets, user)
	os.Remove(r.socketHome + "/" + user)
}

func (r *Restartd) listen(user string) {
	for {
		conn, err := r.sockets[user].socket.AcceptUnix()
		if err != nil {
			panic(err)
		}
		var buff [1024]byte
		n, err := conn.Read(buff[:])
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			conn.Close()
		}
		fmt.Printf("%s: %s\n", user, string(buff[:n]))
		conn.Write(r.sockets[user].handler.Handle(buff[:n]))
	}
}

func (r *Restartd) Listen() {
	defer r.Cleanup()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		r.bus <- true
	}()
	<-r.bus
}

func (r *Restartd) Cleanup() {
	for user, _ := range r.sockets {
		r.RemoveUser(user)
	}
	fmt.Println("bye")
}

type Echo struct {
}

func (e Echo) Handle(req []byte) []byte {
	return req
}
