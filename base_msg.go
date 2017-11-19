package diameter

import (
	"bytes"
	"fmt"
	"net"
)

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
	OriginHost    Identity
	OriginRealm   Identity
	HostIPAddress []net.IP
	VendorID      uint32
	ProductName   string
	OriginStateID uint32
	ApplicationID map[uint32][]uint32
	// []InbandSecurityID
	FirmwareRevision uint32
}

func (v CER) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "%sOrigin-Host     =%s\n", Indent, v.OriginHost)
	fmt.Fprintf(w, "%sOrigin-Realm    =%s\n", Indent, v.OriginRealm)
	fmt.Fprintf(w, "%sHost-IP-Address =", Indent)
	for _, ip := range v.HostIPAddress {
		fmt.Fprintf(w, "%s ", ip)
	}
	fmt.Fprint(w, "\n")
	fmt.Fprintf(w, "%sVendor-ID       =%d\n", Indent, v.VendorID)
	fmt.Fprintf(w, "%sProduct-Name    =%s\n", Indent, v.ProductName)
	fmt.Fprintf(w, "%sOrigin-State-ID =0x%x\n", Indent, v.OriginStateID)
	fmt.Fprintf(w, "%sSupported-Application-ID =\n", Indent)
	for vID, aIDs := range v.ApplicationID {
		fmt.Fprintf(w, "%s%sVendor-ID =%d\n", Indent, Indent, vID)
		for _, aID := range aIDs {
			fmt.Fprintf(w, "%s%s%sApplication-ID =%d\n", Indent, Indent, Indent, aID)
		}
	}
	fmt.Fprintf(w, "%sFirmware-Revision=%d", Indent, v.FirmwareRevision)

	return w.String()
}

// ToRaw return RawMsg struct of this value
func (v CER) ToRaw(s string) RawMsg {
	m := RawMsg{
		Ver:  DiaVer,
		FlgR: true, FlgP: false, FlgE: false, FlgT: false,
		Code: 257, AppID: 0,
		AVP: make([]RawAVP, 0, 20)}

	m.AVP = append(m.AVP, SetOriginHost(v.OriginHost))
	m.AVP = append(m.AVP, SetOriginRealm(v.OriginRealm))
	for _, ip := range v.HostIPAddress {
		m.AVP = append(m.AVP, setHostIPAddress(ip))
	}
	m.AVP = append(m.AVP, setVendorID(v.VendorID))
	m.AVP = append(m.AVP, setProductName(v.ProductName))
	if v.OriginStateID != 0 {
		m.AVP = append(m.AVP, setOriginStateID(v.OriginStateID))
	}
	for vID, aIDs := range v.ApplicationID {
		if vID == 0 {
			for _, aID := range aIDs {
				m.AVP = append(m.AVP, setAuthAppID(aID))
			}
		} else {
			m.AVP = append(m.AVP, setSupportedVendorID(vID))
			for _, aID := range aIDs {
				m.AVP = append(m.AVP, SetVendorSpecAppID(vID, aID))
			}
		}
	}
	if v.FirmwareRevision != 0 {
		m.AVP = append(m.AVP, setFirmwareRevision(v.FirmwareRevision))
	}
	return m
}

