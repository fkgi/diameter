package msg

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/fkgi/abnf"
)

const (
	idFQDN      = iota
	idSCHEME    = iota
	idPORT      = iota
	idTRANSPORT = iota
	idPROTOCOL  = iota
)

// DiameterIdentity is identity of Diameter protocol
type DiameterIdentity string

// ParseDiameterIdentity parse Diamter identity form string
func ParseDiameterIdentity(str string) (id DiameterIdentity, e error) {
	t := abnf.ParseString(str, _identity())
	if t == nil {
		e = fmt.Errorf("Invalid id text")
	} else {
		id = DiameterIdentity(t.Child(idFQDN).V)
	}
	return
}

// Compare compares two Diameter identity
func Compare(id1, id2 DiameterIdentity) int {
	s1 := strings.ToLower(string(id1))
	s2 := strings.ToLower(string(id2))

	l := 0
	r := 0
	if len(s1) > len(s2) {
		l = len(s2)
		r = 1
	} else if len(s1) < len(s2) {
		l = len(s1)
		r = -1
	}

	for i := 0; i < l; i++ {
		if s1[i] > s2[i] {
			return 1
		} else if s1[i] < s2[i] {
			return -1
		}
	}
	return r
}

// DiameterURI is URI of Diameter protocol
type DiameterURI struct {
	Scheme    string
	Fqdn      DiameterIdentity
	Port      int
	Transport string
	Protocol  string
}

// ParseDiameterURI parse Diamter URI form string
func ParseDiameterURI(str string) (uri DiameterURI, e error) {
	t := abnf.ParseString(str, _uri())
	if t == nil {
		e = fmt.Errorf("Invalid id text")
	} else {
		uri.Scheme = string(t.Child(idSCHEME).V)
		uri.Fqdn = DiameterIdentity(t.Child(idFQDN).V)
		p, _ := strconv.ParseInt(string(t.Child(idPORT).V), 10, 32)
		uri.Port = int(p)
		uri.Transport = string(t.Child(idTRANSPORT).V)
		uri.Protocol = string(t.Child(idPROTOCOL).V)
	}
	return
}

func (d DiameterURI) String() string {
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
	return abnf.C(_scheme(), _fqdn(), abnf.O(_port()), abnf.O(_transport()), abnf.O(_protocol()), abnf.EOF())
}

func _scheme() abnf.Rule {
	return abnf.A(
		abnf.C(abnf.K(abnf.VS("aaa"), idSCHEME), abnf.VS("://")),
		abnf.C(abnf.K(abnf.VS("aaas"), idSCHEME), abnf.VS("://")))
}

func _identity() abnf.Rule {
	return abnf.C(_fqdn(), abnf.EOF())
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
	return abnf.C(abnf.VS(";transport="), abnf.K(abnf.VSL("tcp", "sctp", "udp"), idTRANSPORT))
}

func _protocol() abnf.Rule {
	return abnf.C(abnf.VS(";protocol="), abnf.K(abnf.VSL("diameter", "radius", "tacacs+"), idPROTOCOL))
}
