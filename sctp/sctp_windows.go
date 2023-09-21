package sctp

import (
	"unsafe"
)

func sockOpenV4() (int, error) {
	return 0, nil
}

func sockOpenV6() (int, error) {
	return 0, nil
}

func sockListen(fd int) error {
	return nil
}

func sockAccept(fd int) (int, error) {
	return 0, nil
}

func sockClose(fd int) error {
	return nil
}

func sctpBindx(fd int, addr []byte) error {
	return nil
}

func sctpConnectx(fd int, addr []byte) error {
	return nil
}

func sctpSend(fd int, b []byte) (int, error) {
	return 0, nil
}

func sctpRecvmsg(fd int, b []byte) (int, error) {
	return 0, nil
}

func sctpGetladdrs(fd int) (unsafe.Pointer, int, error) {
	return nil, 0, nil
}

func sctpFreeladdrs(addr unsafe.Pointer) {
}

func sctpGetpaddrs(fd int) (unsafe.Pointer, int, error) {
	return nil, 0, nil
}

func sctpFreepaddrs(addr unsafe.Pointer) {
}
