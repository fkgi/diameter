package diameter

import (
	"bytes"
	"fmt"
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
		 * [ Inband-Security-Id ]             // not supported (not recommended)
		 * [ Acct-Application-Id ]            // not supported
		 * [ Vendor-Specific-Application-Id ] // only support auth
		   [ Firmware-Revision ]
		 * [ AVP ]                            // no any other AVP

Capabilities-Exchange-Answer message
 <CEA> ::= < Diameter Header: 257 >
		   { Result-Code }
		   { Origin-Host }
		   { Origin-Realm }
		1* { Host-IP-Address }
		   { Vendor-Id }
		   { Product-Name }
		   [ Origin-State-Id ]
		   [ Error-Message ]                  // ignored
		   [ Failed-AVP ]                     // ignored
		 * [ Supported-Vendor-Id ]
		 * [ Auth-Application-Id ]
		 * [ Inband-Security-Id ]             // not supported (not recommended)
		 * [ Acct-Application-Id ]            // not supported
		 * [ Vendor-Specific-Application-Id ] // only support auth
		   [ Firmware-Revision ]              // ignored
		 * [ AVP ]                            // ignored
*/

type eventRcvCER struct {
	m Message
}

func (eventRcvCER) String() string {
	return "Rcv-CER"
}

func (v eventRcvCER) exec(c *Connection) error {
	RxReq++
	if c.state != waitCER {
		RejectReq++
		return notAcceptableEvent{e: v, s: c.state}
	}
	if TraceMessage != nil {
		TraceMessage(v.m, Rx, nil)
	}

	var oHost Identity
	var oRealm Identity
	var hostIP = make([]net.IP, 0, 2)
	var venID uint32
	var prodName string
	var oState uint32
	var supportVendor = []uint32{}
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
				supportVendor = append(supportVendor, vid)
			}
		case 258:
			var aid uint32
			if aid, err = GetAuthAppID(a); err == nil {
				authApps[aid] = 0
			}
		case 299:
			_, err = getInbandSecurityID(a)
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
		case 267:
			_, err = getFirmwareRevision(a)
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
		for _, vid := range supportVendor {
			if vid == 0 {
				continue
			}
			// ToDo: verify common supported AVP vendor
		}
	}

	if err == nil {
		if _, ok := authApps[0xffffffff]; ok {
			if len(applications) != 0 {
				for aid, app := range applications {
					c.commonApp[aid] = app
				}
			}
		} else if len(applications) == 0 {
			for aid, vid := range authApps {
				c.commonApp[aid] = application{
					venID:    vid,
					handlers: make(map[uint32]Handler)}
			}
		} else {
			for laid, lapp := range applications {
				if pvid, ok := authApps[laid]; ok && pvid == lapp.venID {
					c.commonApp[laid] = lapp
					break
				}
			}
			if len(c.commonApp) == 0 {
				rap := "required applications are "
				for k := range authApps {
					rap = fmt.Sprintf("%s, %d", rap, k)
				}
				err = InvalidMessage{
					Code:   ApplicationUnsupported,
					ErrMsg: rap}
			}
		}
	}

	result := Success
	if v.m.FlgP || v.m.FlgT {
		result = InvalidHdrBits
		err = InvalidMessage{
			Code: result, ErrMsg: "CER must not enable P and T flag"}
	} else if iavp, ok := err.(InvalidAVP); ok {
		result = iavp.Code
	} else if imsg, ok := err.(InvalidMessage); ok {
		result = imsg.Code
	} else if len(oHost) == 0 {
		result = MissingAvp
		err = InvalidAVP{Code: result, AVP: SetOriginHost("")}
	} else if len(oRealm) == 0 {
		result = MissingAvp
		err = InvalidAVP{Code: result, AVP: SetOriginRealm("")}
	} else if c.Host != "" && oHost != c.Host {
		result = UnknownPeer
		err = InvalidMessage{
			Code: result,
			ErrMsg: fmt.Sprintf(
				"peer host %s is not match with %s",
				oHost, c.Host)}
	} else if c.Realm != "" && oRealm != c.Realm {
		result = UnknownPeer
		err = InvalidMessage{
			Code: result,
			ErrMsg: fmt.Sprintf(
				"peer realm %s is not match with %s",
				oRealm, c.Realm)}
	} else if len(hostIP) == 0 {
		result = MissingAvp
		err = InvalidAVP{Code: result, AVP: setHostIPAddress(net.IPv4zero)}
		/*
			} else if venID == 0 {
				result = MissingAvp
				err = InvalidAVP{Code: result, AVP: SetVendorID(0)}
			} else if len(prodName) == 0 {
				result = MissingAvp
				err = InvalidAVP{Code: result, AVP: setProductName("")}
		*/
	} else {
		if oState != 0 {
			c.stateID = oState
		}
		c.Host = oHost
		c.Realm = oRealm
	}

	buf := new(bytes.Buffer)
	SetResultCode(result).MarshalTo(buf)
	SetOriginHost(Host).MarshalTo(buf)
	SetOriginRealm(Realm).MarshalTo(buf)

	if len(OverwriteAddr) != 0 {
		for _, h := range OverwriteAddr {
			setHostIPAddress(h).MarshalTo(buf)
		}
	} else {
		h, _, _ := net.SplitHostPort(c.conn.LocalAddr().String())
		for _, h := range strings.Split(h, "/") {
			setHostIPAddress(net.ParseIP(h)).MarshalTo(buf)
		}
	}

	SetVendorID(VendorID).MarshalTo(buf)
	setProductName(ProductName).MarshalTo(buf)
	if stateID != 0 {
		setOriginStateID(stateID).MarshalTo(buf)
	}
	if len(c.commonApp) == 0 {
		SetAuthAppID(0xffffffff).MarshalTo(buf)
	} else {
		vmap := make(map[uint32]interface{})
		for aid, app := range c.commonApp {
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

	if e := cea.MarshalTo(c.conn); e != nil {
		TxAnsFail++
		c.conn.Close()
		err = e
	} else if err == nil {
		CountTxCode(result)
		c.state = open
		// wdTimer.Stop()
		c.wdTimer = time.AfterFunc(WDInterval, func() {
			c.notify <- eventWatchdog{}
		})
		if ConnectionUpNotify != nil {
			ConnectionUpNotify(c)
		}
	} else {
		CountTxCode(result)
	}

	if TraceMessage != nil {
		TraceMessage(cea, Tx, err)
	}
	return err
}

// RcvCEA
type eventRcvCEA struct {
	m Message
}

func (eventRcvCEA) String() string {
	return "Rcv-CEA"
}

func (v eventRcvCEA) exec(c *Connection) error {
	// verify diameter header
	if v.m.FlgP {
		InvalidAns++
		return InvalidMessage{
			Code: InvalidHdrBits, ErrMsg: "CEA must not enable P flag"}
	}
	if c.state != waitCEA {
		InvalidAns++
		return notAcceptableEvent{e: v, s: c.state}
	}
	if _, ok := c.sndQueue[v.m.HbHID]; !ok {
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
	var errorMsg string
	var failedAVP []AVP
	var supportVendor = []uint32{}
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
				supportVendor = append(supportVendor, vid)
			}
		case 258:
			var aid uint32
			if aid, err = GetAuthAppID(a); err == nil {
				authApps[aid] = 0
			}
		case 299:
			_, err = getInbandSecurityID(a)
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
		case 281:
			errorMsg, err = GetErrorMessage(a)
		case 279:
			failedAVP, err = getFailedAVP(a)
		case 267:
			_, err = getFirmwareRevision(a)
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
		for _, vid := range supportVendor {
			if vid == 0 {
				continue
			}
			// ToDo: verify common supported AVP vendor
		}
	}
	if err == nil {
		if _, ok := authApps[0xffffffff]; ok {
			if len(applications) != 0 {
				for aid, app := range applications {
					c.commonApp[aid] = app
				}
			}
		} else if len(applications) == 0 {
			for aid, vid := range authApps {
				c.commonApp[aid] = application{
					venID:    vid,
					handlers: make(map[uint32]Handler)}
			}
		} else {
			for laid, lapp := range applications {
				if pvid, ok := authApps[laid]; ok && pvid == lapp.venID {
					c.commonApp[laid] = lapp
					break
				}
			}
			if len(c.commonApp) == 0 {
				rap := "required applications are "
				for k := range authApps {
					rap = fmt.Sprintf("%s, %d", rap, k)
				}
				err = InvalidMessage{
					Code:   ApplicationUnsupported,
					ErrMsg: rap}
			}
		}
	}

	if v.m.FlgE && result == Success {
		err = InvalidMessage{
			Code:   InvalidHdrBits,
			ErrMsg: "error flag is true but success response code"}
	} else if result == 0 {
		err = InvalidAVP{Code: MissingAvp, AVP: SetResultCode(0)}
	} else if len(oHost) == 0 {
		err = InvalidAVP{Code: MissingAvp, AVP: SetOriginHost("")}
	} else if len(oRealm) == 0 {
		err = InvalidAVP{Code: MissingAvp, AVP: SetOriginRealm("")}
	} else if oHost != c.Host && oHost != Host {
		err = InvalidMessage{
			Code: UnknownPeer,
			ErrMsg: fmt.Sprintf(
				"peer host %s is not match with %s or %s",
				oHost, c.Host, Host)}
	} else if oRealm != c.Realm && oRealm != Realm {
		err = InvalidMessage{
			Code: UnknownPeer,
			ErrMsg: fmt.Sprintf(
				"peer realm %s is not match with %s or %s",
				oRealm, c.Realm, Host)}
	} else if result != Success {
		err = FailureAnswer{Code: result, ErrMsg: errorMsg, Avps: failedAVP}
	} else if err != nil {
		// invalid AVP value
	} else if len(hostIP) == 0 {
		err = InvalidAVP{Code: MissingAvp, AVP: setHostIPAddress(net.IPv4zero)}
		/*
			} else if venID == 0 {
				err = InvalidAVP{Code: MissingAvp, AVP: SetVendorID(0)}
			} else if len(prodName) == 0 {
				err = InvalidAVP{Code: MissingAvp, AVP: setProductName("")}
		*/
	} else {
		if oState != 0 {
			c.stateID = oState
		}

		c.state = open
		c.wdTimer.Stop()
		c.wdTimer = time.AfterFunc(WDInterval, func() {
			c.notify <- eventWatchdog{}
		})
		delete(c.sndQueue, v.m.HbHID)
		//ch <- v.m
		if ConnectionUpNotify != nil {
			ConnectionUpNotify(c)
		}
	}
	CountRxCode(result)
	if TraceMessage != nil {
		TraceMessage(v.m, Rx, err)
	}

	if err != nil {
		c.wdTimer.Stop()
		delete(c.sndQueue, v.m.HbHID)
		// close(ch)
		c.conn.Close()
	}
	return err
}
