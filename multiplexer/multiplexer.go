package multiplexer

import (
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
	dest int
)

func ListenAndServe(la string) (err error) {
	scheme, host, realm, ips, port, err := connector.ResolveIdentiry(la)
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
