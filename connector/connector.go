package connector

import (
	"errors"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"

	"github.com/fkgi/diameter"
	"github.com/fkgi/diameter/sctp"
)

var (
	TermSignals []os.Signal              // Signals for closing diameter connection.
	TermCause   diameter.Enumerated = -1 // Cause value for termination

	con diameter.Connection // default Diameter connection
)

// DialAndServe start diameter connection handling process as initiator.
// Inputs are string of local(la) and peer(pa) host information with format for ResolveIdentiry.
func DialAndServe(la, pa string) (err error) {
	lscheme, host, realm, lips, lport, err := ResolveIdentity(la)
	if err != nil {
		return
	}
	diameter.Host = host
	diameter.Realm = realm

	pscheme, host, realm, pips, pport, err := ResolveIdentity(pa)
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

	if len(TermSignals) != 0 {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, TermSignals...)
		go func() {
			<-sigc
			if TermCause == diameter.Rebooting ||
				TermCause == diameter.DoNotWantToTalkToYou {
				con.Close(TermCause)
			} else {
				con.Close(diameter.Rebooting)
			}
		}()
	}

	return con.DialAndServe(c)
}

// ListenAndServe start diameter connection handling process as responder.
// Inputs are string of local(la) host information with format for ResolveIdentiry.
func ListenAndServe(la, pa string) (err error) {
	lscheme, host, realm, lips, lport, err := ResolveIdentity(la)
	if err != nil {
		return
	}
	diameter.Host = host
	diameter.Realm = realm

	pscheme, host, realm, pips, pport, err := ResolveIdentity(pa)
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

	var l net.Listener
	switch lscheme {
	case "sctp":
		l, err = sctp.ListenSCTP(&sctp.SCTPAddr{IP: lips, Port: lport})
	default:
		l, err = net.ListenTCP("tcp", &net.TCPAddr{IP: lips[0], Port: lport})
	}
	if err != nil {
		return
	}

	var c net.Conn
	for {
		c, err = l.Accept()
		if err != nil {
			c.Close()
			return
		}
		ra := c.RemoteAddr()
		hs, p, err := net.SplitHostPort(ra.String())
		if err != nil {
			c.Close()
			continue
		}
		if strconv.Itoa(pport) != p {
			c.Close()
			continue
		}
		if !checkIP(strings.Split(hs, "/"), pips) {
			c.Close()
			continue
		}
		break
	}
	l.Close()

	if len(TermSignals) != 0 {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, TermSignals...)
		go func() {
			<-sigc
			if TermCause == diameter.Rebooting ||
				TermCause == diameter.DoNotWantToTalkToYou {
				con.Close(TermCause)
			} else {
				con.Close(diameter.DoNotWantToTalkToYou)
			}
		}()
	}

	return con.ListenAndServe(c)
}

func checkIP(addrs []string, ips []net.IP) bool {
	for _, h := range addrs {
		i := net.ParseIP(h)
		if i == nil {
			return false
		}
		ok := false
		for _, ip := range ips {
			if ip.Equal(i) {
				ok = true
				break
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

var DefaultRouter diameter.Router = func(diameter.Message) *diameter.Connection {
	return &con
}
