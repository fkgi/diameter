package ts29338

import (
	"github.com/fkgi/diameter/msg"
	"github.com/fkgi/sms"
	"github.com/fkgi/teldata"
)

/*
var (
	s6cAppID = rfc6733.VendorSpecificApplicationID{
		VendorID:          rfc6733.VendorID(10415),
		AuthApplicationID: rfc6733.AuthApplicationID(16777312)}
	sessionState = rfc6733.StateNotMaintained
)
*/

/*
SRR is Send-Routing-Info-For-SM-Request message.
 <SRR> ::= < Diameter Header: 8388647, REQ, PXY, 16777312 >
           < Session-Id >
		   [ DRMP ]  // not supported
           [ Vendor-Specific-Application-Id ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
           [ Destination-Host ]
           { Destination-Realm }
           [ MSISDN ]
           [ User-Name ]
           [ SMSMI-Correlation-ID ] // not supported
         * [ Supported-Features ] // not supported
           [ SC-Address ]
           [ SM-RP-MTI ]
           [ SM-RP-SMEA ]
           [ SRR-Flags ]
           [ SM-Delivery-Not-Intended ]
         * [ AVP ]
         * [ Proxy-Info ] // not supported
         * [ Route-Record ]
*/
type SRR struct {
	// DRMP
	OriginHost       msg.DiameterIdentity
	OriginRealm      msg.DiameterIdentity
	DestinationHost  msg.DiameterIdentity
	DestinationRealm msg.DiameterIdentity
	MSISDN           teldata.TBCD
	UserName         string
	// SMSMICorrelationID
	// []SupportedFeatures
	SCAddress             string
	SMRPMTI               bool
	SMRPSMEA              sms.Address
	GPRSIndicator         bool
	SMRPPRI               bool
	SingleAttempt         bool
	SMDeliveryNotIntended int
	// []ProxyInfo
}

/*

// ToRaw return msg.RawMsg struct of this value
func (v SRR) ToRaw() msg.RawMsg {
	m := msg.RawMsg{
		Ver:  msg.DiaVer,
		FlgR: true, FlgP: true, FlgE: false, FlgT: false,
		Code: 8388647, AppID: 16777312,
		AVP: make([]msg.RawAVP, 0, 15)}

	session := sock.NextSession()
	m.AVP = append(m.AVP, session.ToRaw())
	m.AVP = append(m.AVP, s6cAppID.ToRaw())
	m.AVP = append(m.AVP, sessionState.ToRaw())

	m.AVP = append(m.AVP, v.OriginHost.ToRaw())
	m.AVP = append(m.AVP, v.OriginRealm.ToRaw())
	if len(v.DestinationHost) != 0 {
		m.AVP = append(m.AVP, v.DestinationHost.ToRaw())
	}
	m.AVP = append(m.AVP, v.DestinationRealm.ToRaw())
	if len(v.MSISDN) != 0 {
		m.AVP = append(m.AVP, v.MSISDN.ToRaw())
	}
	if len(v.UserName) != 0 {
		m.AVP = append(m.AVP, v.UserName.ToRaw())
	}
	if len(v.SCAddress) != 0 {
		m.AVP = append(m.AVP, v.SCAddress.ToRaw())
	}
	if v.SMRPMTI != 0 {
		m.AVP = append(m.AVP, v.SMRPMTI.ToRaw())
	}
	if v.SMRPSMEA.Addr != nil {
		m.AVP = append(m.AVP, v.SMRPSMEA.ToRaw())
	}
	if v.SRRFlags.GprsIndicator ||
		v.SRRFlags.SingleAttempt ||
		v.SRRFlags.SMRPPRI {
		m.AVP = append(m.AVP, v.SRRFlags.ToRaw())
	}
	if v.SMDeliveryNotIntended != 0 {
		m.AVP = append(m.AVP, v.SMDeliveryNotIntended.ToRaw())
	}

	rt := rfc6733.RouteRecord(v.OriginHost)
	m.AVP = append(m.AVP, rt.ToRaw())
	return m
}

// FromRaw make this value from msg.RawMsg struct
func (SRR) FromRaw(m msg.RawMsg) (msg.Request, error) {
	e := m.Validate(16777312, 8388647, true, true, false, false)
	if e != nil {
		return nil, e
	}

	v := SRR{}
	for _, a := range m.AVP {
		if a.VenID == 0 && a.Code == 264 {
			e = v.OriginHost.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 296 {
			e = v.OriginRealm.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 293 {
			e = v.DestinationHost.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 283 {
			e = v.DestinationRealm.FromRaw(a)

		} else if a.VenID == 10415 && a.Code == 701 {
			e = v.MSISDN.FromRaw(a)
		} else if a.VenID == 0 && a.Code == 1 {
			e = v.UserName.FromRaw(a)
		} else if a.VenID == 10415 && a.Code == 3300 {
			e = v.SCAddress.FromRaw(a)
		} else if a.VenID == 10415 && a.Code == 3308 {
			e = v.SMRPMTI.FromRaw(a)
		} else if a.VenID == 10415 && a.Code == 3308 {
			e = v.SMRPMTI.FromRaw(a)
		} else if a.VenID == 10415 && a.Code == 3309 {
			e = v.SMRPSMEA.FromRaw(a)
		} else if a.VenID == 10415 && a.Code == 3310 {
			e = v.SRRFlags.FromRaw(a)
		} else if a.VenID == 10415 && a.Code == 3311 {
			e = v.SMDeliveryNotIntended.FromRaw(a)
		}

		if e != nil {
			return nil, e
		}
	}

	if len(v.OriginHost) == 0 ||
		len(v.OriginRealm) == 0 ||
		len(v.DestinationRealm) == 0 {
		e = msg.NoMandatoryAVP{}
	}
	return v, e
}

// Failed make error message for timeout
func (v SRR) Failed(c uint32, s string) msg.Answer {
	return SRA{
		ResultCode:    ResultCode(c),
		OriginHost:    v.OriginHost,
		OriginRealm:   v.OriginRealm,
		HostIPAddress: v.HostIPAddress,
		VendorID:      v.VendorID,
		ProductName:   v.ProductName,
		ErrorMessage:  ErrorMessage(s)}
}

/*
SRA is SendRoutingInfoForSMAnswer message.
 <SRA> ::= < Diameter Header: 8388647, PXY, 16777312 >
           < Session-Id >
		   [ DRMP ] // not supported
           [ Vendor-Specific-Application-Id ]
           [ Result-Code ]
           [ Experimental-Result ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
           [ User-Name ]
         * [ Supported-Features ] // not supported
           [ Serving-Node ]
           [ Additional-Serving-Node ]
           [ LMSI ]
           [ User-Identifier ]
           [ MWD-Status ]
           [ MME-Absent-User-Diagnostic-SM ]
           [ MSC-Absent-User-Diagnostic-SM ]
           [ SGSN-Absent-User-Diagnostic-SM ]
         * [ AVP ]
         * [ Failed-AVP ]
         * [ Proxy-Info ] // not supported
         * [ Route-Record ]
*/
type SRA struct {
	// DRMP
	ResultCode
	ExperimentalResult
	OriginHost  msg.DiameterIdentity
	OriginRealm msg.DiameterIdentity
	UserName
	// []SupportedFeatures
	ServingNode
	AdditionalServingNode
	LMSI
	UserIdentifier
	MWDStatus
	MMEAbsentUserDiagnosticSM
	MSCAbsentUserDiagnosticSM
	SGSNAbsentUserDiagnosticSM
	FailedAVP []FailedAVP
	// []ProxyInfo
}

