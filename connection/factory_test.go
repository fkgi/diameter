package connection

import (
	"testing"
	"time"

	"net"

	"github.com/fkgi/diameter/msg"
)

func localnode() (*Connection, net.Listener, net.Conn) {
	l := LocalNode{
		Realm:   msg.DiameterIdentity("realm.com"),
		Host:    msg.DiameterIdentity("host.com"),
		StateID: 1,
		Properties: Properties{
			Tw:       time.Second * time.Duration(30),
			Ew:       3,
			Ts:       time.Millisecond * time.Duration(100),
			Cp:       3,
			Tp:       time.Second * time.Duration(30),
			AuthApps: map[msg.VendorID][]msg.ApplicationID{}}}
	l.InitIDs()

	nl, _ := net.Listen("tcp", "localhost:3868")
	nc, _ := net.Dial("tcp", "localhost:3868")

	c := Connection{
		local:    &l,
		peer:     nil,
		notify:   make(chan stateEvent),
		state:    closed,
		con:      nc,
		rcvstack: make(chan *msg.Message, MsgStackLen),
		sndstack: make(map[uint32]chan *msg.Message),
		openNtfy: make(chan bool, 1),
	}
	return &c, nl, nc
}

func printCERAAVPs(avp msg.GroupedAVP, t *testing.T) {
	if v, ok := msg.GetOriginHost(avp); ok {
		t.Log("Origin-Host:", v)
	}
	if v, ok := msg.GetOriginRealm(avp); ok {
		t.Log("Origin-Realm:", v)
	}
	for _, v := range msg.GetHostIPAddresses(avp) {
		t.Log("Host-IP-Address:", (net.IP(v)).String())
	}
	if v, ok := msg.GetVendorID(avp); ok {
		t.Log("Vendor-Id:", v)
	}
	if v, ok := msg.GetProductName(avp); ok {
		t.Log("Product-Name:", v)
	}
	if v, ok := msg.GetOriginStateID(avp); ok {
		t.Log("Origin-State-Id:", v)
	}
	for _, v := range msg.GetSupportedVendorIDs(avp) {
		t.Log("Supported-Vendor-Id:", v)
	}
	for _, v := range msg.GetAuthApplicationIDs(avp) {
		t.Log("Auth-Application-Id:", v)
	}
	for _, g := range msg.GetVendorSpecificApplicationIDs(avp) {
		t.Log("Vendor-Specific-Application-Id: ven:", g.VendorID, "app:", g.App)
	}
	if v, ok := msg.GetFirmwareRevision(avp); ok {
		t.Log("Firmware-Revision:", v)
	}
}

func TestMakeCER(t *testing.T) {
	c, nl, nc := localnode()

	c.local.AuthApps[0] = []msg.ApplicationID{
		msg.AuthApplicationID(0xffffffff),
		msg.AuthApplicationID(0)}
	c.local.AuthApps[10415] = []msg.ApplicationID{
		msg.AuthApplicationID(16777313),
		msg.AuthApplicationID(16777312)}

	m := c.makeCER()
	// s := new(bytes.Buffer)
	// m.PrintStack(s)
	// t.Log(s)

	if avp, e := m.Decode(); e == nil {
		printCERAAVPs(avp, t)
	}

	nc.Close()
	nl.Close()
}

func TestMakeCEArelay(t *testing.T) {
	c, nl, nc := localnode()
	c.local.AuthApps[10415] = []msg.ApplicationID{
		msg.AuthApplicationID(16777313),
		msg.AuthApplicationID(16777312)}

	r := c.makeCER()
	delete(c.local.AuthApps, 10415)
	c.local.AuthApps[0] = []msg.ApplicationID{
		msg.AuthApplicationID(0xffffffff)}

	p := &PeerNode{}
	if avp, e := r.Decode(); e == nil {
		if t, ok := msg.GetOriginHost(avp); ok {
			p.Host = msg.DiameterIdentity(t)
		}
		if t, ok := msg.GetOriginRealm(avp); ok {
			p.Realm = msg.DiameterIdentity(t)
		}
		getProvidedAuthApp(p, avp)
	}

	m, _ := c.makeCEA(r, p)

	if avp, e := m.Decode(); e == nil {
		if v, ok := msg.GetResultCode(avp); ok {
			t.Log("Result-Code:", v)
		}
		printCERAAVPs(avp, t)
	}

	nc.Close()
	nl.Close()
}

func TestMakeCEAhostfail(t *testing.T) {
	c, nl, nc := localnode()
	c.local.AuthApps[10415] = []msg.ApplicationID{
		msg.AuthApplicationID(16777313),
		msg.AuthApplicationID(16777312)}

	r := c.makeCER()
	delete(c.local.AuthApps, 10415)
	c.local.AuthApps[0] = []msg.ApplicationID{
		msg.AuthApplicationID(0xffffffff)}

	p := &PeerNode{}
	if avp, e := r.Decode(); e == nil {
		if t, ok := msg.GetOriginHost(avp); ok {
			p.Host = msg.DiameterIdentity(t)
		}
		if t, ok := msg.GetOriginRealm(avp); ok {
			p.Realm = msg.DiameterIdentity(t)
		}
		getProvidedAuthApp(p, avp)
	}

	c.peer = &PeerNode{
		Realm: msg.DiameterIdentity("realm.com"),
		Host:  msg.DiameterIdentity("validhost.com")}
	m, _ := c.makeCEA(r, p)

	if avp, e := m.Decode(); e == nil {
		if v, ok := msg.GetResultCode(avp); ok {
			t.Log("Result-Code:", v)
		}
		printCERAAVPs(avp, t)
	}

	nc.Close()
	nl.Close()
}

