package diameter

import (
	"bytes"
	"net"
	"strings"
	"time"
)

/*
Capabilities-Exchange-Request message
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
		 * [ AVP ] // ignore any other AVP

Capabilities-Exchange-Answer message
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
		 * [ AVP ] // ignore any other AVP
*/

// RcvCER
type eventRcvCER struct {
	m Message
}

func (eventRcvCER) String() string {
	return "Rcv-CER"
}

func (v eventRcvCER) exec() error {
	RxReq++
	if state != waitCER {
		RejectReq++
		return notAcceptableEvent{e: v, s: state}
	}
	TraceMessage(v.m, Rx, nil)

	var oHost Identity
	var oRealm Identity
	var hostIP = make([]net.IP, 0, 2)
	var venID uint32
	var prodName string
	var oState uint32
	var suuportVendor = []uint32{}
	var authApps = make(map[uint32]uint32)
	// var firmwareRevision uint32
	var err error

	for rdr := bytes.NewReader(v.m.AVPs); rdr.Len() != 0; {
		a := AVP{}
		if err = a.wrapedUnmarshalFrom(rdr); err != nil {
			break
		}
		if a.VendorID != 0 {
			if a.Mandatory {
				err = InvalidAVP{Code: AvpUnsupported, AVP: a}
				break
			}
			continue
		}

		switch a.Code {
		case 264:
			if len(oHost) != 0 {
				err = InvalidAVP{Code: AvpOccursTooManyTimes, AVP: a}
			} else {
				oHost, err = GetOriginHost(a)
			}
		case 296:
			if len(oRealm) != 0 {
				err = InvalidAVP{Code: AvpOccursTooManyTimes, AVP: a}
			} else {
				oRealm, err = GetOriginRealm(a)
			}
		case 257:
			if ip, e := getHostIPAddress(a); e == nil {
				hostIP = append(hostIP, ip)
			} else {
				err = e
			}
		case 266:
			if venID != 0 {
				err = InvalidAVP{Code: AvpOccursTooManyTimes, AVP: a}
			} else {
				venID, err = GetVendorID(a)
			}
		case 269:
			if len(prodName) != 0 {
				err = InvalidAVP{Code: AvpOccursTooManyTimes, AVP: a}
			} else {
				prodName, err = getProductName(a)
			}
		case 265:
			var vid uint32
			if vid, err = getSupportedVendorID(a); err == nil {
				suuportVendor = append(suuportVendor, vid)
			}
		case 258:
			var aid uint32
			if aid, err = GetAuthAppID(a); err == nil {
				authApps[aid] = 0
			}
		case 260:
			var vid, aid uint32
			if vid, aid, err = GetVendorSpecAppID(a); err == nil {
				authApps[aid] = vid
			}
		case 278:
			if oState != 0 {
				err = InvalidAVP{Code: AvpOccursTooManyTimes, AVP: a}
			} else {
				oState, err = getOriginStateID(a)
			}
		// case 267:
		//	firmwareRevision, e = getFirmwareRevision(a)
		default:
			if a.Mandatory {
				err = InvalidAVP{Code: AvpUnsupported, AVP: a}
			}
		}

		if err != nil {
			break
		}
	}

	if err == nil && len(authApps) == 0 {
		err = InvalidAVP{Code: MissingAvp, AVP: SetAuthAppID(0)}
	}
	if err == nil {
		for _, vid := range authApps {
			if vid == 0 {
				continue
			}
			err = InvalidAVP{Code: MissingAvp, AVP: SetVendorID(0)}
			for _, v := range suuportVendor {
				if vid == v {
					err = nil
					break
				}
			}
			if err != nil {
				break
			}
		}
	}
	if err == nil {
		if _, ok := authApps[0xffffffff]; ok {
		} else if len(applications) == 0 {
			for aid, vid := range authApps {
				applications[aid] = application{
					venID:    vid,
					handlers: make(map[uint32]func(bool, []byte) (bool, []byte))}
			}
		} else {
			var commonApp = make(map[uint32]application)
			for laid, lapp := range applications {
				if pvid, ok := authApps[laid]; ok && pvid == lapp.venID {
					commonApp[laid] = lapp
					break
				}
			}
			if len(commonApp) == 0 {
				err = InvalidMessage(ApplicationUnsupported)
			} else {
				applications = commonApp
			}
		}
	}

	result := Success
	if v.m.FlgP || v.m.FlgT {
		result = InvalidHdrBits
		err = InvalidMessage(result)
	} else if iavp, ok := err.(InvalidAVP); ok {
		result = iavp.Code
	} else if imsg, ok := err.(InvalidMessage); ok {
		result = uint32(imsg)
	} else if len(oHost) == 0 {
		result = MissingAvp
		err = InvalidAVP{Code: result, AVP: SetOriginHost("")}
	} else if len(oRealm) == 0 {
		result = MissingAvp
		err = InvalidAVP{Code: result, AVP: SetOriginRealm("")}
	} else if Peer.Host != "" && oHost != Peer.Host {
		result = UnknownPeer
		err = InvalidMessage(result)
	} else if Peer.Realm != "" && oRealm != Peer.Realm {
		result = UnknownPeer
		err = InvalidMessage(result)
	} else if len(hostIP) == 0 {
		result = MissingAvp
		err = InvalidAVP{Code: result, AVP: setHostIPAddress(net.IPv4zero)}
	} else if venID == 0 {
		result = MissingAvp
		err = InvalidAVP{Code: result, AVP: SetVendorID(0)}
	} else if len(prodName) == 0 {
		result = MissingAvp
		err = InvalidAVP{Code: result, AVP: setProductName("")}
	} else {
		if oState != 0 {
			Peer.state = oState
		}
		Peer.Host = oHost
		Peer.Realm = oRealm
	}

	buf := new(bytes.Buffer)
	SetResultCode(result).MarshalTo(buf)
	SetOriginHost(Local.Host).MarshalTo(buf)
	SetOriginRealm(Local.Realm).MarshalTo(buf)

	if len(OverwriteAddr) != 0 {
		for _, h := range OverwriteAddr {
			setHostIPAddress(h).MarshalTo(buf)
		}
	} else {
		h, _, _ := net.SplitHostPort(conn.LocalAddr().String())
		for _, h := range strings.Split(h, "/") {
			setHostIPAddress(net.ParseIP(h)).MarshalTo(buf)
		}
	}

	SetVendorID(VendorID).MarshalTo(buf)
	setProductName(ProductName).MarshalTo(buf)
	if Local.state != 0 {
		setOriginStateID(Local.state).MarshalTo(buf)
	}
	if len(applications) == 0 {
		SetAuthAppID(0xffffffff).MarshalTo(buf)
	} else {
		vmap := make(map[uint32]interface{})
		for aid, app := range applications {
			if app.venID == 0 {
				SetAuthAppID(aid).MarshalTo(buf)
			} else if _, ok := vmap[app.venID]; ok {
				SetVendorSpecAppID(app.venID, aid).MarshalTo(buf)
			} else {
				setSupportedVendorID(app.venID).MarshalTo(buf)
				SetVendorSpecAppID(app.venID, aid).MarshalTo(buf)
				vmap[app.venID] = nil
			}
		}
	}
	if iavp, ok := err.(InvalidAVP); ok {
		setFailedAVP([]AVP{iavp.AVP}).MarshalTo(buf)
	}
	setFirmwareRevision(FirmwareRev).MarshalTo(buf)

	cea := Message{
		FlgR: false, FlgP: false, FlgE: result != Success, FlgT: false,
		Code: 257, AppID: 0,
		HbHID: v.m.HbHID, EtEID: v.m.EtEID,
		AVPs: buf.Bytes()}

	conn.SetWriteDeadline(time.Now().Add(TxTimeout))
	if e := cea.MarshalTo(conn); e != nil {
		TxAnsFail++
		conn.Close()
		err = e
	} else if err == nil {
		countTxCode(result)
		state = open
		// wdTimer.Stop()
		wdTimer = time.AfterFunc(WDInterval, func() {
			notify <- eventWatchdog{}
		})
	} else {
		countTxCode(result)
	}

	TraceMessage(cea, Tx, err)
	return err
}

