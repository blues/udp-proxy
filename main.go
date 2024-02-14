// Copyright 2024 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// Main service entry point
func main() {

	// Register endpoint for udp-proxy.net lookups, which are
	// performed when notehub starts up so that it knows
	// what IP:PORT to issue to devices.
	http.HandleFunc("/", httpProxyLookupHandler)

	// Register the UDP proxy handlers, which are used by devices
	// to send messages to the notehub - proxied to HTTP by us here.
	go udpProxyHandlers()

	// Register AWS health check endpoint
	http.HandleFunc("/ping", httpPingHandler)
	go func() {
		if err := http.ListenAndServe(":80", nil); err != nil {
			fmt.Println("Error starting HTTP listener: ", err)
			os.Exit(41)
		}
	}()

	// Spawn our signal handler
	go signalHandler()

	// Handle console input so we can manually quit and relaunch
	inputHandler()

}

// Ping handler, for AWS health checks
func httpPingHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte(time.Now().UTC().Format("2006-01-02T15:04:05Z")))
	fmt.Println("Ping from ", r.RemoteAddr)
}

func inputHandler() {

	scanner := bufio.NewScanner(os.Stdin)

	// Loops indefinitely, waiting for input unless Stdin gets an error or EOF
	for scanner.Scan() {
		message := scanner.Text()

		args := strings.Split(message, " ")

		switch args[0] {
		case "q":
			os.Exit(0)
		case "":
			// just re-prompt
		default:
			fmt.Printf("Unrecognized: '%s'\n", message)
		}

		fmt.Print("\n> ")

	}

}

// Our app's signal handler
func signalHandler() {
	ch := make(chan os.Signal, 100)
	signal.Notify(ch, syscall.SIGTERM)
	signal.Notify(ch, syscall.SIGINT)
	signal.Notify(ch, syscall.SIGSEGV)
	for {
		switch <-ch {
		case syscall.SIGINT:
			fmt.Printf("*** Exiting because of SIGNAL \n")
			os.Exit(0)
		case syscall.SIGTERM:
			return
		}
	}
}
