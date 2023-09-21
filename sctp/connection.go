package sctp

import (
	"fmt"
	"net"
	"syscall"
	"time"
)

const (
	ProtocolID uint32 = 46 // Diameter
)

// SCTPConn is an implementation of the Conn interface for SCTP network connections.
type SCTPConn struct {
	sock int
}

// DialSCTP connects from the local address laddr
// to the remote address raddr.
func DialSCTP(laddr, raddr *SCTPAddr) (c *SCTPConn, e error) {
	c = &SCTPConn{}

	if laddr == nil {
		e = fmt.Errorf("nil local address")
	} else if laddr.IP[0].To4() != nil {
		c.sock, e = sockOpenV4()
	} else if laddr.IP[0].To16() != nil {
		c.sock, e = sockOpenV6()
	} else {
		e = &net.AddrError{
			Err:  "unknown address format",
			Addr: laddr.String()}
	}
	if e == nil {
		if e = sctpBindx(c.sock, laddr.rawBytes()); e != nil {
			_ = sockClose(c.sock)
		}
	}

	if e != nil {
	} else if raddr == nil {
		e = fmt.Errorf("nil peer address")
	} else {
		if e = sctpConnectx(c.sock, raddr.rawBytes()); e != nil {
			_ = sockClose(c.sock)
		}
	}

	if e != nil {
		e = &net.OpError{
			Op: "dial", Net: "sctp",
			Source: laddr, Addr: raddr, Err: e}
	}
	return
}

func (c *SCTPConn) Read(b []byte) (n int, e error) {
	if n, e = sctpRecvmsg(c.sock, b); e != nil {
		e = &net.OpError{
			Op: "read", Net: "sctp",
			Source: c.LocalAddr(), Addr: c.RemoteAddr(), Err: e}
	}
	return
}

func (c *SCTPConn) Write(b []byte) (n int, e error) {
	buf := make([]byte, len(b))
	copy(buf, b)

	if n, e = sctpSend(c.sock, b); e != nil {
		e = &net.OpError{
			Op: "write", Net: "sctp",
			Source: c.LocalAddr(), Addr: c.RemoteAddr(), Err: e}
	}
	return
}

// Close closes the connection.
func (c *SCTPConn) Close() (e error) {
	if e = sockClose(c.sock); e != nil {
		e = &net.OpError{
			Op: "close", Net: "sctp",
			Source: c.LocalAddr(), Addr: c.RemoteAddr(), Err: e}
	}
	return e
}

// LocalAddr returns the local network address.
func (c *SCTPConn) LocalAddr() net.Addr {
	ptr, n, e := sctpGetladdrs(c.sock)
	defer sctpFreeladdrs(ptr)
	if e != nil {
		return nil
	}
	return resolveFromRawAddr(ptr, n)
}

// RemoteAddr returns the remote network address.
func (c *SCTPConn) RemoteAddr() net.Addr {
	ptr, n, e := sctpGetpaddrs(c.sock)
	defer sctpFreepaddrs(ptr)
	if e != nil {
		return nil
	}
	return resolveFromRawAddr(ptr, n)
}

// SetDeadline implements the Conn SetDeadline method.
func (c *SCTPConn) SetDeadline(t time.Time) error {
	return syscall.EOPNOTSUPP
}

// SetReadDeadline implements the Conn SetReadDeadline method.
func (c *SCTPConn) SetReadDeadline(t time.Time) error {
	return syscall.EOPNOTSUPP
}

// SetWriteDeadline implements the Conn SetWriteDeadline method.
func (c *SCTPConn) SetWriteDeadline(t time.Time) error {
	return syscall.EOPNOTSUPP
}

/*
// SetRtoInfo set retransmit timer options
func (c *SCTPConn) SetRtoInfo(ini, min, max int) error {
	type opt struct {
		assocID assocT
		ini     uint32
		max     uint32
		min     uint32
	}
	attr := opt{
		assocID: c.id,
		ini:     uint32(ini),
		max:     uint32(max),
		min:     uint32(min)}
	l := unsafe.Sizeof(attr)
	p := unsafe.Pointer(&attr)

	return setSockOpt(c.l.sock, sctpRtoInfo, p, l)
}

// SetAssocinfo set association parameter
func (c *SCTPConn) SetAssocinfo(pRwnd, lRwnd, cLife, assocMaxRxt, numPeerDest int) error {
	type opt struct {
		assocID     assocT
		pRwnd       uint32
		lRwnd       uint32
		cLife       uint32
		assocMaxRxt uint16
		numPeerDest uint16
	}
	attr := opt{
		assocID:     c.id,
		pRwnd:       uint32(pRwnd),
		lRwnd:       uint32(lRwnd),
		cLife:       uint32(cLife),
		assocMaxRxt: uint16(assocMaxRxt),
		numPeerDest: uint16(numPeerDest)}
	l := unsafe.Sizeof(attr)
	p := unsafe.Pointer(&attr)

	return setSockOpt(c.l.sock, sctpAssocInfo, p, l)
}

// SetNodelay set delay answer or not
func (c *SCTPConn) SetNodelay(attr bool) error {
	l := unsafe.Sizeof(attr)
	p := unsafe.Pointer(&attr)

	return setSockOpt(c.l.sock, sctpNodelay, p, l)
}

const (
	HB_ENABLE         = uint32(C.SPP_HB_ENABLE)
	HB_DISABLE        = uint32(C.SPP_HB_DISABLE)
	HB_DEMAND         = uint32(C.SPP_HB_DEMAND)
	PMTUD_ENABLE      = uint32(C.SPP_PMTUD_ENABLE)
	PMTUD_DISABLE     = uint32(C.SPP_PMTUD_DISABLE)
	SACKDELAY_ENABLE  = uint32(C.SPP_SACKDELAY_ENABLE)
	SACKDELAY_DISABLE = uint32(C.SPP_SACKDELAY_DISABLE)
	HB_TIME_IS_ZERO   = uint32(C.SPP_HB_TIME_IS_ZERO)
)

func (c *SCTPConn) SetPeerAddrParams(
	hbinterval uint32, pathmaxrxt uint16, pathmtu, sackdelay, flags uint32) error {
	attr := C.struct_sctp_paddrparams{}
	l := C.socklen_t(unsafe.Sizeof(attr))

	attr.spp_assoc_id = c.id
	attr.spp_hbinterval = C.__u32(hbinterval)
	attr.spp_pathmaxrxt = C.__u16(pathmaxrxt)
	//attr.spp_pathmtu = C.__u32(pathmtu)
	//attr.spp_sackdelay = C.__u32(sackdelay)
	//attr.spp_flags = C.__u32(flags)

	p := unsafe.Pointer(&attr)
	i, e := C.setsockopt(c.sock, C.SOL_SCTP, C.SCTP_PEER_ADDR_PARAMS, p, l)
	if int(i) < 0 {
		return e
	}
	return nil
}
*/