// FromRaw make this value from RawMsg struct
func (CER) FromRaw(m RawMsg) (Request, string, error) {
	e := m.Validate(0, 257, true, false, false, false)
	if e != nil {
		return nil, "", e
	}

	v := CER{
		HostIPAddress: make([]net.IP, 0, 2),
		ApplicationID: make(map[uint32][]uint32, 5)}

	for _, a := range m.AVP {
		if a.VenID == 0 && a.Code == 264 {
			v.OriginHost, e = GetOriginHost(a)
		} else if a.VenID == 0 && a.Code == 296 {
			v.OriginRealm, e = GetOriginRealm(a)
		} else if a.VenID == 0 && a.Code == 257 {
			if t, e2 := getHostIPAddress(a); e2 == nil {
				v.HostIPAddress = append(v.HostIPAddress, t)
			} else {
				e = e2
			}
		} else if a.VenID == 0 && a.Code == 266 {
			v.VendorID, e = getVendorID(a)
		} else if a.VenID == 0 && a.Code == 269 {
			v.ProductName, e = getProductName(a)
		} else if a.VenID == 0 && a.Code == 278 {
			v.OriginStateID, e = getOriginStateID(a)
		} else if a.VenID == 0 && a.Code == 265 {
			if t, e2 := getSupportedVendorID(a); e2 != nil {
				e = e2
			} else if _, ok := v.ApplicationID[t]; !ok {
				v.ApplicationID[t] = []uint32{}
			}
		} else if a.VenID == 0 && a.Code == 258 {
			if t, e2 := getAuthAppID(a); e2 != nil {
				e = e2
			} else if _, ok := v.ApplicationID[0]; !ok {
				v.ApplicationID[0] = []uint32{t}
			} else {
				v.ApplicationID[0] = append(v.ApplicationID[0], t)
			}
		} else if a.VenID == 0 && a.Code == 260 {
			if vi, ai, e2 := GetVendorSpecAppID(a); e2 != nil {
				e = e2
			} else if _, ok := v.ApplicationID[vi]; !ok {
				v.ApplicationID[vi] = []uint32{ai}
			} else {
				v.ApplicationID[vi] = append(v.ApplicationID[vi], ai)
			}
		} else if a.VenID == 0 && a.Code == 267 {
			v.FirmwareRevision, e = getFirmwareRevision(a)
		}

		if e != nil {
			return nil, "", e
		}
	}

	if len(v.OriginHost) == 0 ||
		len(v.OriginRealm) == 0 ||
		len(v.HostIPAddress) == 0 ||
		v.VendorID == 0 ||
		len(v.ProductName) == 0 {
		e = NoMandatoryAVP{}
	}
	return v, "", e
}

// Failed make error message for timeout
func (v CER) Failed(c uint32) Answer {
	return CEA{
		ResultCode:    c,
		OriginHost:    v.OriginHost,
		OriginRealm:   v.OriginRealm,
		HostIPAddress: v.HostIPAddress,
		VendorID:      v.VendorID,
		ProductName:   v.ProductName}
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
	ResultCode    uint32
	OriginHost    Identity
	OriginRealm   Identity
	HostIPAddress []net.IP
	VendorID      uint32
	ProductName   string
	OriginStateID uint32
	ErrorMessage  string
	FailedAVP     []RawAVP
	ApplicationID map[uint32][]uint32
	// []InbandSecurityID
	FirmwareRevision uint32
}

func (v CEA) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "%sResult-Code     =%d\n", Indent, v.ResultCode)
	fmt.Fprintf(w, "%sOrigin-Host     =%s\n", Indent, v.OriginHost)
	fmt.Fprintf(w, "%sOrigin-Realm    =%s\n", Indent, v.OriginRealm)
	fmt.Fprintf(w, "%sHost-IP-Address =", Indent)
	for _, ip := range v.HostIPAddress {
		fmt.Fprintf(w, "%s ", ip)
	}
	fmt.Fprint(w, "\n")
	fmt.Fprintf(w, "%sVendor-ID       =%d\n", Indent, v.VendorID)
	fmt.Fprintf(w, "%sProduct-Name    =%s\n", Indent, v.ProductName)
	fmt.Fprintf(w, "%sOrigin-State-ID =0x%x\n", Indent, v.OriginStateID)
	for _, avp := range v.FailedAVP {
		fmt.Fprintf(w, "%sFailed-AVP      =\n%s", Indent, avp)
	}
	fmt.Fprintf(w, "%sSupported-Application-ID =\n", Indent)
	for vID, aIDs := range v.ApplicationID {
		fmt.Fprintf(w, "%s%sVendor-ID =%d\n", Indent, Indent, vID)
		for _, aID := range aIDs {
			fmt.Fprintf(w, "%s%s%sApplication-ID =%d\n", Indent, Indent, Indent, aID)
		}
	}
	fmt.Fprintf(w, "%sFirmware-Revision=%d", Indent, v.FirmwareRevision)

	return w.String()
}

