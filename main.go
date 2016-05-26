package main

import (
	"./listen"
)

func main() {
	r := listen.New("/tmp")
	err := r.AddUser("pim", listen.Echo{})
	if err != nil {
		panic(err)
	}
	err = r.AddUser("pam", listen.Echo{})
	if err != nil {
		panic(err)
	}
	err = r.AddUser("poum", listen.Echo{})
	if err != nil {
		panic(err)
	}
	r.Listen()
}
