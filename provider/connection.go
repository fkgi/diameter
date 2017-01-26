package provider

import (
	"fmt"
	"net"
	"time"

	"github.com/fkgi/diameter/msg"
)

// Connection is Diameter connection
type Connection struct {
	Peer  *PeerNode
	Local *LocalNode
	conn  net.Conn
}

// IsAvailable returns availability of Connection
func (c *Connection) IsAvailable() bool {
	return c.conn != nil
}

// Close close Diameter connection
func (c *Connection) Close() (e error) {
	var la, pa net.Addr
	if c.conn == nil {
		e = fmt.Errorf("connection has been closed")
	} else {
		la = c.conn.LocalAddr()
		pa = c.conn.RemoteAddr()
		tmp := c.conn
		c.conn = nil
		e = tmp.Close()
	}

	// output logs
	if Notificator != nil {
		Notificator(&TransportStateChange{
			Open: false, Local: c.Local, Peer: c.Peer, LAddr: la, PAddr: pa, Err: e})
	}
	return
}

// Low-level message writer
func (c *Connection) Write(s time.Duration, m msg.Message) (e error) {
	if c.conn == nil {
		e = fmt.Errorf("connection has been closed")
	} else {
		t := time.Time{}
		if s != time.Duration(0) {
			t = time.Now().Add(s)
		}
		c.conn.SetWriteDeadline(t)
		_, e = m.WriteTo(c.conn)
	}

	if Notificator != nil {
		Notificator(&MessageTransfer{
			Tx: true, Local: c.Local, Peer: c.Peer, Err: e, dump: m.PrintStack})
	}
	return
}

// Low-level message reader
func (c *Connection) Read(s time.Duration) (m msg.Message, e error) {
	if c.conn == nil {
		e = fmt.Errorf("connection has been closed")
	} else {
		t := time.Time{}
		if s != time.Duration(0) {
			t = time.Now().Add(s)
		}
		c.conn.SetReadDeadline(t)
		_, e = m.ReadFrom(c.conn)
		/*
			if ne, ok := e.(net.Error); ok && ne.Timeout() {
			}
		*/
	}

	if Notificator != nil {
		Notificator(&MessageTransfer{
			Tx: false, Local: c.Local, Peer: c.Peer, Err: e, dump: m.PrintStack})
	}
	return
}

/*
func (c *Connection) hostnames() (lh, ph string) {
	if c.Local == nil {
		lh = "unknown"
	} else {
		lh = string(c.Local.Host)
	}
	if c.Peer == nil {
		ph = "unknown"
	} else {
		ph = string(c.Peer.Host)
	}
	return
}
*/
