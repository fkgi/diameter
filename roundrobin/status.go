package main

import (
	"encoding/json"
	"net/http"

	"github.com/fkgi/diameter"
)

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

var statistics = struct {
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
}{}

func statsHandler(w http.ResponseWriter, r *http.Request) {
	if jd, e := json.Marshal(statistics); e != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jd)
	}
}
