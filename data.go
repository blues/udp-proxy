// Copyright 2024 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package main

var proxyData = map[string][]struct{ Key, Value string }{

	// Scott's dev
	"api.ray.blues.tools/udp": {
		{Key: "udp_ipv4", Value: "216.245.146.112"},
		{Key: "udp_port", Value: "8086"},
	},

	// Ray's dev
	"api.scott.blues.tools/udp": {
		{Key: "udp_ipv4", Value: "216.245.146.112"},
		{Key: "udp_port", Value: "8087"},
	},

	// Staging
	"api.staging.blues.tools/udp": {
		{Key: "udp_ipv4", Value: "216.245.146.112"},
		{Key: "udp_port", Value: "8088"},
	},

	// Production
	"api.notefile.net/udp": {
		{Key: "udp_ipv4", Value: "216.245.146.112"},
		{Key: "udp_port", Value: "8089"},
	},
}
