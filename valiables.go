package diameter

import (
	"fmt"
	"math/rand"
	"net"
	"time"
)

const (
	VendorID    uint32 = 41102         // VendorID of this code
	ProductName        = "round-robin" // ProductName of this code
	FirmwareRev uint32 = 230619001     // FirmwareRev of this code

	avpBufferSize = 10
)

var (
	WDInterval = time.Second * 30 // WDInterval is watchdog send interval time
	WDMaxSend  = 3                // WDMaxSend is watchdog expired count

	Host    Identity // Local diameter hostname
	Realm   Identity // Local diameter realm
	stateID uint32   // Local diameter state ID

	OverwriteAddr []net.IP // Overwrite IP addresses of local host in CER

	// Acceptable Application-ID and commands of the application.
	// Empty map indicate that accept any application.
	applications = make(map[uint32]application)

	hbhID     = make(chan uint32, 1) // Hop-by-Hop ID source
	eteID     = make(chan uint32, 1) // End-to-End ID source
	sessionID = make(chan uint32, 1) // Session-ID source

	sndQueue = make(map[uint32]chan Message, 65535) // Sending Request message queue
	rcvQueue = make(chan Message, 65535)            // Receiving Request message queue
)

func init() {
	ut := time.Now().Unix()
	hbhID <- rand.Uint32()
	eteID <- (uint32(ut^0xFFF) << 20) | (rand.Uint32() ^ 0xFFFFF)
	sessionID <- rand.Uint32()
	stateID = uint32(ut)
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
		h = Host.String()
	}
	return fmt.Sprintf("%s;%d;%d;0", h, time.Now().Unix()+2208988800, ret)
}
