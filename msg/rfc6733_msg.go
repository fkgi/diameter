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
	*OriginStateID
	SupportedVendorID []SupportedVendorID
	ApplicationID     []ApplicationID
	// []InbandSecurityID
	VendorSpecificApplicationID []VendorSpecificApplicationID
	*FirmwareRevision
}

// Encode return Message struct of this value
func (v *CER) Encode() Message {
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
	for _, id := range v.ApplicationID {
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
func (v *CER) Decode(m Message) error {
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
	for _, a := range GetAuthApplicationIDs(avp) {
		v.ApplicationID = append(v.ApplicationID, a)
	}
	for _, a := range GetAcctApplicationIDs(avp) {
		v.ApplicationID = append(v.ApplicationID, a)
	}
	v.VendorSpecificApplicationID = GetVendorSpecificApplicationIDs(avp)
	if tmp, ok := GetFirmwareRevision(avp); ok {
		v.FirmwareRevision = &tmp
	} else {
		v.FirmwareRevision = nil
	}
	return nil
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
type CEA struct {
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
	ApplicationID     []ApplicationID
	// []InbandSecurityID
	// []AcctApplicationId
	VendorSpecificApplicationID []VendorSpecificApplicationID
	*FirmwareRevision
}

// Encode return Message struct of this value
func (v *CEA) Encode() Message {
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
	for _, id := range v.ApplicationID {
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
func (v *CEA) Decode(m Message) error {
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
	for _, a := range GetAuthApplicationIDs(avp) {
		v.ApplicationID = append(v.ApplicationID, a)
	}
	for _, a := range GetAcctApplicationIDs(avp) {
		v.ApplicationID = append(v.ApplicationID, a)
	}
	v.VendorSpecificApplicationID = GetVendorSpecificApplicationIDs(avp)
	if tmp, ok := GetFirmwareRevision(avp); ok {
		v.FirmwareRevision = &tmp
	} else {
		v.FirmwareRevision = nil
	}
	return nil
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

// Encode return Message struct of this value
func (v *DPR) Encode() Message {
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
func (v *DPR) Decode(m Message) error {
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
	*ErrorMessage
	*FailedAVP
}

// Encode return Message struct of this value
func (v *DPA) Encode() Message {
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
func (v *DPA) Decode(m Message) error {
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
	*OriginStateID
}

// Encode return Message struct of this value
func (v *DWR) Encode() Message {
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
func (v *DWR) Decode(m Message) error {
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
	*ErrorMessage
	*FailedAVP
	*OriginStateID
}

// Encode return Message struct of this value
func (v *DWA) Encode() Message {
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
func (v *DWA) Decode(m Message) error {
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
