package sock

import (
	"net"
	"strings"

	"github.com/fkgi/diameter/msg"
)

// HandleMSG is diameter request handler
var HandleMSG = func(r msg.Message, c *Conn) {
	go func() {
		a := c.peer.Handler(r)
		a.HbHID = r.HbHID
		a.EtEID = r.EtEID
		c.notify <- eventSndMsg{a}
	}()
}

// MakeCER returns new CER
var MakeCER = func(c *Conn) msg.CER {
	r := msg.CER{
		OriginHost:  msg.OriginHost(c.local.Host),
		OriginRealm: msg.OriginRealm(c.local.Realm),
		//HostIPAddress: make([]msg.HostIPAddress, 0),
		VendorID:    VendorID,
		ProductName: ProductName,
		// *OriginStateID:
		SupportedVendorID:           make([]msg.SupportedVendorID, 0),
		ApplicationID:               make([]msg.ApplicationID, 0),
		VendorSpecificApplicationID: make([]msg.VendorSpecificApplicationID, 0),
		FirmwareRevision:            &FirmwareRevision}

	switch c.con.LocalAddr().Network() {
	case "tcp":
		s := c.con.LocalAddr().String()
		s = s[:strings.LastIndex(s, ":")]
		r.HostIPAddress = []msg.HostIPAddress{msg.HostIPAddress(net.ParseIP(s))}
	case "sctp":
		s := c.con.LocalAddr().String()
		s = s[:strings.LastIndex(s, ":")]
		r.HostIPAddress = []msg.HostIPAddress{}
		for _, i := range strings.Split(s, "/") {
			r.HostIPAddress = append(r.HostIPAddress, msg.HostIPAddress(net.ParseIP(i)))
		}
	}
	if c.local.StateID != 0 {
		r.OriginStateID = &c.local.StateID
	}
	for v, a := range c.local.AuthApps {
		if v != 0 {
			r.SupportedVendorID = append(
				r.SupportedVendorID, msg.SupportedVendorID(v))
			for _, i := range a {
				r.VendorSpecificApplicationID = append(
					r.VendorSpecificApplicationID,
					msg.VendorSpecificApplicationID{
						VendorID: v,
						App:      i})
			}
		} else {
			for _, i := range a {
				r.ApplicationID = append(r.ApplicationID, i)
			}
		}
	}
	return r
}

// HandleCER is CER handler function
var HandleCER = func(r msg.CER, c *Conn) msg.CEA {
	cea := msg.CEA{
		ResultCode:                  msg.DiameterSuccess,
		OriginHost:                  msg.OriginHost(c.local.Host),
		OriginRealm:                 msg.OriginRealm(c.local.Realm),
		VendorID:                    VendorID,
		ProductName:                 ProductName,
		SupportedVendorID:           make([]msg.SupportedVendorID, 0),
		ApplicationID:               make([]msg.ApplicationID, 0),
		VendorSpecificApplicationID: make([]msg.VendorSpecificApplicationID, 0),
		FirmwareRevision:            &FirmwareRevision}
	switch c.con.LocalAddr().Network() {
	case "tcp":
		s := c.con.LocalAddr().String()
		s = s[:strings.LastIndex(s, ":")]
		cea.HostIPAddress = []msg.HostIPAddress{msg.HostIPAddress(net.ParseIP(s))}
	case "sctp":
		s := c.con.LocalAddr().String()
		s = s[:strings.LastIndex(s, ":")]
		cea.HostIPAddress = []msg.HostIPAddress{}
		for _, i := range strings.Split(s, "/") {
			cea.HostIPAddress = append(cea.HostIPAddress, msg.HostIPAddress(net.ParseIP(i)))
		}
	}
	if c.local.StateID != 0 {
		cea.OriginStateID = &c.local.StateID
	}

	if c.peer == nil {
		c.peer = &Peer{
			Host:  msg.DiameterIdentity(r.OriginHost),
			Realm: msg.DiameterIdentity(r.OriginRealm)}
	} else if msg.DiameterIdentity(r.OriginHost) != c.peer.Host ||
		msg.DiameterIdentity(r.OriginRealm) != c.peer.Realm {
		cea.ResultCode = msg.DiameterUnknownPeer
	}
	if cea.ResultCode == msg.DiameterSuccess {
		app := map[msg.VendorID][]msg.ApplicationID{}
		for _, i := range r.SupportedVendorID {
			app[msg.VendorID(i)] = []msg.ApplicationID{}
		}
		for _, a := range r.VendorSpecificApplicationID {
			if _, ok := app[a.VendorID]; !ok {
				app[a.VendorID] = []msg.ApplicationID{}
			}
			app[a.VendorID] = append(app[a.VendorID], a.App)
		}
		if len(r.ApplicationID) != 0 {
			if _, ok := app[0]; !ok {
				app[0] = []msg.ApplicationID{}
			}
			for _, i := range r.ApplicationID {
				app[0] = append(app[0], i)
			}
		}

		if c.peer.AuthApps == nil {
			relay := msg.AuthApplicationID(0xffffffff)
			for _, id := range c.local.AuthApps[0] {
				if relay.Equals(id) {
					c.peer.AuthApps = app
					break
				}
			}
			if c.peer.AuthApps == nil {
				c.peer.AuthApps = map[msg.VendorID][]msg.ApplicationID{}
				for key, ids := range app {
					for _, rid := range ids {
						if _, ok := c.local.AuthApps[key]; !ok {
							continue
						}
						for _, lid := range c.local.AuthApps[key] {
							if rid.Equals(lid) {
								if _, ok := c.peer.AuthApps[key]; !ok {
									c.peer.AuthApps[key] = []msg.ApplicationID{}
								}
								c.peer.AuthApps[key] = append(c.peer.AuthApps[key], rid)
							}
						}
					}
				}
				if len(c.peer.AuthApps) == 0 {
					cea.ResultCode = msg.DiameterApplicationUnsupported
				}
			}
		} else {
			a := map[msg.VendorID][]msg.ApplicationID{}
			for key, ids := range app {
				for _, rid := range ids {
					if _, ok := c.peer.AuthApps[key]; !ok {
						continue
					}
					for _, lid := range c.peer.AuthApps[key] {
						if rid.Equals(lid) {
							if _, ok := a[key]; !ok {
								a[key] = []msg.ApplicationID{}
							}
							a[key] = append(a[key], rid)
						}
					}
				}
			}
			if len(a) == 0 {
				cea.ResultCode = msg.DiameterApplicationUnsupported
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
				cea.SupportedVendorID, msg.SupportedVendorID(v))
			for _, i := range a {
				cea.VendorSpecificApplicationID = append(
					cea.VendorSpecificApplicationID,
					msg.VendorSpecificApplicationID{
						VendorID: v,
						App:      i})
			}
		} else {
			for _, i := range a {
				cea.ApplicationID = append(cea.ApplicationID, i)
			}
		}
	}
	return cea
}

