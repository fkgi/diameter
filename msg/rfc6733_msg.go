package msg

/*
CapabilitiesExchangeRequest is CER message
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
type CapabilitiesExchangeRequest struct {
	OriginHost
	OriginRealm
	HostIPAddress []HostIPAddress
	VendorID
	ProductName
	*OriginStateID
	SupportedVendorID []SupportedVendorID
	AuthApplicationID []AuthApplicationID
	// []InbandSecurityID
	// []AcctApplicationId
	VendorSpecificApplicationID []VendorSpecificApplicationID
	*FirmwareRevision
}

// Encode return Message struct of this value
func (v *CapabilitiesExchangeRequest) Encode() Message {
	m := Message{
		Ver:  DiaVer,
		FlgR: true, FlgP: false, FlgE: false, FlgT: false,
		Code: 257, AppID: 0}

	var avps []Avp
	avps = append(avps, v.OriginHost.Encode())
	avps = append(avps, v.OriginRealm.Encode())
	for _, ip := range v.HostIPAddress {
		avps = append(avps, ip.Encode())
	}
	avps = append(avps, v.VendorID.Encode())
	avps = append(avps, v.ProductName.Encode())
	if v.OriginStateID != nil {
		avps = append(avps, v.OriginStateID.Encode())
	}
	for _, id := range v.SupportedVendorID {
		avps = append(avps, id.Encode())
	}
	for _, id := range v.AuthApplicationID {
		avps = append(avps, id.Encode())
	}
	for _, id := range v.VendorSpecificApplicationID {
		avps = append(avps, id.Encode())
	}
	if v.FirmwareRevision != nil {
		avps = append(avps, v.FirmwareRevision.Encode())
	}

	m.Encode(avps)
	return m
}

// Decode make this value from Message struct
func (v *CapabilitiesExchangeRequest) Decode(m Message) error {
	if m.AppID != 0 || m.Code != 257 || !m.FlgR {
		return InvalidMessageError{}
	}
	avp, e := m.Decode()
	if e != nil {
		return e
	}
	var ok bool
	if v.OriginHost, ok = GetOriginHost(avp); !ok {
		return NoMandatoryAVPError{}
	}
	if v.OriginRealm, ok = GetOriginRealm(avp); !ok {
		return NoMandatoryAVPError{}
	}
	v.HostIPAddress = GetHostIPAddresses(avp)
	if v.VendorID, ok = GetVendorID(avp); !ok {
		return NoMandatoryAVPError{}
	}
	if v.ProductName, ok = GetProductName(avp); !ok {
		return NoMandatoryAVPError{}
	}
	if tmp, ok := GetOriginStateID(avp); ok {
		v.OriginStateID = &tmp
	} else {
		v.OriginStateID = nil
	}
	v.SupportedVendorID = GetSupportedVendorIDs(avp)
	v.AuthApplicationID = GetAuthApplicationIDs(avp)
	v.VendorSpecificApplicationID = GetVendorSpecificApplicationIDs(avp)
	if tmp, ok := GetFirmwareRevision(avp); ok {
		v.FirmwareRevision = &tmp
	} else {
		v.FirmwareRevision = nil
	}
	return nil
}

/*
CapabilitiesExchangeAnswer is CEA message
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
type CapabilitiesExchangeAnswer struct {
	ResultCode
	OriginHost
	OriginRealm
	HostIPAddress []HostIPAddress
	VendorID
	ProductName
	*OriginStateID
	*ErrorMessage
	*FailedAVP
	SupportedVendorID []SupportedVendorID
	AuthApplicationID []AuthApplicationID
	// []InbandSecurityID
	// []AcctApplicationId
	VendorSpecificApplicationID []VendorSpecificApplicationID
	*FirmwareRevision
}

// Encode return Message struct of this value
func (v *CapabilitiesExchangeAnswer) Encode() Message {
	m := Message{
		Ver:  DiaVer,
		FlgR: false, FlgP: false, FlgE: false, FlgT: false,
		Code: 257, AppID: 0}
	m.FlgE = v.ResultCode != DiameterSuccess

	var avps []Avp
	avps = append(avps, v.ResultCode.Encode())
	avps = append(avps, v.OriginHost.Encode())
	avps = append(avps, v.OriginRealm.Encode())
	for _, ip := range v.HostIPAddress {
		avps = append(avps, ip.Encode())
	}
	avps = append(avps, v.VendorID.Encode())
	avps = append(avps, v.ProductName.Encode())
	if v.OriginStateID != nil {
		avps = append(avps, v.OriginStateID.Encode())
	}
	if v.ErrorMessage != nil {
		avps = append(avps, v.ErrorMessage.Encode())
	}
	if v.FailedAVP != nil {
		avps = append(avps, v.FailedAVP.Encode())
	}
	for _, id := range v.SupportedVendorID {
		avps = append(avps, id.Encode())
	}
	for _, id := range v.AuthApplicationID {
		avps = append(avps, id.Encode())
	}
	for _, id := range v.VendorSpecificApplicationID {
		avps = append(avps, id.Encode())
	}
	if v.FirmwareRevision != nil {
		avps = append(avps, v.FirmwareRevision.Encode())
	}

	m.Encode(avps)
	return m
}

// Decode make this value from Message struct
func (v *CapabilitiesExchangeAnswer) Decode(m Message) error {
	if m.AppID != 0 || m.Code != 257 || m.FlgR {
		return InvalidMessageError{}
	}
	avp, e := m.Decode()
	if e != nil {
		return e
	}
	var ok bool
	if v.ResultCode, ok = GetResultCode(avp); !ok {
		return NoMandatoryAVPError{}
	}
	if v.OriginHost, ok = GetOriginHost(avp); !ok {
		return NoMandatoryAVPError{}
	}
	if v.OriginRealm, ok = GetOriginRealm(avp); !ok {
		return NoMandatoryAVPError{}
	}
	v.HostIPAddress = GetHostIPAddresses(avp)
	if v.VendorID, ok = GetVendorID(avp); !ok {
		return NoMandatoryAVPError{}
	}
	if v.ProductName, ok = GetProductName(avp); !ok {
		return NoMandatoryAVPError{}
	}
	if tmp, ok := GetOriginStateID(avp); ok {
		v.OriginStateID = &tmp
	} else {
		v.OriginStateID = nil
	}
	if tmp, ok := GetErrorMessage(avp); ok {
		v.ErrorMessage = &tmp
	} else {
		v.ErrorMessage = nil
	}
	if tmp, ok := GetFailedAVP(avp); ok {
		v.FailedAVP = &tmp
	} else {
		v.FailedAVP = nil
	}
	v.SupportedVendorID = GetSupportedVendorIDs(avp)
	v.AuthApplicationID = GetAuthApplicationIDs(avp)
	v.VendorSpecificApplicationID = GetVendorSpecificApplicationIDs(avp)
	if tmp, ok := GetFirmwareRevision(avp); ok {
		v.FirmwareRevision = &tmp
	} else {
		v.FirmwareRevision = nil
	}
	return nil
}

/*
DisconnectPeerRequest is DPR message
 <DPR>  ::= < Diameter Header: 282, REQ >
			{ Origin-Host }
			{ Origin-Realm }
			{ Disconnect-Cause }
		  * [ AVP ]
*/
type DisconnectPeerRequest struct {
	OriginHost
	OriginRealm
	DisconnectCause
}

