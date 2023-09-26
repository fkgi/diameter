package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/fkgi/diameter"
	"github.com/fkgi/diameter/dictionary"
)

var rxPath string

func handleRx(m diameter.Message) diameter.Message {
	path, e := dictionary.DecodeMessage(m)
	if e != nil {
		log.Println(e)
		return makeError(m, diameter.UnableToDeliver)
	}
	avps, e := m.GetAVP()
	if e != nil {
		log.Println(e)
		return makeError(m, diameter.InvalidAvpValue)
	}
	data, e := formatAVPs(avps)
	if e != nil {
		log.Println(e)
		return makeError(m, diameter.InvalidAvpValue)
	}
	jsondata, e := json.Marshal(data)
	if e != nil {
		log.Println(e)
		return makeError(m, diameter.InvalidAvpValue)
	}
	r, e := http.Post(rxPath+path, "application/json", bytes.NewBuffer(jsondata))
	if e != nil {
		log.Println(e)
		return makeError(m, diameter.UnableToDeliver)
	}

	b, e := io.ReadAll(r.Body)
	defer r.Body.Close()
	if e != nil {
		log.Println(e)
		return makeError(m, diameter.UnableToDeliver)
	}
	data = make(map[string]any)
	if e = json.Unmarshal(b, &data); e != nil {
		log.Println(e)
		return makeError(m, diameter.UnableToComply)
	}
	avps, e = parseAVPs(data)
	if e != nil {
		log.Println(e)
		return makeError(m, diameter.UnableToComply)
	}

	m.FlgR = false
	m.SetAVP(avps)
	return m
}

func makeError(req diameter.Message, code uint32) (ans diameter.Message) {
	ans.AppID = req.AppID
	ans.Code = req.Code
	ans.FlgE = true
	ans.FlgP = req.FlgP
	ans.FlgR = false
	ans.FlgT = false
	ans.SetAVP([]diameter.AVP{
		diameter.SetResultCode(code),
		diameter.SetOriginHost(diameter.Host),
		diameter.SetOriginRealm(diameter.Realm)})
	return
}

func formatAVPs(avps []diameter.AVP) (map[string]any, error) {
	result := make(map[string]any)
	for _, a := range avps {
		n, v, e := dictionary.DecodeAVP(a)
		if e != nil {
			return nil, e
		}
		result[n] = v
	}
	return result, nil
}
