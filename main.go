// Copyright 2022 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package main

import (
	"net/http"
	"time"
)

// Main service entry point
func main() {

	// Register endpoint for udp-proxy.net lookups
	http.HandleFunc("/", httpProxyLookupHandler)

	// Register stock endpoints
	http.HandleFunc("/github", httpGithubHandler) // Re-launch upon github push of repo
	http.HandleFunc("/ping", httpPingHandler)     // AWS health check
	go http.ListenAndServe(":80", nil)

	// Spawn the console input handler
	go inputHandler()

	// Perform periodic housekeeping, if any
	for {
		time.Sleep(1 * time.Minute)
	}

}

// Ping handler, for AWS health checks
func httpPingHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(time.Now().UTC().Format("2006-01-02T15:04:05Z")))
}
