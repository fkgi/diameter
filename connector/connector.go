package connector

import (
	"errors"
	"net"
	"os"
	"os/signal"
	"strings"

	"github.com/fkgi/diameter"
	"github.com/fkgi/diameter/sctp"
)

const UndefinedCause diameter.Enumerated = -1

var (
	TermSignals []os.Signal                          // Signals for closing diameter connection.
	TermCause   diameter.Enumerated = UndefinedCause // Cause value for termination

	con diameter.Connection // default Diameter connection
)

// DialAndServe start diameter connection handling process as initiator.
// Inputs are string of local(la) and peer(pa) host information with format for ResolveIdentity.
func DialAndServe(la, pa string) (err error) {
	src, dst, err := resolveAddresses(la, pa)
	if err != nil {
		return
	}

	var c net.Conn
	switch src.Network() {
	case "sctp":
		ssrc := src.(*sctp.SCTPAddr)
		sdst := dst.(*sctp.SCTPAddr)
		c, err = sctp.DialSCTP(ssrc, sdst)
	default:
		tsrc := src.(*net.TCPAddr)
		tdst := dst.(*net.TCPAddr)
		c, err = net.DialTCP("tcp", tsrc, tdst)
	}
	if err != nil {
		return
	}

	if TransportUpNotify != nil {
		TransportUpNotify(c.LocalAddr(), c.RemoteAddr())
	}
	registerTermSignals(diameter.Rebooting)
	return con.DialAndServe(c)
}

// ListenAndServe start diameter connection handling process as responder.
// Inputs are string of local(la) and peer(pa) host information with format for ResolveIdentity.
// Connection request will be rejected if peer address is not same as pa.
func ListenAndServe(la, pa string) (err error) {
	src, dst, err := resolveAddresses(la, pa)
	if err != nil {
		return
	}

	ehs, ep, err := net.SplitHostPort(dst.String())
	if err != nil {
		return
	}

	var l net.Listener
	switch src.Network() {
	case "sctp":
		ssrc := src.(*sctp.SCTPAddr)
		l, err = sctp.ListenSCTP(ssrc)
	default:
		tsrc := src.(*net.TCPAddr)
		l, err = net.ListenTCP("tcp", tsrc)
	}
	if err != nil {
		return
	}

	c, err := l.Accept()
	l.Close()
	if err != nil {
		return
	}

	if TransportUpNotify != nil {
		TransportUpNotify(c.LocalAddr(), c.RemoteAddr())
	}
	ahs, ap, err := net.SplitHostPort(c.RemoteAddr().String())
	if err != nil {
	} else if ep != "0" && ep != ap {
		err = errors.New("connection peer transport port is invalid")
	} else {
		for _, eh := range strings.Split(ehs, "/") {
			ok := false
			for _, ah := range strings.Split(ahs, "/") {
				if eh == ah {
					ok = true
					break
				}
			}
			if !ok {
				err = errors.New("connection peer transport address is invalid")
				break
			}
		}
	}
	if err != nil {
		c.Close()
		return
	}

	registerTermSignals(diameter.DoNotWantToTalkToYou)
	return con.ListenAndServe(c)
}

var DefaultRouter diameter.Router = func(diameter.Message) *diameter.Connection {
	return &con
}

var TransportInfoNotify func(src, dst net.Addr) = nil
var TransportUpNotify func(src, dst net.Addr) = nil

func resolveAddresses(la, pa string) (src, dst net.Addr, err error) {
	scheme, host, realm, pips, pport, err := ResolveIdentity(pa)
	if err != nil {
		return
	}
	con = diameter.Connection{Host: host, Realm: realm}

	_, host, realm, lips, lport, err := ResolveIdentity(la)
	if err != nil {
		return
	}
	diameter.Host = host
	diameter.Realm = realm

	switch scheme {
	case "sctp":
		src = &sctp.SCTPAddr{IP: lips, Port: lport}
		dst = &sctp.SCTPAddr{IP: pips, Port: pport}
		if TransportInfoNotify != nil {
			TransportInfoNotify(src, dst)
		}
	case "tcp", "":
		src = &net.TCPAddr{IP: lips[0], Port: lport}
		dst = &net.TCPAddr{IP: pips[0], Port: pport}
		if TransportInfoNotify != nil {
			TransportInfoNotify(src, dst)
		}
	default:
		err = errors.New("invalid scheme")
	}
	return
}

func registerTermSignals(cause diameter.Enumerated) {
	if len(TermSignals) != 0 {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, TermSignals...)
		go func() {
			<-sigc
			if TermCause == diameter.Rebooting ||
				TermCause == diameter.DoNotWantToTalkToYou {
				con.Close(TermCause)
			} else {
				con.Close(cause)
			}
		}()
	}
}

// RxQueue returns length of Rx queue
func RxQueue() int { return con.RxQueue() }

// TxQueue returns length of Tx queue
func TxQueue() int { return con.TxQueue() }

// LocalAddr returns local address of transport connection
func LocalAddr() net.Addr { return con.LocalAddr() }

// PeerAddr returns transport connection
func PeerAddr() net.Addr { return con.PeerAddr() }

// PeerName returns hostname of diameter connection
func PeerName() diameter.Identity { return con.Host }

// PeerRealm returns realm of diameter connection
func PeerRealm() diameter.Identity { return con.Realm }

// State returns state machine state
func State() string { return con.State() }

// AvailableApplications returns supported application list
func AvailableApplications() []uint32 { return con.AvailableApplications() }
