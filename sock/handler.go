package sock

import (
	"net"
	"strings"

	"github.com/fkgi/diameter/msg"
	"github.com/fkgi/diameter/rfc6733"
)

// HandleMSG is diameter request handler
var HandleMSG = func(m msg.Request) msg.Answer {
	return nil
}

// MakeCER returns new CER
var MakeCER = func(c *Conn) rfc6733.CER {
	ips := make([]net.IP, 0, 2)
	s := c.con.LocalAddr().String()
	s = s[:strings.LastIndex(s, ":")]
	for _, i := range strings.Split(s, "/") {
		ips = append(ips, net.ParseIP(i))
	}

	return rfc6733.CER{
		OriginHost:       Host,
		OriginRealm:      Realm,
		HostIPAddress:    ips,
		VendorID:         VendorID,
		ProductName:      ProductName,
		OriginStateID:    StateID,
		ApplicationID:    getSupportedApps(),
		FirmwareRevision: FirmwareRevision}
}

// HandleCER is CER handler function
var HandleCER = func(r rfc6733.CER, c *Conn) rfc6733.CEA {
	ips := make([]net.IP, 0, 2)
	s := c.con.LocalAddr().String()
	s = s[:strings.LastIndex(s, ":")]
	for _, i := range strings.Split(s, "/") {
		ips = append(ips, net.ParseIP(i))
	}

	result := rfc6733.DiameterSuccess
	if c.peer == nil {
		c.peer = &Peer{Host: r.OriginHost, Realm: r.OriginRealm}
	} else if r.OriginHost != c.peer.Host || r.OriginRealm != c.peer.Realm {
		result = rfc6733.DiameterUnknownPeer
	}

	if result == rfc6733.DiameterSuccess {
		a := make(map[uint32][]uint32)
		apps := c.peer.AuthApps
		if apps == nil {
			apps = getSupportedApps()
		}
		for vID, aIDs := range r.ApplicationID {
			if _, ok := apps[vID]; !ok {
				continue
			}
			for _, aID := range match(apps[vID], aIDs) {
				if _, ok := a[vID]; !ok {
					a[vID] = []uint32{aID}
				} else {
					a[vID] = append(a[vID], aID)
				}
			}
		}
		if len(a) == 0 {
			result = rfc6733.DiameterApplicationUnsupported
		}
		c.peer.AuthApps = a
	}

	if c.peer.WDInterval == 0 {
		c.peer.WDInterval = WDInterval
	}
	if c.peer.WDExpired == 0 {
		c.peer.WDExpired = WDExpired
	}
	if c.peer.SndTimeout == 0 {
		c.peer.SndTimeout = SndTimeout
	}

	return rfc6733.CEA{
		ResultCode:       result,
		OriginHost:       Host,
		OriginRealm:      Realm,
		HostIPAddress:    ips,
		VendorID:         VendorID,
		ProductName:      ProductName,
		OriginStateID:    StateID,
		ApplicationID:    c.peer.AuthApps,
		FirmwareRevision: FirmwareRevision}
}

func match(a, b []uint32) []uint32 {
	r := make([]uint32, 0, len(a))
	for _, va := range a {
		for _, vb := range b {
			if va == vb {
				r = append(r, va)
			}
		}
	}
	return r
}

// HandleCEA is CEA handler function
var HandleCEA = func(m rfc6733.CEA, c *Conn) {
	c.peer.AuthApps = m.ApplicationID
}

// MakeDWR returns new DWR
var MakeDWR = func(c *Conn) rfc6733.DWR {
	dwr := rfc6733.DWR{
		OriginHost:    Host,
		OriginRealm:   Realm,
		OriginStateID: StateID}
	return dwr
}

// HandleDWR is DWR handler function
var HandleDWR = func(r rfc6733.DWR, c *Conn) rfc6733.DWA {
	dwa := rfc6733.DWA{
		ResultCode:    rfc6733.DiameterSuccess,
		OriginHost:    Host,
		OriginRealm:   Realm,
		OriginStateID: StateID}
	if c.peer.Host != r.OriginHost || c.peer.Realm != r.OriginRealm {
		dwa.ResultCode = rfc6733.DiameterUnknownPeer
	}

	return dwa
}

// HandleDWA is DWA handler function
var HandleDWA = func(r rfc6733.DWA, c *Conn) {
}

// MakeDPR returns new DWR
var MakeDPR = func(c *Conn) rfc6733.DPR {
	return rfc6733.DPR{
		OriginHost:      Host,
		OriginRealm:     Realm,
		DisconnectCause: rfc6733.Rebooting}
}

// HandleDPR is DPR handler function
var HandleDPR = func(r rfc6733.DPR, c *Conn) rfc6733.DPA {
	dpa := rfc6733.DPA{
		ResultCode:  rfc6733.DiameterSuccess,
		OriginHost:  Host,
		OriginRealm: Realm}
	if c.peer.Host != r.OriginHost || c.peer.Realm != r.OriginRealm {
		dpa.ResultCode = rfc6733.DiameterUnknownPeer
	}
	return dpa
}

// HandleDPA is DPA handler function
var HandleDPA = func(r rfc6733.DPA, c *Conn) {
}
