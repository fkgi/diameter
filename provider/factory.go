package provider

import (
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
func (c *Connection) makeCER() (m msg.Message) {
	m = msg.Message{}
	m.Ver = msg.DiaVer
	m.FlgR = true
	m.FlgP = false
	m.FlgE = false
	m.FlgT = false
	m.Code = uint32(257)
	m.AppID = uint32(0)

	m.EtEID = c.Local.NextEtE()

	var avps []msg.Avp
	avps = append(avps, msg.OriginHost(c.Local.Host))
	avps = append(avps, msg.OriginRealm(c.Local.Realm))
	for _, ip := range c.Local.Addr {
		avps = append(avps, msg.HostIPAddress(ip))
	}
	avps = append(avps, msg.VendorID(VendorID))
	avps = append(avps, msg.ProductName(ProductName))

	for _, app := range c.Peer.SupportedApps {
		if app[0] != 0 {
			avps = append(avps, msg.SupportedVendorID(app[0]))
			avps = append(avps, msg.VendorSpecificApplicationID(app[0], app[1], true))
		}
		avps = append(avps, msg.AuthApplicationID(app[1]))
	}

	avps = append(avps, msg.FirmwareRevision(FirmwareRevision))

	m.Encode(avps)

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
func (c *Connection) makeCEA(r msg.Message) (m msg.Message, i int) {
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
	avps = append(avps, msg.ResultCode(uint32(2001)))
	i = 2001
	avps = append(avps, msg.OriginHost(c.Local.Host))
	avps = append(avps, msg.OriginRealm(c.Local.Realm))
	for _, ip := range c.Local.Addr {
		avps = append(avps, msg.HostIPAddress(ip))
	}
	avps = append(avps, msg.VendorID(VendorID))
	avps = append(avps, msg.ProductName(ProductName))

	for _, app := range c.Peer.SupportedApps {
		if app[0] != 0 {
			avps = append(avps, msg.SupportedVendorID(app[0]))
			avps = append(avps, msg.VendorSpecificApplicationID(app[0], app[1], true))
		}
		avps = append(avps, msg.AuthApplicationID(app[1]))
	}

	avps = append(avps, msg.FirmwareRevision(FirmwareRevision))

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
func (c *Connection) makeDPR(i msg.Enumerated) (m msg.Message) {
	m = msg.Message{}
	m.Ver = msg.DiaVer
	m.FlgR = true
	m.FlgP = false
	m.FlgE = false
	m.FlgT = false
	m.Code = uint32(282)
	m.AppID = uint32(0)

	m.EtEID = c.Local.NextEtE()

	var avps []msg.Avp
	avps = append(avps, msg.OriginHost(c.Local.Host))
	avps = append(avps, msg.OriginRealm(c.Local.Realm))
	avps = append(avps, msg.DisconnectCause(i))

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
func (c *Connection) makeDPA(r msg.Message) (m msg.Message, i int) {
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
	avps = append(avps, msg.ResultCode(uint32(2001)))
	i = 2001
	avps = append(avps, msg.OriginHost(c.Local.Host))
	avps = append(avps, msg.OriginRealm(c.Local.Realm))

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
func (c *Connection) makeDWR() (m msg.Message) {
	m = msg.Message{}
	m.Ver = msg.DiaVer
	m.FlgR = true
	m.FlgP = false
	m.FlgE = false
	m.FlgT = false
	m.Code = uint32(280)
	m.AppID = uint32(0)

	m.EtEID = c.Local.NextEtE()

	var avps []msg.Avp
	avps = append(avps, msg.OriginHost(c.Local.Host))
	avps = append(avps, msg.OriginRealm(c.Local.Realm))

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
func (c *Connection) makeDWA(r msg.Message) (m msg.Message, i int) {
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
	avps = append(avps, msg.ResultCode(uint32(2001)))
	i = 2001
	avps = append(avps, msg.OriginHost(c.Local.Host))
	avps = append(avps, msg.OriginRealm(c.Local.Realm))

	m.Encode(avps)

	return
}
