package sock

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fkgi/diameter/msg"
)

var (
	// TransportTimeout is transport packet send timeout
	TransportTimeout = time.Second
	// WDInterval is watchdog send interval time
	WDInterval = time.Second * time.Duration(30)
	// WDExpired is watchdog expired count
	WDExpired = 3
	// SndTimeout is message send timeout time
	SndTimeout = time.Second * time.Duration(30)
)

func init() {
	rand.Seed(time.Now().Unix())
}

// Local is local node of Diameter
type Local struct {
	Realm msg.DiameterIdentity
	Host  msg.DiameterIdentity
	//Addr  net.Addr

	StateID msg.OriginStateID

	hbHID     chan uint32
	etEID     chan uint32
	sessionID chan uint32

	AuthApps map[msg.VendorID][]msg.ApplicationID
	// for Vendor-Specific-Application-Id,
	//     Auth-Application-Id, Supported-Vendor-Id AVP
}

func (l *Local) String() string {
	if l == nil {
		return "<nil>"
	}
	return string(l.Host)
}

// NextHbH make HbH ID
func (l *Local) NextHbH() uint32 {
	if l.hbHID == nil {
		l.hbHID = make(chan uint32, 1)
		l.hbHID <- rand.Uint32()
	}
	ret := <-l.hbHID
	l.hbHID <- ret + 1
	return ret
}

// NextEtE make EtE ID
func (l *Local) NextEtE() uint32 {
	if l.etEID == nil {
		tmp := uint32(time.Now().Unix() ^ 0xFFF)
		tmp = (tmp << 20) | (rand.Uint32() ^ 0x000FFFFF)
		l.etEID = make(chan uint32, 1)
		l.etEID <- tmp
	}
	ret := <-l.etEID
	l.etEID <- ret + 1
	return ret
}

// NextSession make session ID
func (l *Local) NextSession() msg.SessionID {
	if l.sessionID == nil {
		l.sessionID = make(chan uint32, 1)
		l.sessionID <- rand.Uint32()
	}
	ret := <-l.sessionID
	l.sessionID <- ret + 1
	return msg.SessionID(fmt.Sprintf("%s;%d;%d;0",
		l.Host, time.Now().Unix()+2208988800, ret))
}

// Peer is peer node of Diameter
type Peer struct {
	Realm msg.DiameterIdentity
	Host  msg.DiameterIdentity
	//Addr  net.Addr

	WDInterval time.Duration
	WDExpired  int
	SndTimeout time.Duration

	Handler func(msg.Message) msg.Message

	AuthApps map[msg.VendorID][]msg.ApplicationID
}

func (p *Peer) String() string {
	if p == nil {
		return "<nil>"
	}
	return string(p.Host)
}
