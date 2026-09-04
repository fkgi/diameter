package dictionary

import (
	_ "embed"
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/fkgi/diameter"
)

//go:embed base.xml
var baseXML []byte

func init() {
	if e := AppendDictionary(baseXML); e != nil {
		panic(e)
	}
}

type AppInfo struct {
	Name    string
	Command map[uint32]string
}

type AvpInfo struct {
	Name string
	Type string
	M    bool
	P    bool
	R    bool
	Enum map[int32]string
}

type VendorInfo struct {
	Name        string
	Application map[uint32]AppInfo
	AVP         map[uint32]AvpInfo
}

var (
	encAVPs    = map[string]func(any) (diameter.AVP, error){}
	decAVPs    = map[uint64]func(diameter.AVP) (string, any, error){}
	encCommand = map[string]uint64{}
	decCommand = map[uint64]string{}

	Data = map[uint32]VendorInfo{} // Data of current dictionary data
)

// GetApplicationName returns the application name associated with an ID.
func GetApplicationName(id uint32) string {
	for _, v := range Data {
		for i, app := range v.Application {
			if id == i {
				return app.Name
			}
		}
	}
	return fmt.Sprintf("UNKNOWN(%d)", id)
}

// AppendDictionary adds vendor, application, command, and AVP definitions from XML data.
func AppendDictionary(data []byte) error {
	xd := struct {
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
	}{}
	if e := xml.Unmarshal(data, &xd); e != nil {
		return e
	}

	for _, v := range xd.V {
		vnd := Data[v.I]
		if vnd.Application == nil {
			vnd.Application = map[uint32]AppInfo{}
		}
		if vnd.AVP == nil {
			vnd.AVP = map[uint32]AvpInfo{}
		}
		vnd.Name = v.N

		for _, p := range v.P {
			app := vnd.Application[p.I]
			if app.Command == nil {
				app.Command = map[uint32]string{}
			}
			app.Name = p.N
			for _, c := range p.C {
				app.Command[c.I] = c.N
			}
			vnd.Application[p.I] = app
		}

		for _, a := range v.V {
			avp := vnd.AVP[a.I]
			avp.Name = a.N
			avp.Type = a.T
			avp.M = a.M
			avp.P = a.P
			avp.R = a.R
			if len(a.E) != 0 && avp.Enum == nil {
				avp.Enum = map[int32]string{}
			}
			for _, e := range a.E {
				avp.Enum[e.I] = e.V
			}
			vnd.AVP[a.I] = avp
		}
		Data[v.I] = vnd
	}

	return nil
}

// LoadDictionary builds the encoders and decoders from the current dictionary data.
func LoadDictionary() error {
	encAVPs = map[string]func(any) (diameter.AVP, error){}
	decAVPs = map[uint64]func(diameter.AVP) (string, any, error){}
	encCommand = map[string]uint64{}
	decCommand = map[uint64]string{}

	for vid, vnd := range Data {
		for aid, app := range vnd.Application {
			for cid, cmd := range app.Command {
				p := vnd.Name + "/" + app.Name + "/" + cmd
				i := (uint64(aid) << 32) | uint64(cid)
				encCommand[p] = i
				decCommand[i] = p
			}
		}

		for aid, avp := range vnd.AVP {
			if _, ok := encAVPs[avp.Name]; ok {
				return fmt.Errorf("duplicated AVP definition: %s", avp.Name)
			}
			if _, ok := decAVPs[(uint64(vid)<<32)|uint64(aid)]; ok {
				return fmt.Errorf("duplicated AVP definition: %s", avp.Name)
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
				m1 := make(map[string]int32, len(avp.Enum))
				m2 := make(map[int32]string, len(avp.Enum))
				for ei, ev := range avp.Enum {
					if _, ok := m1[ev]; ok {
						return fmt.Errorf("duplicate ENUM entry in %s", avp.Name)
					}
					m1[ev] = ei

					if _, ok := m2[ei]; ok {
						return fmt.Errorf("duplicate ENUM entry in %s", avp.Name)
					}
					m2[ei] = ev
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
				return errors.New("invalid AVP type: " + avp.Name)
			}

			c := uint32(aid)
			v := uint32(vid)
			m := avp.M
			p := avp.P
			r := avp.R
			encAVPs[avp.Name] = func(data any) (diameter.AVP, error) {
				a := diameter.AVP{
					Code:      c,
					VendorID:  v,
					Mandatory: m,
					Protected: p,
					Reserved:  [5]bool{r, r, r, r, r}}
				e := encf(data, &a)
				return a, e
			}

			n := avp.Name
			decAVPs[(uint64(vid)<<32)|uint64(aid)] =
				func(a diameter.AVP) (string, any, error) {
					data, e := decf(&a)
					return n, data, e
				}
		}
	}
	return nil
}
