//go:build linux && !386

package main

import (
	"fmt"
	"net/http"
	"syscall"
)

func sysmetricsHandler(w http.ResponseWriter) {
	var usage syscall.Rusage
	if e := syscall.Getrusage(syscall.RUSAGE_SELF, &usage); e == nil {
		fmt.Fprintln(w, "# HELP cpu_seconds_total Total user and system CPU time spent in seconds")
		fmt.Fprintln(w, "# TYPE cpu_seconds_total counter")
		fmt.Fprintf(w, "cpu_seconds_total %.4f\n",
			float64(usage.Utime.Sec)+float64(usage.Utime.Usec)/1e6+
				float64(usage.Stime.Sec)+float64(usage.Stime.Usec)/1e6)
	}
}
