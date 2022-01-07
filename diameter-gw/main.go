package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"syscall"

	dm "github.com/fkgi/diameter"
)

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.Printf("booting Diameter-gateway <%s REV.%d>...", dm.ProductName, dm.FirmwareRev)

	tmp, err := os.Hostname()
	if err != nil {
		tmp = "diameter-hub.internal"
	}
	lDiaAddr := flag.String("l", tmp+":3868", "diameter local host:port")
	pDiaAddr := flag.String("p", "", "diameter peer host:port")
	api := flag.String("m", ":8080", "http local host:port")
	sctp := flag.Bool("sctp", false, "flag for sctp")
	flag.Parse()

	log.Printf("local address = %s\n", *lDiaAddr)
	log.Printf("peer address = %s\n", *pDiaAddr)
	if *sctp {
		log.Println("transport = sctp")
	} else {
		log.Println("transport = tcp")
	}

	go func() {
		sigc := make(chan os.Signal, 2)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sigc

		dm.Close(dm.Rebooting)
	}()
	go func() {
		http.ListenAndServe(*api, http.Handler(apiHandler))
	}()

	if len(*pDiaAddr) == 0 {
		log.Println("listening...")
		err = dm.ListenAndServe(*lDiaAddr, *sctp)
	} else {
		log.Println("connecting...")
		err = dm.DialAndServe(*lDiaAddr, *pDiaAddr, *sctp)
	}

	log.Println("closed, error=", err)
}

var apiHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	if p == "" {
		p = "/"
	} else {
		if p[0] != '/' {
			p = "/" + p
		}
		np := path.Clean(p)
		if p[len(p)-1] == '/' && np != "/" {
			if len(p) == len(np)+1 && strings.HasPrefix(p, np) {
				np = p
			} else {
				np += "/"
			}
		}
		p = np
	}

	code, _ := strconv.ParseUint(strings.Split(p, "/")[3], 10, 32)
	appId, _ := strconv.ParseUint(strings.Split(p, "/")[2], 10, 32)
	vendorId, _ := strconv.ParseUint(strings.Split(p, "/")[1], 10, 32)

	d := message{}
	b, e := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if e = json.Unmarshal(b, &d); e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	avp := make([]dm.AVP, 0, 10)
	avp = append(avp, dm.SetVendorSpecAppID(uint32(vendorId), uint32(appId)))
	if d.SessionID == "" {
		avp = append(avp, dm.SetSessionID("hogehogesession"))
		avp = append(avp, dm.SetAuthSessionState(false))
	} else {
		avp = append(avp, dm.SetSessionID(d.SessionID))
		avp = append(avp, dm.SetAuthSessionState(true))
	}
	host, _ := dm.ParseIdentity(d.Host)
	avp = append(avp, dm.SetOriginHost(host))
	realm, _ := dm.ParseIdentity(d.Realm)
	avp = append(avp, dm.SetOriginRealm(realm))
	for _, a := range d.AVP {
		a2 := dm.AVP{
			Code:      a.Code,
			VendorID:  a.VendorID,
			Mandatory: a.Mandatory}
		switch a.Type {
		case "ip":
			str, ok := a.Data.(string)
			if ok {
				tmp := net.ParseIP(str)
				a2.Encode(tmp)
			}
		case "time":
		case "identity":
		case "uri":
		case "enumerated":
		case "string":
			tmp := a.Data.(string)
			a2.Encode(tmp)
		case "group":
		case "byte":
			tmp := a.Data.([]byte)
			a2.Encode(tmp)
		case "int32":
			var tmp int32 = int32(a.Data.(float64))
			a2.Encode(tmp)
		case "int64":
			var tmp int64 = int64(a.Data.(float64))
			a2.Encode(tmp)
		case "uint32":
			var tmp uint32 = uint32(a.Data.(float64))
			a2.Encode(tmp)
		case "uint64":
			var tmp uint64 = uint64(a.Data.(float64))
			a2.Encode(tmp)
		case "float32":
			var tmp float32 = float32(a.Data.(float64))
			a2.Encode(tmp)
		case "float64":
			var tmp float64 = float64(a.Data.(float64))
			a2.Encode(tmp)
		}
		//a2.Encode(a.Data)
		avp = append(avp, a2)
	}
	buf := new(bytes.Buffer)
	for _, a := range avp {
		a.MarshalTo(buf)
	}
	dm.Send(uint32(code), uint32(appId), 0, false, buf.Bytes())

	w.WriteHeader(http.StatusNoContent)
})

