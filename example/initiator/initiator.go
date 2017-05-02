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

	logger.Println("initiator sample booting ...")

	isock, osock, conf := common.GeneratePath("initiator")
	ln := connection.LocalNode{}
	pn := connection.PeerNode{}
	_, pa := common.LoadConfig(conf, &ln, &pn)

	il, ol := common.ListenUnixSock(isock, osock)

	logger.Println("connecting ...")
	c, e := net.Dial(pa.Network(), pa.String())
	if e != nil {
		logger.Fatalln(e)
	}
	con := ln.Dial(&pn, c)
	common.RunUnixsockRelay(con, il, ol)

	common.CloseUnixSock(isock, osock)
}
