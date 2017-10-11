package ts29338

import (
	"time"

	"github.com/fkgi/diameter/msg"
)

// SMRPMTI AVP contain the RP-Message Type Indicator of the Short Message.
type SMRPMTI msg.Enumerated

const (
	// SmDeliver is Enumerated value 0
	SmDeliver SMRPMTI = 0
	// SmStatusReport is Enumerated value 1
	SmStatusReport SMRPMTI = 1
)

// ToRaw return AVP struct of this value
func (v *SMRPMTI) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 3308, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(msg.Enumerated(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *SMRPMTI) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(10415, 3308, true, true, false); e != nil {
		return
	}
	s := new(msg.Enumerated)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = SMRPMTI(*s)
	return
}

// SMRPSMEA AVP contain the RP-Originating SME-address of
// the Short Message Entity that has originated the SM.
// It shall be formatted according to the formatting rules of
// the address fields described in 3GPP TS 23.040.
type SMRPSMEA []byte

// ToRaw return AVP struct of this value
func (v *SMRPSMEA) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 3309, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode([]byte(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *SMRPSMEA) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(10415, 3309, true, true, false); e != nil {
		return
	}
	s := new([]byte)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = SMRPSMEA(*s)
	return
}

// SRRFlags AVP contain a bit mask.
// GprsIndicator shall be ture if the SMS-GMSC supports receiving
// of two serving nodes addresses from the HSS.
// SmRpPri shall be true if the delivery of the short message shall
// be attempted when a service centre address is already contained
// in the Message Waiting Data file.
// SingleAttempt if true indicates that only one delivery attempt
// shall be performed for this particular SM.
type SRRFlags struct {
	GprsIndicator bool
	SMRPPRI       bool
	SingleAttempt bool
}

// ToRaw return AVP struct of this value
func (v *SRRFlags) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 3310, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}

	if v != nil {
		i := uint32(0)
		if v.GprsIndicator {
			i = i | 0x00000001
		}
		if v.SMRPPRI {
			i = i | 0x00000002
		}
		if v.SingleAttempt {
			i = i | 0x00000004
		}
		a.Encode(i)
	}
	return a
}

// FromRaw get AVP value
func (v *SRRFlags) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(10415, 3310, true, true, false); e != nil {
		return
	}
	s := new(uint32)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = SRRFlags{
		GprsIndicator: (*s)&0x00000001 == 0x00000001,
		SMRPPRI:       (*s)&0x00000002 == 0x00000002,
		SingleAttempt: (*s)&0x00000004 == 0x00000004}
	return
}

// SMDeliveryNotIntended AVP indicate by its presence
// that delivery of a short message is not intended.
type SMDeliveryNotIntended msg.Enumerated

const (
	// OnlyImsiRequested is Enumerated value 0
	OnlyImsiRequested SMDeliveryNotIntended = 0
	// OnlyMccMncRequested is Enumerated value 1
	OnlyMccMncRequested SMDeliveryNotIntended = 1
)

// ToRaw return AVP struct of this value
func (v *SMDeliveryNotIntended) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 3311, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(msg.Enumerated(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *SMDeliveryNotIntended) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(10415, 3311, true, true, false); e != nil {
		return
	}
	s := new(msg.Enumerated)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = SMDeliveryNotIntended(*s)
	return
}

// MWDStatus AVP contain a bit mask.
// SCAddrNotIncluded shall indicate the presence of
// the SC Address in the Message Waiting Data in the HSS.
// MNRFSet shall indicate that the MNRF flag is set in the HSS.
// MCEFSet shall indicate that the MCEF flag is set in the HSS.
// MNRGSet shall indicate that the MNRG flag is set in the HSS.
type MWDStatus struct {
	SCAddrNotIncluded bool
	MNRFSet           bool
	MCEFSet           bool
	MNRGSet           bool
}

// ToRaw return AVP struct of this value
func (v *MWDStatus) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 3312, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}

	if v != nil {
		i := uint32(0)
		if v.SCAddrNotIncluded {
			i = i | 0x00000001
		}
		if v.MNRFSet {
			i = i | 0x00000002
		}
		if v.MCEFSet {
			i = i | 0x00000004
		}
		if v.MNRGSet {
			i = i | 0x00000008
		}
		a.Encode(i)
	}
	return a
}

// FromRaw get AVP value
func (v *MWDStatus) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(10415, 3312, true, true, false); e != nil {
		return
	}
	s := new(uint32)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = MWDStatus{
		SCAddrNotIncluded: (*s)&0x00000001 == 0x00000001,
		MNRFSet:           (*s)&0x00000002 == 0x00000002,
		MCEFSet:           (*s)&0x00000004 == 0x00000004,
		MNRGSet:           (*s)&0x00000008 == 0x00000008}
	return
}

// MMEAbsentUserDiagnosticSM AVP shall indicate the diagnostic
// explaining the absence of the user given by the MME.
type MMEAbsentUserDiagnosticSM uint32

// ToRaw return AVP struct of this value
func (v *MMEAbsentUserDiagnosticSM) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 3313, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(uint32(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *MMEAbsentUserDiagnosticSM) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(10415, 3313, true, true, false); e != nil {
		return
	}
	s := new(uint32)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = MMEAbsentUserDiagnosticSM(*s)
	return
}

// MSCAbsentUserDiagnosticSM AVP shall indicate the diagnostic
// explaining the absence of the user given by the MSC.
type MSCAbsentUserDiagnosticSM uint32

// ToRaw return AVP struct of this value
func (v *MSCAbsentUserDiagnosticSM) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 3314, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(uint32(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *MSCAbsentUserDiagnosticSM) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(10415, 3314, true, true, false); e != nil {
		return
	}
	s := new(uint32)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = MSCAbsentUserDiagnosticSM(*s)
	return
}

// SGSNAbsentUserDiagnosticSM AVP shall indicate the diagnostic
// explaining the absence of the user given by the SGSN.
type SGSNAbsentUserDiagnosticSM uint32

// ToRaw return AVP struct of this value
func (v *SGSNAbsentUserDiagnosticSM) ToRaw() msg.RawAVP {
	a := msg.RawAVP{Code: 3315, VenID: 10415,
		FlgV: true, FlgM: true, FlgP: false}
	if v != nil {
		a.Encode(uint32(*v))
	}
	return a
}

// FromRaw get AVP value
func (v *SGSNAbsentUserDiagnosticSM) FromRaw(a msg.RawAVP) (e error) {
	if e = a.Validate(10415, 3315, true, true, false); e != nil {
		return
	}
	s := new(uint32)
	if e = a.Decode(s); e != nil {
		return
	}
	*v = SGSNAbsentUserDiagnosticSM(*s)
	return
}

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