// ToRaw return RawMsg struct of this value
func (v CEA) ToRaw(s string) RawMsg {
	m := RawMsg{
		Ver:  DiaVer,
		FlgR: false, FlgP: false, FlgE: false, FlgT: false,
		Code: 257, AppID: 0,
		AVP: make([]RawAVP, 0, 20)}
	m.FlgE = v.ResultCode != DiameterSuccess

	m.AVP = append(m.AVP, SetResultCode(v.ResultCode))
	m.AVP = append(m.AVP, SetOriginHost(v.OriginHost))
	m.AVP = append(m.AVP, SetOriginRealm(v.OriginRealm))
	for _, ip := range v.HostIPAddress {
		m.AVP = append(m.AVP, setHostIPAddress(ip))
	}
	m.AVP = append(m.AVP, setVendorID(v.VendorID))
	m.AVP = append(m.AVP, setProductName(v.ProductName))
	if v.OriginStateID != 0 {
		m.AVP = append(m.AVP, setOriginStateID(v.OriginStateID))
	}
	if len(v.ErrorMessage) != 0 {
		m.AVP = append(m.AVP, setErrorMessage(v.ErrorMessage))
	}
	if len(v.FailedAVP) != 0 {
		m.AVP = append(m.AVP, setFailedAVP(v.FailedAVP))
	}
	for vID, aIDs := range v.ApplicationID {
		if vID == 0 {
			for _, aID := range aIDs {
				m.AVP = append(m.AVP, setAuthAppID(aID))
			}
		} else {
			m.AVP = append(m.AVP, setSupportedVendorID(vID))
			for _, aID := range aIDs {
				m.AVP = append(m.AVP, SetVendorSpecAppID(vID, aID))
			}
		}
	}
	if v.FirmwareRevision != 0 {
		m.AVP = append(m.AVP, setFirmwareRevision(v.FirmwareRevision))
	}
	return m
}

// FromRaw make this value from RawMsg struct
func (CEA) FromRaw(m RawMsg) (Answer, string, error) {
	e := m.Validate(0, 257, false, false, false, false)
	if e != nil {
		return nil, "", e
	}

	v := CEA{
		HostIPAddress: make([]net.IP, 0, 2),
		ApplicationID: make(map[uint32][]uint32, 5)}

	for _, a := range m.AVP {
		if a.VenID == 0 && a.Code == 268 {
			v.ResultCode, e = GetResultCode(a)
		} else if a.VenID == 0 && a.Code == 264 {
			v.OriginHost, e = GetOriginHost(a)
		} else if a.VenID == 0 && a.Code == 296 {
			v.OriginRealm, e = GetOriginRealm(a)
		} else if a.VenID == 0 && a.Code == 257 {
			if t, e2 := getHostIPAddress(a); e2 == nil {
				v.HostIPAddress = append(v.HostIPAddress, t)
			} else {
				e = e2
			}
		} else if a.VenID == 0 && a.Code == 266 {
			v.VendorID, e = getVendorID(a)
		} else if a.VenID == 0 && a.Code == 269 {
			v.ProductName, e = getProductName(a)
		} else if a.VenID == 0 && a.Code == 278 {
			v.OriginStateID, e = getOriginStateID(a)
		} else if a.VenID == 0 && a.Code == 281 {
			v.ErrorMessage, e = getErrorMessage(a)
		} else if a.VenID == 0 && a.Code == 279 {
			v.FailedAVP, e = getFailedAVP(a)
		} else if a.VenID == 0 && a.Code == 265 {
			if t, e2 := getSupportedVendorID(a); e2 != nil {
				e = e2
			} else if _, ok := v.ApplicationID[t]; !ok {
				v.ApplicationID[t] = []uint32{}
			}
		} else if a.VenID == 0 && a.Code == 258 {
			if t, e2 := getAuthAppID(a); e2 != nil {
				e = e2
			} else if _, ok := v.ApplicationID[0]; !ok {
				v.ApplicationID[0] = []uint32{t}
			} else {
				v.ApplicationID[0] = append(v.ApplicationID[0], t)
			}
		} else if a.VenID == 0 && a.Code == 260 {
			if vi, ai, e2 := GetVendorSpecAppID(a); e2 != nil {
				e = e2
			} else if _, ok := v.ApplicationID[vi]; !ok {
				v.ApplicationID[vi] = []uint32{ai}
			} else {
				v.ApplicationID[vi] = append(v.ApplicationID[vi], ai)
			}
		} else if a.VenID == 0 && a.Code == 267 {
			v.FirmwareRevision, e = getFirmwareRevision(a)
		}

		if e != nil {
			return nil, "", e
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

	return v, "", e
}

// Result returns result-code
func (v CEA) Result() uint32 {
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
	OriginHost      Identity
	OriginRealm     Identity
	DisconnectCause Enumerated
}

func (v DPR) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "%sOrigin-Host     =%s\n", Indent, v.OriginHost)
	fmt.Fprintf(w, "%sOrigin-Realm    =%s\n", Indent, v.OriginRealm)
	fmt.Fprintf(w, "%sDisconnect-Cause=%d\n", Indent, v.DisconnectCause)

	return w.String()
}

