package dictionary

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/fkgi/diameter"
)

type Post func(path string, hdr http.Header, body io.Reader) (resp *http.Response, err error)

func (d XDictionary) RegisterHandler(p Post, path string, rt diameter.Router) {
	for _, vnd := range d.V {
		if vnd.I == 0 {
			continue
		}
		for _, app := range vnd.P {
			for _, cmd := range app.C {
				registerHandler(p, path+vnd.N+"/"+app.N+"/"+cmd.N,
					cmd.I, app.I, vnd.I, rt)
			}
		}
	}
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		httpErr("not found", "invalid URI path", http.StatusNotFound, w)
	})
}

func registerHandler(p Post, path string, cid, aid, vid uint32, rt diameter.Router) {
	serveDiameter := func(retry bool, avps []diameter.AVP) (bool, []diameter.AVP) {
		sid := ""
		for _, a := range avps {
			if a.Code == 263 {
				a.Decode(&sid)
				break
			}
		}

		data, e := DecodeAVPs(avps)
		if e != nil {
			return diameterErr(avps, diameter.InvalidAvpValue,
				"unable to decode Diameter AVP by dictionary: "+e.Error())
		}
		jsondata, e := json.Marshal(data)
		if e != nil {
			return diameterErr(avps, diameter.InvalidAvpValue,
				"unable to marshal AVPs to JSON: "+e.Error())
		}

		hdr := http.Header{}
		if retry {
			hdr.Add("X-Retry", "true")
		}
		r, e := p(path, hdr, bytes.NewBuffer(jsondata))
		if e != nil {
			return diameterErr(avps, diameter.UnableToDeliver,
				"unable to send HTTP request to backend: "+e.Error())
		}
		switch r.StatusCode {
		case http.StatusOK:
		case http.StatusServiceUnavailable:
			return true, nil
		default:
			return diameterErr(avps, diameter.UnableToComply,
				"error in HTTP")
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
		avps, e = EncodeAVPs(data)
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
		avps, e := EncodeAVPs(data)
		if e != nil {
			httpErr("unable to encode Diameter AVP by dictionary", e.Error(),
				http.StatusBadRequest, w)
			return
		}

		var route diameter.Identity
		for i := range avps {
			switch avps[i].Code {
			case 263: // Session-ID
				if len(avps[i].Data) == 0 {
					avps[i].Encode(diameter.NextSession(diameter.Host.String()))
				}
			case 264: // Origin-Host
				if len(avps[i].Data) == 0 {
					avps[i].Encode(diameter.Host)
				} else if e = avps[i].Decode(&route); e != nil {
					route = ""
				}
			case 296: // Origin-Realm
				if len(avps[i].Data) == 0 {
					avps[i].Encode(diameter.Realm)
				}
			}
		}
		if route != "" {
			avps = append(avps, diameter.SetRouteRecord(route))
		}

		retry := false
		if r.Header.Get("X-Retry") == "true" {
			retry = true
		}
		_, avps = handleTx(retry, avps)

		if data, e = DecodeAVPs(avps); e != nil {
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
	if NotifyHandlerError != nil {
		NotifyHandlerError("HTTP", title+": "+detail)
	}

	data, _ := json.Marshal(struct {
		T string `json:"title"`
		D string `json:"detail"`
	}{T: title, D: detail})

	w.Header().Add("Content-Type", "application/problem+json")
	w.WriteHeader(code)
	w.Write(data)
}

func diameterErr(avp []diameter.AVP, code uint32, err string) (bool, []diameter.AVP) {
	if NotifyHandlerError != nil {
		NotifyHandlerError("Diameter", err)
	}

	ret := []diameter.AVP{}
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
	ret = append(ret, diameter.SetResultCode(code))
	ret = append(ret, diameter.SetOriginHost(diameter.Host))
	ret = append(ret, diameter.SetOriginRealm(diameter.Realm))
	ret = append(ret, diameter.SetErrorMessage(err))

	return true, ret
}
