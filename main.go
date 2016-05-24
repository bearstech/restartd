package main

import (
	"./listen"
)

func main() {
	r := listen.New("/tmp", listen.Echo{})
	err := r.AddUser("pim")
	if err != nil {
		panic(err)
	}
	err = r.AddUser("pam")
	if err != nil {
		panic(err)
	}
	err = r.AddUser("poum")
	if err != nil {
		panic(err)
	}
	r.Listen()
}
