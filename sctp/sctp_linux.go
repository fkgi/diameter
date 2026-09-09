//go:build linux && !386

package sctp

import (
	"bytes"
	"encoding/binary"
	"io"
	"syscall"
	"time"
	"unsafe"
)

var epfd int
var cmap = make(chan map[int32]*SCTPConn, 1)

func init() {
	var err error
	if epfd, err = syscall.EpollCreate1(syscall.EPOLL_CLOEXEC); err != nil {
		panic(err)
	}
	// defer syscall.Close(epfd)

	cmap <- map[int32]*SCTPConn{}
	go func() {
		events := make([]syscall.EpollEvent, 128)

		for {
			n, e := syscall.EpollWait(epfd, events, -1)
			if e != nil && e == syscall.EINTR {
				continue
			} else if e != nil {
				panic(e)
			}

			for i := range n {
				ev := events[i]
				cs := <-cmap
				if c, ok := cs[ev.Fd]; ok {
					// fmt.Println(ev.Fd, ev.Events, len(c.rPoll), len(c.wPoll))

					if ev.Events&syscall.EPOLLIN != 0 && len(c.rPoll) == 0 {
						c.rPoll <- nil
					}
					if ev.Events&syscall.EPOLLOUT != 0 && len(c.wPoll) == 0 {
						c.wPoll <- nil
					}
					if ev.Events&(syscall.EPOLLHUP|
						syscall.EPOLLRDHUP|
						syscall.EPOLLERR) != 0 {
						// syscall.EpollCtl(epfd, syscall.EPOLL_CTL_DEL, int(ev.Fd), nil)
						delete(cs, ev.Fd)
						close(c.rPoll)
						close(c.wPoll)
					}
				}
				cmap <- cs
			}
		}
	}()
	/*
		go func() {
			t := time.Tick(time.Second)
			for range t {
				buf := new(bytes.Buffer)
				cs := <-cmap
				for k, v := range cs {
					fmt.Fprint(buf, k, ",", len(v.wPoll), ",", len(v.wPoll), ",")
				}
				cmap <- cs
				fmt.Println(buf.String())
			}

		}()
	*/
}

func eventIN(fd int) *syscall.EpollEvent {
	return &syscall.EpollEvent{
		Events: syscall.EPOLLIN |
			syscall.EPOLLRDHUP |
			0x80000000, //syscall.EPOLLET
		Fd: int32(fd)}
}
func eventINOUT(fd int) *syscall.EpollEvent {
	return &syscall.EpollEvent{
		Events: syscall.EPOLLOUT |
			syscall.EPOLLIN |
			syscall.EPOLLRDHUP |
			0x80000000, //syscall.EPOLLET
		Fd: int32(fd)}
}

func registerPoll(c *SCTPConn) (e error) {
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
		dataIo:          0,
		association:     1,
		address:         0,
		sendFailed:      0,
		peerError:       0,
		shutdown:        0,
		partialDelivery: 0,
		adaptationLayer: 0,
		authentication:  0,
		senderDry:       0}

	if _, _, ne := syscall.Syscall6(
		syscall.SYS_SETSOCKOPT,
		uintptr(c.sock),
		132, // SOL_SCTP
		11,  // SCTP_EVENTS
		uintptr(unsafe.Pointer(&event)),
		uintptr(unsafe.Sizeof(event)),
		0); ne != 0 {
	} else if e = syscall.EpollCtl(
		epfd,
		syscall.EPOLL_CTL_ADD,
		c.sock,
		eventINOUT(c.sock)); e != nil {
	} else {
		c.wPoll = make(chan any, 128)
		c.rPoll = make(chan any, 128)
		cs := <-cmap
		cs[int32(c.sock)] = c
		cmap <- cs
	}
	return
}

func sockOpenV4(seq bool) (int, error) {
	if seq {
		return syscall.Socket(
			syscall.AF_INET,
			syscall.SOCK_SEQPACKET|syscall.SOCK_CLOEXEC,
			syscall.IPPROTO_SCTP)
	} else {
		return syscall.Socket(
			syscall.AF_INET,
			syscall.SOCK_STREAM|syscall.SOCK_CLOEXEC,
			syscall.IPPROTO_SCTP)
	}
}

func sockOpenV6(seq bool) (int, error) {
	if seq {
		return syscall.Socket(
			syscall.AF_INET6,
			syscall.SOCK_SEQPACKET|syscall.SOCK_CLOEXEC,
			syscall.IPPROTO_SCTP)
	} else {
		return syscall.Socket(
			syscall.AF_INET6,
			syscall.SOCK_STREAM|syscall.SOCK_CLOEXEC,
			syscall.IPPROTO_SCTP)
	}
}

