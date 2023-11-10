// Copyright 2023 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package main

import (
	"net/http"
)

var headerIndex = map[string][]struct{ Header, Value string }{

	// Ray's dev
	"/api.ray.blues.tools/udp": {
		{Header: "udp_ipv4", Value: "44.209.181.127"},
		{Header: "udp_port", Value: "8087"},
	},

	// Staging
	"/api.staging.blues.tools/udp": {
		{Header: "udp_ipv4", Value: "44.209.181.127"},
		{Header: "udp_port", Value: "8088"},
	},

	// Production
	"/api.notefile.net/udp": {
		{Header: "udp_ipv4", Value: "44.209.181.127"},
		{Header: "udp_port", Value: "8089"},
	},
}

// Lookup the proxy for a given server
func httpProxyLookupHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" || r.URL.Path == "/favicon.ico" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	headers, present := headerIndex[r.URL.Path]
	if !present {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for _, v := range headers {
		w.Header().Set(v.Header, v.Value)
	}

	w.WriteHeader(http.StatusOK)

}
