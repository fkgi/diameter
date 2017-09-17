package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fkgi/diameter/msg"
	"github.com/fkgi/diameter/sock"
)

func main() {
	log.Println("initiator sample booting ...")

	ln := sock.Local{
		Host:  "init.test.com",
		Realm: "test.com",
		AuthApps: map[msg.VendorID][]msg.ApplicationID{
			0: []msg.ApplicationID{msg.AuthApplicationID(0xffffffff)}}}
	pn := sock.Peer{
		Host:       "hub.test.com",
		Realm:      "test.com",
		WDInterval: time.Second * 30,
		WDExpired:  3,
		SndTimeout: time.Second * 10,
		Handler:    handler,
		AuthApps: map[msg.VendorID][]msg.ApplicationID{
			0: []msg.ApplicationID{msg.AuthApplicationID(0xffffffff)}}}

	log.Println("connecting ...")
	c, e := net.Dial("tcp", "localhost:3868")
	if e != nil {
		log.Fatalln(" | invalid address:", e)
	}
	con, _ := ln.Dial(&pn, c)

	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-sigc
	log.Println("shutdown ...")
	con.Close(time.Millisecond * 100)
	for con.State() != "close" {
		time.Sleep(time.Millisecond * 100)
	}
}

func handler(m msg.Message) msg.Message {
	return msg.Message{}
}
