package ts29338

import (
	"bytes"
	"fmt"
	"time"

	dia "github.com/fkgi/diameter"
	"github.com/fkgi/sms"
	"github.com/fkgi/teldata"
)

/*
TFR is MT-Forward-Short-Message-Request message.
 <TFR> ::= < Diameter Header: 8388646, REQ, PXY, 16777313 >
           < Session-Id >
           [ DRMP ]  // not supported
           [ Vendor-Specific-Application-Id ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
           { Destination-Host }
           { Destination-Realm }
           { User-Name }
         * [ Supported-Features ]  // not supported
           [ SMSMI-Correlation-ID ]  // not supported
           { SC-Address }
           { SM-RP-UI }
           [ MME-Number-for-MT-SMS ]
           [ SGSN-Number ]
           [ TFR-Flags ]
           [ SM-Delivery-Timer ]
           [ SM-Delivery-Start-Time ]
           [ Maximum-Retransmission-Time ]
           [ SMS-GMSC-Address ]
         * [ AVP ]
         * [ Proxy-Info ]  // not supported
         * [ Route-Record ]
*/
type TFR struct {
	OriginHost       dia.Identity
	OriginRealm      dia.Identity
	DestinationHost  dia.Identity
	DestinationRealm dia.Identity

	teldata.IMSI
	SCAddress teldata.E164
	SMSPDU    sms.Deliver
	Flags     struct {
		MMS bool
	}

	MMEAddress  teldata.E164
	SGSNAddress teldata.E164

	DeliveryTimer     uint32
	DeliveryStartTime time.Time
	MaxRetransTime    time.Time
	SMSGMSCAddress    teldata.E164
}

func (v TFR) String() string {
	w := new(bytes.Buffer)

	fmt.Fprintf(w, "%sOrigin-Host       =%s\n", dia.Indent, v.OriginHost)
	fmt.Fprintf(w, "%sOrigin-Realm      =%s\n", dia.Indent, v.OriginRealm)
	fmt.Fprintf(w, "%sDestination-Host  =%s\n", dia.Indent, v.DestinationHost)
	fmt.Fprintf(w, "%sDestination-Realm =%s\n", dia.Indent, v.DestinationRealm)

	fmt.Fprintf(w, "%sIMSI              =%s\n", dia.Indent, v.IMSI)
	fmt.Fprintf(w, "%sSC Address        =%s\n", dia.Indent, v.SCAddress)
	fmt.Fprintf(w, "%sSMS Data Unit     =%s\n", dia.Indent, v.SMSPDU.String())
	fmt.Fprintf(w, "%sMore SM to sent   =%t\n", dia.Indent, v.Flags.MMS)

	fmt.Fprintf(w, "%sMME Address       =%s\n", dia.Indent, v.MMEAddress)
	fmt.Fprintf(w, "%sSGSN Address      =%s\n", dia.Indent, v.SGSNAddress)

	fmt.Fprintf(w, "%sDelivery Timer     =%d\n", dia.Indent, v.DeliveryTimer)
	fmt.Fprintf(w, "%sDelivery Start Time=%s\n", dia.Indent, v.DeliveryStartTime)
	fmt.Fprintf(w, "%sMax Retransmit Time=%s\n", dia.Indent, v.MaxRetransTime)
	fmt.Fprintf(w, "%sSMS-GMSC Address   =%s\n", dia.Indent, v.SMSGMSCAddress)

	return w.String()
}

