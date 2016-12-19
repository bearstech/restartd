package systemd

import (
	"fmt"
	"testing"
)

func TestService(t *testing.T) {
	u, err := GetStatus("rsyslog")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(u)
	if u.Name != "rsyslog.service" {
		t.Fatal(fmt.Errorf("Bad name:", u.Name))
	}
}
