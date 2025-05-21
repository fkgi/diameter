package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fkgi/diameter"
	"github.com/fkgi/diameter/connector"
)

var upLink diameter.Identity

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "multiplexer.internal"
	}
	dlocal := flag.String("l", hostname, "Diameter local host. `[(tcp|sctp)://][realm/]hostname[:port]`")
	hlocal := flag.String("i", ":12001", "HTTP local interface address. `[host]:port`")
	to := flag.Int("t", int(diameter.WDInterval/time.Second), "Message timeout timer [s]")
	help := flag.Bool("h", false, "Print usage")
	flag.Parse()

	diameter.WDInterval = time.Duration(*to) * time.Second

	upLink, err = diameter.ParseIdentity(flag.Arg(0))
	if *help || err != nil || upLink == "" {
		if err != nil {
			fmt.Println("invalid uplink peer hostname:", err)
			fmt.Println()
		} else if upLink == "" {
			fmt.Println("no uplink peer hostname")
			fmt.Println()
		}

		fmt.Printf("Usage: %s [OPTION]... UPLINK_PEER\n", os.Args[0])
		fmt.Println("UPLINK_PEER format is diameter hostname FQDN")
		fmt.Println()
		flag.PrintDefaults()
		return
	}

	log.Printf("[INFO] booting multiplexer for Round-Robin <%s REV.%d>...",
		diameter.ProductName, diameter.FirmwareRev)
	log.Printf("[INFO] uplink peer hostname is %s", upLink)

	http.HandleFunc("/diastate/v1/connection", conStateHandler)
	log.Println("[INFO] listening HTTP local port:", *hlocal)
	go func() {
		err := http.ListenAndServe(*hlocal, nil)
		if err != nil {
			log.Println("[ERROR]", "failed to listen HTTP, API is not available:",
				err)
		}
	}()

	diameter.DefaultRxHandler = rxhandler

	log.Println("[INFO]", "listening Diameter...")
	var l net.Listener
	l, err = connector.Listen(*dlocal)
	if err != nil {
		log.Fatalln("[ERROR]", err)
	}
	log.Println("[INFO]", "local host/realm:", diameter.Host, "/", diameter.Realm)

	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sigc

		l.Close()
	}()

	for {
		var c net.Conn
		if c, err = l.Accept(); err != nil {
			l.Close()
			break
		}
		buf := new(strings.Builder)
		fmt.Fprint(buf, "transport connection up")
		fmt.Fprintf(buf, "\n| local: %s://%s", c.LocalAddr().Network(), c.LocalAddr().String())
		fmt.Fprintf(buf, "\n| peer : %s://%s", c.RemoteAddr().Network(), c.RemoteAddr().String())
		log.Println("[INFO]", buf)

		con := newConnection(c)
		go func() {
			con.ListenAndServe(c)
			delConnection(c)
		}()
	}

	for _, con := range refConnection() {
		con.Close(diameter.Rebooting)
	}
	wait()
	log.Println("[INFO]", "closed")
}
