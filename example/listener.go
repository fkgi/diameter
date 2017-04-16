package provider

/*
// Listener is Diameter server listener
type Listener struct {
	local *LocalNode
	provs map[*PeerNode]*Provider
}

// Listen make Listener with LocalNode
func Listen(n *LocalNode) (l *Listener) {
	//	n.InitIDs()
	l = &Listener{}
	l.local = n
	l.provs = make(map[*PeerNode]*Provider)
	return
}

// Bind bind local address to Listener
func (l *Listener) Bind(lnr net.Listener) (e error) {
		var lnr net.Listener
		if laddr == nil {
			return fmt.Errorf("Local address is nil")
		} else if laddr.Network() == "sctp" {
			a, ok := laddr.(*extnet.SCTPAddr)
			if !ok {
				return fmt.Errorf("invalid sctp address")
			}
			dialer := extnet.SCTPDialer{
				PPID:      46,
				Unordered: true,
				LocalAddr: a}
			if lnr, e = dialer.Listen(); e != nil {
				return e
			}
		} else if laddr.Network() == "tcp" {
			if a, ok := laddr.(*net.TCPAddr); !ok {
				return fmt.Errorf("invalid tcp address")
			} else if lnr, e = net.ListenTCP(laddr.Network(), a); e != nil {
				return e
			}
		} else {
			return fmt.Errorf("invalid address")
		}

	go func() {
		for {
			if lnr == nil {
				break
			}
			c, e := lnr.Accept()
			if e != nil {
				break
			}

			// output logs
			if Notificator != nil {
				Notificator(&TransportStateChange{
					Open: true, Local: l.local, Peer: nil,
					LAddr: lnr.Addr(), PAddr: c.RemoteAddr(),
					Err: e})
			}
			go l.bindProvider(c)
		}
	}()
	return
}
func (l *Listener) bindProvider(c net.Conn) {
	// R-Accept
	if c == nil {
		return
	}
	m := msg.Message{}
	c.SetReadDeadline(time.Time{})
	_, e := m.ReadFrom(c)

	if Notificator != nil {
		Notificator(&MessageTransfer{
			Tx: false, Local: l.local, Peer: nil, Err: e, dump: m.PrintStack})
	}

	if e == nil && !(m.AppID == 0 && m.Code == 257 && m.FlgR) {
		e = fmt.Errorf("not CER")
	}
	if e != nil {
		// output logs
		if Notificator != nil {
			Notificator(&CapabilityExchangeEvent{
				Tx: false, Req: true, Local: l.local, Peer: nil, Err: e})
		}
		c.Close()
		return
	}

	avp, e := m.Decode()
	var h, r msg.DiameterIdentity
	for _, a := range avp {
		if a.Code == uint32(264) && a.VenID == uint32(0) {
			a.Decode(&h)
		}
		if a.Code == uint32(296) && a.VenID == uint32(0) {
			a.Decode(&r)
		}
	}

	for k, v := range l.provs {
		if k.Host == h && k.Realm == r {
			// c.Peer = k
			v.notify <- eventRConnCER{m, c}
			return
		}
	}

	if Notificator != nil {
		e = fmt.Errorf("CER from unknown peer")
		Notificator(&CapabilityExchangeEvent{
			Tx: false, Req: true, Local: l.local, Peer: nil, Err: e})
	}
	c.Close()
}

// AddPeer add new PeerNode to Listener
func (l *Listener) AddPeer(n *PeerNode) (p *Provider) {
	p = &Provider{}

	p.Notify = make(chan stateEvent)

	p.state = shutdown

	p.rcvstack = make(chan *msg.Message, MsgStackLen)
	p.sndstack = make(map[uint32]chan *msg.Message)

	p.local = l.local
	p.peer = n

	go p.run()

	l.provs[n] = p
	return
}

// Dial add new PeerNode to Listener
func (l *Listener) Dial(n *PeerNode, con net.Conn) (e error) {
	p := l.provs[n]
	if p == nil {
		return
	}
	if p.state == shutdown {
		p.state = closed
	}
	p.Notify <- eventStart{con, l.local, n}
	return
}
*/
