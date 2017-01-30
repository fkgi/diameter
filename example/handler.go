package example

import (
	"bytes"
	"flag"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/fkgi/diameter/msg"
)

// Handler is message handler
type Handler struct {
	sock      *string
	OrigHost  msg.DiameterIdentity
	OrigRealm msg.DiameterIdentity
	DestHost  msg.DiameterIdentity
	DestRealm msg.DiameterIdentity
	SessionID string
}

func genHostRealm(fqdn string) (host, realm msg.DiameterIdentity) {
	fqdn = strings.TrimSpace(fqdn)
	var e error
	host, e = msg.ParseDiameterIdentity(fqdn)
	if e != nil {
		log.Fatalln("invalid host name:", e)
	}

	realm, e = msg.ParseDiameterIdentity(fqdn[strings.Index(fqdn, ".")+1:])
	if e != nil {
		log.Fatalln("invalid host realm:", e)
	}
	return
}

// Init initialize each parameter
func (h *Handler) Init(sockname, ohost, dhost string) {
	log.Println("initializing...")

	// get option flag
	h.sock = flag.String("s", sockname, "UNIX socket name")
	flag.Parse()

	// create path
	if wdir, e := os.Getwd(); e != nil {
		log.Fatalln(e)
	} else {
		*h.sock = wdir + string(os.PathSeparator) + *h.sock
	}

	// generate host/realm
	h.OrigHost, h.OrigRealm = genHostRealm(ohost)
	h.DestHost, h.DestRealm = genHostRealm(dhost)

	// generate session id
	rand.Seed(time.Now().Unix())
	h.SessionID = string(ohost) + ";"
	h.SessionID += strconv.FormatInt(time.Now().Unix()+2208988800, 10) + ";"
	h.SessionID += strconv.FormatInt(int64(rand.Uint32()), 10) + ";"
	h.SessionID += "0"
}

// Push send message
func (h *Handler) Push(f func() *msg.Message) {
	// open UNIX socket
	c, e := net.Dial("unix", *h.sock)
	if e != nil {
		log.Fatalln(e)
	}
	defer c.Close()

	// create message
	log.Println("generating...")
	m := f()

	// handle message
	buf := new(bytes.Buffer)
	m.PrintStack(buf)
	log.Println("send message\n" + buf.String())
	if _, e = m.WriteTo(c); e != nil {
		log.Fatalln(e)
	}

	log.Println("waiting message...")
	if _, e = m.ReadFrom(c); e != nil {
		log.Fatalln(e)
	}
	buf = new(bytes.Buffer)
	m.PrintStack(buf)
	log.Println("receive message\n" + buf.String())
}

// Pull get message
func (h *Handler) Pull(f func(*msg.Message) *msg.Message) {
	// open UNIX socket
	c, e := net.Dial("unix", *h.sock)
	if e != nil {
		log.Fatalln(e)
	}
	defer c.Close()

	// create message
	m := &msg.Message{}

	// handle message
	log.Println("waiting message...")
	if _, e = m.ReadFrom(c); e != nil {
		log.Fatalln(e)
	}
	buf := new(bytes.Buffer)
	m.PrintStack(buf)
	log.Println("receive message\n" + buf.String())

	if avp, e := m.Decode(); e != nil {
		log.Fatalln(e)
	} else {
		for _, a := range avp {
			if a.Code == 263 {
				a.Decode(&h.SessionID)
				break
			}
		}
	}
	m = f(m)

	buf = new(bytes.Buffer)
	m.PrintStack(buf)
	log.Println("send message\n" + buf.String())
	if _, e = m.WriteTo(c); e != nil {
		log.Fatalln(e)
	}

	log.Println("end")
}
