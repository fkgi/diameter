package connection

import (
	"net"
	"strings"

	"github.com/fkgi/diameter/msg"
)

/*
 <CER> ::= < Diameter Header: 257, REQ >
		   { Origin-Host }
		   { Origin-Realm }
		1* { Host-IP-Address }
		   { Vendor-Id }
		   { Product-Name }
		   [ Origin-State-Id ]
		 * [ Supported-Vendor-Id ]
		 * [ Auth-Application-Id ]
		 * [ Inband-Security-Id ]   // not supported (not recommended)
		 * [ Acct-Application-Id ]  // not supported
		 * [ Vendor-Specific-Application-Id ] // only support auth
		   [ Firmware-Revision ]
		 * [ AVP ]
*/
func (c *Connection) makeCER() (m msg.Message) {
	m = msg.Message{
		Ver:  msg.DiaVer,
		FlgR: true, FlgP: false, FlgE: false, FlgT: false,
		Code: uint32(257), AppID: uint32(0)}

	m.EtEID = c.local.NextEtE()

	avps := []msg.Avp{
		msg.OriginHost(c.local.Host).Encode(),
		msg.OriginRealm(c.local.Realm).Encode()}
	for _, ip := range getIP(c.con) {
		avps = append(avps, msg.HostIPAddress(ip).Encode())
	}
	avps = append(avps, VendorID.Encode())
	avps = append(avps, ProductName.Encode())
	if c.local.StateID != 0 {
		avps = append(avps, msg.OriginStateID(c.local.StateID).Encode())
	}
	for ven, app := range c.local.AuthApps {
		if ven != 0 {
			avps = append(avps, msg.SupportedVendorID(ven).Encode())
			for _, id := range app {
				avps = append(avps, msg.VendorSpecificApplicationID{
					VendorID: ven,
					App:      id}.Encode())
			}
		} else {
			for _, id := range app {
				avps = append(avps, id.Encode())
			}
		}
	}
	avps = append(avps, FirmwareRevision.Encode())

	m.Encode(avps)
	return
}

func getIP(c net.Conn) (ip []net.IP) {
	addr := c.LocalAddr()
	switch addr.Network() {
	case "tcp":
		s := addr.String()
		ip = append(ip, net.ParseIP(s[:strings.LastIndex(s, ":")]))
	case "sctp":
		s := addr.String()
		ips := s[:strings.LastIndex(s, ":")]
		for _, i := range strings.Split(ips, "/") {
			ip = append(ip, net.ParseIP(i))
		}
	}
	return
}

