package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fkgi/diameter"
	"github.com/fkgi/diameter/connector"
	"github.com/fkgi/diameter/sctp"
)

var (
	cons      map[net.Conn]*diameter.Connection
	defaultRt diameter.Identity
)

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "spleader.internal"
	}
	dlocal := flag.String("l", hostname,
		"Diameter local host. `[(tcp|sctp)://][realm/]hostname[:port]`")
	droute := flag.String("d", "",
		"Diameter default route host. `hostname`")
	flag.Parse()

	defaultRt, err = diameter.ParseIdentity(*droute)
	if err != nil {
		log.Fatalln("invalid default route host:", err)
	}
	log.Printf("defoult route peer hostname is %s", defaultRt)

	log.Printf("booting spleader for Round-Robin <%s REV.%d>...",
		diameter.ProductName, diameter.FirmwareRev)

	diameter.DefaultRxHandler = rxhandler
	connector.TermSignals = []os.Signal{syscall.SIGINT, syscall.SIGTERM, os.Interrupt}
	connector.TermCause = diameter.Rebooting

	connector.TransportInfoNotify = func(src, dst net.Addr) {
		buf := new(strings.Builder)
		fmt.Fprintln(buf, "Detected transport address")
		if src != nil {
			fmt.Fprintf(buf, "| local: %s://%s\n", src.Network(), src.String())
		}
		if dst != nil {
			fmt.Fprintf(buf, "| peer : %s://%s\n", dst.Network(), dst.String())
		}
		log.Print(buf)
	}
	connector.TransportUpNotify = func(src, dst net.Addr) {
		buf := new(strings.Builder)
		fmt.Fprintln(buf, "Transport connection up")
		fmt.Fprintf(buf, "| local: %s://%s\n", src.Network(), src.String())
		fmt.Fprintf(buf, "| peer : %s://%s\n", dst.Network(), dst.String())
		log.Print(buf)
	}

	diameter.ConnectionUpNotify = func(c *diameter.Connection) {
		buf := new(strings.Builder)
		fmt.Fprintln(buf, "Diameter connection up")
		fmt.Fprintln(buf, "| local host/realm:", diameter.Host, "/", diameter.Realm)
		fmt.Fprintln(buf, "| peer  host/realm:", c.Host, "/", c.Realm)
		fmt.Fprintf(buf, "| applications:    %d\n", c.AvailableApplications())
		log.Print(buf)
	}
	diameter.TraceEvent = func(old, new, event string, err error) {
		if old != new || err != nil {
			log.Println("Diameter state update:",
				old, "->", new, "by event", event, "with error", err)
		}
	}
	diameter.TraceMessage = func(msg diameter.Message, dct diameter.Direction, err error) {
		if msg.AppID != 0 {
			log.Printf("%s diameter message handling: peer=%s, error=%v\n",
				dct, msg.PeerName, err)
		}
	}
	log.Println("listening Diameter...")
	log.Println("closed, error=", ListenAndServe(*dlocal))
}

func ListenAndServe(la string) (err error) {
	scheme, host, realm, ips, port, err := connector.ResolveIdentity(la)
	if err != nil {
		return
	}
	diameter.Host = host
	diameter.Realm = realm
	cons = make(map[net.Conn]*diameter.Connection)

	var l net.Listener
	switch scheme {
	case "sctp":
		src := &sctp.SCTPAddr{IP: ips, Port: port}
		if connector.TransportInfoNotify != nil {
			connector.TransportInfoNotify(src, nil)
		}
		l, err = sctp.ListenSCTP(src)
	default:
		src := &net.TCPAddr{IP: ips[0], Port: port}
		if connector.TransportInfoNotify != nil {
			connector.TransportInfoNotify(src, nil)
		}
		l, err = net.ListenTCP("tcp", src)
	}
	if err != nil {
		return
	}

	if len(connector.TermSignals) != 0 {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, connector.TermSignals...)
		go func() {
			<-sigc
			l.Close()
		}()
	}

	for {
		var c net.Conn
		c, err = l.Accept()
		if err != nil {
			break
		}
		if connector.TransportUpNotify != nil {
			connector.TransportUpNotify(c.LocalAddr(), c.RemoteAddr())
		}

		con := diameter.Connection{}
		cons[c] = &con
		go func() {
			con.ListenAndServe(c)
			delete(cons, c)
		}()
	}
	l.Close()

	cause := diameter.DoNotWantToTalkToYou
	if connector.TermCause == diameter.Rebooting {
		cause = diameter.Rebooting
	}
	for _, con := range cons {
		con.Close(cause)
	}

	for len(cons) != 0 {
		time.Sleep(time.Millisecond * 100)
	}

	return
}

func rxhandler(m diameter.Message) diameter.Message {
	var dHost diameter.Identity
	for rdr := bytes.NewReader(m.AVPs); rdr.Len() != 0; {
		a := diameter.AVP{}
		e := a.UnmarshalFrom(rdr)
		if e != nil {
			continue
		}
		if a.VendorID != 0 {
			continue
		}
		if a.Code == 293 {
			dHost, _ = diameter.GetDestinationHost(a)
		}
	}

	dcon := []*diameter.Connection{}
	if m.PeerName == defaultRt {
		if dHost == defaultRt {
			return m.GenerateAnswerBy(diameter.UnableToDeliver)
		}
		for _, con := range cons {
			if con.Host == dHost {
				dcon = append(dcon, con)
			}
		}
		if len(dcon) == 0 {
			for _, con := range cons {
				if con.Host != defaultRt {
					dcon = append(dcon, con)
				}
			}
		}
	} else {
		for _, con := range cons {
			if con.Host == defaultRt {
				dcon = append(dcon, con)
			}
		}
	}

	if len(dcon) == 0 {
		return m.GenerateAnswerBy(diameter.UnableToDeliver)
	}

	buf := bytes.NewBuffer(m.AVPs)
	diameter.SetRouteRecord(diameter.Host).MarshalTo(buf)
	m.AVPs = buf.Bytes()
	return dcon[rand.Intn(len(dcon))].DefaultTxHandler(m)
}
