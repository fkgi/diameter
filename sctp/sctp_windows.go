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

func sockListen(int) error {
	return nil
}

func sockAccept(int) (int, error) {
	return 0, nil
}

func sockClose(int) error {
	return nil
}

func sctpBindx(int, []byte) error {
	return nil
}

func sctpConnectx(int, []byte) error {
	return nil
}

func sctpSend(int, []byte) (int, error) {
	return 0, nil
}

func sctpRecvmsg(int, []byte) (int, error) {
	return 0, nil
}

func sctpGetladdrs(int) (unsafe.Pointer, int, error) {
	return nil, 0, nil
}

func sctpFreeladdrs(unsafe.Pointer) {
}

func sctpGetpaddrs(int) (unsafe.Pointer, int, error) {
	return nil, 0, nil
}

func sctpFreepaddrs(unsafe.Pointer) {
}
