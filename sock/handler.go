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
	r := rfc6733.CER{
		OriginHost:  rfc6733.OriginHost(Host),
		OriginRealm: rfc6733.OriginRealm(Realm),
		//HostIPAddress: make([]rfc6733.HostIPAddress, 0),
		VendorID:                    VendorID,
		ProductName:                 ProductName,
		OriginStateID:               StateID,
		SupportedVendorID:           make([]rfc6733.SupportedVendorID, 0),
		AuthApplicationID:           make([]rfc6733.AuthApplicationID, 0),
		VendorSpecificApplicationID: make([]rfc6733.VendorSpecificApplicationID, 0),
		FirmwareRevision:            &FirmwareRevision}

	switch c.con.LocalAddr().Network() {
	case "tcp":
		s := c.con.LocalAddr().String()
		s = s[:strings.LastIndex(s, ":")]
		r.HostIPAddress = []rfc6733.HostIPAddress{rfc6733.HostIPAddress(net.ParseIP(s))}
	case "sctp":
		s := c.con.LocalAddr().String()
		s = s[:strings.LastIndex(s, ":")]
		r.HostIPAddress = []rfc6733.HostIPAddress{}
		for _, i := range strings.Split(s, "/") {
			r.HostIPAddress = append(r.HostIPAddress,
				rfc6733.HostIPAddress(net.ParseIP(i)))
		}
	}
	for v, a := range getSupportedApps() {
		if v != 0 {
			r.SupportedVendorID = append(
				r.SupportedVendorID, rfc6733.SupportedVendorID(v))
			for _, i := range a {
				r.VendorSpecificApplicationID = append(
					r.VendorSpecificApplicationID,
					rfc6733.VendorSpecificApplicationID{
						VendorID:          v,
						AuthApplicationID: i})
			}
		} else {
			for _, i := range a {
				r.AuthApplicationID = append(r.AuthApplicationID, i)
			}
		}
	}
	return r
}

