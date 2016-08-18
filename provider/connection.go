package provider

import (
	"fmt"
	"monitor"
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
	if c.conn == nil {
		e = fmt.Errorf("connection is closed")
	} else {
		tmp := c.conn
		c.conn = nil
		e = tmp.Close()
	}

	// output logs
	if e == nil {
		lh, ph := c.hostnames()
		monitor.Notify(monitor.Info, "close connection", lh, ph)
	} else {
		monitor.Notify(monitor.Major, "close connection failed", e.Error())
	}
	return
}

// Low-level message writer
func (c *Connection) Write(s time.Duration, m msg.Message) (e error) {
	if c.conn == nil {
		e = fmt.Errorf("connection is closed")
	} else {
		t := time.Time{}
		if s != time.Duration(0) {
			t = time.Now().Add(s)
		}
		c.conn.SetWriteDeadline(t)
		_, e = m.WriteTo(c.conn)
	}

	if e == nil {
		lh, ph := c.hostnames()
		monitor.Notify(monitor.Trace, "write message data", lh, ph)
		monitor.Dump("== Diameter message stack ==", m.PrintStack)
	} else {
		monitor.Notify(monitor.Minor, "write message failed", e.Error())
	}
	return
}

// Low-level message reader
func (c *Connection) Read(s time.Duration) (m msg.Message, e error) {
	if c.conn == nil {
		e = fmt.Errorf("connection is closed")
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

	if e == nil {
		lh, ph := c.hostnames()
		monitor.Notify(monitor.Trace, "read message data", lh, ph)
		monitor.Dump("== Diameter message stack ==", m.PrintStack)
	} else {
		monitor.Notify(monitor.Minor, "read message failed", e.Error())
	}
	return
}

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
