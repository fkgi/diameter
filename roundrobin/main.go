package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fkgi/diameter"
	"github.com/fkgi/diameter/connector"
	"github.com/fkgi/diameter/dictionary"
)

const apiPath = "/diamsg/v1/"

var (
	dicData dictionary.XDictionary
	con     diameter.Connection
)

func main() {
	host, err := os.Hostname()
	if err != nil {
		host = "roundrobin.internal"
	}
	dlocal := flag.String("l", host, "Diameter local host. `[realm/]hostname[:port]`")
	hlocal := flag.String("i", ":8080", "HTTP local interface address. `[host]:port`")
	hpeer := flag.String("b", "localhost", "HTTP backend host address. `host[:port]`")
	dict := flag.String("d", "dictionary.xml", "Diameter dictionary file `path`.")
	to := flag.Int("t", int(diameter.WDInterval/time.Second), "Message timeout timer [s]")
	verbose := flag.Bool("v", false, "Verbose log output")
	help := flag.Bool("h", false, "Print usage")
	flag.Parse()

	dpeer := flag.Arg(0)
	if *help || dpeer == "" {
		fmt.Printf("usage: %s [OPTION]... DIAMETER_PEER\n", os.Args[0])
		fmt.Println("DIAMETER_PEER format is [(tcp|sctp)://][realm/]hostname[:port]")
		fmt.Println()
		flag.PrintDefaults()
		return
	}

	log.Printf("[INFO] booting Round-Robin debugger for Diameter <%s REV.%d>...",
		diameter.ProductName, diameter.FirmwareRev)

	if !(*verbose) {
		diameter.TraceEvent = nil
		diameter.TraceMessage = nil
	}

	log.Println("[INFO]", "loading dictionary file", *dict)
	if data, err := os.ReadFile(*dict); err != nil {
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

	diameter.WDInterval = time.Duration(*to) * time.Second

	rxPath := "http://" + *hpeer
	_, err = url.Parse(rxPath)
	if err != nil {
		log.Println("[WARN]", "invalid HTTP backend host, Rx request will be rejected")
		rxPath = ""
	} else {
		log.Println("[INFO]", "HTTP backend:", rxPath)
	}

	var dt *http.Transport
	if t, ok := http.DefaultTransport.(*http.Transport); ok {
		dt = t.Clone()
	} else {
		dt = &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			ForceAttemptHTTP2:     false,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
	}
	dt.MaxIdleConns = 0
	dt.MaxIdleConnsPerHost = 1000
	client := http.Client{
		Transport: dt,
		Timeout:   diameter.WDInterval}
	defer client.CloseIdleConnections()

	dicData.RegisterHandler(
		func(path string, hdr http.Header, body io.Reader) (*http.Response, error) {
			if rxPath == "" {
				return nil, fmt.Errorf("no HTTP backend is defined")
			}

			req, _ := http.NewRequest("POST", rxPath+path, body)
			for k, l := range hdr {
				for _, v := range l {
					req.Header.Add(k, v)
				}
			}
			req.Header.Set("Content-Type", "application/json")
			return client.Do(req)
		},
		apiPath,
		func(diameter.Message) *diameter.Connection {
			return &con
		})

	http.HandleFunc("/diastate/v1/connection", conStateHandler)
	http.HandleFunc("/diastate/v1/statistics", statsHandler)
	log.Println("[INFO]", "listening HTTP...\n | local port:", *hlocal)
	go func() {
		err := http.ListenAndServe(*hlocal, nil)
		if err != nil {
			log.Println("[WARN]", "failed to listen HTTP, Tx request is not available:", err)
		}
	}()

	log.Println("[INFO]", "connecting Diameter...")
	var c net.Conn
	con = diameter.Connection{}
	c, con.Host, con.Realm, err = connector.Dial(*dlocal, dpeer)
	if err != nil {
		log.Fatalln("[ERROR]", err)
	}

	buf := new(strings.Builder)
	fmt.Fprint(buf, "transport connection up")
	fmt.Fprintf(buf, "\n| local: %s://%s", c.LocalAddr().Network(), c.LocalAddr().String())
	fmt.Fprintf(buf, "\n| peer : %s://%s", c.RemoteAddr().Network(), c.RemoteAddr().String())
	log.Println("[INFO]", buf)

	go func() {
		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
		<-sigc

		con.Close(diameter.Rebooting)
	}()

	log.Println("[INFO]", "closed, error=", con.DialAndServe(c))
}