// ToRaw return RawMsg struct of this value
func (v DPR) ToRaw(s string) RawMsg {
	m := RawMsg{
		Ver:  DiaVer,
		FlgR: true, FlgP: false, FlgE: false, FlgT: false,
		Code: 282, AppID: 0,
		AVP: make([]RawAVP, 0, 3)}

	m.AVP = append(m.AVP, SetOriginHost(v.OriginHost))
	m.AVP = append(m.AVP, SetOriginRealm(v.OriginRealm))
	m.AVP = append(m.AVP, setDisconnectCause(v.DisconnectCause))
	return m
}

// FromRaw make this value from RawMsg struct
func (DPR) FromRaw(m RawMsg) (Request, string, error) {
	e := m.Validate(0, 282, true, false, false, false)
	if e != nil {
		return nil, "", e
	}

	v := DPR{
		DisconnectCause: -1}
	for _, a := range m.AVP {
		if a.VenID == 0 && a.Code == 264 {
			v.OriginHost, e = GetOriginHost(a)
		} else if a.VenID == 0 && a.Code == 296 {
			v.OriginRealm, e = GetOriginRealm(a)
		} else if a.VenID == 0 && a.Code == 273 {
			v.DisconnectCause, e = getDisconnectCause(a)
		}

		if e != nil {
			return nil, "", e
		}
	}

	if len(v.OriginHost) == 0 ||
		len(v.OriginRealm) == 0 ||
		v.DisconnectCause < 0 {
		e = NoMandatoryAVP{}
	}
	return v, "", e
}

// Failed make error message for timeout
func (v DPR) Failed(c uint32) Answer {
	return DPA{
		ResultCode:  c,
		OriginHost:  v.OriginHost,
		OriginRealm: v.OriginRealm}
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
	ResultCode   uint32
	OriginHost   Identity
	OriginRealm  Identity
	ErrorMessage string
	FailedAVP    []RawAVP
}

func (v DPA) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "%sResult-Code     =%d\n", Indent, v.ResultCode)
	fmt.Fprintf(w, "%sOrigin-Host     =%s\n", Indent, v.OriginHost)
	fmt.Fprintf(w, "%sOrigin-Realm    =%s\n", Indent, v.OriginRealm)
	fmt.Fprintf(w, "%sError-Message   =%s\n", Indent, v.ErrorMessage)
	for _, avp := range v.FailedAVP {
		fmt.Fprintf(w, "%sFailed-AVP      =\n%s", Indent, avp)
	}

	return w.String()
}

