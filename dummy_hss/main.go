package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	dm "github.com/fkgi/diameter"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("booting Dummy-HSS <%s REV.%d>...", dm.ProductName, dm.FirmwareRev)

	tmp, err := os.Hostname()
	if err != nil {
		tmp = "hss.internal"
	}
	lDiaAddr := flag.String("l", tmp+":3868", "diameter local host:port")
	sctp := flag.Bool("sctp", false, "flag for sctp")
	flag.Parse()

	log.Printf("local address = %s\n", *lDiaAddr)
	if *sctp {
		log.Println("transport = sctp")
	} else {
		log.Println("transport = tcp")
	}

	dm.Handle(316, 16777251, 10415, handleULR)
	// dm.Handle(317, 16777251, 10415, handleCLR)
	dm.Handle(318, 16777251, 10415, handleAIR)
	// dm.Handle(319, 16777251, 10415, handleIDR)
	// dm.Handle(320, 16777251, 10415, handleDSR)
	dm.Handle(321, 16777251, 10415, handlePUR)
	// dm.Handle(322, 16777251, 10415, handleRSR)
	dm.Handle(323, 16777251, 10415, handleNOR)

	go func() {
		sigc := make(chan os.Signal, 2)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sigc

		dm.Close(dm.Busy)
	}()

	log.Println("listening...")
	dm.ListenAndServe(*lDiaAddr, *sctp)
}

func handleULR(retry bool, avp []byte) (bool, []byte) {
	buf := new(bytes.Buffer)
	dm.SetResultCode(dm.Success).MarshalTo(buf)
	dm.SetOriginHost(dm.Local.Host).MarshalTo(buf)
	dm.SetOriginRealm(dm.Local.Realm).MarshalTo(buf)
	return true, buf.Bytes()
}

func handleAIR(retry bool, avp []byte) (bool, []byte) {
	buf := new(bytes.Buffer)
	dm.SetResultCode(dm.Success).MarshalTo(buf)
	dm.SetOriginHost(dm.Local.Host).MarshalTo(buf)
	dm.SetOriginRealm(dm.Local.Realm).MarshalTo(buf)
	return true, buf.Bytes()
}

func handlePUR(retry bool, avp []byte) (bool, []byte) {
	buf := new(bytes.Buffer)
	dm.SetResultCode(dm.Success).MarshalTo(buf)
	dm.SetOriginHost(dm.Local.Host).MarshalTo(buf)
	dm.SetOriginRealm(dm.Local.Realm).MarshalTo(buf)
	return true, buf.Bytes()
}

func handleNOR(retry bool, avp []byte) (bool, []byte) {
	buf := new(bytes.Buffer)
	dm.SetResultCode(dm.Success).MarshalTo(buf)
	dm.SetOriginHost(dm.Local.Host).MarshalTo(buf)
	dm.SetOriginRealm(dm.Local.Realm).MarshalTo(buf)
	return true, buf.Bytes()
}
