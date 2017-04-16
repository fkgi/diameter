package connection

import "log"

// Notificator is called when error or trace event are occured
var Notificator func(n Notify)

// Notify is notification information from connection
type Notify interface {
	// Log generage log text of this notify
	Log(l *log.Logger)
}

// StateUpdate notify event
type StateUpdate struct {
	OldState string
	NewState string
	Event    string
	Local    *LocalNode
	Peer     *PeerNode
	Err      error
}

// Log generate log text of this notify
func (e *StateUpdate) Log(l *log.Logger) {
	if e == nil {
		l.Println("<nil>")
	} else if e.Err == nil {
		l.Printf(
			"state change %s to %s with event %s on Connection %s -> %s",
			e.OldState, e.NewState, e.Event, e.Local, e.Peer)
	} else {
		l.Printf(
			"state change %s to %s with event %s on Connection %s -> %s failed: %s",
			e.OldState, e.NewState, e.Event, e.Local, e.Peer, e.Err)
	}
}

// CapabilityExchangeEvent notify capability exchange related event
type CapabilityExchangeEvent struct {
	Tx    bool
	Req   bool
	Local *LocalNode
	Peer  *PeerNode
	Err   error
}

// Log generate log text of this notify
func (e *CapabilityExchangeEvent) Log(l *log.Logger) {
	if e == nil {
		l.Println("<nil>")
	} else if e.Tx && e.Req && e.Err == nil {
		l.Printf("-> CER (%s -> %s)", e.Local, e.Peer)
	} else if e.Tx && e.Req && e.Err != nil {
		l.Printf("-X CER (%s -> %s), %s", e.Local, e.Peer, e.Err)
	} else if e.Tx && !e.Req && e.Err == nil {
		l.Printf("-> CEA (%s -> %s)", e.Local, e.Peer)
	} else if e.Tx && !e.Req && e.Err != nil {
		l.Printf("-X CEA (%s -> %s), %s", e.Local, e.Peer, e.Err)
	} else if !e.Tx && e.Req && e.Err == nil {
		l.Printf("<- CER (%s -> %s)", e.Local, e.Peer)
	} else if !e.Tx && e.Req && e.Err != nil {
		l.Printf("X- CER (%s -> %s), %s", e.Local, e.Peer, e.Err)
	} else if !e.Tx && !e.Req && e.Err == nil {
		l.Printf("<- CEA (%s -> %s)", e.Local, e.Peer)
	} else if !e.Tx && !e.Req && e.Err != nil {
		l.Printf("X- CEA (%s -> %s), %s", e.Local, e.Peer, e.Err)
	}
}

// WatchdogEvent notify watchdog related event
type WatchdogEvent struct {
	Tx    bool
	Req   bool
	Local *LocalNode
	Peer  *PeerNode
	Err   error
}

// Log generate log text of this notify
func (e *WatchdogEvent) Log(l *log.Logger) {
	if e == nil {
		l.Println("<nil>")
	} else if e.Tx && e.Req && e.Err == nil {
		l.Printf("-> DWR (%s -> %s)", e.Local, e.Peer)
	} else if e.Tx && e.Req && e.Err != nil {
		l.Printf("-X DWR (%s -> %s), %s", e.Local, e.Peer, e.Err)
	} else if e.Tx && !e.Req && e.Err == nil {
		l.Printf("-> DWA (%s -> %s)", e.Local, e.Peer)
	} else if e.Tx && !e.Req && e.Err != nil {
		l.Printf("-X DWA (%s -> %s), %s", e.Local, e.Peer, e.Err)
	} else if !e.Tx && e.Req && e.Err == nil {
		l.Printf("<- DWR (%s -> %s)", e.Local, e.Peer)
	} else if !e.Tx && e.Req && e.Err != nil {
		l.Printf("X- DWR (%s -> %s), %s", e.Local, e.Peer, e.Err)
	} else if !e.Tx && !e.Req && e.Err == nil {
		l.Printf("<- DWA (%s -> %s)", e.Local, e.Peer)
	} else if !e.Tx && !e.Req && e.Err != nil {
		l.Printf("X- DWA (%s -> %s), %s", e.Local, e.Peer, e.Err)
	}
}

// MessageEvent notify diameter message related event
type MessageEvent struct {
	Tx    bool
	Req   bool
	Local *LocalNode
	Peer  *PeerNode
	Err   error
}

// Log generate log text of this notify
func (e *MessageEvent) Log(l *log.Logger) {
	if e == nil {
		l.Println("<nil>")
	} else if e.Tx && e.Req && e.Err == nil {
		l.Printf("-> REQ (%s -> %s)", e.Local, e.Peer)
	} else if e.Tx && e.Req && e.Err != nil {
		l.Printf("-X REQ (%s -> %s), %s", e.Local, e.Peer, e.Err)
	} else if e.Tx && !e.Req && e.Err == nil {
		l.Printf("-> ANS (%s -> %s)", e.Local, e.Peer)
	} else if e.Tx && !e.Req && e.Err != nil {
		l.Printf("-X ANS (%s -> %s), %s", e.Local, e.Peer, e.Err)
	} else if !e.Tx && e.Req && e.Err == nil {
		l.Printf("<- REQ (%s -> %s)", e.Local, e.Peer)
	} else if !e.Tx && e.Req && e.Err != nil {
		l.Printf("X- REQ (%s -> %s), %s", e.Local, e.Peer, e.Err)
	} else if !e.Tx && !e.Req && e.Err == nil {
		l.Printf("<- ANS (%s -> %s)", e.Local, e.Peer)
	} else if !e.Tx && !e.Req && e.Err != nil {
		l.Printf("X- ANS (%s -> %s), %s", e.Local, e.Peer, e.Err)
	}
}

// PurgeEvent notify diameter purge related event
type PurgeEvent struct {
	Tx    bool
	Req   bool
	Local *LocalNode
	Peer  *PeerNode
	Err   error
}

// Log generate log text of this notify
func (e *PurgeEvent) Log(l *log.Logger) {
	if e == nil {
		l.Println("<nil>")
	} else if e.Tx && e.Req && e.Err == nil {
		l.Printf("-> DPR (%s -> %s)", e.Local, e.Peer)
	} else if e.Tx && e.Req && e.Err != nil {
		l.Printf("-X DPR (%s -> %s), %s", e.Local, e.Peer, e.Err)
	} else if e.Tx && !e.Req && e.Err == nil {
		l.Printf("-> DPA (%s -> %s)", e.Local, e.Peer)
	} else if e.Tx && !e.Req && e.Err != nil {
		l.Printf("-X DPA (%s -> %s), %s", e.Local, e.Peer, e.Err)
	} else if !e.Tx && e.Req && e.Err == nil {
		l.Printf("<- DPR (%s -> %s)", e.Local, e.Peer)
	} else if !e.Tx && e.Req && e.Err != nil {
		l.Printf("X- DPR (%s -> %s), %s", e.Local, e.Peer, e.Err)
	} else if !e.Tx && !e.Req && e.Err == nil {
		l.Printf("<- DPA (%s -> %s)", e.Local, e.Peer)
	} else if !e.Tx && !e.Req && e.Err != nil {
		l.Printf("X- DPA (%s -> %s), %s", e.Local, e.Peer, e.Err)
	}
}
