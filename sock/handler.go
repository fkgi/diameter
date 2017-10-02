package sock

import (
	"net"
	"strings"

	"github.com/fkgi/diameter/msg"
)

// HandleMSG is diameter request handler
var HandleMSG = func(m msg.Request) msg.Answer {
	return nil
}

// MakeUnsupportedAnswer make answer for unsupported message
var MakeUnsupportedAnswer = func(m msg.RawMsg) msg.RawMsg {
	a := msg.RawMsg{}
	a.Ver = m.Ver
	a.FlgP = m.FlgP
	a.Code = m.Code
	a.AppID = m.AppID
	a.HbHID = m.HbHID
	a.EtEID = m.EtEID

	host := msg.OriginHost(Host)
	realm := msg.OriginRealm(Realm)
	result := msg.DiameterApplicationUnsupported
	a.AVP = []msg.RawAVP{
		result.ToRaw(),
		host.ToRaw(),
		realm.ToRaw()}

	return a
}

// MakeCER returns new CER
var MakeCER = func(c *Conn) msg.CER {
	r := msg.CER{
		OriginHost:  msg.OriginHost(Host),
		OriginRealm: msg.OriginRealm(Realm),
		//HostIPAddress: make([]msg.HostIPAddress, 0),
		VendorID:                    VendorID,
		ProductName:                 ProductName,
		OriginStateID:               StateID,
		SupportedVendorID:           make([]msg.SupportedVendorID, 0),
		AuthApplicationID:           make([]msg.AuthApplicationID, 0),
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
			r.HostIPAddress = append(r.HostIPAddress,
				msg.HostIPAddress(net.ParseIP(i)))
		}
	}
	for v, a := range getSupportedApps() {
		if v != 0 {
			r.SupportedVendorID = append(
				r.SupportedVendorID, msg.SupportedVendorID(v))
			for _, i := range a {
				r.VendorSpecificApplicationID = append(
					r.VendorSpecificApplicationID,
					msg.VendorSpecificApplicationID{
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
var HandleCER = func(r msg.CER, c *Conn) msg.CEA {
	cea := msg.CEA{
		ResultCode:                  msg.DiameterSuccess,
		OriginHost:                  msg.OriginHost(Host),
		OriginRealm:                 msg.OriginRealm(Realm),
		VendorID:                    VendorID,
		ProductName:                 ProductName,
		OriginStateID:               StateID,
		SupportedVendorID:           make([]msg.SupportedVendorID, 0),
		AuthApplicationID:           make([]msg.AuthApplicationID, 0),
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

	if c.peer == nil {
		c.peer = &Peer{
			Host:  msg.DiameterIdentity(r.OriginHost),
			Realm: msg.DiameterIdentity(r.OriginRealm)}
	} else if msg.DiameterIdentity(r.OriginHost) != c.peer.Host ||
		msg.DiameterIdentity(r.OriginRealm) != c.peer.Realm {
		cea.ResultCode = msg.DiameterUnknownPeer
	}
	if cea.ResultCode == msg.DiameterSuccess {
		app := map[msg.VendorID][]msg.AuthApplicationID{}
		for _, i := range r.SupportedVendorID {
			app[msg.VendorID(i)] = []msg.AuthApplicationID{}
		}
		for _, a := range r.VendorSpecificApplicationID {
			if _, ok := app[a.VendorID]; !ok {
				app[a.VendorID] = []msg.AuthApplicationID{}
			}
			app[a.VendorID] = append(app[a.VendorID], a.AuthApplicationID)
		}
		if len(r.AuthApplicationID) != 0 {
			if _, ok := app[0]; !ok {
				app[0] = []msg.AuthApplicationID{}
			}
			for _, i := range r.AuthApplicationID {
				app[0] = append(app[0], i)
			}
		}

		if c.peer.AuthApps == nil {
			relay := msg.AuthApplicationID(0xffffffff)
			for _, id := range getSupportedApps()[0] {
				if relay == id {
					c.peer.AuthApps = app
					break
				}
			}
			if c.peer.AuthApps == nil {
				c.peer.AuthApps = map[msg.VendorID][]msg.AuthApplicationID{}
				for key, ids := range app {
					for _, rid := range ids {
						if _, ok := getSupportedApps()[key]; !ok {
							continue
						}
						for _, lid := range getSupportedApps()[key] {
							if rid == lid {
								if _, ok := c.peer.AuthApps[key]; !ok {
									c.peer.AuthApps[key] = []msg.AuthApplicationID{}
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
			a := map[msg.VendorID][]msg.AuthApplicationID{}
			for key, ids := range app {
				for _, rid := range ids {
					if _, ok := c.peer.AuthApps[key]; !ok {
						continue
					}
					for _, lid := range c.peer.AuthApps[key] {
						if rid == lid {
							if _, ok := a[key]; !ok {
								a[key] = []msg.AuthApplicationID{}
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
var HandleCEA = func(m msg.CEA, c *Conn) {
	app := map[msg.VendorID][]msg.AuthApplicationID{}
	for _, i := range m.SupportedVendorID {
		app[msg.VendorID(i)] = []msg.AuthApplicationID{}
	}
	for _, a := range m.VendorSpecificApplicationID {
		if _, ok := app[a.VendorID]; !ok {
			app[a.VendorID] = []msg.AuthApplicationID{}
		}
		app[a.VendorID] = append(app[a.VendorID], a.AuthApplicationID)
	}
	if len(m.AuthApplicationID) != 0 {
		if _, ok := app[0]; !ok {
			app[0] = []msg.AuthApplicationID{}
		}
		for _, i := range m.AuthApplicationID {
			app[0] = append(app[0], i)
		}
	}
	c.peer.AuthApps = app
}

// MakeDWR returns new DWR
var MakeDWR = func(c *Conn) msg.DWR {
	dwr := msg.DWR{
		OriginHost:    msg.OriginHost(Host),
		OriginRealm:   msg.OriginRealm(Realm),
		OriginStateID: StateID}
	return dwr
}

// HandleDWR is DWR handler function
var HandleDWR = func(r msg.DWR, c *Conn) msg.DWA {
	dwa := msg.DWA{
		ResultCode:  msg.DiameterSuccess,
		OriginHost:  msg.OriginHost(Host),
		OriginRealm: msg.OriginRealm(Realm),
		//		*ErrorMessage
		//		*FailedAVP
		OriginStateID: StateID}
	if c.peer.Host != msg.DiameterIdentity(r.OriginHost) {
		dwa.ResultCode = msg.DiameterUnknownPeer
	}
	if c.peer.Realm != msg.DiameterIdentity(r.OriginRealm) {
		dwa.ResultCode = msg.DiameterUnknownPeer
	}

	return dwa
}

// HandleDWA is DWA handler function
var HandleDWA = func(r msg.DWA, c *Conn) {
}

// MakeDPR returns new DWR
var MakeDPR = func(c *Conn) msg.DPR {
	r := msg.DPR{
		OriginHost:      msg.OriginHost(Host),
		OriginRealm:     msg.OriginRealm(Realm),
		DisconnectCause: msg.DisconnectCause(msg.Rebooting)}
	return r
}

// HandleDPR is DPR handler function
var HandleDPR = func(r msg.DPR, c *Conn) msg.DPA {
	dpa := msg.DPA{
		ResultCode:  msg.DiameterSuccess,
		OriginHost:  msg.OriginHost(Host),
		OriginRealm: msg.OriginRealm(Realm),
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
