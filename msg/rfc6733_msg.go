package msg

/*
CER is Capabilities-Exchange-Request message
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
type CER struct {
	OriginHost
	OriginRealm
	HostIPAddress []HostIPAddress
	VendorID
	ProductName
	OriginStateID
	SupportedVendorID []SupportedVendorID
	AuthApplicationID []AuthApplicationID
	// []InbandSecurityID
	VendorSpecificApplicationID []VendorSpecificApplicationID
	*FirmwareRevision
}

// ToRaw return RawMsg struct of this value
func (v CER) ToRaw() RawMsg {
	m := RawMsg{
		Ver:  DiaVer,
		FlgR: true, FlgP: false, FlgE: false, FlgT: false,
		Code: 257, AppID: 0,
		AVP: make([]RawAVP, 0, 20)}

	m.AVP = append(m.AVP, v.OriginHost.ToRaw())
	m.AVP = append(m.AVP, v.OriginRealm.ToRaw())
	for _, ip := range v.HostIPAddress {
		m.AVP = append(m.AVP, ip.ToRaw())
	}
	m.AVP = append(m.AVP, v.VendorID.ToRaw())
	m.AVP = append(m.AVP, v.ProductName.ToRaw())
	if v.OriginStateID != 0 {
		m.AVP = append(m.AVP, v.OriginStateID.ToRaw())
	}
	for _, id := range v.SupportedVendorID {
		m.AVP = append(m.AVP, id.ToRaw())
	}
	for _, id := range v.AuthApplicationID {
		m.AVP = append(m.AVP, id.ToRaw())
	}
	for _, id := range v.VendorSpecificApplicationID {
		m.AVP = append(m.AVP, id.ToRaw())
	}
	if v.FirmwareRevision != nil {
		m.AVP = append(m.AVP, v.FirmwareRevision.ToRaw())
	}
	return m
}

// FromRaw make this value from RawMsg struct
func (CER) FromRaw(m RawMsg) (Request, error) {
	e := m.Validate(0, 257, true, false, false, false)
	if e != nil {
		return nil, e
	}

	v := CER{}
	v.HostIPAddress = make([]HostIPAddress, 0, 2)
	v.SupportedVendorID = make([]SupportedVendorID, 0, 5)
	v.AuthApplicationID = make([]AuthApplicationID, 0, 5)
	v.VendorSpecificApplicationID = make([]VendorSpecificApplicationID, 0, 5)

	for _, a := range m.AVP {
		if a.VenID == 0 && a.Code == 264 {
			e = v.OriginHost.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 296 {
			e = v.OriginRealm.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 257 {
			t := HostIPAddress{}
			e = t.FromRaw(a)
			v.HostIPAddress = append(v.HostIPAddress, t)
		} else if a.VenID == 0 && a.Code == 266 {
			e = v.VendorID.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 269 {
			e = v.ProductName.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 278 {
			e = v.OriginStateID.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 265 {
			t := SupportedVendorID(0)
			e = t.FromRaw(a)
			v.SupportedVendorID = append(v.SupportedVendorID, t)
		} else if a.VenID == 0 && a.Code == 258 {
			t := AuthApplicationID(0)
			e = t.FromRaw(a)
			v.AuthApplicationID = append(v.AuthApplicationID, t)
		} else if a.VenID == 0 && a.Code == 260 {
			t := VendorSpecificApplicationID{}
			e = t.FromRaw(a)
			v.VendorSpecificApplicationID = append(v.VendorSpecificApplicationID, t)
		} else if a.VenID == 0 && a.Code == 267 {
			e = v.FirmwareRevision.FromRaw(a)
		}

		if e != nil {
			return nil, e
		}
	}

	if len(v.OriginHost) == 0 ||
		len(v.OriginRealm) == 0 ||
		len(v.HostIPAddress) == 0 ||
		v.VendorID == 0 ||
		len(v.ProductName) == 0 {
		e = NoMandatoryAVP{}
	}
	return v, e
}

// TimeoutMsg make error message for timeout
func (v CER) TimeoutMsg() Answer {
	return CEA{
		ResultCode:                  DiameterUnableToDeliver,
		OriginHost:                  v.OriginHost,
		OriginRealm:                 v.OriginRealm,
		HostIPAddress:               v.HostIPAddress,
		VendorID:                    v.VendorID,
		ProductName:                 v.ProductName,
		OriginStateID:               v.OriginStateID,
		ErrorMessage:                "no response from peer node",
		SupportedVendorID:           v.SupportedVendorID,
		AuthApplicationID:           v.AuthApplicationID,
		VendorSpecificApplicationID: v.VendorSpecificApplicationID,
		FirmwareRevision:            v.FirmwareRevision}
}

/*
CEA is Capabilities-Exchange-Answer message
 <CEA> ::= < Diameter Header: 257 >
		   { Result-Code }
		   { Origin-Host }
		   { Origin-Realm }
		1* { Host-IP-Address }
		   { Vendor-Id }
		   { Product-Name }
		   [ Origin-State-Id ]
		   [ Error-RawMsg ]
		   [ Failed-AVP ]
		 * [ Supported-Vendor-Id ]
		 * [ Auth-Application-Id ]
		 * [ Inband-Security-Id ]   // not supported (not recommended)
		 * [ Acct-Application-Id ]  // not supported
		 * [ Vendor-Specific-Application-Id ] // only support auth
		   [ Firmware-Revision ]
		 * [ AVP ]
*/
type CEA struct {
	ResultCode
	OriginHost
	OriginRealm
	HostIPAddress []HostIPAddress
	VendorID
	ProductName
	OriginStateID
	ErrorMessage
	FailedAVP
	SupportedVendorID []SupportedVendorID
	AuthApplicationID []AuthApplicationID
	// []InbandSecurityID
	VendorSpecificApplicationID []VendorSpecificApplicationID
	*FirmwareRevision
}

