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
func (p *Connection) makeCER() (m msg.Message) {
	m = msg.Message{
		Ver:  msg.DiaVer,
		FlgR: true, FlgP: false, FlgE: false, FlgT: false,
		Code: uint32(257), AppID: uint32(0)}

	m.EtEID = p.local.NextEtE()
	var avps []msg.Avp
	avps = append(avps, msg.OriginHost(p.local.Host).Encode())
	avps = append(avps, msg.OriginRealm(p.local.Realm).Encode())
	for _, ip := range getIP(p.con) {
		avps = append(avps, msg.HostIPAddress(ip).Encode())
	}
	avps = append(avps, VendorID.Encode())
	avps = append(avps, ProductName.Encode())
	if p.local.StateID != 0 {
		avps = append(avps, msg.OriginStateID(p.local.StateID).Encode())
	}
	for ven, app := range p.local.AuthApps {
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
func (p *Connection) makeCEA(
	r msg.Message, peer *PeerNode) (m msg.Message, i msg.ResultCode) {

	i = msg.DiameterSuccess
	if p.peer == nil {
		p.peer = &PeerNode{
			Host:  peer.Host,
			Realm: peer.Realm}
	} else if peer.Host != p.peer.Host || peer.Realm != p.peer.Realm {
		i = msg.DiameterUnknownPeer
	}
	if i == msg.DiameterSuccess {
		relay := msg.AuthApplicationID(0xffffffff)
		for _, id := range p.local.AuthApps[0] {
			if relay.Equals(id) {
				relay = msg.AuthApplicationID(0)
				break
			}
		}
		if relay != 0 {

		}
	}

	m = msg.Message{
		Ver:  msg.DiaVer,
		FlgR: false, FlgP: r.FlgP, FlgE: false, FlgT: false,
		HbHID: r.HbHID, EtEID: r.EtEID,
		Code: r.Code, AppID: r.AppID}

	var avps []msg.Avp
	avps = append(avps, msg.DiameterSuccess.Encode())
	i = msg.DiameterSuccess

	avps = append(avps, msg.OriginHost(p.local.Host).Encode())
	avps = append(avps, msg.OriginRealm(p.local.Realm).Encode())
	for _, ip := range getIP(p.con) {
		avps = append(avps, msg.HostIPAddress(ip).Encode())
	}
	avps = append(avps, VendorID.Encode())
	avps = append(avps, ProductName.Encode())

	for ven, app := range p.local.AuthApps {
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

/*
 <DPR>  ::= < Diameter Header: 282, REQ >
			{ Origin-Host }
			{ Origin-Realm }
			{ Disconnect-Cause }
		  * [ AVP ]
*/
func (p *Connection) makeDPR(i msg.Enumerated) (m msg.Message) {
	m = msg.Message{
		Ver:  msg.DiaVer,
		FlgR: true, FlgP: false, FlgE: false, FlgT: false,
		Code: uint32(282), AppID: uint32(0)}

	m.EtEID = p.local.NextEtE()

	var avps []msg.Avp
	avps = append(avps, msg.OriginHost(p.local.Host).Encode())
	avps = append(avps, msg.OriginRealm(p.local.Realm).Encode())
	avps = append(avps, msg.DisconnectCause(i).Encode())

	m.Encode(avps)

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
func (p *Connection) makeDPA(r msg.Message) (m msg.Message, i msg.ResultCode) {
	m = msg.Message{
		Ver:  msg.DiaVer,
		FlgR: false, FlgP: r.FlgP, FlgE: false, FlgT: false,
		HbHID: r.HbHID, EtEID: r.EtEID,
		Code: r.Code, AppID: r.AppID}

	var avps []msg.Avp
	avps = append(avps, msg.DiameterSuccess.Encode())
	i = msg.DiameterSuccess
	avps = append(avps, msg.OriginHost(p.local.Host).Encode())
	avps = append(avps, msg.OriginRealm(p.local.Realm).Encode())

	m.Encode(avps)

	return
}

/*
 <DWR>  ::= < Diameter Header: 280, REQ >
			{ Origin-Host }
			{ Origin-Realm }
			[ Origin-State-Id ]
		  * [ AVP ]
*/
func (p *Connection) makeDWR() (m msg.Message) {
	m = msg.Message{
		Ver:  msg.DiaVer,
		FlgR: true, FlgP: false, FlgE: false, FlgT: false,
		Code: uint32(280), AppID: uint32(0)}

	m.EtEID = p.local.NextEtE()

	var avps []msg.Avp
	avps = append(avps, msg.OriginHost(p.local.Host).Encode())
	avps = append(avps, msg.OriginRealm(p.local.Realm).Encode())

	m.Encode(avps)

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
func (p *Connection) makeDWA(r msg.Message) (m msg.Message, i msg.ResultCode) {
	m = msg.Message{
		Ver:  msg.DiaVer,
		FlgR: false, FlgP: r.FlgP, FlgE: false, FlgT: false,
		HbHID: r.HbHID, EtEID: r.EtEID,
		Code: r.Code, AppID: r.AppID}

	var avps []msg.Avp
	avps = append(avps, msg.DiameterSuccess.Encode())
	i = msg.DiameterSuccess
	avps = append(avps, msg.OriginHost(p.local.Host).Encode())
	avps = append(avps, msg.OriginRealm(p.local.Realm).Encode())

	m.Encode(avps)

	return
}

func (p *Connection) makeUnableToDeliver(r msg.Message) (m msg.Message) {
	m = msg.Message{
		Ver:  msg.DiaVer,
		FlgR: false, FlgP: r.FlgP, FlgE: true, FlgT: false,
		HbHID: r.HbHID, EtEID: r.EtEID,
		Code: r.Code, AppID: r.AppID}

	var avps []msg.Avp
	avps = append(avps, msg.DiameterUnableToDeliver.Encode())
	avps = append(avps, msg.OriginHost(p.local.Host).Encode())
	avps = append(avps, msg.OriginRealm(p.local.Realm).Encode())

	m.Encode(avps)

	return
}