func TestMakeCEArealmfail(t *testing.T) {
	c, nl, nc := localnode()
	c.local.AuthApps[10415] = []msg.ApplicationID{
		msg.AuthApplicationID(16777313),
		msg.AuthApplicationID(16777312)}

	r := c.makeCER()
	delete(c.local.AuthApps, 10415)
	c.local.AuthApps[0] = []msg.ApplicationID{
		msg.AuthApplicationID(0xffffffff)}

	p := &PeerNode{}
	if avp, e := r.Decode(); e == nil {
		if t, ok := msg.GetOriginHost(avp); ok {
			p.Host = msg.DiameterIdentity(t)
		}
		if t, ok := msg.GetOriginRealm(avp); ok {
			p.Realm = msg.DiameterIdentity(t)
		}
		getProvidedAuthApp(p, avp)
	}

	c.peer = &PeerNode{
		Realm: msg.DiameterIdentity("validrealm.com"),
		Host:  msg.DiameterIdentity("host.com")}
	m, _ := c.makeCEA(r, p)

	if avp, e := m.Decode(); e == nil {
		if v, ok := msg.GetResultCode(avp); ok {
			t.Log("Result-Code:", v)
		}
		printCERAAVPs(avp, t)
	}

	nc.Close()
	nl.Close()
}

func TestMakeCEAappfail(t *testing.T) {
	c, nl, nc := localnode()
	c.local.AuthApps[10415] = []msg.ApplicationID{
		msg.AuthApplicationID(16777313),
		msg.AuthApplicationID(16777312)}

	r := c.makeCER()
	delete(c.local.AuthApps, 10415)
	c.local.AuthApps[0] = []msg.ApplicationID{
		msg.AuthApplicationID(0xffffffff)}

	p := &PeerNode{}
	if avp, e := r.Decode(); e == nil {
		if t, ok := msg.GetOriginHost(avp); ok {
			p.Host = msg.DiameterIdentity(t)
		}
		if t, ok := msg.GetOriginRealm(avp); ok {
			p.Realm = msg.DiameterIdentity(t)
		}
		getProvidedAuthApp(p, avp)
	}

	c.peer = &PeerNode{
		Realm: msg.DiameterIdentity("realm.com"),
		Host:  msg.DiameterIdentity("host.com"),
		Properties: Properties{
			AuthApps: map[msg.VendorID][]msg.ApplicationID{
				10415: []msg.ApplicationID{
					msg.AuthApplicationID(16777315),
					msg.AuthApplicationID(16777314)}}}}
	m, _ := c.makeCEA(r, p)

	if avp, e := m.Decode(); e == nil {
		if v, ok := msg.GetResultCode(avp); ok {
			t.Log("Result-Code:", v)
		}
		printCERAAVPs(avp, t)
	}

	nc.Close()
	nl.Close()
}

func TestMakeCEAappmatch(t *testing.T) {
	c, nl, nc := localnode()
	c.local.AuthApps[10415] = []msg.ApplicationID{
		msg.AuthApplicationID(16777313),
		msg.AuthApplicationID(16777312)}

	r := c.makeCER()

	p := &PeerNode{}
	if avp, e := r.Decode(); e == nil {
		if t, ok := msg.GetOriginHost(avp); ok {
			p.Host = msg.DiameterIdentity(t)
		}
		if t, ok := msg.GetOriginRealm(avp); ok {
			p.Realm = msg.DiameterIdentity(t)
		}
		getProvidedAuthApp(p, avp)
	}

	c.peer = &PeerNode{
		Realm: msg.DiameterIdentity("realm.com"),
		Host:  msg.DiameterIdentity("host.com"),
		Properties: Properties{
			AuthApps: map[msg.VendorID][]msg.ApplicationID{
				10415: []msg.ApplicationID{
					msg.AuthApplicationID(16777313),
					msg.AuthApplicationID(16777312)}}}}
	m, _ := c.makeCEA(r, p)

	if avp, e := m.Decode(); e == nil {
		if v, ok := msg.GetResultCode(avp); ok {
			t.Log("Result-Code:", v)
		}
		printCERAAVPs(avp, t)
	}

	nc.Close()
	nl.Close()
}

func TestMakeCEAsubappmatch(t *testing.T) {
	c, nl, nc := localnode()
	c.local.AuthApps[10415] = []msg.ApplicationID{
		msg.AuthApplicationID(16777313),
		msg.AuthApplicationID(16777312)}

	r := c.makeCER()
	delete(c.local.AuthApps, 10415)
	c.local.AuthApps[10415] = []msg.ApplicationID{
		msg.AuthApplicationID(16777313)}

	p := &PeerNode{}
	if avp, e := r.Decode(); e == nil {
		if t, ok := msg.GetOriginHost(avp); ok {
			p.Host = msg.DiameterIdentity(t)
		}
		if t, ok := msg.GetOriginRealm(avp); ok {
			p.Realm = msg.DiameterIdentity(t)
		}
		getProvidedAuthApp(p, avp)
	}

	/*
		c.peer = &PeerNode{
			Realm: msg.DiameterIdentity("realm.com"),
			Host:  msg.DiameterIdentity("host.com"),
			Properties: Properties{
				AuthApps: map[msg.VendorID][]msg.ApplicationID{
					10415: []msg.ApplicationID{
						msg.AuthApplicationID(16777313),
						msg.AuthApplicationID(16777312)}}}}
	*/
	m, _ := c.makeCEA(r, p)

	if avp, e := m.Decode(); e == nil {
		if v, ok := msg.GetResultCode(avp); ok {
			t.Log("Result-Code:", v)
		}
		printCERAAVPs(avp, t)
	}

	nc.Close()
	nl.Close()
}
