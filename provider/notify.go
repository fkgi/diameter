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
	if d, ok := e.(MessageTransfer); ok {
		d.Dump(os.Stderr)
	}
}

// MessageTransfer indicate Diameter message sent
type MessageTransfer struct {
	Local *LocalNode
	Peer  *PeerNode
	Tx    bool
	Err   error
	dump  func(io.Writer)
}

func (e *MessageTransfer) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Tx && e.Err == nil {
		return fmt.Sprintf(
			"write Diameter message from %s to %s",
			e.Local, e.Peer)
	}
	if e.Tx && e.Err != nil {
		return fmt.Sprintf(
			"write Diameter message failed from %s to %s, reason is %s",
			e.Local, e.Peer, e.Err)
	}
	if !e.Tx && e.Err == nil {
		return fmt.Sprintf(
			"read Diameter message from %s to %s",
			e.Local, e.Peer)
	}
	if !e.Tx && e.Err != nil {
		return fmt.Sprintf(
			"read Diameter message failed from %s to %s, reason is %s",
			e.Local, e.Peer, e.Err)
	}
	return ""
}

// Dump provide Diameter message dump
func (e *MessageTransfer) Dump(w io.Writer) {
	if e.dump != nil {
		e.dump(w)
	}
}

// TransportStateChange indicate transport link open
type TransportStateChange struct {
	Open  bool
	Local *LocalNode
	Peer  *PeerNode
	LAddr net.Addr
	PAddr net.Addr
	Err   error
}

func (e *TransportStateChange) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Open && e.Err == nil {
		return fmt.Sprintf(
			"transport link from %s(%s) to %s(%s) open",
			e.Local, e.LAddr, e.Peer, e.PAddr)
	}
	if e.Open && e.Err != nil {
		return fmt.Sprintf(
			"transport link from %s to %s open failed, reason is %s",
			e.Local, e.Peer, e.Err)
	}
	if !e.Open && e.Err == nil {
		return fmt.Sprintf(
			"transport link from %s(%s) to %s(%s) close",
			e.Local, e.LAddr, e.Peer, e.PAddr)
	}
	if !e.Open && e.Err != nil {
		return fmt.Sprintf(
			"transport link from %s to %s close failed, reason is %s",
			e.Local, e.Peer, e.Err)
	}
	return ""
}

// TransportBind indicate transport listener binded to local node
type TransportBind struct {
	Local *LocalNode
	LAddr net.Addr
	Err   error
}

func (e *TransportBind) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Err == nil {
		return "transport listener binded"
	}
	return fmt.Sprintf(
		"transport listener bind failed, reason is %s",
		e.Err)
}

// StateUpdate notify event
type StateUpdate struct {
	State string
	Local *LocalNode
	Peer  *PeerNode
	Err   error
}

func (e *StateUpdate) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Err != nil {
		return fmt.Sprintf(
			"provider state invalid update to %s, %s", e.State, e.Err)
	}
	return fmt.Sprintf("provider state update to %s", e.State)
}

// ConnectionStateChange notify event
type ConnectionStateChange struct {
	Open  bool
	Local *LocalNode
	Peer  *PeerNode
}

func (e *ConnectionStateChange) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Open {
		return "Diameter connection open"
	}
	return "Diameter connection close"
}

// CapabilityExchangeEvent notify capability exchange related event
type CapabilityExchangeEvent struct {
	Tx    bool
	Req   bool
	Local *LocalNode
	Peer  *PeerNode
	Err   error
}

func (e *CapabilityExchangeEvent) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Tx && e.Req && e.Err == nil {
		return fmt.Sprintf("-> CER")
	}
	if e.Tx && e.Req && e.Err != nil {
		return fmt.Sprintf("-X CER failed, %s", e.Err)
	}
	if e.Tx && !e.Req && e.Err == nil {
		return fmt.Sprintf("-> CEA")
	}
	if e.Tx && !e.Req && e.Err != nil {
		return fmt.Sprintf("-X CEA failed, %s", e.Err)
	}
	if !e.Tx && e.Req && e.Err == nil {
		return fmt.Sprintf("<- CER")
	}
	if !e.Tx && e.Req && e.Err != nil {
		return fmt.Sprintf("X- CER failed, %s", e.Err)
	}
	if !e.Tx && !e.Req && e.Err == nil {
		return fmt.Sprintf("<- CEA")
	}
	if !e.Tx && !e.Req && e.Err != nil {
		return fmt.Sprintf("X- CEA failed, %s", e.Err)
	}
	return ""
}