// ToRaw return RawMsg struct of this value
func (v DPA) ToRaw(s string) RawMsg {
	m := RawMsg{
		Ver:  DiaVer,
		FlgR: false, FlgP: false, FlgE: false, FlgT: false,
		Code: 282, AppID: 0,
		AVP: make([]RawAVP, 0, 5)}
	m.FlgE = v.ResultCode != DiameterSuccess

	m.AVP = append(m.AVP, SetResultCode(v.ResultCode))
	m.AVP = append(m.AVP, SetOriginHost(v.OriginHost))
	m.AVP = append(m.AVP, SetOriginRealm(v.OriginRealm))
	if len(v.ErrorMessage) != 0 {
		m.AVP = append(m.AVP, setErrorMessage(v.ErrorMessage))
	}
	if len(v.FailedAVP) != 0 {
		m.AVP = append(m.AVP, setFailedAVP(v.FailedAVP))
	}
	return m
}

// FromRaw make this value from RawMsg struct
func (DPA) FromRaw(m RawMsg) (Answer, string, error) {
	e := m.Validate(0, 282, false, false, false, false)
	if e != nil {
		return nil, "", e
	}

	v := DPA{}
	for _, a := range m.AVP {
		if a.VenID == 0 && a.Code == 268 {
			v.ResultCode, e = GetResultCode(a)
		} else if a.VenID == 0 && a.Code == 264 {
			v.OriginHost, e = GetOriginHost(a)
		} else if a.VenID == 0 && a.Code == 296 {
			v.OriginRealm, e = GetOriginRealm(a)
		} else if a.VenID == 0 && a.Code == 281 {
			v.ErrorMessage, e = getErrorMessage(a)
		} else if a.VenID == 0 && a.Code == 279 {
			v.FailedAVP, e = getFailedAVP(a)
		}

		if e != nil {
			return nil, "", e
		}
	}
	if v.ResultCode == 0 ||
		len(v.OriginHost) == 0 ||
		len(v.OriginRealm) == 0 {
		e = NoMandatoryAVP{}
	}
	return v, "", e
}

// Result returns result-code
func (v DPA) Result() uint32 {
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
	OriginHost    Identity
	OriginRealm   Identity
	OriginStateID uint32
}

func (v DWR) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "%sOrigin-Host     =%s\n", Indent, v.OriginHost)
	fmt.Fprintf(w, "%sOrigin-Realm    =%s\n", Indent, v.OriginRealm)
	fmt.Fprintf(w, "%sOrigin-State-ID =%d\n", Indent, v.OriginStateID)

	return w.String()
}

// ToRaw return RawMsg struct of this value
func (v DWR) ToRaw(s string) RawMsg {
	m := RawMsg{
		Ver:  DiaVer,
		FlgR: true, FlgP: false, FlgE: false, FlgT: false,
		Code: 280, AppID: 0,
		AVP: make([]RawAVP, 0, 3)}

	m.AVP = append(m.AVP, SetOriginHost(v.OriginHost))
	m.AVP = append(m.AVP, SetOriginRealm(v.OriginRealm))
	if v.OriginStateID != 0 {
		m.AVP = append(m.AVP, setOriginStateID(v.OriginStateID))
	}
	return m
}

// FromRaw make this value from RawMsg struct
func (DWR) FromRaw(m RawMsg) (Request, string, error) {
	e := m.Validate(0, 280, true, false, false, false)
	if e != nil {
		return nil, "", e
	}

	v := DWR{}
	for _, a := range m.AVP {
		if a.VenID == 0 && a.Code == 264 {
			v.OriginHost, e = GetOriginHost(a)
		} else if a.VenID == 0 && a.Code == 296 {
			v.OriginRealm, e = GetOriginRealm(a)
		} else if a.VenID == 0 && a.Code == 278 {
			v.OriginStateID, e = getOriginStateID(a)
		}

		if e != nil {
			return nil, "", e
		}
	}
	if len(v.OriginHost) == 0 ||
		len(v.OriginRealm) == 0 {
		e = NoMandatoryAVP{}
	}
	return v, "", e
}

