package provider

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// Notify is called when error or trace event are occured
var Notify = func(e error) {
	log.Println(e)
	d, ok := e.(Dump)
	if ok {
		d.f(os.Stderr)
	}
}

// TxMessage indicate Diameter message sent
type TxMessage struct {
	Local string
	Peer  string
	Err   error
	dump  func(io.Writer)
}

func (e *TxMessage) Error() (s string) {
	if e == nil {
		s = "<nil>"
	} else if e.Err == nil {
		s = fmt.Sprintf(
			"write Diameter message from %s to %s",
			e.Local, e.Peer)
	} else {
		s = fmt.Sprintf(
			"write Diameter message failed from %s to %s, reason is %s",
			e.Local, e.Peer, e.Err)
	}
	return s
}

// Dump provide Diameter message dump
func (e *TxMessage) Dump(w io.Writer) {
	if e.dump != nil {
		e.dump(w)
	}
}

// RxMessage indicate Diameter message received
type RxMessage struct {
	Local string
	Peer  string
	Err   error
	dump  func(io.Writer)
}

func (e *RxMessage) Error() (s string) {
	if e == nil {
		s = "<nil>"
	} else if e.Err == nil {
		s = fmt.Sprintf(
			"read Diameter message from %s to %s",
			e.Local, e.Peer)
	} else {
		s = fmt.Sprintf(
			"read Diameter message failed from %s to %s, reason is %s",
			e.Local, e.Peer, e.Err)
	}
	return s
}

// Dump provide Diameter message dump
func (e *RxMessage) Dump(w io.Writer) {
	if e.dump != nil {
		e.dump(w)
	}
}

// TransportStateChange indicate transport link open
type TransportStateChange struct {
	Local string
	Peer  string
	LAddr net.Addr
	PAddr net.Addr
	Err   error
	Open  bool
}

func (e *TransportStateChange) Error() (s string) {
	if e == nil {
		s = "<nil>"
	} else if e.Open && e.Err == nil {
		s = fmt.Sprintf(
			"transport link from %s(%s) to %s(%s) open",
			e.Local, e.LAddr, e.Peer, e.PAddr)
	} else if e.Open {
		s = fmt.Sprintf(
			"transport link from %s to %s open failed, reason is %s",
			e.Local, e.Peer, e.Err)
	} else if !e.Open && e.Err == nil {
		s = fmt.Sprintf(
			"transport link from %s(%s) to %s(%s) close",
			e.Local, e.LAddr, e.Peer, e.PAddr)
	} else {
		s = fmt.Sprintf(
			"transport link from %s to %s close failed, reason is %s",
			e.Local, e.Peer, e.Err)
	}
	return
}

// TransportBind indicate transport listener binded to local node
type TransportBind struct {
	Local *LocalNode
	LAddr net.Addr
	Err   error
}

func (e *TransportBind) Error() (s string) {
	if e == nil {
		s = "<nil>"
	} else if e.Err == nil {
		s = fmt.Sprintf(
			"transport listener binded")
	} else {
		s = fmt.Sprintf(
			"transport listener bind failed, reason is %s", e.Err)
	}
	return
}

// InvalidEvent notify event
type InvalidEvent struct {
	Err   error
	State string
}

func (e *InvalidEvent) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf(
		"not acceptable event for provider status %s, %s", e.State, e.Err)
}

// StateUpdate notify event
type StateUpdate struct {
	State string
}

func (e *StateUpdate) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("provider state update, %s", e.State)
}

// DiameterStateChange notify event
type DiameterStateChange struct {
	Open bool
}

func (e *DiameterStateChange) Error() (s string) {
	if e == nil {
		s = "<nil>"
	} else if Open {
		s = "Diameter connection open"
	} else {
		s = "Diameter connection close"
	}
	return
}

// WatchdogEvent notify watchdog related event
type WatchdogEvent struct {
	Tx  bool
	Req bool
	Err bool
}

func (e *WatchdogEvent) Error() (s string) {
	if e == nil {
		s = "<nil>"
	} else if Tx && Req {
		s = fmt.Sprintf("-> DWR")
	} else if Tx && !Req {
		s = fmt.Sprintf("-> DWA")
	} else if !Tx && Req {
		s = fmt.Sprintf("<- DWR")
	} else if !Tx && !Req {
		s = fmt.Sprintf("<- DWA")
	}
	return s
}

// NoWatchdogAns notify event
type NoWatchdogAns struct {
}

func (e *NoWatchdogAns) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("no DWA response")
}

// RxWatchdogReq notify event
type RxWatchdogReq struct {
}

func (e *RxWatchdogReq) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("<- DWR")
}

// TxWatchdogAns notify event
type TxWatchdogAns struct {
}

func (e *TxWatchdogAns) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("-> DWA")
}

// TxExchangeReq notify event
type TxExchangeReq struct {
}

func (e *TxExchangeReq) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("-> CER")
}

// RxExchangeAns notify event
type RxExchangeAns struct {
}

func (e *RxExchangeAns) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("<- CEA")
}

// NoExchangeAns notify event
type NoExchangeAns struct {
}

func (e *NoExchangeAns) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("no CEA response")
}

// RxExchangeReq notify event
type RxExchangeReq struct {
}

func (e *RxExchangeReq) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("<- CER")
}

// TxExchangeAns notify event
type TxExchangeAns struct {
}

func (e *TxExchangeAns) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("-> CEA")
}

// TxDataReq notify event
type TxDataReq struct {
}

func (e *TxDataReq) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("-> REQ")
}

// RxDataAns notify event
type RxDataAns struct {
}

func (e *RxDataAns) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("<- ANS")
}

// NoDataAns notify event
type NoDataAns struct {
	Retry bool
}

func (e *NoDataAns) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Retry {
		return fmt.Sprintf("no response retry")
	}
	return fmt.Sprintf("no response")
}

// RxDataReq notify event
type RxDataReq struct {
}

func (e *RxDataReq) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("<- REQ")
}

// TxDataAns notify event
type TxDataAns struct {
}

func (e *TxDataAns) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("-> ANS")
}