// WatchdogEvent notify watchdog related event
type WatchdogEvent struct {
	Tx    bool
	Req   bool
	Local *LocalNode
	Peer  *PeerNode
	Err   error
}

func (e *WatchdogEvent) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Tx && e.Req && e.Err == nil {
		return fmt.Sprintf("-> DWR")
	}
	if e.Tx && e.Req && e.Err != nil {
		return fmt.Sprintf("-X DWR failed, %s", e.Err)
	}
	if e.Tx && !e.Req && e.Err == nil {
		return fmt.Sprintf("-> DWA")
	}
	if e.Tx && !e.Req && e.Err != nil {
		return fmt.Sprintf("-X DWA failed, %s", e.Err)
	}
	if !e.Tx && e.Req && e.Err == nil {
		return fmt.Sprintf("<- DWR")
	}
	if !e.Tx && e.Req && e.Err != nil {
		return fmt.Sprintf("X- DWR failed, %s", e.Err)
	}
	if !e.Tx && !e.Req && e.Err == nil {
		return fmt.Sprintf("<- DWA")
	}
	if !e.Tx && !e.Req && e.Err != nil {
		return fmt.Sprintf("X- DWA failed, %s", e.Err)
	}
	return ""
}

// MessageEvent notify diameter message related event
type MessageEvent struct {
	Tx    bool
	Req   bool
	Local *LocalNode
	Peer  *PeerNode
	Err   error
}

func (e *MessageEvent) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Tx && e.Req && e.Err == nil {
		return fmt.Sprintf("-> REQ")
	}
	if e.Tx && e.Req && e.Err != nil {
		return fmt.Sprintf("-X REQ failed, %s", e.Err)
	}
	if e.Tx && !e.Req && e.Err == nil {
		return fmt.Sprintf("-> ANS")
	}
	if e.Tx && !e.Req && e.Err != nil {
		return fmt.Sprintf("-X ANS failed, %s", e.Err)
	}
	if !e.Tx && e.Req && e.Err == nil {
		return fmt.Sprintf("<- REQ")
	}
	if !e.Tx && e.Req && e.Err != nil {
		return fmt.Sprintf("X- REQ failed, %s", e.Err)
	}
	if !e.Tx && !e.Req && e.Err == nil {
		return fmt.Sprintf("<- ANS")
	}
	if !e.Tx && !e.Req && e.Err != nil {
		return fmt.Sprintf("X- ANS failed, %s", e.Err)
	}
	return ""
}

// PurgeEvent notify diameter purge related event
type PurgeEvent struct {
	Tx    bool
	Req   bool
	Local *LocalNode
	Peer  *PeerNode
	Err   error
}

func (e *PurgeEvent) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Tx && e.Req && e.Err == nil {
		return fmt.Sprintf("-> DPR")
	}
	if e.Tx && e.Req && e.Err != nil {
		return fmt.Sprintf("-X DPR failed, %s", e.Err)
	}
	if e.Tx && !e.Req && e.Err == nil {
		return fmt.Sprintf("-> DPA")
	}
	if e.Tx && !e.Req && e.Err != nil {
		return fmt.Sprintf("-X DPA failed, %s", e.Err)
	}
	if !e.Tx && e.Req && e.Err == nil {
		return fmt.Sprintf("<- DPR")
	}
	if !e.Tx && e.Req && e.Err != nil {
		return fmt.Sprintf("X- DPR failed, %s", e.Err)
	}
	if !e.Tx && !e.Req && e.Err == nil {
		return fmt.Sprintf("<- DPA")
	}
	if !e.Tx && !e.Req && e.Err != nil {
		return fmt.Sprintf("X- DPA failed, %s", e.Err)
	}
	return ""
}
