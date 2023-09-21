package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"net"
	"os"
	"time"

	"github.com/fkgi/diameter"
)

const pathPrefix = "/msg/v1/"

var (
	encAVPs    = make(map[string]func(any) (diameter.AVP, error))
	decAVPs    = make(map[uint32]map[uint32]func(diameter.AVP) (string, any, error))
	txHandlers = make(map[string]diameter.Handler)
)

func loadDictionary(path string) error {
	log.Println("loading dictionary file", path)

	var dictionary map[string]struct {
		ID   uint32 `json:"id"`
		Apps map[string]struct {
			ID   uint32 `json:"id"`
			Cmds map[string]struct {
				ID uint32 `json:"id"`
			} `json:"command"`
		} `json:"applications"`
		Avps map[string]struct {
			ID    uint32           `json:"id"`
			Mflag bool             `json:"mandatory,omitempty"`
			Type  string           `json:"type"`
			Enum  map[string]int32 `json:"map,omitempty"`
		} `json:"avps"`
	}

	if data, e := os.ReadFile(path); e != nil {
		return errors.Join(
			errors.New("failed to read dictionary file"), e)
	} else if e = json.Unmarshal(data, &dictionary); e != nil {
		return errors.Join(
			errors.New("failed to unmarshal dictionary file"), e)
	}

	for vndname, vnd := range dictionary {
		if _, ok := decAVPs[uint32(vnd.ID)]; ok {
			return errors.New("duplicated vendor definition: " + vndname)
		}
		decAVPs[uint32(vnd.ID)] = make(map[uint32]func(diameter.AVP) (string, any, error))

		log.Println("supported vendor:", vndname, "(", vnd.ID, ")")
		for appname, app := range vnd.Apps {
			log.Println("supported application:", appname, "(", app.ID, ")")
			for cmdname, cmd := range app.Cmds {
				p := pathPrefix + appname + "/" + cmdname
				txHandlers[p] = diameter.Handle(
					uint32(cmd.ID), uint32(app.ID), uint32(vnd.ID),
					func(flg bool, avps []diameter.AVP) (bool, []diameter.AVP) {
						return handleRx(p, flg, avps)
					})
			}
		}

		for name, avp := range vnd.Avps {
			if _, ok := encAVPs[name]; ok {
				return errors.New("duplicated AVP definition: " + name)
			}
			if _, ok := decAVPs[uint32(vnd.ID)][uint32(avp.ID)]; ok {
				return errors.New("duplicated AVP definition: " + name)
			}

			var encf func(any, *diameter.AVP) error
			var decf func(*diameter.AVP) (any, error)
			switch avp.Type {
			case "OctetString":
				encf = encOctetString
				decf = decOctetString
			case "Integer32":
				encf = encInteger32
				decf = decInteger32
			case "Integer64":
				encf = encInteger64
				decf = decInteger64
			case "Unsigned32":
				encf = encUnsigned32
				decf = decUnsigned32
			case "Unsigned64":
				encf = encUnsigned64
				decf = decUnsigned64
			case "Float32":
				encf = encFloat32
				decf = decFloat32
			case "Float64":
				encf = encFloat64
				decf = decFloat64
			case "Grouped":
				encf = encGrouped
				decf = decGrouped
			case "Address":
				encf = encAddress
				decf = decAddress
			case "Time":
				encf = encTime
				decf = decTime
			case "UTF8String":
				encf = encUTF8String
				decf = decUTF8String
			case "DiameterIdentity":
				encf = encDiameterIdentity
				decf = decDiameterIdentity
			case "DiameterURI":
				encf = encDiameterURI
				decf = decDiameterURI
			case "Enumerated":
				enum := avp.Enum
				encf = func(v any, a *diameter.AVP) error {
					return encEnumerated(v, a, enum)
				}
				m := make(map[int32]string)
				for k, v := range avp.Enum {
					m[v] = k
				}
				decf = func(a *diameter.AVP) (any, error) {
					return decEnumerated(a, m)
				}
			case "IPFilterRule":
				encf = encIPFilterRule
				decf = decIPFilterRule
			default:
				return errors.New("infalid AVP type: " + name)
			}

			code := uint32(avp.ID)
			vid := uint32(vnd.ID)
			flg := avp.Mflag
			encAVPs[name] = func(v any) (diameter.AVP, error) {
				a := diameter.AVP{
					Code:      code,
					VendorID:  vid,
					Mandatory: flg}
				e := encf(v, &a)
				return a, e
			}

			n := name
			decAVPs[uint32(vnd.ID)][uint32(avp.ID)] = func(a diameter.AVP) (string, any, error) {
				v, e := decf(&a)
				return n, v, e
			}
		}
	}
	return nil
}

func encOctetString(v any, avp *diameter.AVP) error {
	s, ok := v.(string)
	if !ok {
		return errors.New("not String")
	}
	a, e := hex.DecodeString(s)
	if e != nil {
		return e
	}
	return avp.Encode(a)
}

func decOctetString(avp *diameter.AVP) (any, error) {
	d := new([]byte)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return hex.EncodeToString(*d), nil
}

func encInteger32(v any, avp *diameter.AVP) error {
	d, ok := v.(float64)
	if !ok {
		return errors.New("not Number")
	}
	return avp.Encode(int32(d))
}

