package sctp

import (
	"bytes"
	"encoding/binary"
	"net"
	"strconv"
	"syscall"
	"unsafe"
)

// SCTPAddr represents the address of a SCTP end point.
type SCTPAddr struct {
	IP   []net.IP
	Port int
}

func (a *SCTPAddr) rawBytes() []byte {
	var buf bytes.Buffer

	p := uint16(a.Port<<8) & 0xff00
	p |= uint16(a.Port>>8) & 0x00ff

	if len(a.IP) == 0 {
		raw := syscall.RawSockaddrInet4{
			Family: syscall.AF_INET,
			Port:   p,
		}
		copy(raw.Addr[:], net.IPv4zero)
		binary.Write(&buf, binary.LittleEndian, raw)
	} else if a.IP[0].To4() != nil {
		for _, i := range a.IP {
			i = i.To4()
			if i == nil {
				continue
			}
			raw := syscall.RawSockaddrInet4{
				Family: syscall.AF_INET,
				Port:   p}
			copy(raw.Addr[:], i)
			binary.Write(&buf, binary.LittleEndian, raw)
		}
	} else if a.IP[0].To16() != nil {
		for _, i := range a.IP {
			if i.To4() != nil {
				continue
			}
			i = i.To16()
			raw := syscall.RawSockaddrInet6{
				Family: syscall.AF_INET6,
				Port:   p,
				// Flowinfo:
				// Scope_id:
			}
			copy(raw.Addr[:], i)
			binary.Write(&buf, binary.BigEndian, raw)
		}
	} else {
		raw := syscall.RawSockaddrInet4{
			Family: syscall.AF_INET,
			Port:   p,
		}
		copy(raw.Addr[:], net.IPv4zero)
		binary.Write(&buf, binary.BigEndian, raw)
	}
	return buf.Bytes()
}

func resolveFromRawAddr(ptr unsafe.Pointer, n int) *SCTPAddr {
	addr := &SCTPAddr{}
	p := 0
	addr.IP = make([]net.IP, n)

	switch (*(*syscall.RawSockaddrAny)(ptr)).Addr.Family {
	case syscall.AF_INET:
		p = int((*(*syscall.RawSockaddrInet4)(ptr)).Port)

		for i := 0; i < n; i++ {
			a := *(*syscall.RawSockaddrInet4)(
				unsafe.Pointer(uintptr(ptr) + uintptr(16*i)))
			addr.IP[i] = net.IPv4(a.Addr[0], a.Addr[1], a.Addr[2], a.Addr[3])
		}
	case syscall.AF_INET6:
		p = int((*(*syscall.RawSockaddrInet6)(ptr)).Port)

		for i := 0; i < n; i++ {
			a := *(*syscall.RawSockaddrInet6)(
				unsafe.Pointer(uintptr(ptr) + uintptr(28*i)))
			addr.IP[i] = make([]byte, net.IPv6len)
			for j := 0; j < net.IPv6len; j++ {
				addr.IP[i][j] = a.Addr[j]
			}
		}
	default:
		panic("invalid family of address")
	}

	addr.Port = (p & 0x00ff) << 8
	addr.Port |= (p & 0xff00) >> 8
	return addr
}

func (a *SCTPAddr) String() string {
	var b bytes.Buffer

	for n, i := range a.IP {
		if a.IP[n].To4() != nil {
			b.WriteRune('/')
			b.WriteString(i.String())
		} else if a.IP[n].To16() != nil {
			b.WriteRune('/')
			b.WriteRune('[')
			b.WriteString(i.String())
			b.WriteRune(']')
		}
	}
	b.WriteRune(':')
	b.WriteString(strconv.Itoa(a.Port))

	return b.String()[1:]
}

func (a *SCTPAddr) Network() string { return "sctp" }