/*
 <CEA> ::= < Diameter Header: 257 >
		   { Result-Code }
		   { Origin-Host }
		   { Origin-Realm }
		1* { Host-IP-Address }
		   { Vendor-Id }
		   { Product-Name }
		   [ Origin-State-Id ]
		   [ Error-Message ]
		   [ Failed-AVP ]
		 * [ Supported-Vendor-Id ]
		 * [ Auth-Application-Id ]
		 * [ Inband-Security-Id ]   // not supported (not recommended)
		 * [ Acct-Application-Id ]  // not supported
		 * [ Vendor-Specific-Application-Id ] // only support auth
		   [ Firmware-Revision ]
		 * [ AVP ]
*/
func (c *Connection) makeCEA(
	r msg.Message, p *PeerNode) (a msg.Message, i msg.ResultCode) {

	i = msg.DiameterSuccess
	if c.peer == nil {
		c.peer = &PeerNode{
			Host:  p.Host,
			Realm: p.Realm}
	} else if p.Host != c.peer.Host || p.Realm != c.peer.Realm {
		i = msg.DiameterUnknownPeer
	}
	if i == msg.DiameterSuccess && c.peer.AuthApps == nil {
		i = msg.DiameterApplicationUnsupported
		relay := msg.AuthApplicationID(0xffffffff)
		for _, id := range c.local.AuthApps[0] {
			if relay.Equals(id) {
				c.peer.AuthApps = p.AuthApps
				i = msg.DiameterSuccess
				relay = msg.AuthApplicationID(0)
				break
			}
		}
		if i != msg.DiameterSuccess {
			c.peer.AuthApps = map[msg.VendorID][]msg.ApplicationID{}
			for key, ids := range p.AuthApps {
				for _, rid := range ids {
					if _, ok := c.local.AuthApps[key]; !ok {
						continue
					}
					for _, lid := range c.local.AuthApps[key] {
						if rid.Equals(lid) {
							app, ok := c.peer.AuthApps[key]
							if !ok {
								c.peer.AuthApps[key] = []msg.ApplicationID{}
							}
							c.peer.AuthApps[key] = append(app, rid)
							i = msg.DiameterSuccess
						}
					}
				}
			}
		}
	} else if i == msg.DiameterSuccess {
		apps := map[msg.VendorID][]msg.ApplicationID{}
		for key, ids := range p.AuthApps {
			for _, rid := range ids {
				if _, ok := c.peer.AuthApps[key]; !ok {
					continue
				}
				for _, lid := range c.peer.AuthApps[key] {
					if rid.Equals(lid) {
						if _, ok := apps[key]; !ok {
							apps[key] = []msg.ApplicationID{}
						}
						apps[key] = append(apps[key], rid)
						i = msg.DiameterSuccess
					}
				}
			}
		}
		if len(apps) == 0 {
			i = msg.DiameterApplicationUnsupported
		} else {
			c.peer.AuthApps = apps
		}
	}

	if c.peer.Tw == 0 {
		c.peer.Tw = c.local.Tw
	}
	if c.peer.Ew == 0 {
		c.peer.Ew = c.local.Ew
	}
	if c.peer.Ts == 0 {
		c.peer.Ts = c.local.Ts
	}
	if c.peer.Tp == 0 {
		c.peer.Tp = c.local.Tp
	}
	if c.peer.Cp == 0 {
		c.peer.Cp = c.local.Cp
	}

	a = msg.Message{
		Ver:  msg.DiaVer,
		FlgR: false, FlgP: r.FlgP, FlgE: false, FlgT: false,
		HbHID: r.HbHID, EtEID: r.EtEID,
		Code: r.Code, AppID: r.AppID}

	if i != msg.DiameterSuccess {
		a.FlgE = true
	}
	var avps []msg.Avp
	avps = append(avps, i.Encode())

	avps = append(avps, msg.OriginHost(c.local.Host).Encode())
	avps = append(avps, msg.OriginRealm(c.local.Realm).Encode())
	for _, ip := range getIP(c.con) {
		avps = append(avps, msg.HostIPAddress(ip).Encode())
	}
	avps = append(avps, VendorID.Encode())
	avps = append(avps, ProductName.Encode())
	if c.local.StateID != 0 {
		avps = append(avps, msg.OriginStateID(c.local.StateID).Encode())
	}

	for ven, app := range c.peer.AuthApps {
		if ven != 0 {
			avps = append(avps, msg.SupportedVendorID(ven).Encode())
			for _, id := range app {
				avps = append(avps, msg.VendorSpecificApplicationID{
					VendorID: ven,
					App:      id}.Encode())
			}
		} else {
			for _, id := range app {
				avps = append(avps, id.Encode())
			}
		}
	}

	avps = append(avps, FirmwareRevision.Encode())

	a.Encode(avps)
	return
}

/*
 <DPR>  ::= < Diameter Header: 282, REQ >
			{ Origin-Host }
			{ Origin-Realm }
			{ Disconnect-Cause }
		  * [ AVP ]
*/
func (c *Connection) makeDPR(i msg.Enumerated) (r msg.Message) {
	r = msg.Message{
		Ver:  msg.DiaVer,
		FlgR: true, FlgP: false, FlgE: false, FlgT: false,
		Code: uint32(282), AppID: uint32(0)}

	r.EtEID = c.local.NextEtE()

	avps := []msg.Avp{
		msg.OriginHost(c.local.Host).Encode(),
		msg.OriginRealm(c.local.Realm).Encode(),
		msg.DisconnectCause(i).Encode()}

	r.Encode(avps)

	return
}

