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
	Tx    bool
	Local *LocalNode
	Peer  *PeerNode
	Err   error
	dump  func(io.Writer)
}

func (e *MessageTransfer) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Tx && e.Err == nil {
		return fmt.Sprintf(
			"write Diameter message %s -> %s",
			e.Local, e.Peer)
	}
	if e.Tx && e.Err != nil {
		return fmt.Sprintf(
			"write Diameter message %s -> %s failed: %s",
			e.Local, e.Peer, e.Err)
	}
	if !e.Tx && e.Err == nil {
		return fmt.Sprintf(
			"read Diameter message %s <- %s",
			e.Local, e.Peer)
	}
	if !e.Tx && e.Err != nil {
		return fmt.Sprintf(
			"read Diameter message %s <- %s failed: %s",
			e.Local, e.Peer, e.Err)
	}
	return "invlid state of message transfer"
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
			"transport link %s(%s) -> %s(%s) open",
			e.Local, e.LAddr, e.Peer, e.PAddr)
	}
	if e.Open && e.Err != nil {
		return fmt.Sprintf(
			"transport link %s -> %s open failed: %s",
			e.Local, e.Peer, e.Err)
	}
	if !e.Open && e.Err == nil {
		return fmt.Sprintf(
			"transport link %s(%s) -> %s(%s) closed",
			e.Local, e.LAddr, e.Peer, e.PAddr)
	}
	if !e.Open && e.Err != nil {
		return fmt.Sprintf(
			"transport link %s(%s) -> %s(%s) close failed: %s",
			e.Local, e.LAddr, e.Peer, e.PAddr, e.Err)
	}
	return "invlid state of transport state"
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
		return fmt.Sprintf(
			"transport listener binded on %s(%s)",
			e.Local, e.LAddr)
	}
	return fmt.Sprintf(
		"transport listener bind on %s(%s) failed: %s",
		e.Local, e.LAddr, e.Err)
}

// StateUpdate notify event
type StateUpdate struct {
	State string
	Event string
	Local *LocalNode
	Peer  *PeerNode
	Err   error
}

