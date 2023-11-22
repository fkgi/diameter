package connector

import (
	"errors"
	"net"
	"os"
	"os/signal"

	"github.com/fkgi/diameter"
	"github.com/fkgi/diameter/sctp"
)

var (
	TermSignals        []os.Signal    // Signals for closing diameter connection.
	ConnectionUpNotify func(net.Conn) // ConnectionUpNotify is called when transport connection up.

	con diameter.Connection // default Diameter connection
)

// DialAndServe start diameter connection handling process as initiator.
// Inputs are string of local(la) and peer(pa) host information with format for ResolveIdentiry.
func DialAndServe(la, pa string) (err error) {
	lscheme, host, realm, lips, lport, err := ResolveIdentiry(la)
	if err != nil {
		return
	}
	diameter.Host = host
	diameter.Realm = realm

	pscheme, host, realm, pips, pport, err := ResolveIdentiry(pa)
	if err != nil {
		return
	}
	con = diameter.Connection{Host: host, Realm: realm}

	if lscheme == "" {
		lscheme = "tcp"
	}
	if pscheme != "" && lscheme != pscheme {
		err = errors.New("transport protocol mismatch")
		return
	}

	var c net.Conn
	switch lscheme {
	case "sctp":
		c, err = sctp.DialSCTP(
			&sctp.SCTPAddr{IP: lips, Port: lport},
			&sctp.SCTPAddr{IP: pips, Port: pport})
	default:
		c, err = net.DialTCP("tcp",
			&net.TCPAddr{IP: lips[0], Port: lport},
			&net.TCPAddr{IP: pips[0], Port: pport})
	}
	if err != nil {
		return
	}
	if ConnectionUpNotify != nil {
		ConnectionUpNotify(c)
	}

	go termWithSignals(true)
	return con.DialAndServe(c)
}

// ListenAndServe start diameter connection handling process as responder.
// Inputs are string of local(la) host information with format for ResolveIdentiry.
func ListenAndServe(la string) (err error) {
	scheme, host, realm, ips, port, err := ResolveIdentiry(la)
	if err != nil {
		return
	}
	diameter.Host = host
	diameter.Realm = realm
	con = diameter.Connection{}

	var l net.Listener
	switch scheme {
	case "sctp":
		l, err = sctp.ListenSCTP(&sctp.SCTPAddr{IP: ips, Port: port})
	default:
		l, err = net.ListenTCP("tcp", &net.TCPAddr{IP: ips[0], Port: port})
	}
	if err != nil {
		return
	}

	c, err := l.Accept()
	l.Close()
	if err != nil {
		c.Close()
		return
	}

	/*
		if con.Host == "" {
			names, err := net.LookupAddr(c.RemoteAddr().String())
			if err == nil {
				con.Host, con.Realm, _, _, _ = resolveIdentiry(names[0])
			}
		}
	*/

	if ConnectionUpNotify != nil {
		ConnectionUpNotify(c)
	}

	go termWithSignals(false)
	return con.ListenAndServe(c)
}

func termWithSignals(isTx bool) {
	if len(TermSignals) == 0 {
		return
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, TermSignals...)
	<-sigc

	if isTx {
		con.Close(diameter.DoNotWantToTalkToYou)
	} else {
		con.Close(diameter.Rebooting)
	}
}
