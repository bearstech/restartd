package listen

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/user"
	"os/signal"
	"strconv"
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

func (r *Restartd) AddUser(username string, handler Handler) error {
	// verify the user exists on the system
	User, err := user.Lookup(username)
	if err != nil {
		return err
	}

	os.Remove(r.socketHome + "/" + username)

	// socket path
	sp := r.socketHome + "/" + username

	l, err := net.ListenUnix("unix", &net.UnixAddr{sp, "unix"})
	if err != nil {
		return err
	}

	// get uid user value as int
	uid, err := strconv.Atoi(User.Uid)
	if err != nil {
		return err
	}

	// get gid user value as int
	gid, err := strconv.Atoi(User.Gid)
	if err != nil {
		return err
	}

	// change socket ownsership to username
        err = os.Chown(sp , uid, gid)
	if err != nil {
		return err
	}

	c := channel{
		username,
		l,
		handler,
	}
	r.sockets[username] = c
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
