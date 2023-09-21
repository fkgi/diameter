package diameter

import (
	"bytes"
	"errors"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/fkgi/diameter/sctp"
)

func termWithSignals(isTx bool) {
	if len(TermSignals) == 0 {
		return
	}

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, TermSignals...)
	<-sigc

	if isTx {
		Close(DoNotWantToTalkToYou)
	} else {
		Close(Rebooting)
	}
}

// DialAndServe start diameter connection handling process as initiator.
// Inputs are string of local hostname[:port][/realm] (la),
// peer hostname[:port][/realm] (ra) and bool flag to use SCTP.
func DialAndServe(la, pa string, isSctp bool) (err error) {
	if isSctp {
		tla := &sctp.SCTPAddr{}
		Local.Host, Local.Realm, tla.IP, tla.Port, err = resolveIdentiry(la)
		if err != nil {
			return
		}
		tpa := &sctp.SCTPAddr{}
		Peer.Host, Peer.Realm, tpa.IP, tpa.Port, err = resolveIdentiry(pa)
		if err != nil {
			return
		}
		conn, err = sctp.DialSCTP(tla, tpa)
	} else {
		var ips []net.IP
		tla := &net.TCPAddr{}
		Local.Host, Local.Realm, ips, tla.Port, err = resolveIdentiry(la)
		if err != nil {
			return
		}
		tla.IP = ips[0]
		tpa := &net.TCPAddr{}
		Peer.Host, Peer.Realm, ips, tpa.Port, err = resolveIdentiry(pa)
		if err != nil {
			return
		}
		tpa.IP = ips[0]
		conn, err = net.DialTCP("tcp", tla, tpa)
	}
	if err != nil {
		return
	}

	state = closed
	go termWithSignals(true)
	return serve()
}

// ListenAndServe start diameter connection handling process as responder.
// Inputs are string of local hostname (la) and bool flag to use SCTP.
// If Peer is nil, any peer is accepted.
func ListenAndServe(la string, isSctp bool) (err error) {
	var l net.Listener
	if isSctp {
		tla := &sctp.SCTPAddr{}
		Local.Host, Local.Realm, tla.IP, tla.Port, err = resolveIdentiry(la)
		if err != nil {
			return
		}
		l, err = sctp.ListenSCTP(tla)
	} else {
		var ips []net.IP
		tla := &net.TCPAddr{}
		Local.Host, Local.Realm, ips, tla.Port, err = resolveIdentiry(la)
		if err != nil {
			return
		}
		tla.IP = ips[0]
		l, err = net.ListenTCP("tcp", tla)
	}
	if err != nil {
		return
	}

	t := time.AfterFunc(WDInterval, func() {
		l.Close()
	})

	conn, err = l.Accept()
	t.Stop()
	if err != nil {
		conn.Close()
		l.Close()
		return
	}

	if Peer.Host == "" {
		names, err := net.LookupAddr(conn.RemoteAddr().String())
		if err == nil {
			Peer.Host, Peer.Realm, _, _, _ = resolveIdentiry(names[0])
		}
	}

	state = waitCER
	go termWithSignals(false)
	return serve()
}

func resolveIdentiry(fqdn string) (host, realm Identity, ip []net.IP, port int, err error) {
	f := strings.Split(fqdn, "/")
	h, p, e := net.SplitHostPort(f[0])
	if e != nil {
		err = e
		return
	}
	if p == "" {
		p = "3868"
	}

	if host, err = ParseIdentity(h); err != nil {
		return
	}
	if len(f) > 1 {
		if realm, err = ParseIdentity(f[1]); err != nil {
			return
		}
	} else if i := strings.Index(h, "."); i < 0 {
		err = errors.New("domain part not found in local hostname")
		return
	} else if realm, err = ParseIdentity(h[i+1:]); err != nil {
		return
	}

	a, e := net.LookupHost(h)
	if e != nil {
		err = e
		return
	}
	ip = make([]net.IP, 0)
	for _, s := range a {
		ip = append(ip, net.ParseIP(s))
	}

	port, err = strconv.Atoi(p)
	return
}

func serve() error {
	go func() {
		for {
			m := Message{}

			// conn.SetReadDeadline(time.Time{})
			if err := m.UnmarshalFrom(conn); err != nil {
				notify <- eventPeerDisc{reason: err}
				break
			}

			if m.AppID == 0 && m.Code == 257 && m.FlgR {
				notify <- eventRcvCER{m}
			} else if m.AppID == 0 && m.Code == 257 && !m.FlgR {
				notify <- eventRcvCEA{m}
			} else if m.AppID == 0 && m.Code == 280 && m.FlgR {
				notify <- eventRcvDWR{m}
			} else if m.AppID == 0 && m.Code == 280 && !m.FlgR {
				notify <- eventRcvDWA{m}
			} else if m.AppID == 0 && m.Code == 282 && m.FlgR {
				notify <- eventRcvDPR{m}
			} else if m.AppID == 0 && m.Code == 282 && !m.FlgR {
				notify <- eventRcvDPA{m}
			} else if m.FlgR {
				notify <- eventRcvReq{m}
			} else {
				notify <- eventRcvAns{m}
			}
		}
	}()
	go func() {
	rxNewMsg:
		for req, ok := <-rcvQueue; ok; req, ok = <-rcvQueue {
			if app, ok := applications[req.AppID]; ok {
				if f, ok := app.handlers[req.Code]; ok && f != nil {
					avp := make([]AVP, 0, avpBufferSize)
					for rdr := bytes.NewReader(req.AVPs); rdr.Len() != 0; {
						a := AVP{}
						if e := a.UnmarshalFrom(rdr); e != nil {
							buf := new(bytes.Buffer)
							SetResultCode(InvalidAvpValue).MarshalTo(buf)
							SetOriginHost(Local.Host).MarshalTo(buf)
							SetOriginRealm(Local.Realm).MarshalTo(buf)
							notify <- eventSndMsg{Message{
								FlgR: false, FlgP: req.FlgP, FlgE: req.FlgE, FlgT: false,
								Code: req.Code, AppID: req.AppID,
								HbHID: req.HbHID, EtEID: req.EtEID,
								AVPs: buf.Bytes()}}
							continue rxNewMsg
						}
						avp = append(avp, a)
					}

					req.FlgE, avp = f(req.FlgT, avp)
					buf := new(bytes.Buffer)
					for _, a := range avp {
						a.MarshalTo(buf)
					}
					notify <- eventSndMsg{Message{
						FlgR: false, FlgP: req.FlgP, FlgE: req.FlgE, FlgT: false,
						Code: req.Code, AppID: req.AppID,
						HbHID: req.HbHID, EtEID: req.EtEID,
						AVPs: buf.Bytes()}}
					continue rxNewMsg
				}
			}
			notify <- eventSndMsg{DefaultRxHandler(req)}
		}
	}()

	TraceEvent(shutdown.String(), state.String(), eventInit{}.String(), nil)

	if state != waitCER {
		notify <- eventConnect{}
	}
	for {
		event := <-notify
		old := state
		err := event.exec()
		TraceEvent(old.String(), state.String(), event.String(), err)

		if _, ok := event.(eventPeerDisc); ok {
			break
		}
	}

	return nil
}

// Close Diameter connection and stop state machine.
func Close(cause Enumerated) {
	if state == open {
		notify <- eventLock{}
		for len(rcvQueue) != 0 || len(sndQueue) != 0 {
			time.Sleep(time.Millisecond * 100)
		}
	}
	notify <- eventStop{cause}
}
