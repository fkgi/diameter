package multiplexer

import (
	"bytes"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/fkgi/diameter"
	"github.com/fkgi/diameter/connector"
	"github.com/fkgi/diameter/sctp"
)

var (
	cons map[net.Conn]*diameter.Connection
)

// ListenAndServe start diameter connection handling process as responder.
// Inputs are string of local(la) host information with format for ResolveIdentiry.
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
		l, err = sctp.ListenSCTP(&sctp.SCTPAddr{IP: ips, Port: port})
	default:
		l, err = net.ListenTCP("tcp", &net.TCPAddr{IP: ips[0], Port: port})
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

var DefaultRouter diameter.Router = func(m diameter.Message) *diameter.Connection {
	var dHost diameter.Identity
	for rdr := bytes.NewReader(m.AVPs); rdr.Len() != 0; {
		a := diameter.AVP{}
		if err := a.UnmarshalFrom(rdr); err != nil {
			continue
		}
		if a.VendorID != 0 {
			continue
		}
		switch a.Code {
		case 293:
			dHost, _ = diameter.GetDestinationHost(a)
		}
	}
	for _, con := range cons {
		if con.Host == dHost {
			return con
		}
	}

	t := rand.Intn(len(cons))
	i := 0
	for _, con := range cons {
		if i == t {
			return con
		}
		i++
	}
	return nil
}
