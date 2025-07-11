package dictionary

import (
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"slices"

	"github.com/fkgi/diameter"
)

type XDictionary struct {
	XMLName xml.Name `xml:"dictionary"`
	V       []struct {
		N string `xml:"name,attr"`
		I uint32 `xml:"id,attr"`
		P []struct {
			N string `xml:"name,attr"`
			I uint32 `xml:"id,attr"`
			C []struct {
				N string `xml:"name,attr"`
				I uint32 `xml:"id,attr"`
			} `xml:"command"`
		} `xml:"application"`
		V []struct {
			N string `xml:"name,attr"`
			I uint32 `xml:"id,attr"`
			T string `xml:"type,attr"`
			M bool   `xml:"mandatory,attr"`
			P bool   `xml:"protected,attr"`
			R bool   `xml:"reserved,attr"`
			E []struct {
				I int32  `xml:"value,attr"`
				V string `xml:",chardata"`
			} `xml:"enum"`
		} `xml:"avp"`
	} `xml:"vendor"`
}

var (
	encAVPs    = make(map[string]func(any) (diameter.AVP, error))
	decAVPs    = make(map[uint64]func(diameter.AVP) (string, any, error))
	encCommand = make(map[string]uint64)
	decCommand = make(map[uint64]string)
)

func EncodeAVPs(d map[string]any) ([]diameter.AVP, error) {
	avps := map[uint32][]diameter.AVP{}
	codes := make([]uint32, 0, 20)
	for k, v := range d {
		if l, ok := v.([]any); ok {
			for _, v := range l {
				a, e := EncodeAVP(k, v)
				if e != nil {
					return nil, fmt.Errorf("%s is invalid: %v", k, e)
				}
				if _, ok := avps[a.Code]; ok {
					avps[a.Code] = append(avps[a.Code], a)
				} else {
					avps[a.Code] = []diameter.AVP{a}
					codes = append(codes, a.Code)
				}
			}
		} else {
			a, e := EncodeAVP(k, v)
			if e != nil {
				return nil, fmt.Errorf("%s is invalid: %v", k, e)
			}
			avps[a.Code] = []diameter.AVP{a}
			codes = append(codes, a.Code)
		}
	}
	slices.Sort(codes)

	res := make([]diameter.AVP, 0, 20)
	for _, k := range order {
		if l, ok := avps[k]; ok {
			res = append(res, l...)
			delete(avps, k)
		}
	}
	for _, k := range codes {
		if l, ok := avps[k]; ok {
			res = append(res, l...)
		}
	}

	return res, nil
}

var order = []uint32{263, 301, 260, 268, 298, 277, 264, 296, 293, 283}

func EncodeAVP(name string, value any) (diameter.AVP, error) {
	f, ok := encAVPs[name]
	if !ok {
		return diameter.AVP{}, errors.New("unknown AVP name")
	}
	return f(value)
}

func DecodeAVPs(avps []diameter.AVP) (map[string]any, error) {
	result := make(map[string][]any)
	for _, a := range avps {
		n, v, e := DecodeAVP(a)
		if e != nil {
			return nil, e
		}
		if l, ok := result[n]; ok {
			result[n] = append(l, v)
		} else {
			result[n] = []any{v}
		}
	}

	compat := make(map[string]any, len(result))
	for k, v := range result {
		if len(v) == 1 {
			compat[k] = v[0]
		} else {
			compat[k] = v
		}
	}
	return compat, nil
}

func DecodeAVP(a diameter.AVP) (string, any, error) {
	f, ok := decAVPs[(uint64(a.VendorID)<<32)|uint64(a.Code)]
	if !ok {
		// return "", nil, errors.New("unknown AVP")
		return fmt.Sprintf("UNKNOWN(%d)", a.Code), hex.EncodeToString(a.Data), nil
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
		return fmt.Sprintf("UNKNOWN(%d)", m.Code), nil

	}
	return name, nil
}

func LoadDictionary(data []byte) (XDictionary, error) {
	var xd XDictionary
	if e := xml.Unmarshal(data, &xd); e != nil {
		return xd, e
	}

	for _, vnd := range xd.V {
		for _, app := range vnd.P {
			for _, cmd := range app.C {
				p := vnd.N + "/" + app.N + "/" + cmd.N
				i := (uint64(app.I) << 32) | uint64(cmd.I)
				encCommand[p] = i
				decCommand[i] = p
			}
		}

		for _, avp := range vnd.V {
			if _, ok := encAVPs[avp.N]; ok {
				return xd,
					errors.New("duplicated AVP definition: " + avp.N)
			}
			if _, ok := decAVPs[(uint64(vnd.I)<<32)|uint64(avp.I)]; ok {
				return xd,
					errors.New("duplicated AVP definition: " + avp.N)
			}

			var encf func(any, *diameter.AVP) error
			var decf func(*diameter.AVP) (any, error)
			switch avp.T {
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
				m1 := make(map[string]int32)
				m2 := make(map[int32]string)
				for _, enm := range avp.E {
					m1[enm.V] = enm.I
					m2[enm.I] = enm.V
				}
				encf = func(v any, a *diameter.AVP) error {
					return encEnumerated(v, a, m1)
				}
				decf = func(a *diameter.AVP) (any, error) {
					return decEnumerated(a, m2)
				}
			case "IPFilterRule":
				encf = encIPFilterRule
				decf = decIPFilterRule
			default:
				return xd,
					errors.New("invalid AVP type: " + avp.N)
			}

			code := uint32(avp.I)
			vid := uint32(vnd.I)
			mflg := avp.M
			pflg := avp.P
			rflg := avp.R
			encAVPs[avp.N] = func(v any) (diameter.AVP, error) {
				a := diameter.AVP{
					Code:      code,
					VendorID:  vid,
					Mandatory: mflg,
					Protected: pflg,
					Reserved:  [5]bool{rflg, rflg, rflg, rflg, rflg}}
				e := encf(v, &a)
				return a, e
			}

			n := avp.N
			decAVPs[(uint64(vnd.I)<<32)|uint64(avp.I)] =
				func(a diameter.AVP) (string, any, error) {
					v, e := decf(&a)
					return n, v, e
				}
		}
	}

	return xd, nil
}
