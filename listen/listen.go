package listen

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/user"
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

type Dispatcher struct {
	socketHome string
	sockets    map[string]channel
	bus        chan bool
}

func New(socketHome string) *Dispatcher {
	r := Dispatcher{
		socketHome: socketHome,
		sockets:    make(map[string]channel),
		bus:        make(chan bool),
	}
	return &r
}

func (r *Dispatcher) socket(uzer *user.User) (*net.UnixListener, error) {
	// socket dir
	sd := r.socketHome + "/" + uzer.Username
	_, err := os.Stat(sd)
	if os.IsNotExist(err) {
		err = os.MkdirAll(sd, 0644)
		if err != nil {
			return nil, err
		}
	} else {
		if err != nil {
			return nil, err
		}
	}

	sp := sd + "/" + "restartctl.sock"

	_, err = os.Stat(sp)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	err = os.Remove(sp)
	if err == nil {
		return nil, err
	}

	l, err := net.ListenUnix("unix", &net.UnixAddr{Name: sp, Net: "unix"})
	if err != nil {
		return nil, err
	}

	// get uid user value as int
	uid, err := strconv.Atoi(uzer.Uid)
	if err != nil {
		return nil, err
	}

	// get gid user value as int
	gid, err := strconv.Atoi(uzer.Gid)
	if err != nil {
		return nil, err
	}

	// change socket ownsership to username
	err = os.Chown(sd, uid, gid)
	if err != nil {
		return nil, err
	}

	err = os.Chown(sp, uid, gid)
	if err != nil {
		return nil, err
	}
	return l, nil

}

func (r *Dispatcher) AddUser(username string, handler Handler) error {
	// don't add when it already exist
	if _, ok := r.sockets[username]; ok {
		return nil
	}
	// verify the user exists on the system
	uzer, err := user.Lookup(username)
	if err != nil {
		return err
	}

	l, err := r.socket(uzer)
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

func (r *Dispatcher) RemoveUser(user string) error {
	delete(r.sockets, user)
	return os.Remove(r.socketHome + "/" + user)
}

func (r *Dispatcher) Stop() {
	r.bus <- true
}

func (r *Dispatcher) Listen() {
	defer r.Cleanup()
	<-r.bus
}

func (r *Dispatcher) Cleanup() error {
	for user := range r.sockets {
		err := r.RemoveUser(user)
		if err != nil {
			return err
		}
	}
	fmt.Println("bye")
	return nil
}

type Echo struct {
}

func (e Echo) Handle(req []byte) []byte {
	return req
}
