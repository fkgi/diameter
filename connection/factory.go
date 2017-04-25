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
		 * [ Inband-Security-Id ]
		 * [ Acct-Application-Id ]
		 * [ Vendor-Specific-Application-Id ]
		   [ Firmware-Revision ]
		 * [ AVP ]
*/
func (p *Connection) makeCER(c net.Conn) (m msg.Message) {
	m = msg.Message{}
	m.Ver = msg.DiaVer
	m.FlgR = true
	m.FlgP = false
	m.FlgE = false
	m.FlgT = false
	m.Code = uint32(257)
	m.AppID = uint32(0)

	m.EtEID = p.local.NextEtE()
	var avps []msg.Avp
	avps = append(avps, msg.OriginHost(p.local.Host).Encode())
	avps = append(avps, msg.OriginRealm(p.local.Realm).Encode())
	for _, ip := range getIP(c) {
		avps = append(avps, msg.HostIPAddress(ip).Encode())
	}
	avps = append(avps, VendorID.Encode())
	avps = append(avps, ProductName.Encode())

	for _, app := range p.peer.Apps {
		if app.VendorID != 0 {
			avps = append(avps, msg.SupportedVendorID(app.VendorID).Encode())
			avps = append(avps, msg.VendorSpecificApplicationID{
				VendorID: app.VendorID,
				App:      app.AppID}.Encode())
		}
		avps = append(avps, msg.AuthApplicationID(app.AppID).Encode())
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
		 * [ Inband-Security-Id ]
		 * [ Acct-Application-Id ]
		 * [ Vendor-Specific-Application-Id ]
		   [ Firmware-Revision ]
		 * [ AVP ]
*/
func (p *Connection) makeCEA(r msg.Message, c net.Conn) (m msg.Message, i msg.ResultCode) {
	m = msg.Message{}
	m.Ver = msg.DiaVer
	m.FlgR = false
	m.FlgP = r.FlgP
	m.FlgE = false
	m.FlgT = false
	m.HbHID = r.HbHID
	m.EtEID = r.EtEID
	m.Code = r.Code
	m.AppID = r.AppID

	var avps []msg.Avp
	avps = append(avps, msg.DiameterSuccess.Encode())
	i = msg.DiameterSuccess
	avps = append(avps, msg.OriginHost(p.local.Host).Encode())
	avps = append(avps, msg.OriginRealm(p.local.Realm).Encode())
	for _, ip := range getIP(c) {
		avps = append(avps, msg.HostIPAddress(ip).Encode())
	}
	avps = append(avps, VendorID.Encode())
	avps = append(avps, ProductName.Encode())

	for _, app := range p.peer.Apps {
		if app.VendorID != 0 {
			avps = append(avps, msg.SupportedVendorID(app.VendorID).Encode())
			avps = append(avps, msg.VendorSpecificApplicationID{
				VendorID: app.VendorID,
				App:      app.AppID}.Encode())
		}
		avps = append(avps, msg.AuthApplicationID(app.AppID).Encode())
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
	m = msg.Message{}
	m.Ver = msg.DiaVer
	m.FlgR = true
	m.FlgP = false
	m.FlgE = false
	m.FlgT = false
	m.Code = uint32(282)
	m.AppID = uint32(0)

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
	m = msg.Message{}
	m.Ver = msg.DiaVer
	m.FlgR = false
	m.FlgP = r.FlgP
	m.FlgE = false
	m.FlgT = false
	m.HbHID = r.HbHID
	m.EtEID = r.EtEID
	m.Code = r.Code
	m.AppID = r.AppID

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
	m = msg.Message{}
	m.Ver = msg.DiaVer
	m.FlgR = true
	m.FlgP = false
	m.FlgE = false
	m.FlgT = false
	m.Code = uint32(280)
	m.AppID = uint32(0)

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
	m = msg.Message{}
	m.Ver = msg.DiaVer
	m.FlgR = false
	m.FlgP = r.FlgP
	m.FlgE = false
	m.FlgT = false
	m.HbHID = r.HbHID
	m.EtEID = r.EtEID
	m.Code = r.Code
	m.AppID = r.AppID

	var avps []msg.Avp
	avps = append(avps, msg.DiameterSuccess.Encode())
	i = msg.DiameterSuccess
	avps = append(avps, msg.OriginHost(p.local.Host).Encode())
	avps = append(avps, msg.OriginRealm(p.local.Realm).Encode())

	m.Encode(avps)

	return
}

func (p *Connection) makeUnableToDeliver(r msg.Message) (m msg.Message) {
	m = msg.Message{}
	m.Ver = msg.DiaVer
	m.FlgR = false
	m.FlgP = r.FlgP
	m.FlgE = true
	m.FlgT = false
	m.HbHID = r.HbHID
	m.EtEID = r.EtEID
	m.Code = r.Code
	m.AppID = r.AppID

	var avps []msg.Avp
	avps = append(avps, msg.DiameterUnableToDeliver.Encode())
	avps = append(avps, msg.OriginHost(p.local.Host).Encode())
	avps = append(avps, msg.OriginRealm(p.local.Realm).Encode())

	m.Encode(avps)

	return
}
