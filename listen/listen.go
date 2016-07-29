package listen

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/user"
	"strconv"
	"strings"
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

func (r *Dispatcher) AddUser(username string, handler Handler) error {
	// don't add when it already exist
	if _, ok := r.sockets[username]; ok {
		return nil
	}
	// verify the user exists on the system
	User, err := user.Lookup(username)
	if err != nil {
		return err
	}

	// socket dir
	sd := r.socketHome + "/" + username
	fi, err := os.Stat(sd)
	if err != nil && !strings.HasSuffix(err.Error(), ": no such file or directory") {
		return err
	}

	if fi == nil {
		err = os.Mkdir(sd, 0644)
		if err != nil {
			return err
		}
	}

	sp := sd + "/" + "restartctl.sock"

	fm, err := os.Stat(sp)
	fmt.Println("Adduser fm ", fm, err)
	if err != nil && !strings.HasSuffix(err.Error(), ": no such file or directory") {
		return err
	}
	if fm != nil {
		err = os.Remove(sp)
		if err == nil {
			fmt.Println("Adduser error: ", err)
			return err
		}
	}

	l, err := net.ListenUnix("unix", &net.UnixAddr{Name: sp, Net: "unix"})
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
	err = os.Chown(sd, uid, gid)
	if err != nil {
		return err
	}

	err = os.Chown(sp, uid, gid)
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
