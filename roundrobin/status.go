package main

import (
	"fmt"
	"net/http"

	"github.com/fkgi/diameter"
	"github.com/fkgi/diameter/connector"
)

const constatFmt = `{
	"state": "%s",
	"local": {
		"host": "%s",
		"realm": "%s",
		"address": "%s"
	},
	"peer": {
		"host": "%s",
		"realm": "%s",
		"address": "%s"
	}
}`

func conStateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Add("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(constatFmt,
		connector.State(),
		diameter.Host,
		diameter.Realm,
		connector.LocalAddr(),
		connector.PeerName(),
		connector.PeerRealm(),
		connector.PeerAddr())))
}

var (
	rxReq  uint64
	txDisc uint64
	txAns  [6]uint64

	txReq  uint64
	rxIvld uint64
	rxAns  [6]uint64
)

const statsFmt = `{
	"rx_request": %d,
	"tx_discard": %d,
	"tx_etc": %d,
	"tx_1xxx": %d,
	"tx_2xxx": %d,
	"tx_3xxx": %d,
	"tx_4xxx": %d,
	"tx_5xxx": %d,
	"tx_request": %d,
	"rx_invalid": %d,
	"rx_etc": %d,
	"rx_1xxx": %d,
	"rx_2xxx": %d,
	"rx_3xxx": %d,
	"rx_4xxx": %d,
	"rx_5xxx": %d
}`

func statsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Add("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(statsFmt,
		rxReq, txDisc,
		txAns[0], txAns[1], txAns[2], txAns[3], txAns[4], txAns[5],
		txReq, rxIvld,
		rxAns[0], rxAns[1], rxAns[2], rxAns[3], rxAns[4], rxAns[5])))
}