// Encode return Message struct of this value
func (v *DisconnectPeerRequest) Encode() Message {
	m := Message{
		Ver:  DiaVer,
		FlgR: true, FlgP: false, FlgE: false, FlgT: false,
		Code: 282, AppID: 0}

	avps := []Avp{
		v.OriginHost.Encode(),
		v.OriginRealm.Encode(),
		v.DisconnectCause.Encode()}
	m.Encode(avps)

	return m
}

// Decode make this value from Message struct
func (v *DisconnectPeerRequest) Decode(m Message) error {
	if m.AppID != 0 || m.Code != 282 || !m.FlgR {
		return InvalidMessageError{}
	}
	avp, e := m.Decode()
	if e != nil {
		return e
	}
	var ok bool
	if v.OriginHost, ok = GetOriginHost(avp); !ok {
		return NoMandatoryAVPError{}
	}
	if v.OriginRealm, ok = GetOriginRealm(avp); !ok {
		return NoMandatoryAVPError{}
	}
	if v.DisconnectCause, ok = GetDisconnectCause(avp); !ok {
		return NoMandatoryAVPError{}
	}
	return nil
}

/*
DisconnectPeerAnswer is DPA message
 <DPA>  ::= < Diameter Header: 282 >
			{ Result-Code }
			{ Origin-Host }
			{ Origin-Realm }
			[ Error-Message ]
			[ Failed-AVP ]
		  * [ AVP ]
*/
type DisconnectPeerAnswer struct {
	ResultCode
	OriginHost
	OriginRealm
	*ErrorMessage
	*FailedAVP
}

// Encode return Message struct of this value
func (v *DisconnectPeerAnswer) Encode() Message {
	m := Message{
		Ver:  DiaVer,
		FlgR: false, FlgP: false, FlgE: false, FlgT: false,
		Code: 282, AppID: 0}
	m.FlgE = v.ResultCode != DiameterSuccess

	avps := []Avp{
		v.ResultCode.Encode(),
		v.OriginHost.Encode(),
		v.OriginRealm.Encode()}
	if v.ErrorMessage != nil {
		avps = append(avps, v.ErrorMessage.Encode())
	}
	if v.FailedAVP != nil {
		avps = append(avps, v.FailedAVP.Encode())
	}
	m.Encode(avps)

	return m
}

