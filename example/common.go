package example

import (
	"bytes"
	"flag"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/fkgi/diameter/msg"
)

var (
	osock *string

	VENDOR_ID      uint32
	APPLICATION_ID uint32

	ORIG_HOST  msg.DiameterIdentity
	ORIG_REALM msg.DiameterIdentity
	DEST_HOST  msg.DiameterIdentity
	DEST_REALM msg.DiameterIdentity

	SESSION_ID string
)

func SetOrigin(host, realm string) {
	var e error
	if ORIG_HOST, e = msg.ParseDiameterIdentity(host); e != nil {
		log.Fatalln(e)
	}
	if ORIG_REALM, e = msg.ParseDiameterIdentity(realm); e != nil {
		log.Fatalln(e)
	}
}

func SetDest(host, realm string) {
	var e error
	if DEST_HOST, e = msg.ParseDiameterIdentity(host); e != nil {
		log.Fatalln(e)
	}
	if DEST_REALM, e = msg.ParseDiameterIdentity(realm); e != nil {
		log.Fatalln(e)
	}
}

func SetAppID(ven, app uint32) {
	VENDOR_ID = ven
	APPLICATION_ID = app
}

func Init(sockname string) {
	log.SetFlags(log.Ltime | log.Lmicroseconds)

	// get option flag
	osock = flag.String("s", sockname, "output UNIX socket name")
	flag.Parse()

	// create path
	if wdir, e := os.Getwd(); e != nil {
		log.Fatalln(e)
	} else {
		*osock = wdir + string(os.PathSeparator) + *osock
	}

	rand.Seed(time.Now().Unix())
	SESSION_ID = string(ORIG_HOST) + ";"
	//  SESSION_ID += strconv.FormatInt(1234567890, 10) + ";"
	SESSION_ID += strconv.FormatInt(time.Now().Unix()+2208988800, 10) + ";"
	SESSION_ID += strconv.FormatInt(int64(rand.Uint32()), 10) + ";"
	SESSION_ID += "0"
}

func Push(f func() *msg.Message) {
	// open UNIX socket
	c, e := net.Dial("unix", *osock)
	if e != nil {
		log.Fatalln(e)
	}
	defer c.Close()

	// create message
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
	log.Println("end")
}

func Pull(f func(*msg.Message) *msg.Message) {
	// open UNIX socket
	c, e := net.Dial("unix", *osock)
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
				a.Decode(&SESSION_ID)
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
