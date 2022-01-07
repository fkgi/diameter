package diameter

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/fkgi/abnf"
)

const (
	idFQDN int = iota
	idSCHEME
	idPORT
	idTRANSPORT
	idPROTOCOL
)

// Identity is identity of Diameter protocol
type Identity string

func (i Identity) String() string {
	return string(i)
}

// ParseIdentity parse Diamter identity form string
func ParseIdentity(str string) (id Identity, err error) {
	if t := abnf.ParseString(str, _identity()); t == nil {
		err = fmt.Errorf("invalid id text")
	} else {
		id = Identity(t.Child(idFQDN).V)
	}
	return
}

// URI is URI of Diameter protocol
type URI struct {
	Scheme    string
	Fqdn      Identity
	Port      int
	Transport string
	Protocol  string
}

// ParseURI parse Diamter URI form string
func ParseURI(str string) (uri URI, err error) {
	if t := abnf.ParseString(str, _uri()); t == nil {
		err = fmt.Errorf("invalid id text")
	} else {
		uri.Scheme = string(t.Child(idSCHEME).V)
		uri.Fqdn = Identity(t.Child(idFQDN).V)
		p, _ := strconv.ParseInt(string(t.Child(idPORT).V), 10, 32)
		uri.Port = int(p)
		uri.Transport = string(t.Child(idTRANSPORT).V)
		uri.Protocol = string(t.Child(idPROTOCOL).V)
	}
	return
}

func (d URI) String() string {
	var b bytes.Buffer
	b.WriteString(d.Scheme)
	b.WriteString("://")
	b.WriteString(string(d.Fqdn))
	if d.Port != 0 {
		b.WriteString(":")
		b.WriteString(strconv.Itoa(d.Port))
	}
	if len(d.Transport) != 0 {
		b.WriteString(";transport=")
		b.WriteString(d.Transport)
	}
	if len(d.Protocol) != 0 {
		b.WriteString(";protocol=")
		b.WriteString(d.Protocol)
	}
	return b.String()
}

func _uri() abnf.Rule {
	return abnf.C(
		_scheme(), _fqdn(),
		abnf.O(_port()),
		abnf.O(_transport()),
		abnf.O(_protocol()),
		abnf.ETX())
}

func _scheme() abnf.Rule {
	return abnf.A(
		abnf.C(abnf.K(abnf.VS("aaa"), idSCHEME), abnf.VS("://")),
		abnf.C(abnf.K(abnf.VS("aaas"), idSCHEME), abnf.VS("://")))
}

func _identity() abnf.Rule {
	return abnf.C(_fqdn(), abnf.ETX())
}

func _fqdn() abnf.Rule {
	return abnf.K(abnf.C(_label(), abnf.R0(abnf.C(abnf.V('.'), _label()))), idFQDN)
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

func _transport() abnf.Rule {
	return abnf.C(
		abnf.VS(";transport="),
		abnf.K(abnf.VSL("tcp", "sctp", "udp"), idTRANSPORT))
}

func _protocol() abnf.Rule {
	return abnf.C(
		abnf.VS(";protocol="),
		abnf.K(abnf.VSL("diameter", "radius", "tacacs+"), idPROTOCOL))
}