// Decode make this value from Message struct
func (v *DisconnectPeerAnswer) Decode(m Message) error {
	if m.AppID != 0 || m.Code != 282 || m.FlgR {
		return InvalidMessageError{}
	}
	avp, e := m.Decode()
	if e != nil {
		return e
	}
	var ok bool
	if v.ResultCode, ok = GetResultCode(avp); !ok {
		return NoMandatoryAVPError{}
	}
	if v.OriginHost, ok = GetOriginHost(avp); !ok {
		return NoMandatoryAVPError{}
	}
	if v.OriginRealm, ok = GetOriginRealm(avp); !ok {
		return NoMandatoryAVPError{}
	}
	if tmp, ok := GetErrorMessage(avp); ok {
		v.ErrorMessage = &tmp
	} else {
		v.ErrorMessage = nil
	}
	if tmp, ok := GetFailedAVP(avp); ok {
		v.FailedAVP = &tmp
	} else {
		v.FailedAVP = nil
	}
	return nil
}

/*
DeviceWatchdogRequest is DWR message
 <DWR>  ::= < Diameter Header: 280, REQ >
			{ Origin-Host }
			{ Origin-Realm }
			[ Origin-State-Id ]
		  * [ AVP ]
*/
type DeviceWatchdogRequest struct {
	OriginHost
	OriginRealm
	*OriginStateID
}

// Encode return Message struct of this value
func (v *DeviceWatchdogRequest) Encode() Message {
	m := Message{
		Ver:  DiaVer,
		FlgR: true, FlgP: false, FlgE: false, FlgT: false,
		Code: 280, AppID: 0}

	avps := []Avp{
		v.OriginHost.Encode(),
		v.OriginRealm.Encode()}
	if v.OriginStateID != nil {
		avps = append(avps, v.OriginStateID.Encode())
	}
	m.Encode(avps)

	return m
}

// Decode make this value from Message struct
func (v *DeviceWatchdogRequest) Decode(m Message) error {
	if m.AppID != 0 || m.Code != 280 || !m.FlgR {
		return InvalidMessageError{}
	}
	avp, e := m.Decode()
	if e != nil {
		return e
	}
	var ok bool
	if v.OriginHost, ok = GetOriginHost(avp); !ok {
		return NoMandatoryAVPError{}
	}
	if v.OriginRealm, ok = GetOriginRealm(avp); !ok {
		return NoMandatoryAVPError{}
	}
	if tmp, ok := GetOriginStateID(avp); ok {
		v.OriginStateID = &tmp
	} else {
		v.OriginStateID = nil
	}
	return nil
}

/*
DeviceWatchdogAnswer is DWA message
 <DWA>  ::= < Diameter Header: 280 >
			{ Result-Code }
			{ Origin-Host }
			{ Origin-Realm }
			[ Error-Message ]
			[ Failed-AVP ]
			[ Origin-State-Id ]
		  * [ AVP ]
*/
type DeviceWatchdogAnswer struct {
	ResultCode
	OriginHost
	OriginRealm
	*ErrorMessage
	*FailedAVP
	*OriginStateID
}

// Encode return Message struct of this value
func (v *DeviceWatchdogAnswer) Encode() Message {
	m := Message{
		Ver:  DiaVer,
		FlgR: false, FlgP: false, FlgE: false, FlgT: false,
		Code: 280, AppID: 0}
	m.FlgE = v.ResultCode != DiameterSuccess

	avps := []Avp{
		v.ResultCode.Encode(),
		v.OriginHost.Encode(),
		v.OriginRealm.Encode()}
	if v.ErrorMessage != nil {
		avps = append(avps, v.ErrorMessage.Encode())
	}
	if v.FailedAVP != nil {
		avps = append(avps, v.FailedAVP.Encode())
	}
	if v.OriginStateID != nil {
		avps = append(avps, v.OriginStateID.Encode())
	}
	m.Encode(avps)

	return m
}

// Decode make this value from Message struct
func (v *DeviceWatchdogAnswer) Decode(m Message) error {
	if m.AppID != 0 || m.Code != 280 || m.FlgR {
		return InvalidMessageError{}
	}
	avp, e := m.Decode()
	if e != nil {
		return e
	}
	var ok bool
	if v.ResultCode, ok = GetResultCode(avp); !ok {
		return NoMandatoryAVPError{}
	}
	if v.OriginHost, ok = GetOriginHost(avp); !ok {
		return NoMandatoryAVPError{}
	}
	if v.OriginRealm, ok = GetOriginRealm(avp); !ok {
		return NoMandatoryAVPError{}
	}
	if tmp, ok := GetErrorMessage(avp); ok {
		v.ErrorMessage = &tmp
	} else {
		v.ErrorMessage = nil
	}
	if tmp, ok := GetFailedAVP(avp); ok {
		v.FailedAVP = &tmp
	} else {
		v.FailedAVP = nil
	}
	if tmp, ok := GetOriginStateID(avp); ok {
		v.OriginStateID = &tmp
	} else {
		v.OriginStateID = nil
	}
	return nil
}
