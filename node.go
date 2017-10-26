package diameter

import (
	"fmt"
	"math/rand"
	"time"
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
	Host Identity
	// Realm name for local host
	Realm Identity
	// StateID for local host
	StateID uint32

	// Used for Vendor-Specific-Application-Id, Auth-Application-Id
	// and Supported-Vendor-Id AVP
	supportedApps = make(map[uint32]appSet)

	hbHID     = make(chan uint32, 1)
	etEID     = make(chan uint32, 1)
	sessionID = make(chan uint32, 1)
)

type appSet struct {
	id  uint32
	req map[uint32]Request
	ans map[uint32]Answer
}

func getSupportedApps() map[uint32][]uint32 {
	r := make(map[uint32][]uint32)
	for id, set := range supportedApps {
		if _, ok := r[set.id]; !ok {
			r[set.id] = make([]uint32, 0, 1)
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

	StateID = uint32(ut)

	/*
		vid := rfc6733.VendorID(0)
		aid := rfc6733.AuthApplicationID(0)
		AddSupportedMessage(vid, aid, 257, CER{}, CEA{})
		AddSupportedMessage(vid, aid, 282, DPR{}, DPA{})
		AddSupportedMessage(vid, aid, 280, DWR{}, DWA{})
	*/
}

// AddSupportedMessage add supported application message
func AddSupportedMessage(v, a, c uint32,
	req Request, ans Answer) {

	if _, ok := supportedApps[a]; !ok {
		supportedApps[a] = appSet{
			id:  v,
			req: make(map[uint32]Request),
			ans: make(map[uint32]Answer)}
	}
	supportedApps[a].req[c] = req
	supportedApps[a].ans[c] = ans
}

func nextHbH() uint32 {
	ret := <-hbHID
	hbHID <- ret + 1
	return ret
}

func nextEtE() uint32 {
	ret := <-etEID
	etEID <- ret + 1
	return ret
}

// NextSession returns new SessionID
func NextSession() string {
	ret := <-sessionID
	sessionID <- ret + 1
	return fmt.Sprintf("%s;%d;%d;0",
		Host, time.Now().Unix()+2208988800, ret)
}

// Peer is peer node of Diameter
type Peer struct {
	Realm, Host Identity

	WDInterval time.Duration
	WDExpired  int
	SndTimeout time.Duration

	AuthApps map[uint32][]uint32
}

func (p *Peer) String() string {
	if p == nil {
		return "<nil>"
	}
	return string(p.Host)
}
