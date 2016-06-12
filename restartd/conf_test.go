package main

import (
	"testing"
)

func TestConf(t *testing.T) {
	raw := `
---

user: bob
services:
 - bob-web
 - bob-worker

`
	conf := Conf{}
	err := ReadConf([]byte(raw), &conf)
	if err != nil {
		t.Error(err.Error())
	}
	if conf.User != "bob" {
		t.Error("The user is not Bob : " + conf.User)
	}
}