// ToRaw return RawMsg struct of this value
func (v CEA) ToRaw() RawMsg {
	m := RawMsg{
		Ver:  DiaVer,
		FlgR: false, FlgP: false, FlgE: false, FlgT: false,
		Code: 257, AppID: 0,
		AVP: make([]RawAVP, 0, 20)}
	m.FlgE = v.ResultCode != DiameterSuccess

	m.AVP = append(m.AVP, v.ResultCode.ToRaw())
	m.AVP = append(m.AVP, v.OriginHost.ToRaw())
	m.AVP = append(m.AVP, v.OriginRealm.ToRaw())
	for _, ip := range v.HostIPAddress {
		m.AVP = append(m.AVP, ip.ToRaw())
	}
	m.AVP = append(m.AVP, v.VendorID.ToRaw())
	m.AVP = append(m.AVP, v.ProductName.ToRaw())
	if v.OriginStateID != 0 {
		m.AVP = append(m.AVP, v.OriginStateID.ToRaw())
	}
	if len(v.ErrorMessage) != 0 {
		m.AVP = append(m.AVP, v.ErrorMessage.ToRaw())
	}
	if len(v.FailedAVP) != 0 {
		m.AVP = append(m.AVP, v.FailedAVP.ToRaw())
	}
	for _, id := range v.SupportedVendorID {
		m.AVP = append(m.AVP, id.ToRaw())
	}
	for _, id := range v.AuthApplicationID {
		m.AVP = append(m.AVP, id.ToRaw())
	}
	for _, id := range v.VendorSpecificApplicationID {
		m.AVP = append(m.AVP, id.ToRaw())
	}
	if v.FirmwareRevision != nil {
		m.AVP = append(m.AVP, v.FirmwareRevision.ToRaw())
	}
	return m
}

