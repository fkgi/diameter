package dictionary

import (
	"bytes"
	"encoding/hex"
	"errors"
	"net"
	"time"

	"github.com/fkgi/diameter"
)

func decOctetString(avp *diameter.AVP) (any, error) {
	d := new([]byte)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return hex.EncodeToString(*d), nil
}

func decInteger32(avp *diameter.AVP) (any, error) {
	d := new(int32)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return *d, nil
}

func decInteger64(avp *diameter.AVP) (any, error) {
	d := new(int64)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return *d, nil
}

func decUnsigned32(avp *diameter.AVP) (any, error) {
	d := new(uint32)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return *d, nil
}

func decUnsigned64(avp *diameter.AVP) (any, error) {
	d := new(uint64)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return *d, nil
}

func decFloat32(avp *diameter.AVP) (any, error) {
	d := new(float32)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return *d, nil
}

func decFloat64(avp *diameter.AVP) (any, error) {
	d := new(float64)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return *d, nil
}

func decGrouped(avp *diameter.AVP) (any, error) {
	result := make(map[string]any)
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
		result[n] = v
	}
	return result, nil
}

func decAddress(avp *diameter.AVP) (any, error) {
	d := new(net.IP)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return d.String(), nil
}

func decTime(avp *diameter.AVP) (any, error) {
	d := new(time.Time)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return d.Format(time.RFC3339), nil
}

func decUTF8String(avp *diameter.AVP) (any, error) {
	d := new(string)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return *d, nil
}

func decDiameterIdentity(avp *diameter.AVP) (any, error) {
	d := new(diameter.Identity)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return d.String(), nil
}

func decDiameterURI(avp *diameter.AVP) (any, error) {
	d := new(diameter.URI)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return d.String(), nil
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

func decIPFilterRule(avp *diameter.AVP) (any, error) {
	d := new(diameter.IPFilterRule)
	e := avp.Decode(d)
	if e != nil {
		return nil, e
	}
	return *d, nil
}