// ToRaw return dia.RawMsg struct of this value
func (v TFR) ToRaw(s string) dia.RawMsg {
	m := dia.RawMsg{
		Ver:  dia.DiaVer,
		FlgR: true, FlgP: true, FlgE: false, FlgT: false,
		Code: 8388646, AppID: 16777313,
		AVP: make([]dia.RawAVP, 0, 15)}

	m.AVP = append(m.AVP, dia.SetSessionID(s))
	m.AVP = append(m.AVP, dia.SetVendorSpecAppID(10415, 16777313))
	m.AVP = append(m.AVP, dia.SetAuthSessionState(false))

	m.AVP = append(m.AVP, dia.SetOriginHost(v.OriginHost))
	m.AVP = append(m.AVP, dia.SetOriginRealm(v.OriginRealm))
	if len(v.DestinationHost) != 0 {
		m.AVP = append(m.AVP, dia.SetDestinationHost(v.DestinationHost))
	}
	m.AVP = append(m.AVP, dia.SetDestinationRealm(v.DestinationRealm))

	m.AVP = append(m.AVP, setUserName(v.IMSI))
	m.AVP = append(m.AVP, setSCAddress(v.SCAddress))
	m.AVP = append(m.AVP, setSMRPUI(&v.SMSPDU))

	if v.MMEAddress.Length() != 0 {
		m.AVP = append(m.AVP, setMMENumberForMTSMS(v.MMEAddress))
	}
	if v.SGSNAddress.Length() != 0 {
		m.AVP = append(m.AVP, setSGSNNumber(v.SGSNAddress))
	}
	if v.Flags.MMS {
		m.AVP = append(m.AVP, setTFRFlags(v.Flags.MMS))
	}

	if v.DeliveryTimer != 0 {
		m.AVP = append(m.AVP, setSMDeliveryTimer(v.DeliveryTimer))
	}
	if !v.DeliveryStartTime.IsZero() {
		m.AVP = append(m.AVP, setSMDeliveryStartTime(v.DeliveryStartTime))
	}
	if !v.MaxRetransTime.IsZero() {
		m.AVP = append(m.AVP, setMaximumRetransmissionTime(v.MaxRetransTime))
		m.AVP = append(m.AVP, setSMSGMSCAddress(v.SMSGMSCAddress))
	}

	m.AVP = append(m.AVP, dia.SetRouteRecord(v.OriginHost))
	return m
}

// FromRaw make this value from dia.RawMsg struct
func (TFR) FromRaw(m dia.RawMsg) (dia.Request, string, error) {
	s := ""
	e := m.Validate(16777312, 8388646, true, true, false, false)
	if e != nil {
		return nil, s, e
	}

	v := TFR{}
	for _, a := range m.AVP {
		if a.VenID == 0 && a.Code == 263 {
			s, e = dia.GetSessionID(a)
		} else if a.VenID == 0 && a.Code == 264 {
			v.OriginHost, e = dia.GetOriginHost(a)
		} else if a.VenID == 0 && a.Code == 296 {
			v.OriginRealm, e = dia.GetOriginRealm(a)
		} else if a.VenID == 0 && a.Code == 293 {
			v.DestinationHost, e = dia.GetDestinationHost(a)
		} else if a.VenID == 0 && a.Code == 283 {
			v.DestinationRealm, e = dia.GetDestinationRealm(a)

		} else if a.VenID == 0 && a.Code == 1 {
			v.IMSI, e = getUserName(a)
		} else if a.VenID == 10415 && a.Code == 3300 {
			v.SCAddress, e = getSCAddress(a)
		} else if a.VenID == 10415 && a.Code == 3301 {
			v.SMSPDU, e = getSMRPUIasDeliver(a)
		} else if a.VenID == 10415 && a.Code == 1645 {
			v.MMEAddress, e = getMMENumberForMTSMS(a)
		} else if a.VenID == 10415 && a.Code == 1489 {
			v.SGSNAddress, e = getSGSNNumber(a)
		} else if a.VenID == 10415 && a.Code == 3302 {
			v.Flags.MMS, e = getTFRFlags(a)
		} else if a.VenID == 10415 && a.Code == 3306 {
			v.DeliveryTimer, e = getSMDeliveryTimer(a)
		} else if a.VenID == 10415 && a.Code == 3307 {
			v.DeliveryStartTime, e = getSMDeliveryStartTime(a)
		} else if a.VenID == 10415 && a.Code == 3330 {
			v.MaxRetransTime, e = getMaximumRetransmissionTime(a)
		} else if a.VenID == 10415 && a.Code == 3332 {
			v.SMSGMSCAddress, e = getSMSGMSCAddress(a)
		}

		if e != nil {
			return nil, s, e
		}
	}

	if len(v.OriginHost) == 0 || len(v.OriginRealm) == 0 ||
		len(v.DestinationHost) == 0 || len(v.DestinationRealm) == 0 ||
		v.IMSI.Length() == 0 || v.SCAddress.Length() == 0 ||
		v.SMSPDU.OA.Addr == nil {
		e = dia.NoMandatoryAVP{}
	}
	return v, s, e
}

