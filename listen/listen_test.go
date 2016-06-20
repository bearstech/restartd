package listen

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/user"
	"testing"
)

type HandlerTest struct{}

func (h *HandlerTest) Handle(req io.Reader, resp io.Writer) {
	buf := make([]byte, 4)
	io.ReadFull(req, buf)
	fmt.Println("Buffer: ", buf)
	resp.Write(buf)
}

func TestListener(t *testing.T) {
	fm, err := os.Stat("/tmp")
	if err != nil {
		t.Error(err)
	}
	os.Mkdir("/tmp/test_restartd", fm.Mode())
	//defer os.Remove("/tmp/test_restartd")
	me, err := user.Current()
	if err != nil {
		t.Error(err)
	}
	t.Log("Me: ", me)
	l := New("/tmp/test_restartd")
	defer l.Cleanup()
	err = l.AddUser(me.Username, &HandlerTest{})
	if err != nil {
		t.Error(err)
	}
	conn, err := net.DialUnix("unix", nil, &net.UnixAddr{
		"/tmp/test_restartd/" + me.Username,
		"unix"})
	if err != nil {
		t.Error(err)
	}
	conn.Write([]byte("popo"))
	buf := make([]byte, 4)
	s, err := io.ReadFull(conn, buf)
	if err != nil {
		t.Error(err)
	}
	t.Log("response", s, buf)
	if string(buf) != "popo" {
		t.Error("Bad response")
	}
}
