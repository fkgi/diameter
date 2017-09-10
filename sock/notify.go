package sock

import (
	"fmt"
	"log"
)

// Notify is called when error or trace event are occured
var Notify = func(n Notice) {
	log.Println(n)
}

// Notice is notification information from connection
type Notice interface {
	String() string
}

// StateUpdate notify event
type StateUpdate struct {
	OldState string
	NewState string
	Event    string
	Local    *Local
	Peer     *Peer
	Err      error
}

func (e *StateUpdate) String() string {
	if e == nil {
		return "<nil>"
	} else if e.Err == nil {
		return fmt.Sprintf(
			"state change %s -> %s with event %s on connection %s - %s",
			e.OldState, e.NewState, e.Event, e.Local, e.Peer)
	} else {
		return fmt.Sprintf(
			"state change %s -> %s with event %s on connection %s - %s failed: %s",
			e.OldState, e.NewState, e.Event, e.Local, e.Peer, e.Err)
	}
}

// CapabilityExchangeEvent notify capability exchange related event
type CapabilityExchangeEvent struct {
	Tx    bool
	Req   bool
	Local *Local
	Peer  *Peer
	Err   error
}

func (e *CapabilityExchangeEvent) String() string {
	if e == nil {
		return "<nil>"
	} else if e.Tx && e.Req && e.Err == nil {
		return fmt.Sprintf("-> CER (%s -> %s)", e.Local, e.Peer)
	} else if e.Tx && e.Req && e.Err != nil {
		return fmt.Sprintf("-X CER (%s -> %s), %s", e.Local, e.Peer, e.Err)
	} else if e.Tx && !e.Req && e.Err == nil {
		return fmt.Sprintf("-> CEA (%s -> %s)", e.Local, e.Peer)
	} else if e.Tx && !e.Req && e.Err != nil {
		return fmt.Sprintf("-X CEA (%s -> %s), %s", e.Local, e.Peer, e.Err)
	} else if !e.Tx && e.Req && e.Err == nil {
		return fmt.Sprintf("<- CER (%s -> %s)", e.Local, e.Peer)
	} else if !e.Tx && e.Req && e.Err != nil {
		return fmt.Sprintf("X- CER (%s -> %s), %s", e.Local, e.Peer, e.Err)
	} else if !e.Tx && !e.Req && e.Err == nil {
		return fmt.Sprintf("<- CEA (%s -> %s)", e.Local, e.Peer)
	} else if !e.Tx && !e.Req && e.Err != nil {
		return fmt.Sprintf("X- CEA (%s -> %s), %s", e.Local, e.Peer, e.Err)
	}
	return "invalid event"
}

// WatchdogEvent notify watchdog related event
type WatchdogEvent struct {
	Tx    bool
	Req   bool
	Local *Local
	Peer  *Peer
	Err   error
}

func (e *WatchdogEvent) String() string {
	if e == nil {
		return "<nil>"
	} else if e.Tx && e.Req && e.Err == nil {
		return fmt.Sprintf("-> DWR (%s -> %s)", e.Local, e.Peer)
	} else if e.Tx && e.Req && e.Err != nil {
		return fmt.Sprintf("-X DWR (%s -> %s), %s", e.Local, e.Peer, e.Err)
	} else if e.Tx && !e.Req && e.Err == nil {
		return fmt.Sprintf("-> DWA (%s -> %s)", e.Local, e.Peer)
	} else if e.Tx && !e.Req && e.Err != nil {
		return fmt.Sprintf("-X DWA (%s -> %s), %s", e.Local, e.Peer, e.Err)
	} else if !e.Tx && e.Req && e.Err == nil {
		return fmt.Sprintf("<- DWR (%s -> %s)", e.Local, e.Peer)
	} else if !e.Tx && e.Req && e.Err != nil {
		return fmt.Sprintf("X- DWR (%s -> %s), %s", e.Local, e.Peer, e.Err)
	} else if !e.Tx && !e.Req && e.Err == nil {
		return fmt.Sprintf("<- DWA (%s -> %s)", e.Local, e.Peer)
	} else if !e.Tx && !e.Req && e.Err != nil {
		return fmt.Sprintf("X- DWA (%s -> %s), %s", e.Local, e.Peer, e.Err)
	}
	return "invalid event"
}

// MessageEvent notify diameter message related event
type MessageEvent struct {
	Tx    bool
	Req   bool
	Local *Local
	Peer  *Peer
	Err   error
}

func (e *MessageEvent) String() string {
	if e == nil {
		return "<nil>"
	} else if e.Tx && e.Req && e.Err == nil {
		return fmt.Sprintf("-> REQ (%s -> %s)", e.Local, e.Peer)
	} else if e.Tx && e.Req && e.Err != nil {
		return fmt.Sprintf("-X REQ (%s -> %s), %s", e.Local, e.Peer, e.Err)
	} else if e.Tx && !e.Req && e.Err == nil {
		return fmt.Sprintf("-> ANS (%s -> %s)", e.Local, e.Peer)
	} else if e.Tx && !e.Req && e.Err != nil {
		return fmt.Sprintf("-X ANS (%s -> %s), %s", e.Local, e.Peer, e.Err)
	} else if !e.Tx && e.Req && e.Err == nil {
		return fmt.Sprintf("<- REQ (%s -> %s)", e.Local, e.Peer)
	} else if !e.Tx && e.Req && e.Err != nil {
		return fmt.Sprintf("X- REQ (%s -> %s), %s", e.Local, e.Peer, e.Err)
	} else if !e.Tx && !e.Req && e.Err == nil {
		return fmt.Sprintf("<- ANS (%s -> %s)", e.Local, e.Peer)
	} else if !e.Tx && !e.Req && e.Err != nil {
		return fmt.Sprintf("X- ANS (%s -> %s), %s", e.Local, e.Peer, e.Err)
	}
	return "invalid event"
}

// PurgeEvent notify diameter purge related event
type PurgeEvent struct {
	Tx    bool
	Req   bool
	Local *Local
	Peer  *Peer
	Err   error
}

func (e *PurgeEvent) String() string {
	if e == nil {
		return "<nil>"
	} else if e.Tx && e.Req && e.Err == nil {
		return fmt.Sprintf("-> DPR (%s -> %s)", e.Local, e.Peer)
	} else if e.Tx && e.Req && e.Err != nil {
		return fmt.Sprintf("-X DPR (%s -> %s), %s", e.Local, e.Peer, e.Err)
	} else if e.Tx && !e.Req && e.Err == nil {
		return fmt.Sprintf("-> DPA (%s -> %s)", e.Local, e.Peer)
	} else if e.Tx && !e.Req && e.Err != nil {
		return fmt.Sprintf("-X DPA (%s -> %s), %s", e.Local, e.Peer, e.Err)
	} else if !e.Tx && e.Req && e.Err == nil {
		return fmt.Sprintf("<- DPR (%s -> %s)", e.Local, e.Peer)
	} else if !e.Tx && e.Req && e.Err != nil {
		return fmt.Sprintf("X- DPR (%s -> %s), %s", e.Local, e.Peer, e.Err)
	} else if !e.Tx && !e.Req && e.Err == nil {
		return fmt.Sprintf("<- DPA (%s -> %s)", e.Local, e.Peer)
	} else if !e.Tx && !e.Req && e.Err != nil {
		return fmt.Sprintf("X- DPA (%s -> %s), %s", e.Local, e.Peer, e.Err)
	}
	return "invalid event"
}
