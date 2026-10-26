//go:build !linux || (linux && 386)

package main

import "net/http"

func sysmetricsHandler(w http.ResponseWriter) {
}
