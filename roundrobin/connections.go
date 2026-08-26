package main

import (
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/fkgi/diameter"
)

var (
	connections = make(chan map[net.Conn]*diameter.Connection, 1)
	reference   = []*diameter.Connection{}

	block = make(chan any)
)

func init() {
	connections <- make(map[net.Conn]*diameter.Connection)
}

func appendCon(c net.Conn, con *diameter.Connection, f func(net.Conn) error) {
	list := <-connections
	list[c] = con
	cons := make([]*diameter.Connection, 0, len(list))
	for _, v := range list {
		cons = append(cons, v)
	}
	reference = cons
	connections <- list
	log.Println("[WARN]", "Diameter connection", con.Host, "closed:", f(c))

	list = <-connections
	delete(list, c)
	cons = make([]*diameter.Connection, 0, len(list))
	for _, v := range list {
		cons = append(cons, v)
	}
	reference = cons
	connections <- list
}

func closeCon() {
	close(block)
	list := <-connections
	for _, c := range list {
		c.Close(diameter.Rebooting)
	}
	connections <- list

	for {
		list := <-connections
		l := len(list)
		connections <- list

		if l != 0 {
			time.Sleep(time.Millisecond * 100)
		} else {
			break
		}
	}
}

func selectCon(diameter.Message) *diameter.Connection {
	r := []*diameter.Connection{}
	for _, c := range reference {
		if c.State() == "open" {
			r = append(r, c)
		}
	}
	if l := len(r); l == 0 {
		return nil
	} else if l == 1 {
		return r[0]
	} else {
		return r[rand.Intn(l)]
	}
}
