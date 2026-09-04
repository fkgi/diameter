package dictionary

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"slices"
	"time"

	"github.com/fkgi/diameter"
)

// EncodeMessage creates a Diameter message from a dictionary command name.
func EncodeMessage(name string) (m diameter.Message, e error) {
	if id, ok := encCommand[name]; ok {
		m.AppID = uint32((id & 0xffffffff00000000) >> 32)
		m.Code = uint32(id & 0x00000000ffffffff)
		m.FlgR = true
	} else {
		e = errors.New("unknown command name")
	}
	return
}

// EncodeAVPs encodes a map of dictionary AVP names and values into Diameter AVPs.
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

// EncodeAVP encodes a value using the dictionary definition for an AVP name.
func EncodeAVP(name string, value any) (diameter.AVP, error) {
	if f, ok := encAVPs[name]; ok {
		return f(value)
	}
	return diameter.AVP{}, errors.New("unknown AVP name")
}

func encOctetString(v any, avp *diameter.AVP) error {
	if s, ok := v.(string); !ok {
		return errors.New("not String")
	} else if a, e := hex.DecodeString(s); e != nil {
		return e
	} else {
		return avp.Encode(a)
	}
}

func encInteger32(v any, avp *diameter.AVP) error {
	if d, ok := v.(float64); ok {
		return avp.Encode(int32(d))
	}
	return errors.New("not Number")
}

func encInteger64(v any, avp *diameter.AVP) error {
	if d, ok := v.(float64); ok {
		return avp.Encode(int64(d))
	}
	return errors.New("not Number")
}

func encUnsigned32(v any, avp *diameter.AVP) error {
	if d, ok := v.(float64); ok {
		return avp.Encode(uint32(d))
	}
	return errors.New("not Number")
}

func encUnsigned64(v any, avp *diameter.AVP) error {
	if d, ok := v.(float64); ok {
		return avp.Encode(uint64(d))
	}
	return errors.New("not Number")
}

func encFloat32(v any, avp *diameter.AVP) error {
	if d, ok := v.(float64); ok {
		return avp.Encode(float32(d))
	}
	return errors.New("not Number")
}

func encFloat64(v any, avp *diameter.AVP) error {
	if d, ok := v.(float64); ok {
		return avp.Encode(d)
	}
	return errors.New("not Number")
}

func encGrouped(v any, avp *diameter.AVP) (e error) {
	a, ok := v.(map[string]any)
	if !ok {
		return errors.New("not Grouped")
	}

	avps := map[uint32][]diameter.AVP{}
	codes := make([]uint32, 0, 20)
	for k, v := range a {
		if l, ok := v.([]any); ok {
			for _, v := range l {
				a, e := EncodeAVP(k, v)
				if e != nil {
					return fmt.Errorf("%s is invalid: %v", k, e)
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
				return fmt.Errorf("%s is invalid: %v", k, e)
			}
			avps[a.Code] = []diameter.AVP{a}
			codes = append(codes, a.Code)
		}
	}
	slices.Sort(codes)

	buf := new(bytes.Buffer)
	for _, k := range codes {
		if l, ok := avps[k]; ok {
			for _, a := range l {
				a.MarshalTo(buf)
			}
		}
	}
	avp.Data = buf.Bytes()
	return
}

func encAddress(v any, avp *diameter.AVP) error {
	if s, ok := v.(string); !ok {
		return errors.New("not String")
	} else if a := net.ParseIP(s); a == nil {
		return errors.New("not Address")
	} else {
		return avp.Encode(a)
	}
}

func encTime(v any, avp *diameter.AVP) error {
	if s, ok := v.(string); !ok {
		return errors.New("not String")
	} else if a, e := time.Parse(time.RFC3339, s); e != nil {
		return e
	} else {
		return avp.Encode(a)
	}
}

func encUTF8String(v any, avp *diameter.AVP) error {
	if s, ok := v.(string); ok {
		return avp.Encode(s)
	}
	return errors.New("not String")
}

func encDiameterIdentity(v any, avp *diameter.AVP) error {
	if s, ok := v.(string); !ok {
		return errors.New("not String")
	} else if a, e := diameter.ParseIdentity(s); e != nil {
		return e
	} else {
		return avp.Encode(a)
	}
}

func encDiameterURI(v any, avp *diameter.AVP) error {
	if s, ok := v.(string); !ok {
		return errors.New("not String")
	} else if a, e := diameter.ParseURI(s); e != nil {
		return e
	} else {
		return avp.Encode(a)
	}
}

func encEnumerated(v any, avp *diameter.AVP, enum map[string]int32) error {
	if s, ok := v.(string); !ok {
		return errors.New("not String")
	} else if a, ok := enum[s]; !ok {
		return errors.New("not defined Enumerated")
	} else {
		return avp.Encode(diameter.Enumerated(a))
	}
}

func encIPFilterRule(v any, avp *diameter.AVP) error {
	if s, ok := v.(diameter.IPFilterRule); ok {
		return avp.Encode(s)
	}
	return errors.New("not String")
}
