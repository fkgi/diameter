package main

import (
	"flag"
	"log"
	"net"
	"os"
	"strings"

	"github.com/fkgi/diameter/connection"
	"github.com/fkgi/diameter/msg"
)

var logger = log.New(os.Stderr, "", log.Ltime|log.Lmicroseconds)

func main() {
	logger.Println("booting hub...")

	h, e := os.Hostname()
	if e != nil {
		h = "localhost.localnetwork"
	}
	port := flag.String("p", ":3868", "port")
	fqdn := flag.String("h", h, "diameter host name")
	flag.Parse()

	ln := connection.LocalNode{}
	ln.Host, e = msg.ParseDiameterIdentity(*fqdn)
	if e != nil {
		logger.Fatalln("invalid host name:", e)
	}
	ln.Realm, e = msg.ParseDiameterIdentity((*fqdn)[strings.Index(*fqdn, ".")+1:])
	if e != nil {
		logger.Fatalln("invalid host realm:", e)
	}

	logger.Println("listening ...")
	l, e := net.Listen("tcp", *port)
	if e != nil {
		logger.Fatalln(e)
	}

	for {
		c, e := l.Accept()
		if e != nil {
			logger.Fatalln(e)
		}
		con := ln.Accept(c)
	}
}
