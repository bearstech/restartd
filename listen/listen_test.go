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
	resp.Write(buf)
}

func TestListener(t *testing.T) {
	fm, err := os.Stat("/tmp")
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat("/tmp/test_restartd")
	if err == nil || !os.IsNotExist(err) {
		err = os.RemoveAll("/tmp/test_restartd")
		if err != nil {
			t.Fatal(err)
		}
	} else {
		if err != nil {
			t.Fatal(err)
		}
	}

	err = os.Mkdir("/tmp/test_restartd", fm.Mode())
	if err != nil {
		t.Fatal(err)
	}
	me, err := user.Current()

	err = os.Mkdir("/tmp/test_restartd/"+me.Username, fm.Mode())
	if err != nil {
		t.Fatal(err)
	}
	//defer os.Remove("/tmp/test_restartd")
	t.Log("Me: ", me)
	l := New("/tmp/test_restartd")
	defer l.Cleanup()
	err = l.AddUser(me.Username, &HandlerTest{})
	if err != nil {
		t.Fatal(err)
	}
	socketPath := fmt.Sprintf("/tmp/test_restartd/%s/restartctl.sock", me.Username)
	fm, err = os.Stat(socketPath)
	if err != nil {
		t.Fatal(err)
	}
	conn, err := net.DialUnix("unix", nil, &net.UnixAddr{
		Name: socketPath,
		Net:  "unix"})
	if err != nil {
		t.Fatal(err)
	}
	conn.Write([]byte("popo"))
	buf := make([]byte, 4)
	s, err := io.ReadFull(conn, buf)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("response", s, buf)
	if string(buf) != "popo" {
		t.Error("Bad response")
	}
}
