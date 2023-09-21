package sctp

import (
	"fmt"
	"net"
)

// SCTPListener is a SCTP network listener.
type SCTPListener struct {
	sock int
}

// ListenSCTP announces on the SCTP address laddr
// and returns a SCTP listener.
func ListenSCTP(laddr *SCTPAddr) (l *SCTPListener, e error) {
	l = &SCTPListener{}
	if laddr == nil {
		e = fmt.Errorf("no local address")
	} else if laddr.IP[0].To4() != nil {
		l.sock, e = sockOpenV4()
	} else if laddr.IP[0].To16() != nil {
		l.sock, e = sockOpenV6()
	} else {
		e = &net.AddrError{
			Err:  "unknown address format",
			Addr: laddr.String()}
	}

	/*
		type opt struct {
			o uint16
			i uint16
			a uint16
			t uint16
		}
		attr := opt{
			o: uint16(d.OutStream),
			i: uint16(d.InStream),
			a: uint16(d.MaxAttempts),
			t: uint16(d.InitTimeout)}
		l := unsafe.Sizeof(attr)
		p := unsafe.Pointer(&attr)

		e = setSockOpt(sock, sctpInitMsg, p, l)
		if e != nil {
			sockClose(sock)
			return nil, &net.OpError{
				Op:   "setsockopt",Net:  "sctp",
				Addr: d.LocalAddr,Err:  e}
		}
	*/

	// bind SCTP connection
	if e == nil {
		// ptr, n := laddr.rawAddr()
		if e = sctpBindx(l.sock, laddr.rawBytes()); e != nil {
			_ = sockClose(l.sock)
		}
	}
	if e == nil {
		e = sockListen(l.sock)
	}
	if e != nil {
		return nil, &net.OpError{
			Op: "listen", Net: "sctp",
			Addr: laddr, Err: e}
	}

	return
}

// Accept implements the Accept method in the Listener interface;
// it waits for the next call and returns a generic Conn.
func (l *SCTPListener) Accept() (net.Conn, error) {
	return l.AcceptSCTP()
}

// AcceptSCTP accepts the next incoming call and returns the new connection.
func (l *SCTPListener) AcceptSCTP() (c *SCTPConn, e error) {
	c = &SCTPConn{}
	if c.sock, e = sockAccept(l.sock); e != nil {
		e = &net.OpError{
			Op: "accept", Net: "sctp",
			Addr: l.Addr(), Err: e}
	}
	return
}

// Close stops listening on the SCTP address.
func (l *SCTPListener) Close() (e error) {
	return sockClose(l.sock)
}

// Addr returns the listener's network address, a *SCTPAddr.
func (l *SCTPListener) Addr() net.Addr {
	ptr, n, e := sctpGetladdrs(l.sock)
	if e != nil {
		return nil
	}
	defer sctpFreeladdrs(ptr)

	return resolveFromRawAddr(ptr, n)
}
