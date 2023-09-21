package diameter

import (
	"log"
	"net"
)

// Direction of Diameter message.
// Tx or Rx.
type Direction bool

func (v Direction) String() string {
	if v {
		return "Tx"
	}
	return "Rx"
}

const (
	Tx Direction = true
	Rx Direction = false
)

// TraceMessage is called when Diameter message is receved or sent.
var TraceMessage = func(msg Message, dct Direction, err error) {
	log.Printf("%s diameter message handling: error=%s\n%s", dct, err, msg)
}

// TraceEvent is called when any event is called.
var TraceEvent = func(old, new, event string, err error) {
	log.Println("diameter state update:", old, "->", new, "by event", event, ": error=", err)
}

// RxQueue returns length of Rx queue
func RxQueue() int {
	return len(rcvQueue)
}

// TxQueue returns length of Tx queue
func TxQueue() int {
	return len(sndQueue)
}

// LocalAddr returns transport connection of state machine
func LocalAddr() net.Addr {
	return conn.LocalAddr()
}

// PeerAddr returns transport connection of state machine
func PeerAddr() net.Addr {
	return conn.RemoteAddr()
}

// State returns state machine state
func State() string {
	return state.String()
}

var (
	// Statistics values
	RxReq     uint64
	RejectReq uint64
	TxAnsFail uint64
	Tx1xxx    uint64
	Tx2xxx    uint64
	Tx3xxx    uint64
	Tx4xxx    uint64
	Tx5xxx    uint64
	TxEtc     uint64

	TxReq      uint64
	InvalidAns uint64
	Rx1xxx     uint64
	Rx2xxx     uint64
	Rx3xxx     uint64
	Rx4xxx     uint64
	Rx5xxx     uint64
	RxEtc      uint64
)

func CountRxCode(c uint32) {
	if c < 1000 {
		RxEtc++
	} else if c < 2000 {
		Rx1xxx++
	} else if c < 3000 {
		Rx2xxx++
	} else if c < 4000 {
		Rx3xxx++
	} else if c < 5000 {
		Rx4xxx++
	} else if c < 6000 {
		Rx5xxx++
	} else {
		RxEtc++
	}
}

func CountTxCode(c uint32) {
	if c < 1000 {
		TxEtc++
	} else if c < 2000 {
		Tx1xxx++
	} else if c < 3000 {
		Tx2xxx++
	} else if c < 4000 {
		Tx3xxx++
	} else if c < 5000 {
		Tx4xxx++
	} else if c < 6000 {
		Tx5xxx++
	} else {
		TxEtc++
	}
}
