package provider

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// Notify is called when error or trace event are occured
var Notify = func(e error) {
	log.Println(e)
	d, ok := e.(Dump)
	if ok {
		d.f(os.Stderr)
	}
}

// Dump provide Diameter message dump
type Dump struct {
	f func(io.Writer)
}

func (e *Dump) Error() string {
	if e == nil {
		return "<nil>"
	}
	return "Diameter message dump"
}

// WriteSuccess provide Diameter message dump
type WriteSuccess struct {
	Local string
	Peer  string
}

func (e *WriteSuccess) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf(
		"write Diameter message from %s to %s",
		e.Local, e.Peer)
}

// WriteFail provide Diameter message dump
type WriteFail struct {
	Local string
	Peer  string
	Err   error
}

func (e *WriteFail) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf(
		"write Diameter message failed from %s to %s, reason is %s",
		e.Local, e.Peer, e.Err)
}

// ReadSuccess provide Diameter message dump
type ReadSuccess struct {
	Local string
	Peer  string
}

func (e *ReadSuccess) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf(
		"read Diameter message from %s to %s",
		e.Local, e.Peer)
}

// ReadFail provide Diameter message dump
type ReadFail struct {
	Local string
	Peer  string
	Err   error
}

func (e *ReadFail) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf(
		"read Diameter message failed from %s to %s, reason is %s",
		e.Local, e.Peer, e.Err)
}

// TransportCloseSuccess provide Diameter message dump
type TransportCloseSuccess struct {
	Local string
	Peer  string
}

func (e *TransportCloseSuccess) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf(
		"transport link from %s to %s closed",
		e.Local, e.Peer)
}

// TransportCloseFail provide Diameter message dump
type TransportCloseFail struct {
	Local string
	Peer  string
	Err   error
}

func (e *TransportCloseFail) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf(
		"fail to close transport link from %s to %s, reason is %s",
		e.Local, e.Peer, e.Err)
}

// TransportConnectSuccess provide Diameter message dump
type TransportConnectSuccess struct {
	Local string
	Peer  string
	LAddr net.Addr
	PAddr net.Addr
}

func (e *TransportConnectSuccess) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf(
		"transport link from %s(%s) to %s(%s) connected",
		e.Local, e.LAddr, e.Peer, e.PAddr)
}

// TransportConnectFail provide Diameter message dump
type TransportConnectFail struct {
	Err error
}

func (e *TransportConnectFail) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf(
		"fail to connect transport link, reason is %s", e.Err)
}
