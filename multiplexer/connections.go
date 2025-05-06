package main

import (
	"net"
	"time"

	"github.com/fkgi/diameter"
)

var (
	update    = make(map[net.Conn]*diameter.Connection)
	reference = []*diameter.Connection{}
	lock      = make(chan bool, 1)
)

func init() {
	lock <- true
}

func newConnection(c net.Conn) *diameter.Connection {
	con := new(diameter.Connection)
	<-lock
	update[c] = con
	cons := make([]*diameter.Connection, 0, len(update))
	for _, v := range update {
		cons = append(cons, v)
	}
	reference = cons
	lock <- true
	return con
}

func delConnection(c net.Conn) {
	<-lock
	delete(update, c)
	cons := make([]*diameter.Connection, 0, len(update))
	for _, v := range update {
		cons = append(cons, v)
	}
	reference = cons
	lock <- true
}

func wait() {
	for {
		<-lock
		if len(update) == 0 {
			lock <- true
			break
		}
		lock <- true
		time.Sleep(time.Millisecond * 100)
	}
}

func refConnection() []*diameter.Connection {
	return reference
}
