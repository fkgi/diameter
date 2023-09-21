package diameter

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"
)

const (
	VendorID    uint32 = 41102         // VendorID of this code
	ProductName        = "round-robin" // ProductName of this code
	FirmwareRev uint32 = 230619001     // FirmwareRev of this code
)

var (
	WDInterval = time.Second * 30 // WDInterval is watchdog send interval time
	WDMaxSend  = 3                // WDMaxSend is watchdog expired count

	wdTimer *time.Timer // system message timer
	wdCount = 0         // watchdog expired counter

	Local peer // Local diameter host information
	Peer  peer // Peer diameter host information

	Router        = false     // Router mode add RouteRecord AVP in request message
	OverwriteAddr []net.IP    // Overwrite IP addresses of local host in CER
	TermSignals   []os.Signal // Signals for closing diameter connection

	// Acceptable Application-ID and commands of the application.
	// Empty map indicate that accept any application.
	applications = make(map[uint32]application)

	hbhID     = make(chan uint32, 1) // Hop-by-Hop ID source
	eteID     = make(chan uint32, 1) // End-to-End ID source
	sessionID = make(chan uint32, 1) // Session-ID source

	conn     net.Conn                               // Transport connection
	notify   = make(chan stateEvent, 16)            // state change notification queue
	state    = closed                               // current state
	sndQueue = make(map[uint32]chan Message, 65535) // Sending Request message queue
	rcvQueue = make(chan Message, 65535)            // Receiving Request message queue

	avpBufferSize = 10
)

func init() {
	ut := time.Now().Unix()
	// rand.Seed(ut)

	hbhID <- rand.Uint32()
	eteID <- (uint32(ut^0xFFF) << 20) | (rand.Uint32() ^ 0xFFFFF)
	sessionID <- rand.Uint32()
	Local.state = uint32(ut)
}

type peer struct {
	Host  Identity
	Realm Identity
	state uint32
}

type application struct {
	venID    uint32
	handlers map[uint32]Handler
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

// NextSession generate new session ID data
func NextSession(h string) string {
	ret := <-sessionID
	sessionID <- ret + 1
	if h == "" {
		h = Local.Host.String()
	}
	return fmt.Sprintf("%s;%d;%d;0", h, time.Now().Unix()+2208988800, ret)
}
