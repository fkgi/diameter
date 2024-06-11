package dictionary

import (
	"encoding/json"
	"errors"

	"github.com/fkgi/diameter"
)

type Dictionary map[string]struct {
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

var (
	encAVPs    = make(map[string]func(any) (diameter.AVP, error))
	decAVPs    = make(map[uint64]func(diameter.AVP) (string, any, error))
	encCommand = make(map[string]uint64)
	decCommand = make(map[uint64]string)
)

func EncodeAVP(name string, value any) (diameter.AVP, error) {
	f, ok := encAVPs[name]
	if !ok {
		return diameter.AVP{}, errors.New("unknown AVP name")
	}
	return f(value)
}

func DecodeAVP(a diameter.AVP) (string, any, error) {
	f, ok := decAVPs[(uint64(a.VendorID)<<32)|uint64(a.Code)]
	if !ok {
		return "", nil, errors.New("unknown AVP")
	}
	return f(a)
}

func EncodeMessage(name string) (m diameter.Message, e error) {
	id, ok := encCommand[name]
	if !ok {
		e = errors.New("unknown command name")
	} else {
		m.AppID = uint32((id & 0xffffffff00000000) >> 32)
		m.Code = uint32(id & 0x00000000ffffffff)
		m.FlgR = true
	}
	return
}

func DecodeMessage(m diameter.Message) (string, error) {
	name, ok := decCommand[(uint64(m.AppID)<<32)|uint64(m.Code)]
	if !ok {
		return "", errors.New("unknown command")
	}
	return name, nil
}

func LoadDictionary(data []byte) (Dictionary, error) {
	var dictionary Dictionary

	if e := json.Unmarshal(data, &dictionary); e != nil {
		return dictionary, errors.Join(
			errors.New("failed to unmarshal dictionary file"), e)
	}

	for vndname, vnd := range dictionary {
		for appname, app := range vnd.Apps {
			for cmdname, cmd := range app.Cmds {
				p := vndname + "/" + appname + "/" + cmdname
				i := (uint64(app.ID) << 32) | uint64(cmd.ID)
				encCommand[p] = i
				decCommand[i] = p
			}
		}

		for name, avp := range vnd.Avps {
			if _, ok := encAVPs[name]; ok {
				return dictionary,
					errors.New("duplicated AVP definition: " + name)
			}
			if _, ok := decAVPs[(uint64(vnd.ID)<<32)|uint64(avp.ID)]; ok {
				return dictionary,
					errors.New("duplicated AVP definition: " + name)
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
				return dictionary,
					errors.New("invalid AVP type: " + name)
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
			decAVPs[(uint64(vnd.ID)<<32)|uint64(avp.ID)] =
				func(a diameter.AVP) (string, any, error) {
					v, e := decf(&a)
					return n, v, e
				}
		}
	}
	return dictionary, nil
}
