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
		if id == 0xffffffff {
			continue
		}
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

// EnableRelaySupport add supported application message
func EnableRelaySupport() {
	supportedApps[0xffffffff] = appSet{
		id:  0,
		req: map[uint32]Request{0: GenericReq{}},
		ans: map[uint32]Answer{0: GenericAns{}}}
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

func nextSession() string {
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
	AuthApps   map[uint32][]uint32
}

func (p *Peer) String() string {
	if p == nil {
		return "<nil>"
	}
	return string(p.Host)
}
