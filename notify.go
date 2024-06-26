package diameter

import (
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

var (
	// TraceMessage is called when Diameter message is receved or sent.
	// Inputs are handled message, message direction and occured error while message handling.
	TraceMessage func(Message, Direction, error)

	// TraceEvent is called on event.
	// Inputs are old state, new state, event name and occured error while event handling.
	TraceEvent func(string, string, string, error)

	// ConnectionUpNotify is called when Diameter connection up.
	ConnectionUpNotify func(*Connection)

	// ConnectionDownNotify is called when Diameter connection down.
	ConnectionDownNotify func(*Connection)

	// ConnectionAbortNotify is called when Diameter connection abort.
	ConnectionAbortNotify func(*Connection)
)

// RxQueue returns length of Rx queue
func (c *Connection) RxQueue() int {
	return len(c.rcvQueue)
}

// TxQueue returns length of Tx queue
func (c *Connection) TxQueue() int {
	return len(c.sndQueue)
}

// LocalAddr returns transport connection of state machine
func (c *Connection) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}

// PeerAddr returns transport connection of state machine
func (c *Connection) PeerAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// State returns state machine state
func (c *Connection) State() string {
	return c.state.String()
}

// AvailableApplications returns supported application list
func (c *Connection) AvailableApplications() []uint32 {
	ret := []uint32{}
	for k := range c.commonApp {
		ret = append(ret, k)
	}
	return ret
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
