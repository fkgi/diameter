package main

import (
	"bytes"
	"flag"
	"log"
	"os"

	"github.com/fkgi/diameter"
	"github.com/fkgi/diameter/dictionary"
)

func main() {
	log.Println("load")
	rt := flag.String("r", "route.xml", "route file `path`.")
	flag.Parse()

	if data, err := os.ReadFile(*rt); err != nil {
		log.Fatalln("[ERROR]", "failed to open route file:", err)
	} else if err = loadRoute(data); err != nil {
		log.Fatalln("[ERROR]", "failed to read route file:", err)
	}
	if data, err := os.ReadFile("../dictionary/s6a.xml"); err != nil {
		log.Fatalln("[ERROR]", "failed to open dictionary file:", err)
	} else if _, err = dictionary.LoadDictionary(data); err != nil {
		log.Fatalln("[ERROR]", "failed to read route file:", err)
	}

	m := diameter.Message{
		FlgR:  true,
		Code:  316,
		AppID: 16777251}
	buf := new(bytes.Buffer)
	diameter.SetDestinationHost(diameter.Identity("desthost.localdomain")).MarshalTo(buf)
	m.AVPs = buf.Bytes()

	log.Println("test")
	log.Println(getDestination(m))
	log.Println(getDestination(m))
	log.Println(getDestination(m))
	log.Println(getDestination(m))
	log.Println(getDestination(m))
	log.Println(getDestination(m))
	log.Println(getDestination(m))
	log.Println(getDestination(m))
	log.Println(getDestination(m))
	log.Println(getDestination(m))
}
