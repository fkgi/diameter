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
	ConnectionDownNotify func(*Connection, error)
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
	if c.state == open {
		for k := range c.commonApp {
			ret = append(ret, k)
		}
	}
	return ret
}

// SharedMessagegQueue return lengh of shared queue for recieved stateless message handling.
func SharedMessagegQueue() int {
	return len(sharedQ)
}

// ActiveSharedWorkers return count of active worker for recieved stateless message handling.
func ActiveSharedWorkers() int {
	a := <-activeWorkers
	activeWorkers <- a
	return a
}
