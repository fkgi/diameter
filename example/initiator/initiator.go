package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fkgi/diameter/connection"
	"github.com/fkgi/diameter/msg"
)

func main() {
	logger := log.New(os.Stderr, "", log.Ltime|log.Lmicroseconds)
	connection.Notificator = func(e connection.Notify) { e.Log(logger) }

	// boot log
	log.Println("initiator sample booting ...")

	c, e := net.Dial("tcp", "localhost:3868")
	if e != nil {
		logger.Fatalln(e)
	}
	ln := connection.LocalNode{
		Realm: msg.DiameterIdentity("test.com"),
		Host:  msg.DiameterIdentity("init.test.com"),
		Properties: connection.Properties{
			Tw:   time.Duration(30) * time.Second,
			Ew:   3,
			Ts:   time.Duration(100) * time.Millisecond,
			Tp:   time.Duration(30) * time.Second,
			Cp:   3,
			Apps: [][2]uint32{{0, 0}, {0, 0xffffffff}}}}
	ln.InitIDs()
	pn := connection.PeerNode{
		Realm:      msg.DiameterIdentity("test.com"),
		Host:       msg.DiameterIdentity("resp.test.com"),
		Properties: ln.Properties}
	con := ln.Dial(&pn, c)

	sigc := make(chan os.Signal)
	signal.Notify(sigc, os.Interrupt, os.Kill, syscall.SIGTERM)

	<-sigc
	log.Println("shutdown ...")

	con.Close(msg.Rebooting)
	for con.State() != "Shutdown" {
		time.Sleep(time.Duration(100) * time.Millisecond)
	}
	/*
		// get option flag
		isock, osock, conf := example.GeneratePath("initiator")
		la, ln, pa, pn := example.LoadConfig(conf)

		// open Diameter socket
		log.Println("start connecting Diameter connection")
		pl := provider.Listen(ln)
		prov := pl.AddPeer(pn)

		pl.Dial(pn, la, pa)
		time.Sleep(time.Second)

		// open UNIX socket
		log.Println("start listening on UNIX socket", isock, "and", osock)
		example.RunUnixsockRelay(prov, isock, osock)
		time.Sleep(time.Second)
	*/
}
