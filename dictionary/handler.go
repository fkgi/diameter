package dictionary

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/fkgi/diameter"
)

// Post function is called when recieve Diameter request
type Post func(path string, hdr http.Header, body io.Reader) (resp *http.Response, err error)

// RegisterHandler registers HTTP endpoints for the commands in the current dictionary.
func RegisterHandler(p Post, path string, rt diameter.Router) {
	for vid, vnd := range Data {
		if vid == 0 {
			continue
		}
		for aid, app := range vnd.Application {
			for cid, cmd := range app.Command {
				registerHandler(p, path+vnd.Name+"/"+app.Name+"/"+cmd,
					cid, aid, vid, rt)
			}
		}
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		httpErr(r.URL.Path, nil, http.StatusNotFound, "not found", "invalid URI path", w)
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

		txjson, e := encodeAVPtoJSON(avps)
		if e != nil {
			return diameterErr(path, nil, 0, nil, avps, diameter.InvalidAvpValue, e)
		}

		hdr := http.Header{}
		if retry {
			hdr.Add("X-Retry", "true")
		}
		r, e := p(path, hdr, bytes.NewBuffer(txjson))
		if e != nil {
			e = errors.New("unable to send HTTP request to backend: " + e.Error())
			return diameterErr(path, txjson, 0, nil, avps, diameter.UnableToDeliver, e)
		}
		defer r.Body.Close()

		switch r.StatusCode {
		case http.StatusOK:
		case http.StatusServiceUnavailable:
			if TraceTxHttpRequest != nil {
				TraceTxHttpRequest(path, txjson, r.StatusCode, nil, nil)
			}
			return true, nil
		default:
			e = errors.New("unknown error in HTTP")
			return diameterErr(path, txjson, r.StatusCode, nil, avps, diameter.UnableToComply, e)
		}

		avps, rxjson, e := decodeJSONtoAVP(r.Body)
		if e != nil {
			return diameterErr(path, txjson, r.StatusCode, nil, avps, diameter.UnableToComply, e)
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

		if TraceTxHttpRequest != nil {
			TraceTxHttpRequest(path, txjson, r.StatusCode, rxjson, nil)
		}
		return false, avps
	}
	handleTx := diameter.Handle(cid, aid, vid, serveDiameter, rt)

	serveHttp := func(w http.ResponseWriter, r *http.Request) {
		avps, txJson, e := decodeJSONtoAVP(r.Body)
		r.Body.Close()
		if e != nil {
			httpErr(path, txJson, http.StatusBadRequest,
				"failed to get AVP data", e.Error(), w)
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

		rxJson, e := encodeAVPtoJSON(avps)
		if e != nil {
			httpErr(path, txJson, http.StatusInternalServerError,
				"failed to get JSON data", e.Error(), w)
			return
		}

		if TraceRxHttpRequest != nil {
			TraceRxHttpRequest(path, txJson, http.StatusOK, rxJson, nil)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(rxJson)
	}
	http.HandleFunc("POST "+path, serveHttp)
}

func httpErr(
	path string, txj []byte, hcode int,
	title, detail string, w http.ResponseWriter) {

	data, _ := json.Marshal(struct {
		T string `json:"title"`
		D string `json:"detail"`
	}{T: title, D: detail})

	if TraceRxHttpRequest != nil {
		TraceRxHttpRequest(path, txj, hcode, data, errors.New(title+": "+detail))
	}

	w.Header().Add("Content-Type", "application/problem+json")
	w.WriteHeader(hcode)
	w.Write(data)
}

func diameterErr(
	path string, txj []byte, hcode int, rxj []byte,
	avp []diameter.AVP, dcode uint32, e error) (bool, []diameter.AVP) {

	if TraceTxHttpRequest != nil {
		TraceTxHttpRequest(path, txj, hcode, rxj, e)
	}

	ret := []diameter.AVP{}
	for _, a := range avp {
		if a.VendorID != 0 {
			continue
		}
		switch a.Code {
		case 277: // Auth-Session-State
			ret = append(ret, a)
		case 263: // Session-ID
			ret = append(ret, a)
		}
	}
	ret = append(ret, diameter.SetResultCode(dcode))
	ret = append(ret, diameter.SetOriginHost(diameter.Host))
	ret = append(ret, diameter.SetOriginRealm(diameter.Realm))
	ret = append(ret, diameter.SetErrorMessage(e.Error()))

	return true, ret
}

func encodeAVPtoJSON(avps []diameter.AVP) ([]byte, error) {
	if a, e := decodeAVPs(avps); e != nil {
		return nil, errors.New(
			"unable to decode Diameter AVP by dictionary: " + e.Error())
	} else if j, e := json.Marshal(a); e != nil {
		return nil, errors.New(
			"unable to marshal AVPs to JSON: " + e.Error())
	} else {
		return j, nil
	}
}

func decodeJSONtoAVP(r io.ReadCloser) ([]diameter.AVP, []byte, error) {
	data := make(map[string]any)
	if j, e := io.ReadAll(r); e != nil {
		return nil, nil, errors.New("unable to receive HTTP response: " + e.Error())
	} else if e = json.Unmarshal(j, &data); e != nil {
		return nil, j, errors.New("invalid JSON data of AVP: " + e.Error())
	} else if a, e := encodeAVPs(data); e != nil {
		return nil, j, errors.New("unable to encode Diameter AVP by dictionary: " + e.Error())
	} else {
		return a, j, nil
	}
}
