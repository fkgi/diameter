package sock

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fkgi/diameter/msg"
	"github.com/fkgi/diameter/rfc6733"
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
	StateID rfc6733.OriginStateID

	// Used for Vendor-Specific-Application-Id, Auth-Application-Id
	// and Supported-Vendor-Id AVP
	supportedApps = make(map[rfc6733.AuthApplicationID]appSet)

	hbHID     = make(chan uint32, 1)
	etEID     = make(chan uint32, 1)
	sessionID = make(chan uint32, 1)
)

type appSet struct {
	id  rfc6733.VendorID
	req map[uint32]msg.Request
	ans map[uint32]msg.Answer
}

func getSupportedApps() map[rfc6733.VendorID][]rfc6733.AuthApplicationID {
	r := make(map[rfc6733.VendorID][]rfc6733.AuthApplicationID)
	for id, set := range supportedApps {
		if _, ok := r[set.id]; !ok {
			r[set.id] = make([]rfc6733.AuthApplicationID, 0, 1)
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

	// StateID = rfc6733.OriginStateID(ut)

	/*
		vid := rfc6733.VendorID(0)
		aid := rfc6733.AuthApplicationID(0)
		AddSupportedMessage(vid, aid, 257, msg.CER{}, msg.CEA{})
		AddSupportedMessage(vid, aid, 282, msg.DPR{}, msg.DPA{})
		AddSupportedMessage(vid, aid, 280, msg.DWR{}, msg.DWA{})
	*/
}

// AddSupportedMessage add supported application message
func AddSupportedMessage(
	v rfc6733.VendorID, a rfc6733.AuthApplicationID, c uint32,
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

func nextSession() rfc6733.SessionID {
	ret := <-sessionID
	sessionID <- ret + 1
	return rfc6733.SessionID(fmt.Sprintf("%s;%d;%d;0",
		Host, time.Now().Unix()+2208988800, ret))
}

// Peer is peer node of Diameter
type Peer struct {
	Realm, Host msg.DiameterIdentity

	WDInterval time.Duration
	WDExpired  int
	SndTimeout time.Duration

	Handler func(msg.Request) msg.Answer

	AuthApps map[rfc6733.VendorID][]rfc6733.AuthApplicationID
}

func (p *Peer) String() string {
	if p == nil {
		return "<nil>"
	}
	return string(p.Host)
}
