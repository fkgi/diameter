package connector

import (
	"errors"
	"net"
	"strings"

	"github.com/fkgi/diameter"
	"github.com/fkgi/diameter/sctp"
)

// Dial parse inputs and connect transport connection.
// Inputs are string of local(la) and peer(pa) host information with format for ResolveIdentity.
func Dial(la, pa string) (con net.Conn, host, realm diameter.Identity, err error) {
	var pips, lips []net.IP
	var pport, lport int
	var scheme string

	scheme, host, realm, pips, pport, err = ResolveIdentity(pa)
	if err != nil {
		err = errors.Join(errors.New("invalid peer identity"), err)
		return
	}
	_, diameter.Host, diameter.Realm, lips, lport, err = ResolveIdentity(la)
	if err != nil {
		err = errors.Join(errors.New("invalid local identity"), err)
		return
	}

	switch scheme {
	case "sctp":
		con, err = sctp.DialSCTP(
			&sctp.SCTPAddr{IP: lips, Port: lport}, &sctp.SCTPAddr{IP: pips, Port: pport})
	case "tcp", "":
		con, err = net.DialTCP("tcp",
			&net.TCPAddr{IP: lips[0], Port: lport}, &net.TCPAddr{IP: pips[0], Port: pport})
	default:
		err = errors.New("invalid transport scheme: " + scheme)
		return
	}
	if err != nil {
		err = errors.Join(errors.New("failed to dial transport connection"), err)
	}
	return
}

// Listen parse inputs and listen transport listener.
// Inputs are string of local(la) and peer(pa) host information with format for ResolveIdentity.
func Listen(la string) (l net.Listener, err error) {
	var ips []net.IP
	var port int
	var scheme string

	scheme, diameter.Host, diameter.Realm, ips, port, err = ResolveIdentity(la)
	if err != nil {
		err = errors.Join(errors.New("invalid local identity"), err)
		return
	}

	switch scheme {
	case "sctp":
		l, err = sctp.ListenSCTP(&sctp.SCTPAddr{IP: ips, Port: port})
	case "tcp", "":
		l, err = net.ListenTCP("tcp", &net.TCPAddr{IP: ips[0], Port: port})
	default:
		err = errors.New("invalid transport scheme: " + scheme)
		return
	}
	if err != nil {
		err = errors.Join(errors.New("failed to listen transport connection"), err)
	}
	return
}

// Listen parse inputs and accept transport connection.
// Inputs are string of local(la) and peer(pa) host information with format for ResolveIdentity.
func Accept(la, pa string) (con net.Conn, host, realm diameter.Identity, err error) {
	var pips, lips []net.IP
	var pport, lport int
	var scheme string

	scheme, host, realm, pips, pport, err = ResolveIdentity(pa)
	if err != nil {
		err = errors.Join(errors.New("invalid peer identity"), err)
		return
	}
	_, diameter.Host, diameter.Realm, lips, lport, err = ResolveIdentity(la)
	if err != nil {
		err = errors.Join(errors.New("invalid local identity"), err)
		return
	}

	var l net.Listener
	var dst net.Addr
	switch scheme {
	case "sctp":
		l, err = sctp.ListenSCTP(&sctp.SCTPAddr{IP: lips, Port: lport})
		dst = &sctp.SCTPAddr{IP: pips, Port: pport}
	case "tcp", "":
		l, err = net.ListenTCP("tcp", &net.TCPAddr{IP: lips[0], Port: lport})
		dst = &net.TCPAddr{IP: pips[0], Port: pport}
	default:
		err = errors.New("invalid transport scheme: " + scheme)
	}
	if err != nil {
		err = errors.Join(errors.New("failed to listen transport connection"), err)
		return
	}

	ehs, ep, err := net.SplitHostPort(dst.String())
	if err == nil {
		con, err = l.Accept()
	}
	l.Close()
	if err != nil {
		return
	}

	ahs, ap, err := net.SplitHostPort(con.RemoteAddr().String())
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
		con.Close()
	}
	return
}
