package sock

import (
	"bytes"
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
	oldStat state
	newStat state
	stateEvent
	conn *Conn
	Err  error
}

func (e StateUpdate) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w,
		"state change %s -> %s with event %s on connection %s - %s",
		e.oldStat, e.newStat, e.stateEvent, Host, e.conn.peer)
	if e.Err != nil {
		fmt.Fprintf(w, " failed: %s", e.Err)
	}
	return w.String()
}

func msgHandleLog(x, r bool, c *Conn, e error, req, ans string) string {
	w := new(bytes.Buffer)
	if x {
		fmt.Fprintf(w, "-> ")
	} else {
		fmt.Fprintf(w, "<- ")
	}
	if r {
		fmt.Fprintf(w, req)
	} else {
		fmt.Fprintf(w, ans)
	}
	fmt.Fprintf(w, "(%s -> %s)", Host, c.peer)
	if e != nil {
		fmt.Fprintf(w, " failed: %s", e)
	}
	return w.String()
}

// CapabilityExchangeEvent notify capability exchange related event
type CapabilityExchangeEvent struct {
	tx   bool
	req  bool
	conn *Conn
	Err  error
}

func (e CapabilityExchangeEvent) String() string {
	return msgHandleLog(e.tx, e.req, e.conn, e.Err, "CER", "CEA")
}

// WatchdogEvent notify watchdog related event
type WatchdogEvent struct {
	tx   bool
	req  bool
	conn *Conn
	Err  error
}

func (e WatchdogEvent) String() string {
	return msgHandleLog(e.tx, e.req, e.conn, e.Err, "DWR", "DWA")
}

// MessageEvent notify diameter message related event
type MessageEvent struct {
	tx   bool
	req  bool
	conn *Conn
	Err  error
}

func (e MessageEvent) String() string {
	return msgHandleLog(e.tx, e.req, e.conn, e.Err, "REQ", "ANS")
}

// PurgeEvent notify diameter purge related event
type PurgeEvent struct {
	tx   bool
	req  bool
	conn *Conn
	Err  error
}

func (e PurgeEvent) String() string {
	return msgHandleLog(e.tx, e.req, e.conn, e.Err, "DPR", "DPA")
}