func decInteger32(avp *diameter.AVP) (any, error) {
	d := new(int32)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return *d, nil
}

func encInteger64(v any, avp *diameter.AVP) error {
	d, ok := v.(float64)
	if !ok {
		return errors.New("not Number")
	}
	return avp.Encode(int64(d))
}

func decInteger64(avp *diameter.AVP) (any, error) {
	d := new(int64)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return *d, nil
}

func encUnsigned32(v any, avp *diameter.AVP) error {
	d, ok := v.(float64)
	if !ok {
		return errors.New("not Number")
	}
	return avp.Encode(uint32(d))
}

func decUnsigned32(avp *diameter.AVP) (any, error) {
	d := new(uint32)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return *d, nil
}

func encUnsigned64(v any, avp *diameter.AVP) error {
	d, ok := v.(float64)
	if !ok {
		return errors.New("not Number")
	}
	return avp.Encode(uint64(d))
}

func decUnsigned64(avp *diameter.AVP) (any, error) {
	d := new(uint64)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return *d, nil
}

func encFloat32(v any, avp *diameter.AVP) error {
	d, ok := v.(float64)
	if !ok {
		return errors.New("not Number")
	}
	return avp.Encode(float32(d))
}

func decFloat32(avp *diameter.AVP) (any, error) {
	d := new(float32)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return *d, nil
}

func encFloat64(v any, avp *diameter.AVP) error {
	d, ok := v.(float64)
	if !ok {
		return errors.New("not Number")
	}
	return avp.Encode(d)
}

func decFloat64(avp *diameter.AVP) (any, error) {
	d := new(float64)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return *d, nil
}

func encGrouped(v any, avp *diameter.AVP) (e error) {
	a, ok := v.(map[string]any)
	if !ok {
		return errors.New("not Grouped")
	}

	buf := new(bytes.Buffer)
	for k, v := range a {
		avp, e := encAVPs[k](v)
		if e != nil {
			return errors.Join(errors.New(k+" is invalid"), e)
		}
		avp.MarshalTo(buf)
	}
	avp.Data = buf.Bytes()
	return
}

func decGrouped(avp *diameter.AVP) (any, error) {
	result := make(map[string]any)
	for buf := bytes.NewBuffer(avp.Data); buf.Len() != 0; {
		a := diameter.AVP{}
		e := a.UnmarshalFrom(buf)
		if e != nil {
			return nil, e
		}
		n, v, e := decAVPs[a.VendorID][a.Code](a)
		if e != nil {
			return nil, e
		}
		result[n] = v
	}
	return result, nil
}

func encAddress(v any, avp *diameter.AVP) error {
	s, ok := v.(string)
	if !ok {
		return errors.New("not String")
	}
	a := net.ParseIP(s)
	if a == nil {
		return errors.New("not Address")
	}
	return avp.Encode(a)
}

func decAddress(avp *diameter.AVP) (any, error) {
	d := new(net.IP)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return d.String(), nil
}

func encTime(v any, avp *diameter.AVP) error {
	s, ok := v.(string)
	if !ok {
		return errors.New("not String")
	}
	a, e := time.Parse(time.RFC3339, s)
	if e != nil {
		return e
	}
	return avp.Encode(a)
}

func decTime(avp *diameter.AVP) (any, error) {
	d := new(time.Time)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return d.Format(time.RFC3339), nil
}

func encUTF8String(v any, avp *diameter.AVP) error {
	s, ok := v.(string)
	if !ok {
		return errors.New("not String")
	}
	return avp.Encode(s)
}

func decUTF8String(avp *diameter.AVP) (any, error) {
	d := new(string)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return *d, nil
}

func encDiameterIdentity(v any, avp *diameter.AVP) error {
	s, ok := v.(string)
	if !ok {
		return errors.New("not String")
	}
	a, e := diameter.ParseIdentity(s)
	if e != nil {
		return e
	}
	return avp.Encode(a)
}

func decDiameterIdentity(avp *diameter.AVP) (any, error) {
	d := new(diameter.Identity)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return d.String(), nil
}

func encDiameterURI(v any, avp *diameter.AVP) error {
	s, ok := v.(string)
	if !ok {
		return errors.New("not String")
	}
	a, e := diameter.ParseURI(s)
	if e != nil {
		return e
	}
	return avp.Encode(a)
}

func decDiameterURI(avp *diameter.AVP) (any, error) {
	d := new(diameter.URI)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return d.String(), nil
}

func encEnumerated(v any, avp *diameter.AVP, enum map[string]int32) error {
	s, ok := v.(string)
	if !ok {
		return errors.New("not String")
	}
	a, ok := enum[s]
	if !ok {
		return errors.New("not defined Enumerated")
	}
	return avp.Encode(diameter.Enumerated(a))
}

func decEnumerated(avp *diameter.AVP, enum map[int32]string) (any, error) {
	d := new(diameter.Enumerated)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	a, ok := enum[int32(*d)]
	if !ok {
		return nil, errors.New("not defined Enumerated")
	}
	return a, nil
}

func encIPFilterRule(v any, avp *diameter.AVP) error {
	s, ok := v.(diameter.IPFilterRule)
	if !ok {
		return errors.New("not String")
	}
	return avp.Encode(s)
}

func decIPFilterRule(avp *diameter.AVP) (any, error) {
	d := new(diameter.IPFilterRule)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return *d, nil
}