// RcvCEA
type eventRcvCEA struct {
	m Message
}

func (eventRcvCEA) String() string {
	return "Rcv-CEA"
}

func (v eventRcvCEA) exec() error {
	// verify diameter header
	if v.m.FlgP {
		InvalidAns++
		return InvalidMessage(InvalidHdrBits)
	}
	if state != waitCEA {
		InvalidAns++
		return notAcceptableEvent{e: v, s: state}
	}
	if _, ok := sndStack[v.m.HbHID]; !ok {
		InvalidAns++
		return unknownAnswer(v.m.HbHID)
	}

	// verify diameter AVP
	var result uint32
	var oHost Identity
	var oRealm Identity
	var hostIP = make([]net.IP, 0, 2)
	var venID uint32
	var prodName string
	var oState uint32
	// var errorMsg string
	// var failedAVP []AVP
	var suuportVendor = []uint32{}
	var authApps = make(map[uint32]uint32)
	// var firmwareRevision uint32
	var err error

	for rdr := bytes.NewReader(v.m.AVPs); rdr.Len() != 0; {
		a := AVP{}
		if err = a.wrapedUnmarshalFrom(rdr); err != nil {
			break
		}
		if a.VendorID != 0 {
			if a.Mandatory {
				err = InvalidAVP{Code: AvpUnsupported, AVP: a}
				break
			}
			continue
		}

		switch a.Code {
		case 268:
			if result != 0 {
				err = InvalidAVP{Code: AvpOccursTooManyTimes, AVP: a}
			} else {
				result, err = GetResultCode(a)
			}
		case 264:
			if len(oHost) != 0 {
				err = InvalidAVP{Code: AvpOccursTooManyTimes, AVP: a}
			} else {
				oHost, err = GetOriginHost(a)
			}
		case 296:
			if len(oRealm) != 0 {
				err = InvalidAVP{Code: AvpOccursTooManyTimes, AVP: a}
			} else {
				oRealm, err = GetOriginRealm(a)
			}
		case 257:
			if ip, e := getHostIPAddress(a); e == nil {
				hostIP = append(hostIP, ip)
			} else {
				err = e
			}
		case 266:
			if venID != 0 {
				err = InvalidAVP{Code: AvpOccursTooManyTimes, AVP: a}
			} else {
				venID, err = GetVendorID(a)
			}
		case 269:
			if len(prodName) != 0 {
				err = InvalidAVP{Code: AvpOccursTooManyTimes, AVP: a}
			} else {
				prodName, err = getProductName(a)
			}
		case 265:
			var vid uint32
			if vid, err = getSupportedVendorID(a); err == nil {
				suuportVendor = append(suuportVendor, vid)
			}
		case 258:
			var aid uint32
			if aid, err = GetAuthAppID(a); err == nil {
				authApps[aid] = 0
			}
		case 260:
			var vid, aid uint32
			if vid, aid, err = GetVendorSpecAppID(a); err == nil {
				authApps[aid] = vid
			}
		case 278:
			if oState != 0 {
				err = InvalidAVP{Code: AvpOccursTooManyTimes, AVP: a}
			} else {
				oState, err = getOriginStateID(a)
			}
		// case 281:
		//	errorMsg, e = getErrorMessage(a)
		case 279:
			//	failedAVP, e = getFailedAVP(a)
		// case 267:
		//	firmwareRevision, e = getFirmwareRevision(a)
		default:
			if a.Mandatory {
				err = InvalidAVP{Code: AvpUnsupported, AVP: a}
			}
		}

		if err != nil {
			break
		}
	}

	if len(authApps) == 0 {
		err = InvalidAVP{Code: MissingAvp, AVP: SetAuthAppID(0)}
	}
	if err == nil {
		for _, vid := range authApps {
			if vid == 0 {
				continue
			}
			err = InvalidAVP{Code: MissingAvp, AVP: SetVendorID(0)}
			for _, v := range suuportVendor {
				if vid == v {
					err = nil
					break
				}
			}
			if err != nil {
				break
			}
		}
	}
	if err == nil {
		if _, ok := authApps[0xffffffff]; ok {
		} else if len(applications) == 0 {
			for aid, vid := range authApps {
				applications[aid] = application{
					venID:    vid,
					handlers: make(map[uint32]func(bool, []byte) (bool, []byte))}
			}
		} else {
			var commonApp = make(map[uint32]application)
			for laid, lapp := range applications {
				if pvid, ok := authApps[laid]; ok && pvid == lapp.venID {
					commonApp[laid] = lapp
					break
				}
			}
			if len(commonApp) == 0 {
				err = InvalidMessage(ApplicationUnsupported)
			} else {
				applications = commonApp
			}
		}
	}

	if v.m.FlgE && result == Success {
		err = InvalidMessage(InvalidHdrBits)
	} else if err != nil {
		// invalid AVP value
	} else if result == 0 {
		err = InvalidAVP{Code: MissingAvp, AVP: SetResultCode(0)}
	} else if len(oHost) == 0 {
		err = InvalidAVP{Code: MissingAvp, AVP: SetOriginHost("")}
	} else if len(oRealm) == 0 {
		err = InvalidAVP{Code: MissingAvp, AVP: SetOriginRealm("")}
	} else if oHost != Peer.Host && oHost != Local.Host {
		err = InvalidMessage(UnknownPeer)
	} else if oRealm != Peer.Realm && oRealm != Local.Realm {
		err = InvalidMessage(UnknownPeer)
	} else if len(hostIP) == 0 {
		err = InvalidAVP{Code: MissingAvp, AVP: setHostIPAddress(net.IPv4zero)}
	} else if venID == 0 {
		err = InvalidAVP{Code: MissingAvp, AVP: SetVendorID(0)}
	} else if len(prodName) == 0 {
		err = InvalidAVP{Code: MissingAvp, AVP: setProductName("")}
	} else if result != Success {
		err = FailureAnswer{Code: result}
	} else {
		if oState != 0 {
			Peer.state = oState
		}

		state = open
		wdTimer.Stop()
		wdTimer = time.AfterFunc(WDInterval, func() {
			notify <- eventWatchdog{}
		})
		delete(sndStack, v.m.HbHID)
		//ch <- v.m
	}
	countRxCode(result)
	TraceMessage(v.m, Rx, err)

	if err != nil {
		wdTimer.Stop()
		delete(sndStack, v.m.HbHID)
		// close(ch)
		conn.Close()
	}
	return err
}