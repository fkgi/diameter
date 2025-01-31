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

func (c *Connection) countRxCode(r uint32) {
	if r < 1000 {
		c.RxAns[0]++
	} else if r < 2000 {
		c.RxAns[1]++
	} else if r < 3000 {
		c.RxAns[2]++
	} else if r < 4000 {
		c.RxAns[3]++
	} else if r < 5000 {
		c.RxAns[4]++
	} else if r < 6000 {
		c.RxAns[5]++
	} else {
		c.RxAns[0]++
	}
}

func (c *Connection) countTxCode(r uint32) {
	if r < 1000 {
		c.TxAns[0]++
	} else if r < 2000 {
		c.TxAns[1]++
	} else if r < 3000 {
		c.TxAns[2]++
	} else if r < 4000 {
		c.TxAns[3]++
	} else if r < 5000 {
		c.TxAns[4]++
	} else if r < 6000 {
		c.TxAns[5]++
	} else {
		c.TxAns[0]++
	}
}
