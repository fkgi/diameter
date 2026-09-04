package dictionary

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/fkgi/diameter"
)

// DecodeMessage returns the dictionary command name for a Diameter message.
func DecodeMessage(m diameter.Message) (string, error) {
	if name, ok := decCommand[(uint64(m.AppID)<<32)|uint64(m.Code)]; ok {
		return name, nil
	}
	return fmt.Sprintf("UNKNOWN(%d)", m.Code), nil
}

// DecodeAVPs decodes Diameter AVPs into a map keyed by dictionary AVP names.
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

// DecodeAVP decodes a Diameter AVP into its dictionary name and value.
func DecodeAVP(a diameter.AVP) (string, any, error) {
	if f, ok := decAVPs[(uint64(a.VendorID)<<32)|uint64(a.Code)]; ok {
		return f(a)
	}
	return fmt.Sprintf("UNKNOWN(%d)", a.Code), hex.EncodeToString(a.Data), nil
}

func decOctetString(avp *diameter.AVP) (any, error) {
	d := new([]byte)
	e := avp.Decode(d)
	return hex.EncodeToString(*d), e
}

func decInteger32(avp *diameter.AVP) (any, error) {
	d := new(int32)
	e := avp.Decode(d)
	return *d, e
}

func decInteger64(avp *diameter.AVP) (any, error) {
	d := new(int64)
	e := avp.Decode(d)
	return *d, e
}

func decUnsigned32(avp *diameter.AVP) (any, error) {
	d := new(uint32)
	e := avp.Decode(d)
	return *d, e
}

func decUnsigned64(avp *diameter.AVP) (any, error) {
	d := new(uint64)
	e := avp.Decode(d)
	return *d, e
}

func decFloat32(avp *diameter.AVP) (any, error) {
	d := new(float32)
	e := avp.Decode(d)
	return *d, e
}

func decFloat64(avp *diameter.AVP) (any, error) {
	d := new(float64)
	e := avp.Decode(d)
	return *d, e
}

func decGrouped(avp *diameter.AVP) (any, error) {
	result := make(map[string][]any)
	for buf := bytes.NewBuffer(avp.Data); buf.Len() != 0; {
		a := diameter.AVP{}
		e := a.UnmarshalFrom(buf)
		if e != nil {
			return nil, e
		}
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

func decAddress(avp *diameter.AVP) (any, error) {
	d := new(net.IP)
	e := avp.Decode(d)
	return d.String(), e
}

func decTime(avp *diameter.AVP) (any, error) {
	d := new(time.Time)
	e := avp.Decode(d)
	return d.Format(time.RFC3339), e
}

func decUTF8String(avp *diameter.AVP) (any, error) {
	d := new(string)
	e := avp.Decode(d)
	return *d, e
}

func decDiameterIdentity(avp *diameter.AVP) (any, error) {
	d := new(diameter.Identity)
	e := avp.Decode(d)
	return d.String(), e
}

func decDiameterURI(avp *diameter.AVP) (any, error) {
	d := new(diameter.URI)
	e := avp.Decode(d)
	return d.String(), e
}

func decEnumerated(avp *diameter.AVP, enum map[int32]string) (any, error) {
	d := new(diameter.Enumerated)
	if e := avp.Decode(d); e != nil {
		return nil, e
	}
	if a, ok := enum[int32(*d)]; ok {
		return a, nil
	}
	return nil, errors.New("not defined Enumerated")
}

func decIPFilterRule(avp *diameter.AVP) (any, error) {
	d := new(diameter.IPFilterRule)
	e := avp.Decode(d)
	return *d, e
}