func (e *StateUpdate) Error() string {
	if e == nil {
		return "<nil>"
	}
	if e.Err == nil {
		return fmt.Sprintf(
			"provider %s -> %s state update event %s occured on %s",
			e.Local, e.Peer, e.Event, e.State)
	}
	return fmt.Sprintf(
		"provider %s -> %s state update event %s occured on %s failed: %s",
		e.Local, e.Peer, e.Event, e.State, e.Err)
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
		return fmt.Sprintf(
			"Diameter connection %s -> %s open",
			e.Local, e.Peer)
	}
	return fmt.Sprintf(
		"Diameter connection %s -> %s close",
		e.Local, e.Peer)
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
		return fmt.Sprintf(
			"-> CER (%s -> %s)", e.Local, e.Peer)
	}
	if e.Tx && e.Req && e.Err != nil {
		return fmt.Sprintf(
			"-X CER (%s -> %s), %s", e.Local, e.Peer, e.Err)
	}
	if e.Tx && !e.Req && e.Err == nil {
		return fmt.Sprintf(
			"-> CEA (%s -> %s)", e.Local, e.Peer)
	}
	if e.Tx && !e.Req && e.Err != nil {
		return fmt.Sprintf(
			"-X CEA (%s -> %s), %s", e.Local, e.Peer, e.Err)
	}
	if !e.Tx && e.Req && e.Err == nil {
		return fmt.Sprintf(
			"<- CER (%s -> %s)", e.Local, e.Peer)
	}
	if !e.Tx && e.Req && e.Err != nil {
		return fmt.Sprintf(
			"X- CER (%s -> %s), %s", e.Local, e.Peer, e.Err)
	}
	if !e.Tx && !e.Req && e.Err == nil {
		return fmt.Sprintf(
			"<- CEA (%s -> %s)", e.Local, e.Peer)
	}
	if !e.Tx && !e.Req && e.Err != nil {
		return fmt.Sprintf(
			"X- CEA (%s -> %s), %s", e.Local, e.Peer, e.Err)
	}
	return "invalid capability exchange state"
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
		return fmt.Sprintf(
			"-> DWR (%s -> %s)", e.Local, e.Peer)
	}
	if e.Tx && e.Req && e.Err != nil {
		return fmt.Sprintf(
			"-X DWR (%s -> %s), %s", e.Local, e.Peer, e.Err)
	}
	if e.Tx && !e.Req && e.Err == nil {
		return fmt.Sprintf(
			"-> DWA (%s -> %s)", e.Local, e.Peer)
	}
	if e.Tx && !e.Req && e.Err != nil {
		return fmt.Sprintf(
			"-X DWA (%s -> %s), %s", e.Local, e.Peer, e.Err)
	}
	if !e.Tx && e.Req && e.Err == nil {
		return fmt.Sprintf(
			"<- DWR (%s -> %s)", e.Local, e.Peer)
	}
	if !e.Tx && e.Req && e.Err != nil {
		return fmt.Sprintf(
			"X- DWR (%s -> %s), %s", e.Local, e.Peer, e.Err)
	}
	if !e.Tx && !e.Req && e.Err == nil {
		return fmt.Sprintf(
			"<- DWA (%s -> %s)", e.Local, e.Peer)
	}
	if !e.Tx && !e.Req && e.Err != nil {
		return fmt.Sprintf(
			"X- DWA (%s -> %s), %s", e.Local, e.Peer, e.Err)
	}
	return "invalid watchdog state"
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
		return fmt.Sprintf(
			"-> REQ (%s -> %s)", e.Local, e.Peer)
	}
	if e.Tx && e.Req && e.Err != nil {
		return fmt.Sprintf(
			"-X REQ (%s -> %s), %s", e.Local, e.Peer, e.Err)
	}
	if e.Tx && !e.Req && e.Err == nil {
		return fmt.Sprintf(
			"-> ANS (%s -> %s)", e.Local, e.Peer)
	}
	if e.Tx && !e.Req && e.Err != nil {
		return fmt.Sprintf(
			"-X ANS (%s -> %s), %s", e.Local, e.Peer, e.Err)
	}
	if !e.Tx && e.Req && e.Err == nil {
		return fmt.Sprintf(
			"<- REQ (%s -> %s)", e.Local, e.Peer)
	}
	if !e.Tx && e.Req && e.Err != nil {
		return fmt.Sprintf(
			"X- REQ (%s -> %s), %s", e.Local, e.Peer, e.Err)
	}
	if !e.Tx && !e.Req && e.Err == nil {
		return fmt.Sprintf(
			"<- ANS (%s -> %s)", e.Local, e.Peer)
	}
	if !e.Tx && !e.Req && e.Err != nil {
		return fmt.Sprintf(
			"X- ANS (%s -> %s), %s", e.Local, e.Peer, e.Err)
	}
	return "invalid message handle state"
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
		return fmt.Sprintf(
			"-> DPR (%s -> %s)", e.Local, e.Peer)
	}
	if e.Tx && e.Req && e.Err != nil {
		return fmt.Sprintf(
			"-X DPR (%s -> %s), %s", e.Local, e.Peer, e.Err)
	}
	if e.Tx && !e.Req && e.Err == nil {
		return fmt.Sprintf(
			"-> DPA (%s -> %s)", e.Local, e.Peer)
	}
	if e.Tx && !e.Req && e.Err != nil {
		return fmt.Sprintf(
			"-X DPA (%s -> %s), %s", e.Local, e.Peer, e.Err)
	}
	if !e.Tx && e.Req && e.Err == nil {
		return fmt.Sprintf(
			"<- DPR (%s -> %s)", e.Local, e.Peer)
	}
	if !e.Tx && e.Req && e.Err != nil {
		return fmt.Sprintf(
			"X- DPR (%s -> %s), %s", e.Local, e.Peer, e.Err)
	}
	if !e.Tx && !e.Req && e.Err == nil {
		return fmt.Sprintf(
			"<- DPA (%s -> %s)", e.Local, e.Peer)
	}
	if !e.Tx && !e.Req && e.Err != nil {
		return fmt.Sprintf(
			"X- DPA (%s -> %s), %s", e.Local, e.Peer, e.Err)
	}
	return "invalid purge state"
}
