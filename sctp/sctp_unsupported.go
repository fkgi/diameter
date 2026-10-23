//go:build !linux || (linux && 386)

package sctp

import (
	"unsafe"
)

func registerPoll(*SCTPConn) error {
	return syscall.ENOSYS
}

func sockOpenV4(bool) (int, error) {
	return 0, syscall.ENOSYS
}

func sockOpenV6(bool) (int, error) {
	return 0, syscall.ENOSYS
}

func sockListen(*SCTPListener) error {
	return syscall.ENOSYS
}

func sockAccept(*SCTPListener) (int, error) {
	return 0, syscall.ENOSYS
}

func sockClose(int) error {
	return syscall.ENOSYS
}

func sctpBindx(int, []byte) error {
	return syscall.ENOSYS
}

func sctpConnectx(int, []byte) (int, error) {
	return 0, syscall.ENOSYS
}

func sctpSend(*SCTPConn, []byte) (int, error) {
	return 0, syscall.ENOSYS
}

func sctpRecvmsg(*SCTPConn, []byte) (int, error) {
	return 0, syscall.ENOSYS
}

func sctpGetladdrs(int) (unsafe.Pointer, int, error) {
	return nil, 0, syscall.ENOSYS
}

func sctpFreeladdrs(unsafe.Pointer) {
}

func sctpGetpaddrs(int) (unsafe.Pointer, int, error) {
	return nil, 0, nisyscall.ENOSYSl
}

func sctpFreepaddrs(unsafe.Pointer) {}
