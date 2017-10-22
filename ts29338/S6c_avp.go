package ts29338

import (
	"bytes"
	"encoding/binary"

	"github.com/fkgi/diameter/msg"
	"github.com/fkgi/sms"
	"github.com/fkgi/teldata"
)

func setMSISDN(v teldata.E164) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 701, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(v.Bytes())
	return
}

func getMSISDN(a msg.RawAVP) (v teldata.E164, e error) {
	s := new([]byte)
	if e = a.Validate(10415, 701, true, true, false); e != nil {
	} else if e = a.Decode(s); e == nil {
		v, e = teldata.ToE164(*s)
	}
	return
}

// SM-RP-MTI AVP contain the RP-Message Type Indicator of the Short Message.
func setSMRPMTI(v bool) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 3308, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	if v {
		// SM_DELIVER (0)
		a.Encode(msg.Enumerated(0))
	} else {
		// SM_STATUS_REPORT (1)
		a.Encode(msg.Enumerated(1))
	}
	return
}

func getSMRPMTI(a msg.RawAVP) (v bool, e error) {
	s := new(msg.Enumerated)
	if e = a.Validate(10415, 3308, true, true, false); e != nil {
	} else if e = a.Decode(s); e != nil {
	} else if *s == 0 {
		v = true
	} else if *s == 1 {
		v = false
	} else {
		e = msg.InvalidAVP{}
	}
	return
}

// SM-RP-SMEA AVP contain the RP-Originating SME-address of
// the Short Message Entity that has originated the SM.
// It shall be formatted according to the formatting rules of
// the address fields described in 3GPP TS 23.040.
func setSMRPSMEA(v sms.Address) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 3309, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	_, b := v.Encode()
	a.Encode(b)
	return
}

func getSMRPSMEA(a msg.RawAVP) (v sms.Address, e error) {
	s := new([]byte)
	if e = a.Validate(10415, 3309, true, true, false); e != nil {
	} else if e = a.Decode(s); e == nil {
		v.Decode(byte(len(*s)-3)*2, *s)
	}
	return
}

// SRR-Flags AVP contain a bit mask.
// GPRS-Indicator shall be ture if the SMS-GMSC supports receiving
// of two serving nodes addresses from the HSS.
// SM-RP-PRI shall be true if the delivery of the short message shall
// be attempted when a service centre address is already contained
// in the Message Waiting Data file.
// Single-Attempt if true indicates that only one delivery attempt
// shall be performed for this particular SM.
func setSRRFlags(g, p, s bool) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 3310, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	i := uint32(0)
	if g {
		i = i | 0x00000001
	}
	if p {
		i = i | 0x00000002
	}
	if s {
		i = i | 0x00000004
	}
	a.Encode(i)
	return
}

func getSRRFlags(a msg.RawAVP) (g, p, s bool, e error) {
	v := new(uint32)
	if e = a.Validate(10415, 3310, true, true, false); e != nil {
	} else if e = a.Decode(v); e == nil {
		g = (*v)&0x00000001 == 0x00000001
		p = (*v)&0x00000002 == 0x00000002
		s = (*v)&0x00000004 == 0x00000004
	}
	return
}

// SM-Delivery-Not-Intended AVP indicate by its presence
// that delivery of a short message is not intended.
const (
	// LocationRequested is no SM-Delivery-Not-Intended
	LocationRequested = iota
	// OnlyImsiRequested is ONLY_IMSI_REQUESTED
	OnlyImsiRequested
	// OnlyMccMncRequested is ONLY_MCC_MNC_REQUESTED
	OnlyMccMncRequested
)

func setSMDeliveryNotIntended(v int) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 3311, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(msg.Enumerated(v - 1))
	return
}

func getSMDeliveryNotIntended(a msg.RawAVP) (v int, e error) {
	s := new(msg.Enumerated)
	if e = a.Validate(10415, 3311, true, true, false); e != nil {
	} else if e = a.Decode(s); e != nil {
	} else if *s == 0 {
		v = 1
	} else if *s == 1 {
		v = 2
	} else {
		e = msg.InvalidAVP{}
	}
	return
}

// ServingNode AVP shall contain the information about
// the network node serving the targeted SMS user.

