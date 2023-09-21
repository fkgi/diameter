package main

import (
	"log"
	"testing"
)

func TestDictionary(t *testing.T) {
	if e := loadDictionary("dictionary.json"); e != nil {
		t.Fatal(e)
	}

	var k string = "Auth-Session-State"
	var v any = "NO_STATE_MAINTAINED"

	a, e := encAVPs[k](v)
	if e != nil {
		t.Fatal(e)
	}
	log.Println(a)

	k2, v2, e := decAVPs[a.VendorID][a.Code](a)
	if e != nil {
		t.Fatal(e)
	}
	if k != k2 {
		t.Fatal("k not match", k2)
	}
	if v != v2 {
		t.Fatal("k not match", v2)
	}

}
