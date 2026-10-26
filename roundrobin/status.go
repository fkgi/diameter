package main

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/fkgi/diameter"
)

var txMetrics = make(chan map[string]uint64, 1)
var rxMetrics = make(chan map[string]uint64, 1)

func init() {
	txMetrics <- map[string]uint64{}
	rxMetrics <- map[string]uint64{}
}

func metricsHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")

	sysmetricsHandler(w)

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	fmt.Fprintln(w, "# HELP heap_alloc_bytes Number of heap bytes allocated and still in use")
	fmt.Fprintln(w, "# TYPE heap_alloc_bytes gauge")
	fmt.Fprintf(w, "heap_alloc_bytes %d\n", ms.Alloc)

	fmt.Fprintln(w, "# HELP roundrobin_active_worker_count Count of active worker")
	fmt.Fprintln(w, "# TYPE roundrobin_active_worker gauge")
	fmt.Fprintf(w, "roundrobin_active_worker %d\n", diameter.ActiveSharedWorkers())

	if len(reference) != 0 {
		fmt.Fprintln(w, "# HELP roundrobin_diameter_connection_state Status of Diameter connection(down=0, up=1)")
		fmt.Fprintln(w, "# TYPE roundrobin_diameter_connection_state gauge")
	}
	for _, a := range reference {
		switch a.State() {
		case "open", "locked":
			fmt.Fprintf(w, "roundrobin_diameter_connection_state{peer=\"%s\"} 1\n", a.Host)
		default:
			fmt.Fprintf(w, "roundrobin_diameter_connection_state{peer=\"%s\"} 0\n", a.Host)
		}
	}

	m := <-txMetrics
	if len(m) != 0 {
		fmt.Fprintln(w, "# HELP roundrobin_tx_diam_message Total number of outgoing Diameter message")
		fmt.Fprintln(w, "# TYPE roundrobin_tx_diam_message counter")
	}
	for k, v := range m {
		fmt.Fprintf(w, "roundrobin_tx_diam_message{type=\"%s\"} %d\n", k, v)
	}
	txMetrics <- m

	m = <-rxMetrics
	if len(m) != 0 {
		fmt.Fprintln(w, "# HELP roundrobin_rx_diam_message Total number of incoming Diameter message")
		fmt.Fprintln(w, "# TYPE roundrobin_rx_diam_message counter")
	}
	for k, v := range m {
		fmt.Fprintf(w, "roundrobin_rx_diam_message{type=\"%s\"} %d\n", k, v)
	}
	rxMetrics <- m
}

/*
type statFmt struct {
	L peerFmt   `json:"local"`
	P []peerFmt `json:"peer"`
}
type peerFmt struct {
	S string `json:"state,omitempty"`
	H string `json:"host"`
	R string `json:"realm"`
	A string `json:"address"`
}

func conStateHandler(w http.ResponseWriter, r *http.Request) {
	con := reference
	if len(con) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	st := statFmt{
		L: peerFmt{
			H: diameter.Host.String(),
			R: diameter.Realm.String(),
			A: con[0].LocalAddr().String()},
		P: []peerFmt{}}
	for _, c := range con {
		st.P = append(st.P, peerFmt{
			S: c.State(),
			H: c.Host.String(),
			R: c.Realm.String(),
			A: c.PeerAddr().String()})
	}
	if jd, e := json.Marshal(st); e != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jd)
	}
}
*/
/*
type statistics struct {
	RxReq     uint64 `json:"rx_request"`
	TxDisc    uint64 `json:"tx_discard"`
	TxAnsEtc  uint64 `json:"tx_etc"`
	TxAns1xxx uint64 `json:"tx_1xxx"`
	TxAns2xxx uint64 `json:"tx_2xxx"`
	TxAns3xxx uint64 `json:"tx_3xxx"`
	TxAns4xxx uint64 `json:"tx_4xxx"`
	TxAns5xxx uint64 `json:"tx_5xxx"`
	TxReq     uint64 `json:"tx_request"`
	RxIvld    uint64 `json:"rx_invalid"`
	RxAnsEtc  uint64 `json:"rx_etc"`
	RxAns1xxx uint64 `json:"rx_1xxx"`
	RxAns2xxx uint64 `json:"rx_2xxx"`
	RxAns3xxx uint64 `json:"rx_3xxx"`
	RxAns4xxx uint64 `json:"rx_4xxx"`
	RxAns5xxx uint64 `json:"rx_5xxx"`
}

var stats = make(chan statistics, 1)

func init() {
	stats <- statistics{}
}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	s := <-stats
	stats <- s

	if jd, e := json.Marshal(s); e != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jd)
	}
}
*/

func count(msg diameter.Message, rx bool, err error) {
	if rx {
		if msg.FlgR {
			s := <-rxMetrics
			s["request"] = s["request"] + 1
			rxMetrics <- s

			if _, ok := err.(diameter.RejectRxMessage); ok {
				s := <-txMetrics
				s["discard"] = s["discard"] + 1
				txMetrics <- s
			}
		} else if _, ok := err.(diameter.FailureAnswer); err != nil && !ok {
			s := <-txMetrics
			s["invalid"] = s["invalid"] + 1
			txMetrics <- s
		} else {
			var code uint32
			if avps, e := msg.GetAVP(); e == nil {
				for _, a := range avps {
					switch a.Code {
					case 268:
						code, _ = diameter.GetResultCode(a)
					case 297:
						code, _ = diameter.GetResultCode(a)
						code %= 10000
					}
					if code != 0 {
						break
					}
				}
			}

			s := <-rxMetrics
			if code < 1000 {
				s["etc"] = s["etc"] + 1
			} else if code < 2000 {
				s["1000"] = s["1000"] + 1
			} else if code < 3000 {
				s["2000"] = s["2000"] + 1
			} else if code < 4000 {
				s["3000"] = s["3000"] + 1
			} else if code < 5000 {
				s["4000"] = s["4000"] + 1
			} else if code < 6000 {
				s["5000"] = s["5000"] + 1
			} else {
				s["etc"] = s["etc"] + 1
			}
			rxMetrics <- s
		}
	} else {
		if msg.FlgR {
			s := <-txMetrics
			s["request"] = s["request"] + 1
			txMetrics <- s
		} else {
			var code uint32
			if avps, e := msg.GetAVP(); e == nil {
				for _, a := range avps {
					switch a.Code {
					case 268:
						code, _ = diameter.GetResultCode(a)
					case 297:
						code, _ = diameter.GetResultCode(a)
						code %= 10000
					}
					if code != 0 {
						break
					}
				}
			}

			s := <-txMetrics
			if code < 1000 {
				s["etc"] = s["etc"] + 1
			} else if code < 2000 {
				s["1000"] = s["1000"] + 1
			} else if code < 3000 {
				s["2000"] = s["2000"] + 1
			} else if code < 4000 {
				s["3000"] = s["3000"] + 1
			} else if code < 5000 {
				s["4000"] = s["4000"] + 1
			} else if code < 6000 {
				s["5000"] = s["5000"] + 1
			} else {
				s["etc"] = s["etc"] + 1
			}
			txMetrics <- s
		}
	}
}
