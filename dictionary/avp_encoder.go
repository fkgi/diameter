package dictionary

import (
	"bytes"
	"encoding/hex"
	"errors"
	"net"
	"time"

	"github.com/fkgi/diameter"
)

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
	d, ok := v.(float64)
	if !ok {
		return errors.New("not Number")
	}
	return avp.Encode(int32(d))
}

func encInteger64(v any, avp *diameter.AVP) error {
	d, ok := v.(float64)
	if !ok {
		return errors.New("not Number")
	}
	return avp.Encode(int64(d))
}

func encUnsigned32(v any, avp *diameter.AVP) error {
	d, ok := v.(float64)
	if !ok {
		return errors.New("not Number")
	}
	return avp.Encode(uint32(d))
}

func encUnsigned64(v any, avp *diameter.AVP) error {
	d, ok := v.(float64)
	if !ok {
		return errors.New("not Number")
	}
	return avp.Encode(uint64(d))
}

func encFloat32(v any, avp *diameter.AVP) error {
	d, ok := v.(float64)
	if !ok {
		return errors.New("not Number")
	}
	return avp.Encode(float32(d))
}

func encFloat64(v any, avp *diameter.AVP) error {
	d, ok := v.(float64)
	if !ok {
		return errors.New("not Number")
	}
	return avp.Encode(d)
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

func encUTF8String(v any, avp *diameter.AVP) error {
	s, ok := v.(string)
	if !ok {
		return errors.New("not String")
	}
	return avp.Encode(s)
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

func encIPFilterRule(v any, avp *diameter.AVP) error {
	s, ok := v.(diameter.IPFilterRule)
	if !ok {
		return errors.New("not String")
	}
	return avp.Encode(s)
}
