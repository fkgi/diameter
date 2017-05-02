package main

import (
	"log"
	"net"
	"os"

	"github.com/fkgi/diameter/connection"
	"github.com/fkgi/diameter/example/common"
)

func main() {
	logger := log.New(os.Stderr, "", log.Ltime|log.Lmicroseconds)
	connection.Notificator = func(e connection.Notice) { e.Log(logger) }
	common.Log = logger

	logger.Println("responder sample booting ...")

	isock, osock, conf := common.GeneratePath("responder")
	ln := connection.LocalNode{}
	la, _ := common.LoadConfig(conf, &ln, nil)

	il, ol := common.ListenUnixSock(isock, osock)

	logger.Println("listening ...")
	l, e := net.Listen(la.Network(), la.String())
	if e != nil {
		logger.Fatalln(e)
	}

	c, e := l.Accept()
	if e != nil {
		logger.Fatalln(e)
	}
	con := ln.Accept(c)
	common.RunUnixsockRelay(con, il, ol)

	common.CloseUnixSock(isock, osock)
}
