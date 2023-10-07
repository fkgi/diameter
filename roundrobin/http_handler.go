package main

import (
	"encoding/json"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/fkgi/diameter/connector"
	"github.com/fkgi/diameter/dictionary"
)

func handleTx(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Add("Allow", "POST")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if !strings.HasPrefix(r.URL.Path, apipath) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	msg, e := dictionary.EncodeMessage(r.URL.Path[len(apipath):])
	if e != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	b, e := io.ReadAll(r.Body)
	defer r.Body.Close()
	if e != nil {
		errorAnswer(
			"unable to read HTTP request body", e.Error(),
			http.StatusBadRequest, w)
		return
	}

	data := make(map[string]any)
	if e = json.Unmarshal(b, &data); e != nil {
		errorAnswer(
			"invalid JSON data of AVPs", e.Error(),
			http.StatusBadRequest, w)
		return
	}

	avps, e := parseAVPs(data)
	if e != nil {
		errorAnswer(
			"unable to encode Diameter AVP by dictionary", e.Error(),
			http.StatusBadRequest, w)
		return
	}
	msg.SetAVP(avps)

	msg.EtEID = nextEtE()
	msg.FlgP = true

	msg = connector.DefaultTxHandler(msg)
	if avps, e = msg.GetAVP(); e != nil {
		errorAnswer(
			"unable to get AVPs from message", e.Error(),
			http.StatusBadRequest, w)
		return
	}
	if data, e = formatAVPs(avps); e != nil {
		errorAnswer(
			"unable to decode Diameter AVP by dictionary", e.Error(),
			http.StatusBadRequest, w)
		return
	}
	if b, e = json.Marshal(data); e != nil {
		errorAnswer("unable to marshal AVPs to JSON", e.Error(),
			http.StatusInternalServerError, w)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func errorAnswer(title, detail string, code int, w http.ResponseWriter) {
	log.Println(title+":", detail)
	data, _ := json.Marshal(struct {
		T string `json:"title"`
		D string `json:"detail"`
	}{T: title, D: detail})
	w.Header().Add("Content-Type", "application/problem+json")
	w.WriteHeader(code)
	w.Write(data)
}

var eteID = make(chan uint32, 1) // End-to-End ID source

func init() {
	ut := time.Now().Unix()
	eteID <- (uint32(ut^0xFFF) << 20) | (rand.Uint32() ^ 0xFFFFF)
}

func nextEtE() uint32 {
	ret := <-eteID
	eteID <- ret + 1
	return ret
}
