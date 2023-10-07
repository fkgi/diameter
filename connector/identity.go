package connector

import (
	"errors"
	"net"
	"strconv"
	"strings"

	"github.com/fkgi/abnf"
	"github.com/fkgi/diameter"
)

/*
ResolveIdentiry parse Diameter peer parameter.

[tcp|sctp://][realm/]hostname[:port]

Default transport is tcp.
Default realm is generated from hostname (following text after first ".").
Default port is 3868.
IP addresses are resolved from hostname.
*/
func ResolveIdentiry(uri string) (
	scheme string, host, realm diameter.Identity, ip []net.IP, port int, err error) {
	t := abnf.ParseString(uri, _uri())
	if t == nil {
		err = errors.New("invalid id text")
		return
	}

	if c := t.Child(idSCHEME); c != nil {
		scheme = string(c.V)
	} else {
		scheme = ""
	}
	switch scheme {
	case "tcp", "sctp", "":
	default:
		err = errors.New("invalid transport protocol")
		return
	}

	h := ""
	if c := t.Child(idHOST); c != nil {
		h = string(c.V)
	}
	if h == "" {
		err = errors.New("invalid hostname")
		return
	}
	if host, err = diameter.ParseIdentity(h); err != nil {
		return
	}

	if c := t.Child(idExtREALM); c != nil {
		if c = c.Child(idREALM); c != nil {
			realm = diameter.Identity(c.V)
		} else {
			realm = ""
		}
	} else {
		realm = ""
	}
	if realm == "" {
		if i := strings.Index(h, "."); i < 0 {
			err = errors.New("domain part not found in local hostname")
			return
		} else if realm, err = diameter.ParseIdentity(h[i+1:]); err != nil {
			return
		}
	}

	p := ""
	if c := t.Child(idPORT); c != nil {
		p = string(c.V)
	}
	if p == "" {
		p = "3868"
	}
	port, err = strconv.Atoi(p)
	if err != nil {
		return
	}
	if port > 65535 {
		err = errors.New("invalid port number")
		return
	}

	a, err := net.LookupHost(h)
	if err != nil {
		return
	}
	ip = make([]net.IP, 0)
	for _, s := range a {
		ip = append(ip, net.ParseIP(s))
	}

	return
}

const (
	idSCHEME int = iota
	idREALM
	idExtREALM
	idHOST
	idPORT
)

func _uri() abnf.Rule {
	return abnf.C(
		abnf.O(_scheme()),
		abnf.O(abnf.K(_realm(), idExtREALM)),
		_host(),
		abnf.O(_port()),
		abnf.ETX())
}

func _scheme() abnf.Rule {
	return abnf.A(
		abnf.C(abnf.K(abnf.VS("tcp"), idSCHEME), abnf.VS("://")),
		abnf.C(abnf.K(abnf.VS("sctp"), idSCHEME), abnf.VS("://")))
}

func _realm() abnf.Rule {
	return abnf.C(abnf.K(_fqdn(), idREALM), abnf.V('/'))
}

func _host() abnf.Rule {
	return abnf.K(_fqdn(), idHOST)
}

func _fqdn() abnf.Rule {
	return abnf.C(_label(), abnf.R0(abnf.C(abnf.V('.'), _label())))
}

func _label() abnf.Rule {
	return abnf.C(abnf.ALPHANUM(), abnf.R0(_ldhstr()))
}

func _ldhstr() abnf.Rule {
	return abnf.A(abnf.ALPHANUM(), abnf.C(abnf.V('-'), abnf.ALPHANUM()))
}

func _port() abnf.Rule {
	return abnf.C(abnf.V(':'), abnf.K(abnf.RV(1, -1, abnf.DIGIT()), idPORT))
}
