package main

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strings"

	"github.com/fkgi/diameter"
	"github.com/fkgi/diameter/dictionary"
)

var (
	groups map[diameter.Identity][]diameter.Identity
	routes []route
)

type route struct {
	destination diameter.Identity
	condition   map[*regexp.Regexp][]string
}

func loadRoute(data []byte) (e error) {
	xr := struct {
		XMLName xml.Name `xml:"router"`
		Group   []struct {
			Name string   `xml:"name,attr"`
			Peer []string `xml:"peer"`
		} `xml:"group"`
		Route []struct {
			Dest      string `xml:"destination,attr"`
			Condition []struct {
				Param string `xml:"param,attr"`
				Text  string `xml:",chardata"`
			} `xml:"condition"`
		} `xml:"route"`
	}{}
	if e = xml.Unmarshal(data, &xr); e != nil {
		return errors.Join(
			errors.New("failed to unmarshal route file"), e)
	}

	groups = make(map[diameter.Identity][]diameter.Identity)
	for _, gr := range xr.Group {
		id, e := diameter.ParseIdentity(gr.Name)
		if e != nil {
			return errors.Join(
				errors.New("invalid group name"), e)
		}

		p := make([]diameter.Identity, len(gr.Peer))
		for i := range gr.Peer {
			p[i], e = diameter.ParseIdentity(gr.Peer[i])
			if e != nil {
				return errors.Join(
					errors.New("invalid peer name"), e)
			}
		}
		groups[id] = p
	}

	routes = []route{}
	for _, rt := range xr.Route {
		var id diameter.Identity
		id, e = diameter.ParseIdentity(rt.Dest)
		if e != nil {
			return errors.Join(
				errors.New("invalid destination"), e)
		}

		rm := make(map[*regexp.Regexp][]string)
		for _, c := range rt.Condition {
			var r *regexp.Regexp
			r, e = regexp.Compile(c.Text)
			if e != nil {
				return
			}
			rm[r] = strings.Split(c.Param, "/")
		}
		routes = append(routes, route{destination: id, condition: rm})
	}

	return
}

func getDestination(m diameter.Message) (peers []diameter.Identity) {
	mname, e := dictionary.DecodeMessage(m)
	if e != nil {
		mname = ""
	}

	avps := []diameter.AVP{}
	for rdr := bytes.NewReader(m.AVPs); rdr.Len() != 0; {
		a := diameter.AVP{}
		if e := a.UnmarshalFrom(rdr); e == nil {
			avps = append(avps, a)
		}
	}
	root, _ := dictionary.DecodeAVPs(avps)

	peers = []diameter.Identity{}
	for _, route := range routes {
		match := true
		for regex, path := range route.condition {
			if path[0] == "$command" && !regex.MatchString(mname) {
				match = false
				break
			}
			if !checkAVP(root, path, regex) {
				match = false
				break
			}
		}
		if match {
			if ids, ok := groups[route.destination]; ok {
				l := len(peers)
				peers = append(peers, ids...)
				rand.Shuffle(len(ids), func(i, j int) {
					ids[l+i], ids[l+j] = ids[l+j], ids[l+i]
				})
			} else {
				peers = append(peers, route.destination)
			}
		}
	}
	return
}

func checkAVP(avps map[string]any, path []string, reg *regexp.Regexp) bool {
	if a, ok := avps[path[0]]; !ok {
		return false
	} else if len(path) == 1 {
		return checkValue(a, reg)
	} else if gr, ok := a.(map[string]any); ok {
		return checkAVP(gr, path[1:], reg)
	}
	return false
}

func checkValue(a any, reg *regexp.Regexp) bool {
	switch v := a.(type) {
	case string:
		return reg.MatchString(v)
	case int32, int64, uint32, uint64, float32, float64:
		return reg.MatchString(fmt.Sprintf("%d", v))
	case []any:
		for _, a := range v {
			if checkValue(a, reg) {
				return true
			}
		}
		return false
	default:
		return false
	}
}
