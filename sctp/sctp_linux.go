//go:build linux && !386

package sctp

import (
	"bytes"
	"encoding/binary"
	"io"
	"syscall"
	"unsafe"
)

/*
	func setNotify(fd int) error {
		type opt struct {
			dataIo          uint8
			association     uint8
			address         uint8
			sendFailed      uint8
			peerError       uint8
			shutdown        uint8
			partialDelivery uint8
			adaptationLayer uint8
			authentication  uint8
			senderDry       uint8
		}

		event := opt{
			dataIo:          1,
			association:     1,
			address:         0,
			sendFailed:      0,
			peerError:       0,
			shutdown:        0,
			partialDelivery: 0,
			adaptationLayer: 0,
			authentication:  0,
			senderDry:       0}
		l := unsafe.Sizeof(event)
		p := unsafe.Pointer(&event)

	_, _, e := syscall.Syscall6(syscall.SYS_SETSOCKOPT,
		uintptr(fd),
		132, // SOL_SCTP
		11, // SCTP_EVENTS
		uintptr(p),
		uintptr(l),
		0)
	}
*/

func sockOpenV4(seq bool) (int, error) {
	if seq {
		return syscall.Socket(
			syscall.AF_INET,
			syscall.SOCK_SEQPACKET,
			syscall.IPPROTO_SCTP)
	} else {
		return syscall.Socket(
			syscall.AF_INET,
			syscall.SOCK_STREAM,
			syscall.IPPROTO_SCTP)
	}
}

func sockOpenV6(seq bool) (int, error) {
	if seq {
		return syscall.Socket(
			syscall.AF_INET6,
			syscall.SOCK_SEQPACKET,
			syscall.IPPROTO_SCTP)
	} else {
		return syscall.Socket(
			syscall.AF_INET6,
			syscall.SOCK_STREAM,
			syscall.IPPROTO_SCTP)
	}
}

func sockListen(fd int) error {
	if e := syscall.Listen(fd, 1024); e != nil {
		return e
	}
	return syscall.SetNonblock(fd, true)
}

func sockAccept(fd int) (nfd int, e error) {
	nfd, _, e = syscall.Accept(fd)
	return
}

func sockClose(fd int) error {
	syscall.Shutdown(fd, syscall.SHUT_RDWR)
	return syscall.Close(fd)
}

func sctpBindx(fd int, addr []byte) error {
	if _, _, e := syscall.Syscall6(syscall.SYS_SETSOCKOPT,
		uintptr(fd),
		syscall.IPPROTO_SCTP,
		100, // SCTP_SOCKOPT_BINDX_ADD
		uintptr(unsafe.Pointer(&addr[0])),
		uintptr(len(addr)),
		0); e != 0 {
		return e
	}
	return nil
}

func sctpConnectx(fd int, addr []byte) (int, error) {
	t, _, e := syscall.Syscall6(syscall.SYS_SETSOCKOPT,
		uintptr(fd),
		syscall.IPPROTO_SCTP,
		110, // SCTP_SOCKOPT_CONNECTX
		uintptr(unsafe.Pointer(&addr[0])),
		uintptr(len(addr)),
		0)
	/*
		if e == syscall.EINPROGRESS {
			fdset := &syscall.FdSet{}
			fdset.Bits[fd/64] |= 1 << (uint(fd) % 64)
			to := &syscall.Timeval{Sec: 5, Usec: 0}

			n, e := syscall.Select(fd+1, nil, fdset, nil, to)
			if e != nil {
				return 0, e
			} else if n == 0 {
				return 0, errors.New("timeout")
			} else if n, e = syscall.GetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_ERROR); e != nil {
				return 0, e
			} else if n != 0 {
				return 0, syscall.Errno(n)
			}
		} else if e != 0 {
	*/
	if e != 0 {
		return 0, e
	}

	peel := struct {
		aid int32
		sd  int32
	}{aid: int32(t)}
	l := unsafe.Sizeof(peel)
	if _, _, e := syscall.Syscall6(syscall.SYS_GETSOCKOPT,
		uintptr(fd),
		syscall.IPPROTO_SCTP,
		102, // SCTP_SOCKOPT_PEELOFF
		uintptr(unsafe.Pointer(&peel)),
		uintptr(unsafe.Pointer(&l)),
		0); e != 0 {
		return 0, e
	}
	return int(peel.sd), nil
}

func sctpSend(fd int, b []byte) (int, error) {
	hdr := &syscall.Cmsghdr{
		Level: syscall.IPPROTO_SCTP,
		Type:  2, //SCTP_SNDINFO
	}
	hdr.SetLen(syscall.CmsgSpace(16))

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, hdr)
	binary.Write(buf, binary.LittleEndian, uint16(0))      // stream ID(2 byte)=0
	binary.Write(buf, binary.LittleEndian, uint16(0x0001)) // flag(2 byte)=SCTP_UNORDERED
	binary.Write(buf, binary.BigEndian, uint32(46))        // PPID(4 byte)=diameter(46)
	binary.Write(buf, binary.LittleEndian, uint32(0))      // context(4 byte) = empty
	binary.Write(buf, binary.LittleEndian, uint32(0))      // assoc ID(4 byte)

	return syscall.SendmsgN(fd, b, buf.Bytes(), nil, syscall.MSG_DONTWAIT|syscall.MSG_EOR)
}

func sctpRecvmsg(fd int, b []byte) (int, error) {
	n, on, _, _, e := syscall.Recvmsg(fd, b, make([]byte, syscall.CmsgSpace(32)), 0)
	if e == nil && n <= 0 && on <= 0 {
		e = io.EOF
	}
	return n, e
}

func sctpGetladdrs(fd int) (unsafe.Pointer, int, error) {
	addr := struct {
		_     int32
		num   uint32
		addrs [4096]byte
	}{}
	l := unsafe.Sizeof(addr)
	if _, _, e := syscall.Syscall6(syscall.SYS_GETSOCKOPT,
		uintptr(fd),
		syscall.IPPROTO_SCTP,
		109, // SCTP_GET_LOCAL_ADDRS
		uintptr(unsafe.Pointer(&addr)),
		uintptr(unsafe.Pointer(&l)),
		0); e != 0 {
		return nil, 0, e
	}
	return unsafe.Pointer(&addr.addrs), int(addr.num), nil
}

func sctpFreeladdrs(addr unsafe.Pointer) {}

func sctpGetpaddrs(fd int) (unsafe.Pointer, int, error) {
	addr := struct {
		_     int32
		num   uint32
		addrs [4096]byte
	}{}
	l := unsafe.Sizeof(addr)
	if _, _, e := syscall.Syscall6(syscall.SYS_GETSOCKOPT,
		uintptr(fd),
		syscall.IPPROTO_SCTP,
		108, // SCTP_GET_PEER_ADDRS
		uintptr(unsafe.Pointer(&addr)),
		uintptr(unsafe.Pointer(&l)),
		0); e != 0 {
		return nil, 0, e
	}
	return unsafe.Pointer(&addr.addrs), int(addr.num), nil
}

func sctpFreepaddrs(addr unsafe.Pointer) {}
