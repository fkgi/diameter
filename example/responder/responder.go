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
	log.Println("responder sample booting ...")

	l, e := net.Listen("tcp", "localhost:3868")
	if e != nil {
		logger.Fatalln(e)
	}
	c, e := l.Accept()
	if e != nil {
		logger.Fatalln(e)
	}
	ln := connection.LocalNode{
		Realm: msg.DiameterIdentity("test.com"),
		Host:  msg.DiameterIdentity("resp.test.com"),
		Properties: connection.Properties{
			Tw: time.Duration(30) * time.Second,
			Ew: 3,
			Ts: time.Duration(100) * time.Millisecond,
			Tp: time.Duration(30) * time.Second,
			Cp: 3,
			Apps: []connection.AuthApplication{
				{VendorID: 0, AppID: 0},
				{VendorID: 0, AppID: 0xffffffff}}}}
	ln.InitIDs()
	con := ln.Accept(c)

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
		isock, osock, conf := example.GeneratePath("responder")
		la, ln, _, pn := example.LoadConfig(conf)

		// open Diameter socket
		log.Println("start connecting Diameter connection")
		pl := provider.Listen(ln)
		prov := pl.AddPeer(pn)
		prov.Open()
		if e := pl.Bind(la); e != nil {
			log.Fatalln(e)
		}
		time.Sleep(time.Second)

		// open UNIX socket
		log.Println("start listening on UNIX socket", isock, "and", osock)
		example.RunUnixsockRelay(prov, isock, osock)
		time.Sleep(time.Second)
	*/
}
