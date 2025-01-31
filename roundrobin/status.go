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
	s := connector.Stats()
	w.Write([]byte(fmt.Sprintf(statsFmt,
		s[0], s[1], s[2], s[3], s[4], s[5], s[6], s[7], s[8], s[9], s[10], s[11], s[12], s[13], s[14], s[15])))
}
