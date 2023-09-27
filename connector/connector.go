package connector

import (
	"errors"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/fkgi/diameter"
	"github.com/fkgi/diameter/sctp"
)

var (
	TermSignals []os.Signal // Signals for closing diameter connection

	con diameter.Connection // default Diameter connection
)

// DialAndServe start diameter connection handling process as initiator.
// Inputs are string of local hostname[:port][/realm] (la),
// peer hostname[:port][/realm] (ra) and bool flag to use SCTP.
func DialAndServe(la, pa string, isSctp bool) (err error) {
	con = diameter.Connection{}
	var c net.Conn

	if isSctp {
		tla := &sctp.SCTPAddr{}
		diameter.Host, diameter.Realm, tla.IP, tla.Port, err = resolveIdentiry(la)
		if err != nil {
			return
		}
		tpa := &sctp.SCTPAddr{}
		con.Host, con.Realm, tpa.IP, tpa.Port, err = resolveIdentiry(pa)
		if err != nil {
			return
		}
		c, err = sctp.DialSCTP(tla, tpa)
	} else {
		var ips []net.IP
		tla := &net.TCPAddr{}
		diameter.Host, diameter.Realm, ips, tla.Port, err = resolveIdentiry(la)
		if err != nil {
			return
		}
		tla.IP = ips[0]
		tpa := &net.TCPAddr{}
		con.Host, con.Realm, ips, tpa.Port, err = resolveIdentiry(pa)
		if err != nil {
			return
		}
		tpa.IP = ips[0]
		c, err = net.DialTCP("tcp", tla, tpa)
	}
	if err != nil {
		return
	}

	go termWithSignals(true)
	return con.DialAndServe(c)
}

// ListenAndServe start diameter connection handling process as responder.
// Inputs are string of local hostname (la) and bool flag to use SCTP.
// If Peer is nil, any peer is accepted.
func ListenAndServe(la string, isSctp bool) (err error) {
	con = diameter.Connection{}
	var c net.Conn
	var l net.Listener

	if isSctp {
		tla := &sctp.SCTPAddr{}
		diameter.Host, diameter.Realm, tla.IP, tla.Port, err = resolveIdentiry(la)
		if err != nil {
			return
		}
		l, err = sctp.ListenSCTP(tla)
	} else {
		var ips []net.IP
		tla := &net.TCPAddr{}
		diameter.Host, diameter.Realm, ips, tla.Port, err = resolveIdentiry(la)
		if err != nil {
			return
		}
		tla.IP = ips[0]
		l, err = net.ListenTCP("tcp", tla)
	}
	if err != nil {
		return
	}

	t := time.AfterFunc(diameter.WDInterval, func() {
		l.Close()
	})

	c, err = l.Accept()
	t.Stop()
	if err != nil {
		c.Close()
		l.Close()
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
	go termWithSignals(false)
	return con.ListenAndServe(c)
}

func resolveIdentiry(fqdn string) (host, realm diameter.Identity, ip []net.IP, port int, err error) {
	f := strings.Split(fqdn, "/")
	h, p, e := net.SplitHostPort(f[0])
	if e != nil {
		err = e
		return
	}
	if p == "" {
		p = "3868"
	}

	if host, err = diameter.ParseIdentity(h); err != nil {
		return
	}
	if len(f) > 1 {
		if realm, err = diameter.ParseIdentity(f[1]); err != nil {
			return
		}
	} else if i := strings.Index(h, "."); i < 0 {
		err = errors.New("domain part not found in local hostname")
		return
	} else if realm, err = diameter.ParseIdentity(h[i+1:]); err != nil {
		return
	}

	a, e := net.LookupHost(h)
	if e != nil {
		err = e
		return
	}
	ip = make([]net.IP, 0)
	for _, s := range a {
		ip = append(ip, net.ParseIP(s))
	}

	port, err = strconv.Atoi(p)
	return
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
