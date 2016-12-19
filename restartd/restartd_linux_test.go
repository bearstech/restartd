package restartd

import (
	"fmt"
	"testing"
)

func TestAllStatus(t *testing.T) {
	r := &Restartd{
		PrefixService: false,
		User:          "Bob",
		Services:      []string{"rsyslog"},
	}
	s, err := r.getAllStatus()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(s)
}