// Failed make error message for timeout
func (v TFR) Failed(c uint32) dia.Answer {
	return TFA{
		ResultCode:  c,
		OriginHost:  dia.Host,
		OriginRealm: dia.Realm}
}

/*
TFA is MT-Forward-Short-Message-Answer message.
 <TFA> ::= < Diameter Header: 8388646, PXY, 16777313 >
           < Session-Id >
           [ DRMP ] // not supported
           [ Vendor-Specific-Application-Id ]
           [ Result-Code ]
           [ Experimental-Result ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
         * [ Supported-Features ] // not supported
           [ Absent-User-Diagnostic-SM ]
           [ SM-Delivery-Failure-Cause ]
           [ SM-RP-UI ]
           [ Requested-Retransmission-Time ]
           [ User-Identifier ] // not supported
         * [ AVP ]
         * [ Failed-AVP ]
         * [ Proxy-Info ] // not supported
         * [ Route-Record ]
*/
type TFA struct {
	ResultCode  uint32
	OriginHost  dia.Identity
	OriginRealm dia.Identity

	SMSPDU sms.DeliverReport

	AbsentUserDiag uint32
	DeliveryFailureCause
	ReqRetransTime time.Time

	FailedAVP []dia.RawAVP
}

func (v TFA) String() string {
	w := new(bytes.Buffer)

	if v.ResultCode > dia.ResultOffset {
		fmt.Fprintf(w, "%sExp-Result-Code   =%d\n", dia.Indent, v.ResultCode-dia.ResultOffset)

	} else {
		fmt.Fprintf(w, "%sResult-Code       =%d\n", dia.Indent, v.ResultCode)
	}
	fmt.Fprintf(w, "%sOrigin-Host       =%s\n", dia.Indent, v.OriginHost)
	fmt.Fprintf(w, "%sOrigin-Realm      =%s\n", dia.Indent, v.OriginRealm)

	fmt.Fprintf(w, "%sSMS Data Unit     =%s\n", dia.Indent, v.SMSPDU.String())

	fmt.Fprintf(w, "%sAbsent User Diag  =%d\n", dia.Indent, v.AbsentUserDiag)
	switch v.DeliveryFailureCause {
	case CauseMemoryCapacityExceeded:
		fmt.Fprintf(w, "%sFailure Cause     =MEMORY_CAPACITY_EXCEEDED\n", dia.Indent)
	case CauseEquipmentProtocolError:
		fmt.Fprintf(w, "%sFailure Cause     =EQUIPMENT_PROTOCOL_ERROR\n", dia.Indent)
	case CauseEquipmentNotSMEquipped:
		fmt.Fprintf(w, "%sFailure Cause     =EQUIPMENT_NOT_SM-EQUIPPED\n", dia.Indent)
	case CauseUnknownServiceCenter:
		fmt.Fprintf(w, "%sFailure Cause     =UNKNOWN_SERVICE_CENTRE\n", dia.Indent)
	case CauseSCCongestion:
		fmt.Fprintf(w, "%sFailure Cause     =SC-CONGESTION\n", dia.Indent)
	case CauseInvalidSMEAddress:
		fmt.Fprintf(w, "%sFailure Cause     =INVALID_SME-ADDRESS\n", dia.Indent)
	case CauseUserNotSCUser:
		fmt.Fprintf(w, "%sFailure Cause     =USER_NOT_SC-USER\n", dia.Indent)
	}
	fmt.Fprintf(w, "%sRequested Retrans Time=%s\n", dia.Indent, v.ReqRetransTime)

	return w.String()
}

