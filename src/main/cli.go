package main

import (
	"kvraft"
)

func main() {
	user := raftkv.Command{}
	user.Setup("tmp", false, false, false, -1)
	user.Put("foo", "bar")
	user.Put("foo", "bar1")
	user.Get("foo")
}
