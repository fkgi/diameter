package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/fkgi/diameter"
)

var rxPath string

func handleRx(path string, flag bool, avps []diameter.AVP) (bool, []diameter.AVP) {
	data, e := formatAVPs(avps)
	if e != nil {
		log.Println(e)
		return true, []diameter.AVP{
			diameter.SetResultCode(diameter.InvalidAvpValue),
			diameter.SetOriginHost(diameter.Local.Host),
			diameter.SetOriginRealm(diameter.Local.Realm)}
	}
	jsondata, e := json.Marshal(data)
	if e != nil {
		log.Println(e)
		return true, []diameter.AVP{
			diameter.SetResultCode(diameter.InvalidAvpValue),
			diameter.SetOriginHost(diameter.Local.Host),
			diameter.SetOriginRealm(diameter.Local.Realm)}
	}
	r, e := http.Post(rxPath+path, "application/json", bytes.NewBuffer(jsondata))
	if e != nil {
		log.Println(e)
		return true, []diameter.AVP{
			diameter.SetResultCode(diameter.UnableToDeliver),
			diameter.SetOriginHost(diameter.Local.Host),
			diameter.SetOriginRealm(diameter.Local.Realm)}
	}

	b, e := io.ReadAll(r.Body)
	defer r.Body.Close()
	if e != nil {
		log.Println(e)
		return true, []diameter.AVP{
			diameter.SetResultCode(diameter.UnableToDeliver),
			diameter.SetOriginHost(diameter.Local.Host),
			diameter.SetOriginRealm(diameter.Local.Realm)}
	}
	data = make(map[string]any)
	if e = json.Unmarshal(b, &data); e != nil {
		log.Println(e)
		return true, []diameter.AVP{
			diameter.SetResultCode(diameter.UnableToComply),
			diameter.SetOriginHost(diameter.Local.Host),
			diameter.SetOriginRealm(diameter.Local.Realm)}
	}
	avps, e = parseAVPs(data)
	if e != nil {
		log.Println(e)
		return true, []diameter.AVP{
			diameter.SetResultCode(diameter.UnableToComply),
			diameter.SetOriginHost(diameter.Local.Host),
			diameter.SetOriginRealm(diameter.Local.Realm)}
	}
	return false, avps
}

func formatAVPs(avps []diameter.AVP) (map[string]any, error) {
	result := make(map[string]any)
	for _, a := range avps {
		n, v, e := decAVPs[a.VendorID][a.Code](a)
		if e != nil {
			return nil, e
		}
		result[n] = v
	}
	return result, nil
}
