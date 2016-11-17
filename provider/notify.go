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
	if d, ok := e.(TxMessage); ok {
		d.Dump(os.Stderr)
	} else if d, ok := e.(RxMessage); ok {
		d.Dump(os.Stderr)
	}
}

// MessageTransfer indicate Diameter message sent
type MessageTransfer struct {
	Local string
	Peer  string
	Tx    bool
	Err   error
	dump  func(io.Writer)
}

func (e *MessageTransfer) Error() (s string) {
	if e == nil {
		s = "<nil>"
	} else if e.Tx && e.Err == nil {
		s = fmt.Sprintf(
			"write Diameter message from %s to %s",
			e.Local, e.Peer)
	} else if e.Tx {
		s = fmt.Sprintf(
			"write Diameter message failed from %s to %s, reason is %s",
			e.Local, e.Peer, e.Err)
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
func (e *MessageTransfer) Dump(w io.Writer) {
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
	} else if e.Open {
		s = "Diameter connection open"
	} else {
		s = "Diameter connection close"
	}
	return
}

// CapabilityExchangeEvent notify capability exchange related event
type CapabilityExchangeEvent struct {
	Tx  bool
	Req bool
	Err error
}

func (e *CapabilityExchangeEvent) Error() (s string) {
	if e == nil {
		s = "<nil>"
	} else if e.Tx && e.Req && e.Err == nil {
		s = fmt.Sprintf("-> CER")
	} else if e.Tx && e.Req {
		s = fmt.Sprintf("-X CER failed, %s", e.Err)

	} else if e.Tx && !e.Req && e.Err == nil {
		s = fmt.Sprintf("-> CEA")
	} else if e.Tx && !e.Req {
		s = fmt.Sprintf("-X CEA failed, %s", e.Err)

	} else if !e.Tx && e.Req && e.Err == nil {
		s = fmt.Sprintf("<- CER")
	} else if !e.Tx && e.Req {
		s = fmt.Sprintf("X- CER failed, %s", e.Err)

	} else if !e.Tx && !e.Req && e.Err == nil {
		s = fmt.Sprintf("<- CEA")
	} else if !e.Tx && !e.Req {
		s = fmt.Sprintf("X- CEA failed, %s", e.Err)
	}
	return s
}

// WatchdogEvent notify watchdog related event
type WatchdogEvent struct {
	Tx  bool
	Req bool
	Err error
}

func (e *WatchdogEvent) Error() (s string) {
	if e == nil {
		s = "<nil>"
	} else if e.Tx && e.Req && e.Err == nil {
		s = fmt.Sprintf("-> DWR")
	} else if e.Tx && e.Req {
		s = fmt.Sprintf("-X DWR failed, %s", e.Err)

	} else if e.Tx && !e.Req && e.Err == nil {
		s = fmt.Sprintf("-> DWA")
	} else if e.Tx && !e.Req {
		s = fmt.Sprintf("-X DWA failed, %s", e.Err)

	} else if !e.Tx && e.Req && e.Err == nil {
		s = fmt.Sprintf("<- DWR")
	} else if !e.Tx && e.Req {
		s = fmt.Sprintf("X- DWR failed, %s", e.Err)

	} else if !e.Tx && !e.Req && e.Err == nil {
		s = fmt.Sprintf("<- DWA")
	} else if !e.Tx && !e.Req {
		s = fmt.Sprintf("X- DWA failed, %s", e.Err)
	}
	return s
}

// MessageEvent notify diameter message related event
type MessageEvent struct {
	Tx  bool
	Req bool
	Err error
}

func (e *MessageEvent) Error() (s string) {
	if e == nil {
		s = "<nil>"
	} else if e.Tx && e.Req && e.Err == nil {
		s = fmt.Sprintf("-> REQ")
	} else if e.Tx && e.Req {
		s = fmt.Sprintf("-X REQ failed, %s", e.Err)

	} else if e.Tx && !e.Req && e.Err == nil {
		s = fmt.Sprintf("-> ANS")
	} else if e.Tx && !e.Req {
		s = fmt.Sprintf("-X ANS failed, %s", e.Err)

	} else if !e.Tx && e.Req && e.Err == nil {
		s = fmt.Sprintf("<- REQ")
	} else if !e.Tx && e.Req {
		s = fmt.Sprintf("X- REQ failed, %s", e.Err)

	} else if !e.Tx && !e.Req && e.Err == nil {
		s = fmt.Sprintf("<- ANS")
	} else if !e.Tx && !e.Req {
		s = fmt.Sprintf("X- ANS failed, %s", e.Err)
	}
	return s
}

// PurgeEvent notify diameter message related event
type PurgeEvent struct {
	Tx  bool
	Req bool
	Err error
}

func (e *PurgeEvent) Error() (s string) {
	if e == nil {
		s = "<nil>"
	} else if e.Tx && e.Req && e.Err == nil {
		s = fmt.Sprintf("-> DPR")
	} else if e.Tx && e.Req {
		s = fmt.Sprintf("-X DPR failed, %s", e.Err)

	} else if e.Tx && !e.Req && e.Err == nil {
		s = fmt.Sprintf("-> DPA")
	} else if e.Tx && !e.Req {
		s = fmt.Sprintf("-X DPA failed, %s", e.Err)

	} else if !e.Tx && e.Req && e.Err == nil {
		s = fmt.Sprintf("<- DPR")
	} else if !e.Tx && e.Req {
		s = fmt.Sprintf("X- DPR failed, %s", e.Err)

	} else if !e.Tx && !e.Req && e.Err == nil {
		s = fmt.Sprintf("<- DPA")
	} else if !e.Tx && !e.Req {
		s = fmt.Sprintf("X- DPA failed, %s", e.Err)
	}
	return s
}