/*
 <DPA>  ::= < Diameter Header: 282 >
			{ Result-Code }
			{ Origin-Host }
			{ Origin-Realm }
			[ Error-Message ]
			[ Failed-AVP ]
		  * [ AVP ]
*/
func (c *Connection) makeDPA(r msg.Message) (a msg.Message, i msg.ResultCode) {
	i = msg.DiameterSuccess
	if avp, e := r.Decode(); e != nil {
		i = msg.DiameterInvalidAvpValue
	} else {
		if t, ok := msg.GetOriginHost(avp); !ok {
			i = msg.DiameterInvalidAvpValue
		} else if c.peer.Host != msg.DiameterIdentity(t) {
			i = msg.DiameterUnknownPeer
		}
		if t, ok := msg.GetOriginRealm(avp); !ok {
			i = msg.DiameterInvalidAvpValue
		} else if c.peer.Realm != msg.DiameterIdentity(t) {
			i = msg.DiameterUnknownPeer
		}
	}

	a = msg.Message{
		Ver:  msg.DiaVer,
		FlgR: false, FlgP: r.FlgP, FlgE: false, FlgT: false,
		HbHID: r.HbHID, EtEID: r.EtEID,
		Code: r.Code, AppID: r.AppID}

	if i != msg.DiameterSuccess {
		a.FlgE = true
	}
	avps := []msg.Avp{
		i.Encode(),
		msg.OriginHost(c.local.Host).Encode(),
		msg.OriginRealm(c.local.Realm).Encode()}

	a.Encode(avps)

	return
}

/*
 <DWR>  ::= < Diameter Header: 280, REQ >
			{ Origin-Host }
			{ Origin-Realm }
			[ Origin-State-Id ]
		  * [ AVP ]
*/
func (c *Connection) makeDWR() (a msg.Message) {
	a = msg.Message{
		Ver:  msg.DiaVer,
		FlgR: true, FlgP: false, FlgE: false, FlgT: false,
		Code: uint32(280), AppID: uint32(0)}

	a.EtEID = c.local.NextEtE()

	avps := []msg.Avp{
		msg.OriginHost(c.local.Host).Encode(),
		msg.OriginRealm(c.local.Realm).Encode()}
	if c.local.StateID != 0 {
		avps = append(avps, msg.OriginStateID(c.local.StateID).Encode())
	}

	a.Encode(avps)

	return
}

/*
 <DWA>  ::= < Diameter Header: 280 >
			{ Result-Code }
			{ Origin-Host }
			{ Origin-Realm }
			[ Error-Message ]
			[ Failed-AVP ]
			[ Origin-State-Id ]
		  * [ AVP ]
*/
func (c *Connection) makeDWA(r msg.Message) (a msg.Message, i msg.ResultCode) {
	i = msg.DiameterSuccess
	if avp, e := r.Decode(); e != nil {
		i = msg.DiameterInvalidAvpValue
	} else {
		if t, ok := msg.GetOriginHost(avp); !ok {
			i = msg.DiameterInvalidAvpValue
		} else if c.peer.Host != msg.DiameterIdentity(t) {
			i = msg.DiameterUnknownPeer
		}
		if t, ok := msg.GetOriginRealm(avp); !ok {
			i = msg.DiameterInvalidAvpValue
		} else if c.peer.Realm != msg.DiameterIdentity(t) {
			i = msg.DiameterUnknownPeer
		}
	}

	a = msg.Message{
		Ver:  msg.DiaVer,
		FlgR: false, FlgP: r.FlgP, FlgE: false, FlgT: false,
		HbHID: r.HbHID, EtEID: r.EtEID,
		Code: r.Code, AppID: r.AppID}

	if i != msg.DiameterSuccess {
		a.FlgE = true
	}
	avps := []msg.Avp{
		i.Encode(),
		msg.OriginHost(c.local.Host).Encode(),
		msg.OriginRealm(c.local.Realm).Encode()}
	if c.local.StateID != 0 {
		avps = append(avps, msg.OriginStateID(c.local.StateID).Encode())
	}

	a.Encode(avps)

	return
}

func (c *Connection) makeUnableToDeliver(r msg.Message) (a msg.Message) {
	a = msg.Message{
		Ver:  msg.DiaVer,
		FlgR: false, FlgP: r.FlgP, FlgE: true, FlgT: false,
		HbHID: r.HbHID, EtEID: r.EtEID,
		Code: r.Code, AppID: r.AppID}

	avps := []msg.Avp{
		msg.DiameterUnableToDeliver.Encode(),
		msg.OriginHost(c.local.Host).Encode(),
		msg.OriginRealm(c.local.Realm).Encode()}
	a.Encode(avps)

	return
}
