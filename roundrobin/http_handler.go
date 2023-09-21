package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/fkgi/diameter"
)

func listenAndServeHttp(addr string) {
	go http.ListenAndServe(
		addr,
		http.HandlerFunc(handleTx))
}

func handleTx(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	handler, ok := txHandlers[r.URL.Path]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	b, e := io.ReadAll(r.Body)
	defer r.Body.Close()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data := make(map[string]any)
	if e = json.Unmarshal(b, &data); e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	avps, e := parseAVPs(data)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, avps = handler(false, avps)
	if data, e = formatAVPs(avps); e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if b, e = json.Marshal(data); e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func parseAVPs(d map[string]any) ([]diameter.AVP, error) {
	avps := make([]diameter.AVP, 0, 10)
	for k, v := range d {
		a, e := encAVPs[k](v)
		if e != nil {
			return nil, errors.Join(errors.New(k+" is invalid"), e)
		}
		avps = append(avps, a)
	}
	return avps, nil
}