/*
AlertServiceCentreRequest is ALR message.
 <ALR> ::= < Diameter Header: 8388648, REQ, PXY, 16777312 >
           < Session-Id >
           [ DRMP ]
           [ Vendor-Specific-Application-Id ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
           [ Destination-Host ]
           { Destination-Realm }
           { SC-Address }
           { User-Identifier }
           [ SMSMI-Correlation-ID ]
           [ Maximum-UE-Availability-Time ]
           [ SMS-GMSC-Alert-Event ]
           [ Serving-Node ]
         * [ Supported-Features ]
         * [ AVP ]
         * [ Proxy-Info ]
         * [ Route-Record ]
*/
/*
 AlertServiceCentreAnswer is ALA message.
 <ALA> ::= < Diameter Header: 8388648, PXY, 16777312 >
           < Session-Id >
           [ DRMP ]
           [ Vendor-Specific-Application-Id ]
           [ Result-Code ]
           [ Experimental-Result ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
         * [ Supported-Features ]
         * [ AVP ]
         * [ Failed-AVP ]
         * [ Proxy-Info ]
         * [ Route-Record ]
*/

/*
ReportSMDeliveryStatusRequest is RDR message.
 <RDR> ::= < Diameter Header: 8388649, REQ, PXY, 16777312 >
           < Session-Id >
           [ DRMP ]
           [ Vendor-Specific-Application-Id ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
           [ Destination-Host ]
           { Destination-Realm }
         * [ Supported-Features ]
           { User-Identifier }
           [ SMSMI-Correlation-ID ]
           { SC-Address }
           { SM-Delivery-Outcome }
           [ RDR-Flags ]
         * [ AVP ]
         * [ Proxy-Info ]
         * [ Route-Record ]
*/
/*
ReportSMDeliveryStatusAnswer is RDA message.
 <RDA> ::= < Diameter Header: 8388649, PXY, 16777312 >
           < Session-Id >
           [ DRMP ]
           [ Vendor-Specific-Application-Id ]
           [ Result-Code ]
           [ Experimental-Result ]
           { Auth-Session-State }
           { Origin-Host }
           { Origin-Realm }
         * [ Supported-Features ]
           [ User-Identifier ]
         * [ AVP ]
         * [ Failed-AVP ]
         * [ Proxy-Info ]
         * [ Route-Record ]
*/
