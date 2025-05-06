package main

import (
	"bytes"
	"math/rand"

	"github.com/fkgi/diameter"
)

func rxhandler(m diameter.Message) diameter.Message {
	var dHost diameter.Identity
	for rdr := bytes.NewReader(m.AVPs); rdr.Len() != 0; {
		a := diameter.AVP{}
		e := a.UnmarshalFrom(rdr)
		if e != nil {
			continue
		}
		if a.VendorID != 0 {
			continue
		}
		if a.Code == 293 {
			dHost, _ = diameter.GetDestinationHost(a)
		}
	}

	if dHost == diameter.Host {
		return m.GenerateAnswerBy(diameter.CommandUnspported)
	}

	dcon := []*diameter.Connection{}
	if m.PeerName == upLink {
		if dHost == upLink {
			return m.GenerateAnswerBy(diameter.UnableToDeliver)
		}
		for _, con := range refConnection() {
			if con.Host != dHost {
				continue
			}
			if len(con.AvailableApplications()) == 0 {
				dcon = append(dcon, con)
				continue
			}
			for _, i := range con.AvailableApplications() {
				if i == m.AppID {
					dcon = append(dcon, con)
					continue
				}
			}
		}
		if len(dcon) == 0 {
			for _, con := range refConnection() {
				if con.Host == upLink {
					continue
				}
				if len(con.AvailableApplications()) == 0 {
					dcon = append(dcon, con)
					continue
				}
				for _, i := range con.AvailableApplications() {
					if i == m.AppID {
						dcon = append(dcon, con)
						continue
					}
				}
			}
		}
	} else {
		for _, con := range refConnection() {
			if con.Host == upLink {
				dcon = append(dcon, con)
			}
		}
	}

	if len(dcon) == 0 {
		return m.GenerateAnswerBy(diameter.UnableToDeliver)
	}

	buf := bytes.NewBuffer(m.AVPs)
	diameter.SetRouteRecord(diameter.Host).MarshalTo(buf)
	m.AVPs = buf.Bytes()
	return dcon[rand.Intn(len(dcon))].DefaultTxHandler(m)
}
