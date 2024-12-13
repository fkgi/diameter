package dictionary

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/fkgi/diameter"
)

func (d XDictionary) RegisterHandler(backend, path string, rt diameter.Router) {
	for _, vnd := range d.V {
		if vnd.I == 0 {
			continue
		}
		for _, app := range vnd.P {
			for _, cmd := range app.C {
				registerHandler(backend, path+vnd.N+"/"+app.N+"/"+cmd.N,
					cmd.I, app.I, vnd.I, rt)
			}
		}
	}
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		httpErr("not found", "invalid URI path", http.StatusNotFound, w)
	})
}

func registerHandler(backend, path string, cid, aid, vid uint32, rt diameter.Router) {
	serveDiameter := func(_ bool, avps []diameter.AVP) (bool, []diameter.AVP) {
		if backend == "" {
			return diameterErr(avps, diameter.UnableToDeliver,
				"no HTTP backend is defined")
		}

		sid := ""
		for _, a := range avps {
			if a.Code == 263 {
				a.Decode(&sid)
				break
			}
		}

		data, e := formatAVPs(avps)
		if e != nil {
			return diameterErr(avps, diameter.InvalidAvpValue,
				"unable to decode Diameter AVP by dictionary: "+e.Error())
		}
		jsondata, e := json.Marshal(data)
		if e != nil {
			return diameterErr(avps, diameter.InvalidAvpValue,
				"unable to marshal AVPs to JSON: "+e.Error())
		}
		r, e := http.Post(backend+path, "application/json", bytes.NewBuffer(jsondata))
		if e != nil {
			return diameterErr(avps, diameter.UnableToDeliver,
				"unable to send HTTP request to backend: "+e.Error())
		}

		jsondata, e = io.ReadAll(r.Body)
		defer r.Body.Close()
		if e != nil {
			return diameterErr(avps, diameter.UnableToDeliver,
				"unable to receive HTTP response: "+e.Error())
		}
		data = make(map[string]any)
		if e = json.Unmarshal(jsondata, &data); e != nil {
			return diameterErr(avps, diameter.UnableToComply,
				"invalid JSON data of AVP: "+e.Error())
		}
		avps, e = parseAVPs(data)
		if e != nil {
			return diameterErr(avps, diameter.UnableToComply,
				"unable to encode Diameter AVP by dictionary: "+e.Error())
		}

		for i := range avps {
			if len(avps[i].Data) != 0 {
				continue
			}
			switch avps[i].Code {
			case 263: // Session-ID
				avps[i].Encode(sid)
			case 264: // Origin-Host
				avps[i].Encode(diameter.Host)
			case 296: // Origin-Realm
				avps[i].Encode(diameter.Realm)
			}
		}

		return false, avps
	}
	handleTx := diameter.Handle(cid, aid, vid, serveDiameter, rt)

	serveHttp := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Add("Allow", "POST")
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		jsondata, e := io.ReadAll(r.Body)
		defer r.Body.Close()
		if e != nil {
			httpErr("unable to read HTTP request body", e.Error(),
				http.StatusBadRequest, w)
			return
		}
		data := make(map[string]any)
		if e = json.Unmarshal(jsondata, &data); e != nil {
			httpErr("invalid JSON data of AVPs", e.Error(),
				http.StatusBadRequest, w)
			return
		}
		avps, e := parseAVPs(data)
		if e != nil {
			httpErr("unable to encode Diameter AVP by dictionary", e.Error(),
				http.StatusBadRequest, w)
			return
		}

		proxy := true
		for i := range avps {
			if len(avps[i].Data) == 0 {
				continue
			}
			switch avps[i].Code {
			case 263: // Session-ID
				avps[i].Encode(diameter.NextSession(diameter.Host.String()))
			case 264: // Origin-Host
				avps[i].Encode(diameter.Host)
				proxy = false
			case 296: // Origin-Realm
				avps[i].Encode(diameter.Realm)
			}
		}
		if proxy {
			avps = append(avps, diameter.SetRouteRecord(diameter.Host))
		}

		_, avps = handleTx(true, avps)

		if data, e = formatAVPs(avps); e != nil {
			httpErr("unable to decode Diameter AVP by dictionary", e.Error(),
				http.StatusBadRequest, w)
			return
		}
		if jsondata, e = json.Marshal(data); e != nil {
			httpErr("unable to marshal AVPs to JSON", e.Error(),
				http.StatusInternalServerError, w)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsondata)
	}
	http.HandleFunc(path, serveHttp)
}

func httpErr(title, detail string, code int, w http.ResponseWriter) {
	log.Println(title+":", detail)

	data, _ := json.Marshal(struct {
		T string `json:"title"`
		D string `json:"detail"`
	}{T: title, D: detail})

	w.Header().Add("Content-Type", "application/problem+json")
	w.WriteHeader(code)
	w.Write(data)
}

func diameterErr(avp []diameter.AVP, code uint32, err string) (bool, []diameter.AVP) {
	log.Println(err)

	ret := []diameter.AVP{
		diameter.SetResultCode(code),
		diameter.SetOriginHost(diameter.Host),
		diameter.SetOriginRealm(diameter.Realm),
		diameter.SetErrorMessage(err)}

	for _, a := range avp {
		if a.VendorID != 0 {
			continue
		}
		switch a.Code {
		case 277:
			ret = append(ret, a)
		case 263:
			ret = append(ret, a)
		}
	}

	return true, ret
}

func formatAVPs(avps []diameter.AVP) (map[string]any, error) {
	result := make(map[string]any)
	for _, a := range avps {
		n, v, e := DecodeAVP(a)
		if e != nil {
			return nil, e
		}
		result[n] = v
	}
	return result, nil
}

func parseAVPs(d map[string]any) ([]diameter.AVP, error) {
	avps := make([]diameter.AVP, 0, 10)
	for k, v := range d {
		a, e := EncodeAVP(k, v)
		if e != nil {
			return nil, fmt.Errorf("%s is invalid: %v", k, e)
		}
		avps = append(avps, a)
	}
	return avps, nil
}
