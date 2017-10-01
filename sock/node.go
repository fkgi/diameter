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

	// Host name for local host
	Host msg.DiameterIdentity
	// Realm name for local host
	Realm msg.DiameterIdentity
	// StateID for local host
	StateID msg.OriginStateID

	// Used for Vendor-Specific-Application-Id, Auth-Application-Id
	// and Supported-Vendor-Id AVP
	supportedApps = make(map[msg.AuthApplicationID]appSet)

	hbHID     = make(chan uint32, 1)
	etEID     = make(chan uint32, 1)
	sessionID = make(chan uint32, 1)
)

type appSet struct {
	id  msg.VendorID
	req map[uint32]msg.Request
	ans map[uint32]msg.Answer
}

func getSupportedApps() map[msg.VendorID][]msg.AuthApplicationID {
	r := make(map[msg.VendorID][]msg.AuthApplicationID)
	for id, set := range supportedApps {
		if _, ok := r[set.id]; !ok {
			r[set.id] = make([]msg.AuthApplicationID, 0, 1)
		}
		r[set.id] = append(r[set.id], id)
	}
	return r
}

func init() {
	ut := time.Now().Unix()
	rand.Seed(ut)

	hbHID <- rand.Uint32()

	tmp := uint32(ut ^ 0xFFF)
	tmp = (tmp << 20) | (rand.Uint32() ^ 0x000FFFFF)
	etEID <- tmp

	sessionID <- rand.Uint32()

	// StateID = msg.OriginStateID(ut)

	/*
		vid := msg.VendorID(0)
		aid := msg.AuthApplicationID(0)
		AddSupportedMessage(vid, aid, 257, msg.CER{}, msg.CEA{})
		AddSupportedMessage(vid, aid, 282, msg.DPR{}, msg.DPA{})
		AddSupportedMessage(vid, aid, 280, msg.DWR{}, msg.DWA{})
	*/
}

// AddSupportedMessage add supported application message
func AddSupportedMessage(
	v msg.VendorID, a msg.AuthApplicationID, c uint32,
	req msg.Request, ans msg.Answer) {

	if _, ok := supportedApps[a]; !ok {
		supportedApps[a] = appSet{
			id:  v,
			req: make(map[uint32]msg.Request),
			ans: make(map[uint32]msg.Answer)}
	}
	supportedApps[a].req[c] = req
	supportedApps[a].ans[c] = ans
}

// NextHbH make HbH ID
func NextHbH() uint32 {
	ret := <-hbHID
	hbHID <- ret + 1
	return ret
}

// NextEtE make EtE ID
func NextEtE() uint32 {
	ret := <-etEID
	etEID <- ret + 1
	return ret
}

// NextSession make session ID
func NextSession() msg.SessionID {
	ret := <-sessionID
	sessionID <- ret + 1
	return msg.SessionID(fmt.Sprintf("%s;%d;%d;0",
		Host, time.Now().Unix()+2208988800, ret))
}

// Peer is peer node of Diameter
type Peer struct {
	Realm, Host msg.DiameterIdentity

	WDInterval time.Duration
	WDExpired  int
	SndTimeout time.Duration

	Handler func(msg.Request) msg.Answer

	AuthApps map[msg.VendorID][]msg.AuthApplicationID
}

func (p *Peer) String() string {
	if p == nil {
		return "<nil>"
	}
	return string(p.Host)
}
