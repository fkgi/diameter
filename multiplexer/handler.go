package multiplexer

import (
	"bytes"
	"math/rand"

	"github.com/fkgi/diameter"
)

var DefaultRouter diameter.Router = func(m diameter.Message) *diameter.Connection {
	var dHost diameter.Identity
	for rdr := bytes.NewReader(m.AVPs); rdr.Len() != 0; {
		a := diameter.AVP{}
		if err := a.UnmarshalFrom(rdr); err != nil {
			continue
		}
		if a.VendorID != 0 {
			continue
		}
		switch a.Code {
		case 293:
			dHost, _ = diameter.GetDestinationHost(a)
		}
	}
	for _, con := range cons {
		if con.Host == dHost {
			return con
		}
	}

	t := rand.Intn(len(cons))
	i := 0
	for _, con := range cons {
		if i == t {
			return con
		}
		i++
	}
	return nil
}

// Handle wrap diameter.Handle for multiple connection.
func Handle(code, appID, venID uint32, h diameter.Handler) diameter.Handler {
	return diameter.Handle(code, appID, venID, h, DefaultRouter)
}

// DefaultTxHandler wrap diameter.DefaultTxHandler for multiple connection
func DefaultTxHandler(m diameter.Message) diameter.Message {
	if con := DefaultRouter(m); con != nil {
		return con.DefaultTxHandler(m)
	}

	buf := new(bytes.Buffer)
	diameter.SetResultCode(diameter.UnableToDeliver).MarshalTo(buf)
	diameter.SetOriginHost(diameter.Host).MarshalTo(buf)
	diameter.SetOriginRealm(diameter.Realm).MarshalTo(buf)

	return diameter.Message{
		FlgR: false, FlgP: m.FlgP, FlgE: true, FlgT: false,
		Code: m.Code, AppID: m.AppID,
		HbHID: m.HbHID, EtEID: m.EtEID,
		AVPs: buf.Bytes()}
}
