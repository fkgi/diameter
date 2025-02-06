package main

import (
	"encoding/json"
	"net/http"
)

type constat struct {
	Host  string   `json:"host"`
	Realm string   `json:"realm"`
	Addr  string   `json:"address"`
	State string   `json:"state"`
	Apps  []uint32 `json:"apps"`
}

func conStateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Add("Allow", "GET")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	stats := []constat{}
	for _, c := range cons {
		stats = append(stats, constat{
			Host:  c.Host.String(),
			Realm: c.Realm.String(),
			Addr:  c.PeerAddr().String(),
			State: c.State(),
			Apps:  c.AvailableApplications()})
	}
	if b, e := json.Marshal(stats); e != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}
