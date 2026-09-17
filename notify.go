package diameter

import (
	"net"
)

var (
	// TraceTxMessage is called when Diameter message is sent.
	// Inputs are handled message and occured error while message handling.
	TraceTxMessage func(*Connection, Message, error)

	// TraceRxMessage is called when Diameter message is receved.
	// Inputs are handled message and occured error while message handling.
	TraceRxMessage func(*Connection, Message, error)

	// TraceEvent is called on event.
	// Inputs are old state, event name and occured error while event handling.
	TraceEvent func(*Connection, string, string, error)
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

// EventQueue returns length of event queue
func (c *Connection) EventQueue() int {
	return len(c.notify)
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