// Failed make error message for timeout
func (v DWR) Failed(c uint32) Answer {
	return DWA{
		ResultCode:  c,
		OriginHost:  v.OriginHost,
		OriginRealm: v.OriginRealm}
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
	ResultCode    uint32
	OriginHost    Identity
	OriginRealm   Identity
	ErrorMessage  string
	FailedAVP     []RawAVP
	OriginStateID uint32
}

func (v DWA) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "%sResult-Code     =%d\n", Indent, v.ResultCode)
	fmt.Fprintf(w, "%sOrigin-Host     =%s\n", Indent, v.OriginHost)
	fmt.Fprintf(w, "%sOrigin-Realm    =%s\n", Indent, v.OriginRealm)
	fmt.Fprintf(w, "%sError-Message   =%s\n", Indent, v.ErrorMessage)
	for _, avp := range v.FailedAVP {
		fmt.Fprintf(w, "%sFailed-AVP      =\n%s", Indent, avp)
	}
	fmt.Fprintf(w, "%sOrigin-State-ID =%d\n", Indent, v.OriginStateID)

	return w.String()
}

// ToRaw return RawMsg struct of this value
func (v DWA) ToRaw(s string) RawMsg {
	m := RawMsg{
		Ver:  DiaVer,
		FlgR: false, FlgP: false, FlgE: false, FlgT: false,
		Code: 280, AppID: 0,
		AVP: make([]RawAVP, 0, 6)}
	m.FlgE = v.ResultCode != DiameterSuccess

	m.AVP = append(m.AVP, SetResultCode(v.ResultCode))
	m.AVP = append(m.AVP, SetOriginHost(v.OriginHost))
	m.AVP = append(m.AVP, SetOriginRealm(v.OriginRealm))
	if len(v.ErrorMessage) != 0 {
		m.AVP = append(m.AVP, setErrorMessage(v.ErrorMessage))
	}
	if len(v.FailedAVP) != 0 {
		m.AVP = append(m.AVP, setFailedAVP(v.FailedAVP))
	}
	if v.OriginStateID != 0 {
		m.AVP = append(m.AVP, setOriginStateID(v.OriginStateID))
	}
	return m
}

// FromRaw make this value from RawMsg struct
func (DWA) FromRaw(m RawMsg) (Answer, string, error) {
	e := m.Validate(0, 280, false, false, false, false)
	if e != nil {
		return nil, "", e
	}

	v := DWA{}
	for _, a := range m.AVP {
		if a.VenID == 0 && a.Code == 268 {
			v.ResultCode, e = GetResultCode(a)
		} else if a.VenID == 0 && a.Code == 264 {
			v.OriginHost, e = GetOriginHost(a)
		} else if a.VenID == 0 && a.Code == 296 {
			v.OriginRealm, e = GetOriginRealm(a)
		} else if a.VenID == 0 && a.Code == 281 {
			v.ErrorMessage, e = getErrorMessage(a)
		} else if a.VenID == 0 && a.Code == 279 {
			v.FailedAVP, e = getFailedAVP(a)
		} else if a.VenID == 0 && a.Code == 278 {
			v.OriginStateID, e = getOriginStateID(a)
		}

		if e != nil {
			return nil, "", e
		}
	}
	if v.ResultCode == 0 ||
		len(v.OriginHost) == 0 ||
		len(v.OriginRealm) == 0 {
		e = NoMandatoryAVP{}
	}
	return v, "", e
}

// Result returns result-code
func (v DWA) Result() uint32 {
	return v.ResultCode
}