type message struct {
	Host      string `json:"host"`
	Realm     string `json:"realm"`
	SessionID string `json:"sessionID,omitempty"`
	// UTF8String []utf8stringAVP `json:"utf8string"`
	AVP []avp `json:"avp"`
}

type avp struct {
	Code      uint32      `json:"code"`
	VendorID  uint32      `json:"vendorID,omitempty"`
	Mandatory bool        `json:"mandatory,omitempty"`
	Type      string      `json:"type"`
	Data      interface{} `json:"data"`
}

/*
type utf8stringAVP struct {
	Code      uint32 `json:"code"`
	VendorID  uint32 `json:"vendorID,omitempty"`
	Mandatory bool   `json:"mandatory,omitempty"`
	Data      string `json:"data"`
}
*/

/*
POST http://host/appID/code
{
	"Host": "",
	"Realm": "",
	"SessionID": "",
	"UTF8String": {
		"Code": 1,
		"Vendor": 0,
		"Mandatory": true,
		"Data": "123"
	]
}

func handleSRR(req ts29338.SRR) ts29338.SRA {
	log.Printf("recieve message\n%s", req)

	txt := "searching subscriber\n"
	res := ts29338.SRA{
		ResultCode:  DiameterUnableToComply,
		OriginHost:  hName,
		OriginRealm: Realm}

	x := xmlConfig{}
	data, e := ioutil.ReadFile(*subscriberData)
	if e == nil {
		e = xml.Unmarshal([]byte(data), &x)
	}
	if e != nil {
		log.Println(txt +
			" | subscriber DB reading failed: " + e.Error())
		log.Printf("send answer message\n%s", res)
		return res
	}

	msisdn := req.MSISDN.String()
	for _, s := range x.Subs {
		s.MSISDN = strings.TrimSpace(s.MSISDN)
		if s.MSISDN != msisdn {
			continue
		}
		txt += " | subscriber " + msisdn + " found"

		res.ResultCode = DiameterSuccess
		res.ServingNode[0].NodeType = ts29338.NodeMME
		res.ServingNode[0].Address, e = teldata.ParseE164(s.Address)
		res.ServingNode[0].Host, e = ParseIdentity(s.Host)
		res.ServingNode[0].Realm, e = ParseIdentity(
			s.Host[strings.Index(s.Host, ".")+1:])

		res.IMSI, e = teldata.ParseIMSI(s.IMSI)
		log.Println(txt)
		log.Printf("send answer message\n%s", res)
		return res
	}

	log.Println(txt + " | subscriber " + msisdn + " not found")
	res.ResultCode = ts29338.DiameterErrorUserUnknown
	log.Printf("send answer message\n%s", res)
	return res
}

func handleTFR(req ts29338.TFR) ts29338.TFA {
	log.Printf("recieve message\n%s", req)

	txt := "searching subscriber\n"
	res := ts29338.TFA{
		ResultCode:  DiameterUnableToComply,
		OriginHost:  mName,
		OriginRealm: Realm}

	x := xmlConfig{}
	data, e := ioutil.ReadFile(*subscriberData)
	if e == nil {
		e = xml.Unmarshal([]byte(data), &x)
	}
	if e != nil {
		log.Println(txt + " | subscriber DB reading failed: " + e.Error())
		log.Printf("send answer message\n%s", res)
		return res
	}

	imsi := req.IMSI.String()
	for _, s := range x.Subs {
		s.IMSI = strings.TrimSpace(s.IMSI)
		if s.IMSI != imsi {
			continue
		}
		txt += " | subscriber " + imsi + " found"

		res.ResultCode = DiameterSuccess
		res.SMSPDU = sms.DeliverReport{}

		log.Println(txt)
		log.Printf("send answer message\n%s", res)
		return res
	}

	log.Println(txt + " | subscriber " + imsi + " not found")
	res.ResultCode = ts29338.DiameterErrorIlleagalUser
	log.Printf("send answer message\n%s", res)
	return res
}

*/
