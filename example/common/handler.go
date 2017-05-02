package common

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"net"
	"os"
	"time"

	"github.com/fkgi/diameter/msg"
)

// Handler is message handler
type Handler struct {
	sock *string
	msg.OriginHost
	msg.OriginRealm
	msg.DestinationHost
	msg.DestinationRealm
	msg.SessionID
}

// Init initialize each parameter
func (h *Handler) Init(sockname, ohost, dhost string) {
	Log.Println("initializing...")

	// get option flag
	h.sock = flag.String("s", sockname, "UNIX socket name")
	flag.Parse()

	// create path
	*h.sock = os.TempDir() + string(os.PathSeparator) + *h.sock

	// generate host/realm
	hst, rlm := genHostRealm(ohost)
	h.OriginHost = msg.OriginHost(hst)
	h.OriginRealm = msg.OriginRealm(rlm)
	hst, rlm = genHostRealm(dhost)
	h.DestinationHost = msg.DestinationHost(hst)
	h.DestinationRealm = msg.DestinationRealm(rlm)

	h.SessionID = msg.SessionID(fmt.Sprintf("%s;%d;%d;0",
		h.OriginHost, time.Now().Unix()+2208988800, rand.Uint32()))
}

// Push send message
func (h *Handler) Push(f func() *msg.Message) {
	// open UNIX socket
	c, e := net.Dial("unix", *h.sock)
	if e != nil {
		Log.Fatalln(e)
	}
	defer c.Close()

	// create message
	Log.Println("generating...")
	m := f()

	// handle message
	buf := new(bytes.Buffer)
	m.PrintStack(buf)
	Log.Println("send message\n" + buf.String())
	if _, e = m.WriteTo(c); e != nil {
		Log.Fatalln(e)
	}

	Log.Println("waiting message...")
	if _, e = m.ReadFrom(c); e != nil {
		Log.Fatalln(e)
	}
	buf = new(bytes.Buffer)
	m.PrintStack(buf)
	Log.Println("receive message\n" + buf.String())
}

// Pull get message
func (h *Handler) Pull(f func(*msg.Message) *msg.Message) {
	// open UNIX socket
	c, e := net.Dial("unix", *h.sock)
	if e != nil {
		Log.Fatalln(e)
	}
	defer c.Close()

	// create message
	m := &msg.Message{}

	// handle message
	Log.Println("waiting message...")
	if _, e = m.ReadFrom(c); e != nil {
		Log.Fatalln(e)
	}
	buf := new(bytes.Buffer)
	m.PrintStack(buf)
	Log.Println("receive message\n" + buf.String())

	if avp, e := m.Decode(); e != nil {
		Log.Fatalln(e)
	} else {
		h.SessionID, _ = msg.GetSessionID(avp)
	}
	m = f(m)

	buf = new(bytes.Buffer)
	m.PrintStack(buf)
	Log.Println("send message\n" + buf.String())
	if _, e = m.WriteTo(c); e != nil {
		Log.Fatalln(e)
	}

	Log.Println("end")
}