// FromRaw make this value from RawMsg struct
func (CEA) FromRaw(m RawMsg) (Answer, error) {
	e := m.Validate(0, 257, false, false, false, false)
	if e != nil {
		return nil, e
	}

	v := CEA{}
	v.HostIPAddress = make([]HostIPAddress, 0, 2)
	v.SupportedVendorID = make([]SupportedVendorID, 0, 5)
	v.AuthApplicationID = make([]AuthApplicationID, 0, 5)
	v.VendorSpecificApplicationID = make([]VendorSpecificApplicationID, 0, 5)

	for _, a := range m.AVP {
		if a.VenID == 0 && a.Code == 268 {
			e = v.ResultCode.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 264 {
			e = v.OriginHost.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 296 {
			e = v.OriginRealm.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 257 {
			t := HostIPAddress{}
			e = t.FromRaw(a)
			v.HostIPAddress = append(v.HostIPAddress, t)
		} else if a.VenID == 0 && a.Code == 266 {
			e = v.VendorID.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 269 {
			e = v.ProductName.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 278 {
			e = v.OriginStateID.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 281 {
			e = v.ErrorMessage.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 279 {
			e = v.FailedAVP.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 265 {
			t := SupportedVendorID(0)
			e = t.FromRaw(a)
			v.SupportedVendorID = append(v.SupportedVendorID, t)
		} else if a.VenID == 0 && a.Code == 258 {
			t := AuthApplicationID(0)
			e = t.FromRaw(a)
			v.AuthApplicationID = append(v.AuthApplicationID, t)
		} else if a.VenID == 0 && a.Code == 260 {
			t := VendorSpecificApplicationID{}
			e = t.FromRaw(a)
			v.VendorSpecificApplicationID = append(v.VendorSpecificApplicationID, t)
		} else if a.VenID == 0 && a.Code == 267 {
			e = v.FirmwareRevision.FromRaw(a)
		}

		if e != nil {
			return nil, e
		}
	}
	if v.ResultCode == 0 ||
		len(v.OriginHost) == 0 ||
		len(v.OriginRealm) == 0 ||
		len(v.HostIPAddress) == 0 ||
		v.VendorID == 0 ||
		len(v.ProductName) == 0 {
		e = NoMandatoryAVP{}
	}

	return v, e
}

// Result returns result-code
func (v CEA) Result() ResultCode {
	return v.ResultCode
}

/*
DPR is Disconnect-Peer-Request message
 <DPR>  ::= < Diameter Header: 282, REQ >
			{ Origin-Host }
			{ Origin-Realm }
			{ Disconnect-Cause }
		  * [ AVP ]
*/
type DPR struct {
	OriginHost
	OriginRealm
	DisconnectCause
}

// ToRaw return RawMsg struct of this value
func (v DPR) ToRaw() RawMsg {
	m := RawMsg{
		Ver:  DiaVer,
		FlgR: true, FlgP: false, FlgE: false, FlgT: false,
		Code: 282, AppID: 0,
		AVP: make([]RawAVP, 0, 3)}

	m.AVP = append(m.AVP, v.OriginHost.ToRaw())
	m.AVP = append(m.AVP, v.OriginRealm.ToRaw())
	m.AVP = append(m.AVP, v.DisconnectCause.ToRaw())

	return m
}

// FromRaw make this value from RawMsg struct
func (DPR) FromRaw(m RawMsg) (Request, error) {
	e := m.Validate(0, 282, true, false, false, false)
	if e != nil {
		return nil, e
	}

	v := DPR{}
	v.DisconnectCause = -1
	for _, a := range m.AVP {
		if a.VenID == 0 && a.Code == 264 {
			e = v.OriginHost.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 296 {
			e = v.OriginRealm.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 273 {
			e = v.DisconnectCause.FromRaw(a)
		}

		if e != nil {
			return nil, e
		}
	}

	if len(v.OriginHost) == 0 ||
		len(v.OriginRealm) == 0 ||
		v.DisconnectCause < 0 {
		e = NoMandatoryAVP{}
	}
	return v, e
}

// TimeoutMsg make error message for timeout
func (v DPR) TimeoutMsg() Answer {
	return DPA{
		ResultCode:   DiameterUnableToDeliver,
		OriginHost:   v.OriginHost,
		OriginRealm:  v.OriginRealm,
		ErrorMessage: "no response from peer node"}
}

/*
DPA is Disconnect-Peer-Answer message
 <DPA>  ::= < Diameter Header: 282 >
			{ Result-Code }
			{ Origin-Host }
			{ Origin-Realm }
			[ Error-Message ]
			[ Failed-AVP ]
		  * [ AVP ]
*/
type DPA struct {
	ResultCode
	OriginHost
	OriginRealm
	ErrorMessage
	FailedAVP
}

// ToRaw return RawMsg struct of this value
func (v DPA) ToRaw() RawMsg {
	m := RawMsg{
		Ver:  DiaVer,
		FlgR: false, FlgP: false, FlgE: false, FlgT: false,
		Code: 282, AppID: 0,
		AVP: make([]RawAVP, 0, 5)}
	m.FlgE = v.ResultCode != DiameterSuccess

	m.AVP = append(m.AVP, v.ResultCode.ToRaw())
	m.AVP = append(m.AVP, v.OriginHost.ToRaw())
	m.AVP = append(m.AVP, v.OriginRealm.ToRaw())

	if len(v.ErrorMessage) != 0 {
		m.AVP = append(m.AVP, v.ErrorMessage.ToRaw())
	}
	if len(v.FailedAVP) != 0 {
		m.AVP = append(m.AVP, v.FailedAVP.ToRaw())
	}

	return m
}

// FromRaw make this value from RawMsg struct
func (DPA) FromRaw(m RawMsg) (Answer, error) {
	e := m.Validate(0, 282, false, false, false, false)
	if e != nil {
		return nil, e
	}

	v := DPA{}
	for _, a := range m.AVP {
		if a.VenID == 0 && a.Code == 268 {
			e = v.ResultCode.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 264 {
			e = v.OriginHost.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 296 {
			e = v.OriginRealm.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 281 {
			e = v.ErrorMessage.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 279 {
			e = v.FailedAVP.FromRaw(a)
		}

		if e != nil {
			return nil, e
		}
	}
	if v.ResultCode == 0 ||
		len(v.OriginHost) == 0 ||
		len(v.OriginRealm) == 0 {
		e = NoMandatoryAVP{}
	}
	return v, e
}

// Result returns result-code
func (v DPA) Result() ResultCode {
	return v.ResultCode
}

/*
DWR is DeviceWatchdogRequest message
 <DWR>  ::= < Diameter Header: 280, REQ >
			{ Origin-Host }
			{ Origin-Realm }
			[ Origin-State-Id ]
		  * [ AVP ]
*/
type DWR struct {
	OriginHost
	OriginRealm
	OriginStateID
}

// ToRaw return RawMsg struct of this value
func (v DWR) ToRaw() RawMsg {
	m := RawMsg{
		Ver:  DiaVer,
		FlgR: true, FlgP: false, FlgE: false, FlgT: false,
		Code: 280, AppID: 0,
		AVP: make([]RawAVP, 0, 3)}

	m.AVP = append(m.AVP, v.OriginHost.ToRaw())
	m.AVP = append(m.AVP, v.OriginRealm.ToRaw())

	if v.OriginStateID != 0 {
		m.AVP = append(m.AVP, v.OriginStateID.ToRaw())
	}

	return m
}

// FromRaw make this value from RawMsg struct
func (DWR) FromRaw(m RawMsg) (Request, error) {
	e := m.Validate(0, 280, true, false, false, false)
	if e != nil {
		return nil, e
	}

	v := DWR{}
	for _, a := range m.AVP {
		if a.VenID == 0 && a.Code == 264 {
			e = v.OriginHost.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 296 {
			e = v.OriginRealm.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 278 {
			e = v.OriginStateID.FromRaw(a)
		}

		if e != nil {
			return nil, e
		}
	}
	if len(v.OriginHost) == 0 ||
		len(v.OriginRealm) == 0 {
		e = NoMandatoryAVP{}
	}
	return v, e
}

// TimeoutMsg make error message for timeout
func (v DWR) TimeoutMsg() Answer {
	return DWA{
		ResultCode:    DiameterUnableToDeliver,
		OriginHost:    v.OriginHost,
		OriginRealm:   v.OriginRealm,
		ErrorMessage:  "no response from peer node",
		OriginStateID: v.OriginStateID}
}

/*
DWA Device-Watchdo-gAnswer message
 <DWA>  ::= < Diameter Header: 280 >
			{ Result-Code }
			{ Origin-Host }
			{ Origin-Realm }
			[ Error-Message ]
			[ Failed-AVP ]
			[ Origin-State-Id ]
		  * [ AVP ]
*/
type DWA struct {
	ResultCode
	OriginHost
	OriginRealm
	ErrorMessage
	FailedAVP
	OriginStateID
}

// ToRaw return RawMsg struct of this value
func (v DWA) ToRaw() RawMsg {
	m := RawMsg{
		Ver:  DiaVer,
		FlgR: false, FlgP: false, FlgE: false, FlgT: false,
		Code: 280, AppID: 0,
		AVP: make([]RawAVP, 0, 6)}
	m.FlgE = v.ResultCode != DiameterSuccess

	m.AVP = append(m.AVP, v.ResultCode.ToRaw())
	m.AVP = append(m.AVP, v.OriginHost.ToRaw())
	m.AVP = append(m.AVP, v.OriginRealm.ToRaw())

	if len(v.ErrorMessage) != 0 {
		m.AVP = append(m.AVP, v.ErrorMessage.ToRaw())
	}
	if len(v.FailedAVP) != 0 {
		m.AVP = append(m.AVP, v.FailedAVP.ToRaw())
	}
	if v.OriginStateID != 0 {
		m.AVP = append(m.AVP, v.OriginStateID.ToRaw())
	}

	return m
}

// FromRaw make this value from RawMsg struct
func (DWA) FromRaw(m RawMsg) (Answer, error) {
	e := m.Validate(0, 280, false, false, false, false)
	if e != nil {
		return nil, e
	}

	v := DWA{}
	for _, a := range m.AVP {
		if a.VenID == 0 && a.Code == 268 {
			e = v.ResultCode.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 264 {
			e = v.OriginHost.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 296 {
			e = v.OriginRealm.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 281 {
			e = v.ErrorMessage.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 279 {
			e = v.FailedAVP.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 278 {
			e = v.OriginStateID.FromRaw(a)
		}

		if e != nil {
			return nil, e
		}
	}
	if v.ResultCode == 0 ||
		len(v.OriginHost) == 0 ||
		len(v.OriginRealm) == 0 {
		e = NoMandatoryAVP{}
	}
	return v, e
}

// Result returns result-code
func (v DWA) Result() ResultCode {
	return v.ResultCode
}
