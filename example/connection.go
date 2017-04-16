package provider

// Connection is Diameter connection
/*
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
		//if ne, ok := e.(net.Error); ok && ne.Timeout() {
		//	}
	}

	if Notificator != nil {
		Notificator(&MessageTransfer{
			Tx: false, Local: c.Local, Peer: c.Peer, Err: e, dump: m.PrintStack})
	}
	return
}
*/

// Connect is Low-level diameter connect    laddr, raddr net.Addr, s time.Duration
/*
func (l *LocalNode) Connect(p *PeerNode, con net.Conn) (c *Connection, e error) {
	if p == nil {
		e = fmt.Errorf("Peer node is nil")
	}
	c = &Connection{p, l, con}
			} else if laddr.Network() == "sctp" {
				if a, ok := laddr.(*extnet.SCTPAddr); !ok {
					e = fmt.Errorf("address type mismatch")
				} else {
					dialer := extnet.SCTPDialer{
						InitTimeout: s,
						PPID:        46,
						Unordered:   true,
						LocalAddr:   a}
					var con net.Conn
					if con, e = dialer.Dial(raddr.Network(), raddr.String()); e == nil {
						c = &Connection{p, l, con}
					}
				}

			} else if laddr.Network() == "tcp" {
				dialer := net.Dialer{
					Timeout:   s,
					LocalAddr: laddr}

				var con net.Conn
				if con, e = dialer.Dial(raddr.Network(), raddr.String()); e == nil {
					c = &Connection{p, l, con}
				}
			} else {
				e = fmt.Errorf("invalid address type")
			}
		// output logs
		if Notificator != nil {
			Notificator(&TransportStateChange{
				Open: true, Local: l, Peer: p, LAddr: laddr, PAddr: raddr, Err: e})
		}
	return
}
*/

// Accept is Low-level diameter accept
/*
func (l *LocalNode) Accept(lnr net.Listener) (c *Connection, e error) {
	if lnr == nil {
		e = fmt.Errorf("Local listener is nil")
	} else {
		var con net.Conn
		if con, e = lnr.Accept(); e == nil {
			c = &Connection{nil, l, con}
		}
	}

	// output logs
	if Notificator != nil {
		var pa net.Addr
		if e == nil {
			pa = c.conn.RemoteAddr()
		}
		Notificator(&TransportStateChange{
			Open: true, Local: l, Peer: nil, LAddr: lnr.Addr(), PAddr: pa, Err: e})
	}
	return
}
*/