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
		return makeError(m, diameter.UnableToDeliver,
			"unable to decode Diameter message by dictionary: "+e.Error())
	}
	avps, e := m.GetAVP()
	if e != nil {
		return makeError(m, diameter.InvalidAvpValue,
			"unable to get AVPs from message: "+e.Error())
	}
	data, e := formatAVPs(avps)
	if e != nil {
		return makeError(m, diameter.InvalidAvpValue,
			"unable to decode Diameter AVP by dictionary: "+e.Error())
	}
	jsondata, e := json.Marshal(data)
	if e != nil {
		return makeError(m, diameter.InvalidAvpValue,
			"unable to marshal AVPs to JSON: "+e.Error())
	}
	if rxPath == "" {
		return makeError(m, diameter.UnableToDeliver,
			"no HTTP destination is defined")
	}
	r, e := http.Post(
		rxPath+path,
		"application/json",
		bytes.NewBuffer(jsondata))
	if e != nil {
		return makeError(m, diameter.UnableToDeliver,
			"unable to send HTTP request: "+e.Error())
	}

	b, e := io.ReadAll(r.Body)
	defer r.Body.Close()
	if e != nil {
		return makeError(m, diameter.UnableToDeliver,
			"unable to receive HTTP response: "+e.Error())
	}
	data = make(map[string]any)
	if e = json.Unmarshal(b, &data); e != nil {
		return makeError(m, diameter.UnableToComply,
			"invalid JSON data of AVP: "+e.Error())
	}
	avps, e = parseAVPs(data)
	if e != nil {
		return makeError(m, diameter.UnableToComply,
			"unable to encode Diameter AVP by dictionary: "+e.Error())
	}

	m.FlgR = false
	m.SetAVP(avps)
	return m
}

func makeError(req diameter.Message, code uint32, err string) (ans diameter.Message) {
	log.Println(err)

	ans.AppID = req.AppID
	ans.Code = req.Code
	ans.EtEID = req.EtEID
	ans.FlgE = true
	ans.FlgP = req.FlgP
	ans.FlgR = false
	ans.FlgT = false
	ans.SetAVP([]diameter.AVP{
		diameter.SetResultCode(code),
		diameter.SetOriginHost(diameter.Host),
		diameter.SetOriginRealm(diameter.Realm),
		diameter.SetErrorMessage(err)})
	return
}
