package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/fkgi/diameter"
	"github.com/fkgi/diameter/connector"
	"github.com/fkgi/diameter/dictionary"
	"github.com/fkgi/diameter/sctp"
)

const apiPath = "/diamsg/v1/"

var (
	dicData dictionary.XDictionary
)

func main() {
	log.Printf("[INFO] booting Round-Robin debugger for Diameter <%s REV.%d>...",
		diameter.ProductName, diameter.FirmwareRev)

	if v := os.Getenv("VERBOSE"); v != "yes" {
		if v != "no" {
			log.Println("[INFO]", "parameter VERBOSE is empty or invalid, set to default")
		}
		diameter.TraceEvent = func(old, new, event string, err error) {
			if err != nil {
				log.Printf("[INFO] event %s handling failed: %v", event, err)
			}
		}
		diameter.TraceMessage = func(msg diameter.Message, dct diameter.Direction, err error) {
			count(msg, dct, err)
		}
	}

	dict := os.Getenv("DICTIONARY")
	if dict == "" {
		dict = "dictionary.xml"
	}
	log.Println("[INFO]", "loading dictionary file", dict)
	if data, err := os.ReadFile(dict); err != nil {
		log.Fatalln("[ERROR]", "failed to open dictionary file:", err)
	} else if dicData, err = dictionary.LoadDictionary(data); err != nil {
		log.Fatalln("[ERROR]", "failed to read dictionary file:", err)
	} else {
		for _, vnd := range dicData.V {
			buf := new(strings.Builder)
			fmt.Fprintf(buf, "supported vendor: %s(%d)", vnd.N, vnd.I)
			for _, app := range vnd.P {
				fmt.Fprintf(buf, "\n | application: %s(%d)\n | | command:",
					app.N, app.I)
				for _, cmd := range app.C {
					fmt.Fprintf(buf, " %s(%d),", cmd.N, cmd.I)
				}
			}
			fmt.Fprint(buf, "\n | AVP:")
			for _, avp := range vnd.V {
				fmt.Fprintf(buf, " %s(%d/%s),", avp.N, avp.I, avp.T)
			}
			log.Println("[INFO]", buf)
		}
	}

	if to := os.Getenv("TIMEOUT"); to == "" {
	} else if t, e := strconv.Atoi(to); e != nil {
		log.Printf("[INFO] parameter TIMEOUT is invalid, set to default %fs",
			diameter.WDInterval.Seconds())
	} else {
		diameter.WDInterval = time.Second * time.Duration(t)
	}

	backend := "http://" + os.Getenv("BACKENDAPI_ADDR")
	if u, e := url.Parse(backend); e != nil || u.Host == "" {
		log.Println("[WARN]", "invalid HTTP backend host, Rx request will be rejected")
		dicData.RegisterHandler(
			func(path string, hdr http.Header, body io.Reader) (*http.Response, error) {
				return nil, fmt.Errorf("no HTTP backend is defined")
			},
			apiPath, selectCon)
	} else {
		log.Println("[INFO]", "HTTP backend:", backend)
		t, _ := http.DefaultTransport.(*http.Transport)
		dt := t.Clone()
		dt.MaxIdleConns = 0
		dt.MaxIdleConnsPerHost = 1000
		client := http.Client{Transport: dt, Timeout: diameter.WDInterval}
		defer client.CloseIdleConnections()

		dicData.RegisterHandler(
			func(path string, hdr http.Header, body io.Reader) (*http.Response, error) {
				req, _ := http.NewRequest("POST", backend+path, body)
				for k, l := range hdr {
					for _, v := range l {
						req.Header.Add(k, v)
					}
				}
				req.Header.Set("Content-Type", "application/json")
				return client.Do(req)
			},
			apiPath, selectCon)
	}

	http.HandleFunc("GET /diastate/v1/connection", conStateHandler)
	http.HandleFunc("GET /diastate/v1/statistics", statsHandler)

	frontend := os.Getenv("LOCALAPI_ADDR")
	log.Println("[INFO]", "listening HTTP...\n | local port:", frontend)
	go func() {
		err := http.ListenAndServe(frontend, nil)
		if err != nil {
			log.Println("[WARN]", "failed to listen HTTP, Tx request is not available:", err)
		}
	}()

	log.Println("[INFO]", "connecting Diameter...")
	var e error
	dlocal := os.Getenv("LOCAL_HOSTPORT")
	if dlocal != "" {
	} else if dlocal, e = os.Hostname(); e != nil {
		log.Fatalln("[ERROR]", "failed to detect localhost name:", e)
	}

	var lips []net.IP
	var lport int
	var scheme string
	if scheme, diameter.Host, diameter.Realm, lips, lport, e =
		connector.ResolveIdentity(dlocal); e != nil {
		log.Fatalln("[ERROR]", "invalid local identity:", e)
	}
	log.Printf("[INFO] Diameter local information"+
		"\n | address:    %s://%v:%d"+
		"\n | realm/host: %s/%s",
		scheme, lips, lport, diameter.Realm, diameter.Host)

	dpeer := []string{}
	for i := range 10 {
		if a := os.Getenv(fmt.Sprintf("PEER_HOSTPORT%d", i)); a == "" {
			continue
		} else if _, _, _, _, _, e = connector.ResolveIdentity(a); e != nil {
			log.Fatalln("[ERROR]", "invalid peer identity of", a, ":", e)
		} else {
			dpeer = append(dpeer, a)
		}
	}

	if len(dpeer) == 0 {
		log.Println("[INFO]", "accepting transport connection")
		var l net.Listener
		var e error
		switch scheme {
		case "sctp":
			l, e = sctp.ListenSCTP(&sctp.SCTPAddr{IP: lips, Port: lport})
		case "tcp":
			l, e = net.ListenTCP("tcp", &net.TCPAddr{IP: lips[0], Port: lport})
		}
		if e != nil {
			log.Fatalln("[ERROR]", "failed to listen transport interface", e)
		}
		go func() {
			wait()
			l.Close()
		}()

		c, e := l.Accept()
		for ; e == nil; c, e = l.Accept() {
			log.Printf("[INFO] transport connection up\n| peer : %s://%s",
				c.RemoteAddr().Network(), c.RemoteAddr().String())
			go func(c net.Conn) {
				con := &diameter.Connection{}
				appendCon(c, con, con.ListenAndServe)
			}(c)
		}
		log.Println("[WARN]", "transport listener closed:", e)
	} else {
		switch scheme {
		case "sctp":
			d, e := sctp.NewDaler(&sctp.SCTPAddr{IP: lips, Port: lport})
			if e != nil {
				log.Fatalln("[ERROR]", "failed to bind local SCTP port:", e)
			}
			for _, p := range dpeer {
				_, host, realm, pips, pport, _ := connector.ResolveIdentity(p)
				a := sctp.SCTPAddr{IP: pips, Port: pport}
				go dial(
					func() (net.Conn, error) { return d.Dial(&a) },
					&a, host, realm)
			}
		case "tcp":
			if len(dpeer) != 1 && lport != 0 {
				log.Fatalln("[ERROR]", "local port must 0 for multiple TCP connection")
			}
			l := net.TCPAddr{IP: lips[0], Port: lport}
			for _, p := range dpeer {
				_, host, realm, pips, pport, _ := connector.ResolveIdentity(p)
				a := net.TCPAddr{IP: pips[0], Port: pport}
				go dial(
					func() (net.Conn, error) { return net.DialTCP("tcp", &l, &a) },
					&a, host, realm)
			}
		}
		wait()
	}

	time.AfterFunc(time.Second*30, func() {
		log.Fatalln("[ERROR]", "closing timeout, forcefully stopped")
	})
	closeCon()
	log.Println("[INFO]", "server stopped")
}

func wait() {
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sigc
	log.Println("[INFO]", "interrupted, closing connections")
}

func dial(f func() (net.Conn, error), a net.Addr, host, realm diameter.Identity) {
	for {
		log.Printf("[INFO] connecting transport connection to %s//%s",
			a.Network(), a.String())
		if c, e := f(); e != nil {
			log.Printf(
				"[WARN] failed to connect transport connection to %s//%s",
				a.Network(), a.String())
		} else {
			buf := new(strings.Builder)
			fmt.Fprint(buf, "transport connection up")
			fmt.Fprintf(buf, "\n| local: %s://%s",
				c.LocalAddr().Network(), c.LocalAddr().String())
			fmt.Fprintf(buf, "\n| peer : %s://%s",
				c.RemoteAddr().Network(), c.RemoteAddr().String())
			log.Println("[INFO]", buf)

			con := &diameter.Connection{Host: host, Realm: realm}
			appendCon(c, con, con.DialAndServe)
		}

		select {
		case <-block:
			return
		case <-time.After(time.Second * 30):
		}
	}
}