// NodeType is tyoe of serving node address
type NodeType uint8

const (
	// NodeSGSN present SGSN
	NodeSGSN NodeType = iota
	// NodeMME present MME
	NodeMME
	// NodeMSC present MSC
	NodeMSC
	// NodeIPSMGW present IP-SM-GW
	// NodeIPSMGW
)

func setServingNode(
	t NodeType, d teldata.E164, n, r msg.DiameterIdentity) msg.RawAVP {
	return setSN(2401, t, d, n, r)
}

func setAdditionalServingNode(
	t NodeType, d teldata.E164, n, r msg.DiameterIdentity) msg.RawAVP {
	return setSN(2406, t, d, n, r)
}

func setSN(c uint32, t NodeType,
	d teldata.E164, n, r msg.DiameterIdentity) (a msg.RawAVP) {
	a = msg.RawAVP{Code: c, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	v := make([]msg.RawAVP, 1, 3)
	switch t {
	case NodeSGSN:
		v[0] = msg.RawAVP{Code: 1489, VenID: 10415,
			FlgV: true, FlgM: true, FlgP: false}
		v[0].Encode(d.Bytes())
		if len(n) != 0 {
			nv := msg.RawAVP{Code: 2409, VenID: 10415,
				FlgV: true, FlgM: false, FlgP: false}
			nv.Encode(n)
			v = append(v, nv)
		}
		if len(r) != 0 {
			rv := msg.RawAVP{Code: 2410, VenID: 10415,
				FlgV: true, FlgM: false, FlgP: false}
			rv.Encode(r)
			v = append(v, rv)
		}
	case NodeMME:
		v[0] = msg.RawAVP{Code: 1645, VenID: 10415,
			FlgV: true, FlgM: true, FlgP: false}
		v[0].Encode(d.Bytes())
		if len(n) != 0 {
			nv := msg.RawAVP{Code: 2402, VenID: 10415,
				FlgV: true, FlgM: true, FlgP: false}
			nv.Encode(n)
			v = append(v, nv)
		}
		if len(r) != 0 {
			rv := msg.RawAVP{Code: 2408, VenID: 10415,
				FlgV: true, FlgM: true, FlgP: false}
			rv.Encode(r)
			v = append(v, rv)
		}
	case NodeMSC:
		v[0] = msg.RawAVP{Code: 2403, VenID: 10415,
			FlgV: true, FlgM: true, FlgP: false}
		v[0].Encode(d.Bytes())
		if len(n) != 0 {
			nv := msg.RawAVP{Code: 2402, VenID: 10415,
				FlgV: true, FlgM: true, FlgP: false}
			nv.Encode(n)
			v = append(v, nv)
		}
		if len(r) != 0 {
			rv := msg.RawAVP{Code: 2408, VenID: 10415,
				FlgV: true, FlgM: true, FlgP: false}
			rv.Encode(r)
			v = append(v, rv)
		}
	}

	a.Encode(v)
	return
}

func getServingNode(a msg.RawAVP) (
	t NodeType, d teldata.E164, n, r msg.DiameterIdentity, e error) {
	return getSN(2401, a)
}

func getAdditionalServingNode(a msg.RawAVP) (
	t NodeType, d teldata.E164, n, r msg.DiameterIdentity, e error) {
	return getSN(2406, a)
}

func getSN(c uint32, a msg.RawAVP) (
	t NodeType, d teldata.E164, n, r msg.DiameterIdentity, e error) {
	o := []msg.RawAVP{}
	if e = a.Validate(10415, c, true, true, false); e == nil {
		e = a.Decode(&o)
	}
	for _, a := range o {
		if a.VenID != 10415 {
			continue
		}
		switch a.Code {
		case 1489:
			if e = a.Validate(1489, 10415, true, true, false); e == nil {
				e = a.Decode(&d)
				t = NodeSGSN
			}
		case 1645:
			if e = a.Validate(1645, 10415, true, true, false); e == nil {
				e = a.Decode(&d)
				t = NodeMME
			}
		case 2403:
			if e = a.Validate(2403, 10415, true, true, false); e == nil {
				e = a.Decode(&d)
				t = NodeMSC
			}
		}
	}
	switch t {
	case NodeSGSN:
		for _, a := range o {
			if a.VenID != 10415 {
				continue
			}
			switch a.Code {
			case 2409:
				if e = a.Validate(2409, 10415, true, true, false); e == nil {
					e = a.Decode(&n)
				}
			case 2410:
				if e = a.Validate(2410, 10415, true, true, false); e == nil {
					e = a.Decode(&r)
				}
			}
		}
	case NodeMME, NodeMSC:
		for _, a := range o {
			if a.VenID != 10415 {
				continue
			}
			switch a.Code {
			case 2402:
				if e = a.Validate(2402, 10415, true, true, false); e == nil {
					e = a.Decode(&n)
				}
			case 2408:
				if e = a.Validate(2408, 10415, true, true, false); e == nil {
					e = a.Decode(&r)
				}
			}
		}
	}
	return
}

func setLMSI(v uint32) msg.RawAVP {
	a := msg.RawAVP{Code: 2400, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, v)
	a.Encode(buf.Bytes())
	return a
}

func getLMSI(a msg.RawAVP) (v uint32, e error) {
	s := new([]byte)
	if e = a.Validate(10415, 2400, true, true, false); e != nil {
	} else if e = a.Decode(s); e == nil && len(*s) == 4 {
		binary.Read(bytes.NewBuffer(*s), binary.BigEndian, &v)
	}
	return
}

func setUserIdentifier(v teldata.E164) msg.RawAVP {
	a := msg.RawAVP{Code: 3102, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	t := []msg.RawAVP{
		msg.RawAVP{Code: 701, VenID: 10415,
			FlgV: true, FlgM: true, FlgP: false}}
	t[0].Encode(v.Bytes())
	a.Encode(t)
	return a
}

func getUserIdentifier(a msg.RawAVP) (v teldata.E164, e error) {
	o := []msg.RawAVP{}
	if e = a.Validate(10415, 3102, true, true, false); e != nil {
	} else if e = a.Decode(&o); e == nil {
		for _, a := range o {
			if a.Code == 701 && a.VenID == 10415 {
				if e = a.Validate(10415, 701, true, true, false); e == nil {
					a.Decode(&v)
				}
			}
		}
	}
	return
}

// MWD-Status AVP contain a bit mask.
// SCAddrNotIncluded shall indicate the presence of
// the SC Address in the Message Waiting Data in the HSS.
// MNRF shall indicate that the MNRF flag is set in the HSS.
// MCEF shall indicate that the MCEF flag is set in the HSS.
// MNRG shall indicate that the MNRG flag is set in the HSS.
func setMWDStatus(sca, mnrf, mcef, mnrg bool) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 3312, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}

	i := uint32(0)
	if sca {
		i = i | 0x00000001
	}
	if mnrf {
		i = i | 0x00000002
	}
	if mcef {
		i = i | 0x00000004
	}
	if mnrg {
		i = i | 0x00000008
	}
	a.Encode(i)
	return
}