// HandleCER is CER handler function
var HandleCER = func(r rfc6733.CER, c *Conn) rfc6733.CEA {
	cea := rfc6733.CEA{
		ResultCode:                  rfc6733.DiameterSuccess,
		OriginHost:                  rfc6733.OriginHost(Host),
		OriginRealm:                 rfc6733.OriginRealm(Realm),
		VendorID:                    VendorID,
		ProductName:                 ProductName,
		OriginStateID:               StateID,
		SupportedVendorID:           make([]rfc6733.SupportedVendorID, 0),
		AuthApplicationID:           make([]rfc6733.AuthApplicationID, 0),
		VendorSpecificApplicationID: make([]rfc6733.VendorSpecificApplicationID, 0),
		FirmwareRevision:            &FirmwareRevision}
	switch c.con.LocalAddr().Network() {
	case "tcp":
		s := c.con.LocalAddr().String()
		s = s[:strings.LastIndex(s, ":")]
		cea.HostIPAddress = []rfc6733.HostIPAddress{rfc6733.HostIPAddress(net.ParseIP(s))}
	case "sctp":
		s := c.con.LocalAddr().String()
		s = s[:strings.LastIndex(s, ":")]
		cea.HostIPAddress = []rfc6733.HostIPAddress{}
		for _, i := range strings.Split(s, "/") {
			cea.HostIPAddress = append(cea.HostIPAddress, rfc6733.HostIPAddress(net.ParseIP(i)))
		}
	}

	if c.peer == nil {
		c.peer = &Peer{
			Host:  msg.DiameterIdentity(r.OriginHost),
			Realm: msg.DiameterIdentity(r.OriginRealm)}
	} else if msg.DiameterIdentity(r.OriginHost) != c.peer.Host ||
		msg.DiameterIdentity(r.OriginRealm) != c.peer.Realm {
		cea.ResultCode = rfc6733.DiameterUnknownPeer
	}
	if cea.ResultCode == rfc6733.DiameterSuccess {
		app := map[rfc6733.VendorID][]rfc6733.AuthApplicationID{}
		for _, i := range r.SupportedVendorID {
			app[rfc6733.VendorID(i)] = []rfc6733.AuthApplicationID{}
		}
		for _, a := range r.VendorSpecificApplicationID {
			if _, ok := app[a.VendorID]; !ok {
				app[a.VendorID] = []rfc6733.AuthApplicationID{}
			}
			app[a.VendorID] = append(app[a.VendorID], a.AuthApplicationID)
		}
		if len(r.AuthApplicationID) != 0 {
			if _, ok := app[0]; !ok {
				app[0] = []rfc6733.AuthApplicationID{}
			}
			for _, i := range r.AuthApplicationID {
				app[0] = append(app[0], i)
			}
		}

		if c.peer.AuthApps == nil {
			relay := rfc6733.AuthApplicationID(0xffffffff)
			for _, id := range getSupportedApps()[0] {
				if relay == id {
					c.peer.AuthApps = app
					break
				}
			}
			if c.peer.AuthApps == nil {
				c.peer.AuthApps = map[rfc6733.VendorID][]rfc6733.AuthApplicationID{}
				for key, ids := range app {
					for _, rid := range ids {
						if _, ok := getSupportedApps()[key]; !ok {
							continue
						}
						for _, lid := range getSupportedApps()[key] {
							if rid == lid {
								if _, ok := c.peer.AuthApps[key]; !ok {
									c.peer.AuthApps[key] = []rfc6733.AuthApplicationID{}
								}
								c.peer.AuthApps[key] = append(c.peer.AuthApps[key], rid)
							}
						}
					}
				}
				if len(c.peer.AuthApps) == 0 {
					cea.ResultCode = rfc6733.DiameterApplicationUnsupported
				}
			}
		} else {
			a := map[rfc6733.VendorID][]rfc6733.AuthApplicationID{}
			for key, ids := range app {
				for _, rid := range ids {
					if _, ok := c.peer.AuthApps[key]; !ok {
						continue
					}
					for _, lid := range c.peer.AuthApps[key] {
						if rid == lid {
							if _, ok := a[key]; !ok {
								a[key] = []rfc6733.AuthApplicationID{}
							}
							a[key] = append(a[key], rid)
						}
					}
				}
			}
			if len(a) == 0 {
				cea.ResultCode = rfc6733.DiameterApplicationUnsupported
			} else {
				c.peer.AuthApps = a
			}
		}
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

	for v, a := range c.peer.AuthApps {
		if v != 0 {
			cea.SupportedVendorID = append(
				cea.SupportedVendorID, rfc6733.SupportedVendorID(v))
			for _, i := range a {
				cea.VendorSpecificApplicationID = append(
					cea.VendorSpecificApplicationID,
					rfc6733.VendorSpecificApplicationID{
						VendorID:          v,
						AuthApplicationID: i})
			}
		} else {
			for _, i := range a {
				cea.AuthApplicationID = append(cea.AuthApplicationID, i)
			}
		}
	}
	return cea
}

// HandleCEA is CEA handler function
var HandleCEA = func(m rfc6733.CEA, c *Conn) {
	app := map[rfc6733.VendorID][]rfc6733.AuthApplicationID{}
	for _, i := range m.SupportedVendorID {
		app[rfc6733.VendorID(i)] = []rfc6733.AuthApplicationID{}
	}
	for _, a := range m.VendorSpecificApplicationID {
		if _, ok := app[a.VendorID]; !ok {
			app[a.VendorID] = []rfc6733.AuthApplicationID{}
		}
		app[a.VendorID] = append(app[a.VendorID], a.AuthApplicationID)
	}
	if len(m.AuthApplicationID) != 0 {
		if _, ok := app[0]; !ok {
			app[0] = []rfc6733.AuthApplicationID{}
		}
		for _, i := range m.AuthApplicationID {
			app[0] = append(app[0], i)
		}
	}
	c.peer.AuthApps = app
}

// MakeDWR returns new DWR
var MakeDWR = func(c *Conn) rfc6733.DWR {
	dwr := rfc6733.DWR{
		OriginHost:    rfc6733.OriginHost(Host),
		OriginRealm:   rfc6733.OriginRealm(Realm),
		OriginStateID: StateID}
	return dwr
}

// HandleDWR is DWR handler function
var HandleDWR = func(r rfc6733.DWR, c *Conn) rfc6733.DWA {
	dwa := rfc6733.DWA{
		ResultCode:  rfc6733.DiameterSuccess,
		OriginHost:  rfc6733.OriginHost(Host),
		OriginRealm: rfc6733.OriginRealm(Realm),
		//		*ErrorMessage
		//		*FailedAVP
		OriginStateID: StateID}
	if c.peer.Host != msg.DiameterIdentity(r.OriginHost) {
		dwa.ResultCode = rfc6733.DiameterUnknownPeer
	}
	if c.peer.Realm != msg.DiameterIdentity(r.OriginRealm) {
		dwa.ResultCode = rfc6733.DiameterUnknownPeer
	}

	return dwa
}

// HandleDWA is DWA handler function
var HandleDWA = func(r rfc6733.DWA, c *Conn) {
}

// MakeDPR returns new DWR
var MakeDPR = func(c *Conn) rfc6733.DPR {
	r := rfc6733.DPR{
		OriginHost:      rfc6733.OriginHost(Host),
		OriginRealm:     rfc6733.OriginRealm(Realm),
		DisconnectCause: rfc6733.Rebooting}
	return r
}

// HandleDPR is DPR handler function
var HandleDPR = func(r rfc6733.DPR, c *Conn) rfc6733.DPA {
	dpa := rfc6733.DPA{
		ResultCode:  rfc6733.DiameterSuccess,
		OriginHost:  rfc6733.OriginHost(Host),
		OriginRealm: rfc6733.OriginRealm(Realm),
		//		*ErrorMessage
		//		*FailedAVP
	}
	if c.peer.Host != msg.DiameterIdentity(r.OriginHost) {
		dpa.ResultCode = rfc6733.DiameterUnknownPeer
	}
	if c.peer.Realm != msg.DiameterIdentity(r.OriginRealm) {
		dpa.ResultCode = rfc6733.DiameterUnknownPeer
	}
	return dpa
}

// HandleDPA is DPA handler function
var HandleDPA = func(r rfc6733.DPA, c *Conn) {
}