func sockListen(l *SCTPListener) (e error) {
	if e = syscall.Listen(l.sock, 128); e != nil {
	} else if e = syscall.SetNonblock(l.sock, true); e != nil {
	} else if e = syscall.EpollCtl(
		epfd,
		syscall.EPOLL_CTL_ADD,
		l.sock,
		eventIN(l.sock)); e != nil {
	} else {
		l.rPoll = make(chan any, 128)
		l.cPoll = make(chan any, 1)
		cs := <-cmap
		cs[int32(l.sock)] = &SCTPConn{
			sock: l.sock, rPoll: l.rPoll, wPoll: make(chan any, 1)}
		cmap <- cs
	}
	return
}

func sockAccept(l *SCTPListener) (nfd int, e error) {
	for {
		switch nfd, _, e = syscall.Accept(l.sock); e {
		case syscall.EAGAIN:
			select {
			case <-l.rPoll:
			case <-l.cPoll:
				cs := <-cmap
				delete(cs, int32(l.sock))
				cmap <- cs
			}
		case syscall.EINTR:
		case nil:
			if e = syscall.SetNonblock(nfd, true); e != nil {
				sockClose(nfd)
			}
			return
		default:
			return
		}
	}
}

func sockClose(fd int) error {
	syscall.Shutdown(fd, syscall.SHUT_RDWR)
	return syscall.Close(fd)
}

func sctpBindx(fd int, addr []byte) error {
	if _, _, e := syscall.Syscall6(
		syscall.SYS_SETSOCKOPT,
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
	t, _, e := syscall.Syscall6(
		syscall.SYS_SETSOCKOPT,
		uintptr(fd),
		syscall.IPPROTO_SCTP,
		110, // SCTP_SOCKOPT_CONNECTX
		uintptr(unsafe.Pointer(&addr[0])),
		uintptr(len(addr)),
		0)
	if e != 0 {
		return 0, e
	}

	peel := struct {
		aid  int32
		sd   int32
		flag int
	}{aid: int32(t), flag: syscall.SOCK_NONBLOCK}
	l := unsafe.Sizeof(peel)
	if _, _, e := syscall.Syscall6(
		syscall.SYS_GETSOCKOPT,
		uintptr(fd),
		syscall.IPPROTO_SCTP,
		122, // SCTP_SOCKOPT_PEELOFF_FLAGS
		uintptr(unsafe.Pointer(&peel)),
		uintptr(unsafe.Pointer(&l)),
		0); e != 0 {
		return 0, e
	}
	return int(peel.sd), nil
}

func sctpSend(c *SCTPConn, b []byte) (n int, e error) {
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

	for {
		n, e = syscall.SendmsgN(c.sock, b, buf.Bytes(), nil, syscall.MSG_EOR)
		if e == nil || e != syscall.EAGAIN {
			return
		}

		if c.wDeadline.IsZero() {
			<-c.wPoll
		} else {
			select {
			case <-c.wPoll:
			case <-time.After(time.Until(c.wDeadline)):
				e = busy{}
				return
			}
		}
	}
}

type busy struct{}

func (busy) Error() string {
	return syscall.EAGAIN.Error()
}
func (busy) Timeout() bool {
	return true
}
func (busy) Temporary() bool {
	return true
}

func sctpRecvmsg(c *SCTPConn, b []byte) (n int, e error) {
	var on, f int
	for {
		n, on, f, _, e = syscall.Recvmsg(
			c.sock, b, make([]byte, syscall.CmsgSpace(32)), 0)

		switch e {
		case nil:
			if f&0x8000 != 0 {
				if b[0] == 0x01 && b[1] == 0x80 &&
					(b[8] == 0x03 || b[8] == 0x01) {
					cs := <-cmap
					delete(cs, int32(c.sock))
					cmap <- cs
					e = io.EOF
					return
				}
			} else {
				if n <= 0 && on <= 0 {
					e = io.EOF
				}
				return
			}
		case syscall.EAGAIN:
			if c.rDeadline.IsZero() {
				<-c.rPoll
			} else {
				select {
				case <-c.rPoll:
				case <-time.After(time.Until(c.rDeadline)):
					e = busy{}
					return
				}
			}
		case syscall.EINTR:
		default:
			return
		}
	}
}

func sctpGetladdrs(fd int) (unsafe.Pointer, int, error) {
	addr := struct {
		_     int32
		num   uint32
		addrs [4096]byte
	}{}
	l := unsafe.Sizeof(addr)
	if _, _, e := syscall.Syscall6(
		syscall.SYS_GETSOCKOPT,
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
	if _, _, e := syscall.Syscall6(
		syscall.SYS_GETSOCKOPT,
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