// HandleCEA is CEA handler function
var HandleCEA = func(m msg.CEA, c *Conn) {
	app := map[msg.VendorID][]msg.ApplicationID{}
	for _, i := range m.SupportedVendorID {
		app[msg.VendorID(i)] = []msg.ApplicationID{}
	}
	for _, a := range m.VendorSpecificApplicationID {
		if _, ok := app[a.VendorID]; !ok {
			app[a.VendorID] = []msg.ApplicationID{}
		}
		app[a.VendorID] = append(app[a.VendorID], a.App)
	}
	if len(m.ApplicationID) != 0 {
		if _, ok := app[0]; !ok {
			app[0] = []msg.ApplicationID{}
		}
		for _, i := range m.ApplicationID {
			app[0] = append(app[0], i)
		}
	}
	c.peer.AuthApps = app
}

// MakeDWR returns new DWR
var MakeDWR = func(c *Conn) msg.DWR {
	dwr := msg.DWR{
		OriginHost:  msg.OriginHost(c.local.Host),
		OriginRealm: msg.OriginRealm(c.local.Realm)}
	if c.local.StateID != 0 {
		dwr.OriginStateID = &c.local.StateID
	}
	return dwr
}

// HandleDWR is DWR handler function
var HandleDWR = func(r msg.DWR, c *Conn) msg.DWA {
	dwa := msg.DWA{
		ResultCode:  msg.DiameterSuccess,
		OriginHost:  msg.OriginHost(c.local.Host),
		OriginRealm: msg.OriginRealm(c.local.Realm),
		//		*ErrorMessage
		//		*FailedAVP
		//		*OriginStateID
	}
	if c.peer.Host != msg.DiameterIdentity(r.OriginHost) {
		dwa.ResultCode = msg.DiameterUnknownPeer
	}
	if c.peer.Realm != msg.DiameterIdentity(r.OriginRealm) {
		dwa.ResultCode = msg.DiameterUnknownPeer
	}
	if c.local.StateID != 0 {
		dwa.OriginStateID = &c.local.StateID
	}

	return dwa
}

// HandleDWA is DWA handler function
var HandleDWA = func(r msg.DWA, c *Conn) {
}

// MakeDPR returns new DWR
var MakeDPR = func(c *Conn) msg.DPR {
	r := msg.DPR{
		OriginHost:      msg.OriginHost(c.local.Host),
		OriginRealm:     msg.OriginRealm(c.local.Realm),
		DisconnectCause: msg.DisconnectCause(msg.Rebooting)}
	return r
}

// HandleDPR is DPR handler function
var HandleDPR = func(r msg.DPR, c *Conn) msg.DPA {
	dpa := msg.DPA{
		ResultCode:  msg.DiameterSuccess,
		OriginHost:  msg.OriginHost(c.local.Host),
		OriginRealm: msg.OriginRealm(c.local.Realm),
		//		*ErrorMessage
		//		*FailedAVP
	}
	if c.peer.Host != msg.DiameterIdentity(r.OriginHost) {
		dpa.ResultCode = msg.DiameterUnknownPeer
	}
	if c.peer.Realm != msg.DiameterIdentity(r.OriginRealm) {
		dpa.ResultCode = msg.DiameterUnknownPeer
	}
	return dpa
}

// HandleDPA is DPA handler function
var HandleDPA = func(r msg.DPA, c *Conn) {
}