func getMWDStatus(a msg.RawAVP) (sca, mnrf, mcef, mnrg bool, e error) {
	s := new(uint32)
	if e = a.Validate(10415, 3312, true, true, false); e != nil {
	} else if e = a.Decode(s); e == nil {
		sca = (*s)&0x00000001 == 0x00000001
		mnrf = (*s)&0x00000002 == 0x00000002
		mcef = (*s)&0x00000004 == 0x00000004
		mnrg = (*s)&0x00000008 == 0x00000008
	}
	return
}

/*
0	-	no paging response via the MSC
1	-	IMSI detached
2	-	roaming restriction
3	-	deregistered in the HLR for non GPRS
4	-	MS purged for non GPRS
5	-	no paging response via the SGSN
6	-	GPRS detached
7	-	deregistered in the HLR for GPRS
8	-	MS purged for GPRS
9	-	Unidentified subscriber via the MSC
10	-	Unidentified subscriber via the SGSN
11	-	deregistered in the HSS/HLR for IMS
12	-	no response via the IP-SM-GW
13		the MS is temporarily unavailable
*/

// MMEAbsentUserDiagnosticSM AVP shall indicate the diagnostic
// explaining the absence of the user given by the MME.
func setMMEAbsentUserDiagnosticSM(v uint32) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 3313, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getMMEAbsentUserDiagnosticSM(a msg.RawAVP) (v uint32, e error) {
	if e = a.Validate(10415, 3313, true, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

// MSCAbsentUserDiagnosticSM AVP shall indicate the diagnostic
// explaining the absence of the user given by the MSC.
func setMSCAbsentUserDiagnosticSM(v uint32) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 3314, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getMSCAbsentUserDiagnosticSM(a msg.RawAVP) (v uint32, e error) {
	if e = a.Validate(10415, 3314, true, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

// SGSNAbsentUserDiagnosticSM AVP shall indicate the diagnostic
// explaining the absence of the user given by the SGSN.
func setSGSNAbsentUserDiagnosticSM(v uint32) (a msg.RawAVP) {
	a = msg.RawAVP{Code: 3315, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	a.Encode(v)
	return
}

func getSGSNAbsentUserDiagnosticSM(a msg.RawAVP) (v uint32, e error) {
	if e = a.Validate(10415, 3315, true, true, false); e == nil {
		e = a.Decode(&v)
	}
	return
}

/*
// SMDeliveryOutcome AVP contains the result of the SM delivery.
type SMDeliveryOutcome struct {
	E msg.Enumerated
	I uint32
}

// ToRaw return AVP struct of this value
func (v *SMDeliveryOutcome) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 3316, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	// a.Encode()
	return a
}

// MMESMDeliveryOutcome AVP shall indicate the outcome of
// the SM delivery for setting the message waiting data
// in the HSS when the SM delivery is with an MME.
type MMESMDeliveryOutcome struct {
	SMDeliveryCause        SMDeliveryCause
	AbsentUserDiagnosticSM AbsentUserDiagnosticSM
}

// ToRaw return AVP struct of this value
func (v *MMESMDeliveryOutcome) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 3317, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	if v != nil {
		t := []msg.RawAVP{
			v.SMDeliveryCause.ToRaw(),
			v.AbsentUserDiagnosticSM.ToRaw()}
		a.Encode(msg.GroupedAVP(t))
	}
	return a
}

// MSCSMDeliveryOutcome AVP shall indicate the outcome of
// the SM delivery for setting the message waiting data
// in the HSS when the SM delivery is with an MSC.
type MSCSMDeliveryOutcome struct {
	SMDeliveryCause        SMDeliveryCause
	AbsentUserDiagnosticSM AbsentUserDiagnosticSM
}

// ToRaw return AVP struct of this value
func (v *MSCSMDeliveryOutcome) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 3318, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	if v != nil {
		t := []msg.RawAVP{
			v.SMDeliveryCause.ToRaw(),
			v.AbsentUserDiagnosticSM.ToRaw()}
		a.Encode(msg.GroupedAVP(t))
	}
	return a
}

// SGSNSMDeliveryOutcome AVP shall indicate the outcome of
// the SM delivery for setting the message waiting data
// in the HSS when the SM delivery is with an SGSN.
type SGSNSMDeliveryOutcome struct {
	SMDeliveryCause        SMDeliveryCause
	AbsentUserDiagnosticSM AbsentUserDiagnosticSM
}

// ToRaw return AVP struct of this value
func (v *SGSNSMDeliveryOutcome) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 3319, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	if v != nil {
		t := []msg.RawAVP{
			v.SMDeliveryCause.ToRaw(),
			v.AbsentUserDiagnosticSM.ToRaw()}
		a.Encode(msg.GroupedAVP(t))
	}
	return a
}

// IPSMGWSMDeliveryOutcome AVP shall indicate the outcome of
// the SM delivery for setting the message waiting data
// when the SM delivery is with an IP-SM-GW.
type IPSMGWSMDeliveryOutcome struct {
	SMDeliveryCause        SMDeliveryCause
	AbsentUserDiagnosticSM AbsentUserDiagnosticSM
}

// ToRaw return AVP struct of this value
func (v *IPSMGWSMDeliveryOutcome) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 3320, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	if v != nil {
		t := []msg.RawAVP{
			v.SMDeliveryCause.ToRaw(),
			v.AbsentUserDiagnosticSM.ToRaw()}
		a.Encode(msg.GroupedAVP(t))
	}
	return a
}

// SMDeliveryCause AVP shall indicate the cause of
// the SMP delivery result.
type SMDeliveryCause msg.Enumerated

const (
	// UeMemoryCapacityExceeded is Enumerated value 0
	UeMemoryCapacityExceeded SMDeliveryCause = 0
	// AbsentUser is Enumerated value 1
	AbsentUser SMDeliveryCause = 1
	// SuccessfulTransfer is Enumerated value 2
	SuccessfulTransfer SMDeliveryCause = 2
)

// ToRaw return AVP struct of this value
func (v *SMDeliveryCause) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 3321, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(msg.Enumerated(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *SMDeliveryCause) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(10415, 3321, true, true, false); e != nil {
		return
	}
	s := new(msg.Enumerated)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = SMDeliveryCause(*s)
	return
}

// AbsentUserDiagnosticSM AVP shall indicate the diagnostic
// explaining the absence of the subscriber.
type AbsentUserDiagnosticSM uint32

// ToRaw return AVP struct of this value
func (v *AbsentUserDiagnosticSM) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 3322, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(uint32(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *AbsentUserDiagnosticSM) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(10415, 3322, true, true, false); e != nil {
		return
	}
	s := new(uint32)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = AbsentUserDiagnosticSM(*s)
	return
}

// RDRFlags AVP contain a bit mask.
// SingleAttemptDelivery indicates that only one delivery attempt
// shall be performed for this particular SM.
type RDRFlags struct {
	SingleAttemptDelivery bool
}

// ToRaw return AVP struct of this value
func (v *RDRFlags) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 3323, VenID: 10415,
		FlgV: true, FlgM: false, FlgP: false}
	if v != nil {
		i := uint32(0)
		if v.SingleAttemptDelivery {
			i = i | 0x00000001
		}
		a.Encode(i)
	}
	return a
}

// FromRaw get AVP value
func (v *RDRFlags) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(10415, 3323, true, false, false); e != nil {
		return
	}
	s := new(uint32)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = RDRFlags{
		SingleAttemptDelivery: (*s)&0x00000001 == 0x00000001}
	return
}

// MaximumUEAvailabilityTime AVP shall contain the timestamp (in UTC)
// until which a UE using a power saving mechanism
// (such as extended idle mode DRX) is expected to be reachable
// for SM Delivery.
type MaximumUEAvailabilityTime time.Time

// ToRaw return AVP struct of this value
func (v *MaximumUEAvailabilityTime) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 3329, VenID: 10415,
		FlgV: true, FlgM: false, FlgP: false}
	if v != nil {
		a.Encode(time.Time(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *MaximumUEAvailabilityTime) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(10415, 3329, true, false, false); e != nil {
		return
	}
	s := new(time.Time)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = MaximumUEAvailabilityTime(*s)
	return
}

// SMSGMSCAlertEvent AVP shall contain a bit mask.
// UEAvailableForMTSMS shall indicate that the UE is
// now available for MT SMS
// UEUnderNewServingNode shall indicate that the UE has moved
// under the coverage of another MME or SGSN.
type SMSGMSCAlertEvent struct {
	UEAvailableForMTSMS   bool
	UEUnderNewServingNode bool
}

// ToRaw return AVP struct of this value
func (v *SMSGMSCAlertEvent) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 3333, VenID: 10415,
		FlgV: true, FlgM: false, FlgP: false}
	if v != nil {
		i := uint32(0)
		if v.UEAvailableForMTSMS {
			i = i | 0x00000001
		}
		if v.UEUnderNewServingNode {
			i = i | 0x00000002
		}
		a.Encode(i)
	}
	return a
}

// FromRaw get AVP value
func (v *SMSGMSCAlertEvent) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(10415, 3333, true, false, false); e != nil {
		return
	}
	s := new(uint32)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = SMSGMSCAlertEvent{
		UEAvailableForMTSMS:   (*s)&0x00000001 == 0x00000001,
		UEUnderNewServingNode: (*s)&0x00000002 == 0x00000002}
	return
}
*/
