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



	n, e := C.setsockopt(
			C.int(fd),
			C.SOL_SCTP,
			C.int(opt),
			p,
			C.socklen_t(l))
		if int(n) < 0 {
			return e
		}
		return nil
*/

func sockOpenV4() (int, error) {
	return syscall.Socket(
		syscall.AF_INET,
		syscall.SOCK_STREAM, //syscall.SOCK_SEQPACKET,
		syscall.IPPROTO_SCTP)
}

func sockOpenV6() (int, error) {
	return syscall.Socket(
		syscall.AF_INET6,
		syscall.SOCK_STREAM, //syscall.SOCK_SEQPACKET,
		syscall.IPPROTO_SCTP)
}

func sockListen(fd int) error {
	return syscall.Listen(fd, 1024)
}

func sockAccept(fd int) (nfd int, e error) {
	nfd, _, e = syscall.Accept4(fd, 0)
	return
}

/*
	Stream  uint16
	SSN     uint16
	Flags   uint16
	_       uint16
	PPID    uint32
	Context uint32
	TTL     uint32
	TSN     uint32
	CumTSN  uint32
	AssocID int32
*/

func sockClose(fd int) error {
	/*
		var buf bytes.Buffer
		hdr := &syscall.Cmsghdr{
			Level: syscall.IPPROTO_SCTP,
			Type:  1, //SCTP_CMSG_SNDRCV
		}

		hdr.SetLen(syscall.CmsgSpace(32))
		binary.Write(&buf, binary.LittleEndian, hdr)
		binary.Write(&buf, binary.LittleEndian, [4]byte{})
		binary.Write(&buf, binary.BigEndian, uint16(0x0100)) // flag
		binary.Write(&buf, binary.LittleEndian, [26]byte{})

		info := struct {
			stream     uint16 // Stream No. in SCTP association
			ssn        uint16 // Stream Sequence No. of DATA (recvmsg only), SSN is same number in fragmented data
			flags      uint16 // Bit flags, UNORDERD/ADDR_OVER/ABORT/EOF/SENDALL
			ppid       uint32 // Protocol ID of user data
			context    uint32 // Context for match error response and request
			timetolive uint32 // Available time in msec, 0 for no time-out
			tsn        uint32 // Transmission Sequence No. of DATA (recvmsg only)
			cumtsn     uint32 // Current cumulativa TSN of DATA (recvmsg only)
			assocID    uint32 // Association ID os SCTP association
		}{
			flags: 0x0100,
		}
		hdr.SetLen(syscall.CmsgSpace(30))
		binary.Write(&buf, binary.LittleEndian, hdr)
		binary.Write(&buf, binary.LittleEndian, info)
		syscall.SendmsgN(fd, nil, buf.Bytes(), nil, 0)
	*/

	syscall.Shutdown(fd, syscall.SHUT_RDWR)
	return syscall.Close(fd)
}

func sctpBindx(fd int, addr []byte) error {
	_, _, e := syscall.Syscall6(syscall.SYS_SETSOCKOPT,
		uintptr(fd),
		132, // SOL_SCTP
		100, // SCTP_SOCKOPT_BINDX_ADD
		uintptr(unsafe.Pointer(&addr[0])),
		uintptr(len(addr)),
		0)
	if e != 0 {
		return e
	}
	return nil
}

func sctpConnectx(fd int, addr []byte) error {
	_, _, e := syscall.Syscall6(syscall.SYS_SETSOCKOPT,
		uintptr(fd),
		132, // SOL_SCTP
		110, // SCTP_SOCKOPT_CONNECTX
		uintptr(unsafe.Pointer(&addr[0])),
		uintptr(len(addr)),
		0)
	if e != 0 {
		return e
	}
	return nil
}

func sctpSend(fd int, b []byte) (int, error) {
	buf := new(bytes.Buffer)
	hdr := &syscall.Cmsghdr{
		Level: syscall.IPPROTO_SCTP,
		Type:  2, //SCTP_SNDINFO
	}
	hdr.SetLen(syscall.CmsgSpace(16))
	/*
		__u16 snd_sid;
		__u16 snd_flags;
		__u32 snd_ppid;
		__u32 snd_context;
		sctp_assoc_t(u32) snd_assoc_id;
	*/
	binary.Write(buf, binary.LittleEndian, hdr)
	binary.Write(buf, binary.LittleEndian, uint16(0))      // session ID=0
	binary.Write(buf, binary.LittleEndian, uint16(0x0001)) // flag=SCTP_UNORDERED
	binary.Write(buf, binary.BigEndian, ProtocolID)        // PPID=diameter
	buf.Write(make([]byte, 8))

	return syscall.SendmsgN(fd, b, buf.Bytes(), nil, 0)
}

func sctpRecvmsg(fd int, b []byte) (int, error) {
	n, on, _, _, e := syscall.Recvmsg(fd, b, make([]byte, 256), 0)
	if e == nil && n == 0 && on == 0 {
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
	_, _, e := syscall.Syscall6(syscall.SYS_GETSOCKOPT,
		uintptr(fd),
		132, // SOL_SCTP
		109, // SCTP_GET_LOCAL_ADDRS
		uintptr(unsafe.Pointer(&addr)),
		uintptr(unsafe.Pointer(&l)),
		0)
	if e != 0 {
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
	_, _, e := syscall.Syscall6(syscall.SYS_GETSOCKOPT,
		uintptr(fd),
		132, // SOL_SCTP
		108, // SCTP_GET_PEER_ADDRS
		uintptr(unsafe.Pointer(&addr)),
		uintptr(unsafe.Pointer(&l)),
		0)
	if e != 0 {
		return nil, 0, e
	}
	return unsafe.Pointer(&addr.addrs), int(addr.num), nil
}

func sctpFreepaddrs(addr unsafe.Pointer) {}
