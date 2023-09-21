package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os/exec"
)

var path string

func main() {
	t := flag.String("t", ":8081", "http local host:port")
	s := flag.String("s", "./server.sh", "response generator scripth path")
	flag.Parse()

	path = *s
	http.ListenAndServe(*t, http.Handler(apiHandler))
}

var apiHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	log.Println(r.URL)

	b, e := io.ReadAll(r.Body)
	defer r.Body.Close()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Println(string(b))

	cmd := exec.Command(path, r.URL.Path)
	stdin, e := cmd.StdinPipe()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, e = stdin.Write(b)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	stdin.Close()

	b, e = cmd.Output()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(b)
})
