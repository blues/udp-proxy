// Copyright 2022 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package main

import (
	"time"
)

// Main service entry point
func main() {

	// Spawn the console input handler
	go inputHandler()

	// Spawn the web request inbound server
	go HTTPInboundHandler(":80")

	// Housekeeping
	for {
		time.Sleep(1 * time.Minute)
	}

}
