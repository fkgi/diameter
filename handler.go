package diameter

import (
	"net"
	"strings"
)

// HandleMSG is diameter request handler
var HandleMSG = func(m Request) Answer {
	return nil
}

// MakeCER returns new CER
var MakeCER = func(c *Conn) CER {
	ips := make([]net.IP, 0, 2)
	s := c.con.LocalAddr().String()
	s = s[:strings.LastIndex(s, ":")]
	for _, i := range strings.Split(s, "/") {
		ips = append(ips, net.ParseIP(i))
	}

	return CER{
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
var HandleCER = func(r CER, c *Conn) CEA {
	ips := make([]net.IP, 0, 2)
	s := c.con.LocalAddr().String()
	s = s[:strings.LastIndex(s, ":")]
	for _, i := range strings.Split(s, "/") {
		ips = append(ips, net.ParseIP(i))
	}

	result := DiameterSuccess
	if c.peer == nil {
		c.peer = &Peer{Host: r.OriginHost, Realm: r.OriginRealm}
	} else if r.OriginHost != c.peer.Host || r.OriginRealm != c.peer.Realm {
		result = DiameterUnknownPeer
	}

	if result == DiameterSuccess {
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
			result = DiameterApplicationUnsupported
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

	return CEA{
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
var HandleCEA = func(m CEA, c *Conn) {
	c.peer.AuthApps = m.ApplicationID
}

// MakeDWR returns new DWR
var MakeDWR = func(c *Conn) DWR {
	dwr := DWR{
		OriginHost:    Host,
		OriginRealm:   Realm,
		OriginStateID: StateID}
	return dwr
}

// HandleDWR is DWR handler function
var HandleDWR = func(r DWR, c *Conn) DWA {
	dwa := DWA{
		ResultCode:    DiameterSuccess,
		OriginHost:    Host,
		OriginRealm:   Realm,
		OriginStateID: StateID}
	if c.peer.Host != r.OriginHost || c.peer.Realm != r.OriginRealm {
		dwa.ResultCode = DiameterUnknownPeer
	}

	return dwa
}

// HandleDWA is DWA handler function
var HandleDWA = func(r DWA, c *Conn) {
}

// MakeDPR returns new DWR
var MakeDPR = func(c *Conn) DPR {
	return DPR{
		OriginHost:      Host,
		OriginRealm:     Realm,
		DisconnectCause: Rebooting}
}

// HandleDPR is DPR handler function
var HandleDPR = func(r DPR, c *Conn) DPA {
	dpa := DPA{
		ResultCode:  DiameterSuccess,
		OriginHost:  Host,
		OriginRealm: Realm}
	if c.peer.Host != r.OriginHost || c.peer.Realm != r.OriginRealm {
		dpa.ResultCode = DiameterUnknownPeer
	}
	return dpa
}

// HandleDPA is DPA handler function
var HandleDPA = func(r DPA, c *Conn) {
}
