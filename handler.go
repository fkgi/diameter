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
var MakeCER = defaultMakeCER

func defaultMakeCER(c *Conn) CER {
	ips := make([]net.IP, 0, 2)
	s := c.con.LocalAddr().String()
	s, _, _ = net.SplitHostPort(s)
	for _, i := range strings.Split(s, "/") {
		ips = append(ips, net.ParseIP(i))
		println(net.ParseIP(i) == nil)
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
var HandleCER = defaultHandleCER

func defaultHandleCER(r CER, c *Conn) CEA {
	ips := make([]net.IP, 0, 2)
	s := c.con.LocalAddr().String()
	s, _, _ = net.SplitHostPort(s)
	for _, i := range strings.Split(s, "/") {
		ips = append(ips, net.ParseIP(i))
	}

	result := DiameterSuccess
	if c.Peer == nil {
		c.Peer = &Peer{Host: r.OriginHost, Realm: r.OriginRealm}
	} else if r.OriginHost != c.Peer.Host || r.OriginRealm != c.Peer.Realm {
		result = DiameterUnknownPeer
	}

	if result == DiameterSuccess {
		if _, ok := supportedApps[0xffffffff]; ok && c.Peer.AuthApps == nil {
			c.Peer.AuthApps = getSupportedApps()
		} else {
			apps := c.Peer.AuthApps
			if apps == nil {
				apps = getSupportedApps()
			}
			a := make(map[uint32][]uint32)
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
				c.Peer.AuthApps = apps
			} else {
				c.Peer.AuthApps = a
			}
		}
	}

	if c.Peer.WDInterval == 0 {
		c.Peer.WDInterval = WDInterval
	}
	if c.Peer.WDExpired == 0 {
		c.Peer.WDExpired = WDExpired
	}

	return CEA{
		ResultCode:       result,
		OriginHost:       Host,
		OriginRealm:      Realm,
		HostIPAddress:    ips,
		VendorID:         VendorID,
		ProductName:      ProductName,
		OriginStateID:    StateID,
		ApplicationID:    c.Peer.AuthApps,
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
var HandleCEA = defaultHandleCEA

func defaultHandleCEA(m CEA, c *Conn) {
	c.Peer.AuthApps = m.ApplicationID
}

// MakeDWR returns new DWR
var MakeDWR = defaultMakeDWR

func defaultMakeDWR(c *Conn) DWR {
	dwr := DWR{
		OriginHost:    Host,
		OriginRealm:   Realm,
		OriginStateID: StateID}
	return dwr
}

// HandleDWR is DWR handler function
var HandleDWR = defaultHandleDWR

func defaultHandleDWR(r DWR, c *Conn) DWA {
	dwa := DWA{
		ResultCode:    DiameterSuccess,
		OriginHost:    Host,
		OriginRealm:   Realm,
		OriginStateID: StateID}
	if c.Peer.Host != r.OriginHost || c.Peer.Realm != r.OriginRealm {
		dwa.ResultCode = DiameterUnknownPeer
	}

	return dwa
}

// HandleDWA is DWA handler function
var HandleDWA = defaultHandleDWA

func defaultHandleDWA(r DWA, c *Conn) {
}

// MakeDPR returns new DWR
var MakeDPR = defaultMakeDPR

func defaultMakeDPR(c *Conn) DPR {
	return DPR{
		OriginHost:      Host,
		OriginRealm:     Realm,
		DisconnectCause: Rebooting}
}

// HandleDPR is DPR handler function
var HandleDPR = defaultHandleDPR

func defaultHandleDPR(r DPR, c *Conn) DPA {
	dpa := DPA{
		ResultCode:  DiameterSuccess,
		OriginHost:  Host,
		OriginRealm: Realm}
	if c.Peer.Host != r.OriginHost || c.Peer.Realm != r.OriginRealm {
		dpa.ResultCode = DiameterUnknownPeer
	}
	return dpa
}

// HandleDPA is DPA handler function
var HandleDPA = defaultHandleDPA

func defaultHandleDPA(r DPA, c *Conn) {
}