// ToRaw return dia.RawMsg struct of this value
func (v TFA) ToRaw(s string) dia.RawMsg {
	m := dia.RawMsg{
		Ver:  dia.DiaVer,
		FlgR: false, FlgP: true, FlgE: false, FlgT: false,
		Code: 8388646, AppID: 16777313,
		AVP: make([]dia.RawAVP, 0, 20)}

	m.AVP = append(m.AVP, dia.SetSessionID(s))
	m.AVP = append(m.AVP, dia.SetVendorSpecAppID(10415, 16777312))
	if v.ResultCode > dia.ResultOffset {
		m.AVP = append(m.AVP, dia.SetExperimentalResult(10415, v.ResultCode))
	} else {
		m.AVP = append(m.AVP, dia.SetResultCode(v.ResultCode))
	}
	m.AVP = append(m.AVP, dia.SetAuthSessionState(false))
	m.AVP = append(m.AVP, dia.SetOriginHost(v.OriginHost))
	m.AVP = append(m.AVP, dia.SetOriginRealm(v.OriginRealm))

	switch v.ResultCode {
	case dia.DiameterSuccess:
		m.AVP = append(m.AVP, setSMRPUI(&v.SMSPDU))
	case DiameterErrorAbsentUser:
		m.AVP = append(m.AVP, setAbsentUserDiagnosticSM(v.AbsentUserDiag))
		if !v.ReqRetransTime.IsZero() {
			m.AVP = append(m.AVP, setRequestedRetransmissionTime(v.ReqRetransTime))
		}
	case DiameterErrorSmDeliveryFailure:
		m.AVP = append(m.AVP, setSMDeliveryFailureCause(
			v.DeliveryFailureCause, v.SMSPDU))
	}
	return m
}

// FromRaw make this value from dia.RawMsg struct
func (TFA) FromRaw(m dia.RawMsg) (dia.Answer, string, error) {
	s := ""
	e := m.Validate(16777313, 8388646, false, true, false, false)
	if e != nil {
		return nil, s, e
	}

	v := TFA{}
	for _, a := range m.AVP {
		if a.VenID == 0 && a.Code == 263 {
			s, e = dia.GetSessionID(a)
		} else if a.VenID == 0 && a.Code == 268 {
			v.ResultCode, e = dia.GetResultCode(a)
		} else if a.VenID == 0 && a.Code == 297 {
			_, v.ResultCode, e = dia.GetExperimentalResult(a)
		} else if a.VenID == 0 && a.Code == 264 {
			v.OriginHost, e = dia.GetOriginHost(a)
		} else if a.VenID == 0 && a.Code == 296 {
			v.OriginRealm, e = dia.GetOriginRealm(a)
		}
		if e != nil {
			return nil, s, e
		}
	}
	switch v.ResultCode {
	case dia.DiameterSuccess:
		for _, a := range m.AVP {
			if a.VenID == 10415 && a.Code == 3301 {
				v.SMSPDU, e = getSMRPUIasDeliverReport(a)
			}
			if e != nil {
				return nil, s, e
			}
		}
	case DiameterErrorAbsentUser:
		for _, a := range m.AVP {
			if a.VenID == 10415 && a.Code == 3322 {
				v.AbsentUserDiag, e = getAbsentUserDiagnosticSM(a)
			} else if a.VenID == 10415 && a.Code == 3331 {
				v.ReqRetransTime, e = getRequestedRetransmissionTime(a)
			}
			if e != nil {
				return nil, s, e
			}
		}
	case DiameterErrorSmDeliveryFailure:
		for _, a := range m.AVP {
			if a.VenID == 10415 && a.Code == 3303 {
				v.DeliveryFailureCause, v.SMSPDU, e = getSMDeliveryFailureCause(a)
			}
			if e != nil {
				return nil, s, e
			}
		}
	}
	return v, s, e
}

// Result returns result-code
func (v TFA) Result() uint32 {
	return v.ResultCode
}
