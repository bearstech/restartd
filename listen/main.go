package listen

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
)

type Handler interface {
	Handle(req io.Reader, resp io.Writer)
}

type channel struct {
	user    string
	socket  *net.UnixListener
	handler Handler
}

func (c *channel) listen() {
	for {
		conn, err := c.socket.AcceptUnix()
		if err != nil {
			panic(err)
		}
		c.handler.Handle(conn, conn)
	}
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
	c := channel{
		user,
		l,
		handler,
	}
	r.sockets[user] = c
	go c.listen()
	return nil
}

func (r *Restartd) RemoveUser(user string) {
	delete(r.sockets, user)
	os.Remove(r.socketHome + "/" + user)
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
