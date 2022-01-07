package diameter

import (
	"errors"
	"math/rand"
	"net"
	"strings"
	"time"
)

// constant values
const (
	VendorID    uint32 = 41102
	ProductName        = "yatagarasu"
	FirmwareRev uint32 = 211031001
)

var (
	// TxTimeout is transport packet send timeout
	TxTimeout = time.Microsecond * 100
	// WDInterval is watchdog send interval time
	WDInterval = time.Second * 30
	// WDMaxSend is watchdog expired count
	WDMaxSend = 3

	wdTimer *time.Timer // system message timer
	wdCount = 0         // watchdog expired counter

	Router = true

	Local peer
	Peer  peer

	OverwriteAddr []net.IP // IP addresses of local host
	conn          net.Conn // Transport connection

	// Acceptable Application-ID and commands of the application.
	// Empty map indicate that accept any application.
	applications = make(map[uint32]application)

	hbhID = make(chan uint32, 1)
	eteID = make(chan uint32, 1)
	// sessionID = make(chan uint32, 1)

	notify   = make(chan stateEvent, 16)
	state    = closed
	sndStack = make(map[uint32]chan Message, 65535)
	rcvStack = make(chan Message, 65535)

	// Statistics value
	RxReq     uint64
	RejectReq uint64
	TxAnsFail uint64
	Tx1xxx    uint64
	Tx2xxx    uint64
	Tx3xxx    uint64
	Tx4xxx    uint64
	Tx5xxx    uint64
	TxEtc     uint64

	TxReq      uint64
	InvalidAns uint64
	TxReqFail  uint64
	Rx1xxx     uint64
	Rx2xxx     uint64
	Rx3xxx     uint64
	Rx4xxx     uint64
	Rx5xxx     uint64
	RxEtc      uint64
)

func init() {
	ut := time.Now().Unix()
	rand.Seed(ut)

	hbhID <- rand.Uint32()
	eteID <- (uint32(ut^0xFFF) << 20) | (rand.Uint32() ^ 0xFFFFF)
	// sessionID <- rand.Uint32()
	Local.state = uint32(ut)
}

type peer struct {
	Host  Identity
	Realm Identity
	state uint32
}

type application struct {
	venID    uint32
	handlers map[uint32]func(bool, []byte) (bool, []byte)
}

// RxQueue returns length of Rx queue
func RxQueue() int {
	return len(rcvStack)
}

// TxQueue returns length of Tx queue
func TxQueue() int {
	return len(sndStack)
}

func countRxCode(c uint32) {
	if c < 1000 {
		RxEtc++
	} else if c < 2000 {
		Rx1xxx++
	} else if c < 3000 {
		Rx2xxx++
	} else if c < 4000 {
		Rx3xxx++
	} else if c < 5000 {
		Rx4xxx++
	} else if c < 6000 {
		Rx5xxx++
	} else {
		RxEtc++
	}
}

func countTxCode(c uint32) {
	if c < 1000 {
		TxEtc++
	} else if c < 2000 {
		Tx1xxx++
	} else if c < 3000 {
		Tx2xxx++
	} else if c < 4000 {
		Tx3xxx++
	} else if c < 5000 {
		Tx4xxx++
	} else if c < 6000 {
		Tx5xxx++
	} else {
		TxEtc++
	}
}

func nextHbH() uint32 {
	ret := <-hbhID
	hbhID <- ret + 1
	return ret
}

func nextEtE() uint32 {
	ret := <-eteID
	eteID <- ret + 1
	return ret
}

/*
func nextSession() string {
	ret := <-sessionID
	sessionID <- ret + 1
	return fmt.Sprintf("%s;%d;%d;0",
		LocalHost, time.Now().Unix()+2208988800, ret)
}
*/

// ResolveIdentity get Diameter Host, Realm and address of the FQDN
func ResolveIdentiry(fqdn string) (host, realm Identity, err error) {
	h, _, err := net.SplitHostPort(fqdn)
	if err != nil {
		return
	}

	if host, err = ParseIdentity(h); err != nil {
		return
	}
	if i := strings.Index(h, "."); i < 0 {
		err = errors.New("domain part not found in local hostname")
		return
	} else if realm, err = ParseIdentity(h[i+1:]); err != nil {
		return
	}

	return
}

// LocalAddr returns transport connection of state machine
func LocalAddr() net.Addr {
	return conn.LocalAddr()
}

// PeerAddr returns transport connection of state machine
func PeerAddr() net.Addr {
	return conn.RemoteAddr()
}

// State returns state machine state
func State() string {
	return state.String()
}
