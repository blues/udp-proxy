// Copyright 2023 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
)

const traceIo = true

var headerIndex = map[string][]struct{ Key, Value string }{

	// Scott's dev
	"api.ray.blues.tools/udp": {
		{Key: "udp_ipv4", Value: "44.209.181.127"},
		{Key: "udp_port", Value: "8086"},
	},

	// Ray's dev
	"api.scott.blues.tools/udp": {
		{Key: "udp_ipv4", Value: "44.209.181.127"},
		{Key: "udp_port", Value: "8087"},
	},

	// Staging
	"api.staging.blues.tools/udp": {
		{Key: "udp_ipv4", Value: "44.209.181.127"},
		{Key: "udp_port", Value: "8088"},
	},

	// Production
	"api.notefile.net/udp": {
		{Key: "udp_ipv4", Value: "44.209.181.127"},
		{Key: "udp_port", Value: "8089"},
	},
}

// Lookup the proxy for a given server
func httpProxyLookupHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" || r.URL.Path == "/favicon.ico" {
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	headers, present := headerIndex[strings.TrimPrefix(r.URL.Path, "/")]
	if !present {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for _, v := range headers {
		w.Header().Set(v.Key, v.Value)
	}

	w.WriteHeader(http.StatusOK)

}

// Register a UDP handler for each target
func udpProxyHandlers() {

	for target, headers := range headerIndex {
		for _, header := range headers {
			if header.Key == "udp_port" {
				go udpProxyHandler(target, header.Value)
			}
		}
	}

}

// Register a UDP handler for a single target
func udpProxyHandler(target string, port string) {

	addr, err := net.ResolveUDPAddr("udp", ":"+port)
	if err != nil {
		fmt.Printf("can't resolve UDP port %s for target %s: %v", port, target, err)
		os.Exit(43)
	}

	sock, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("can't listen on UDP port %s for target %s: %v", port, target, err)
		os.Exit(44)
	}
	defer sock.Close()

	targetURL := "https://" + target

	for {

		// Receive a packet
		buf := make([]byte, 65535)
		buflen, addr, err := sock.ReadFrom(buf)
		if err != nil {
			fmt.Printf("error reading from UDP port %s for target %s: %v", port, target, err)
			continue
		}

		// Instantiate a send func as a closure that sends a UDP message back to caller
		sendFunc := func(data []byte) error {

			ip, err := net.ResolveUDPAddr("udp", addr.String())
			if err != nil {
				return fmt.Errorf("udp: resolving return IP: %s", err)
			}

			if traceIo {
				fmt.Printf("rsp %s %d bytes to %s\n", port, len(data), addr.String())
			}
			_, err = sock.WriteTo(data, ip)
			if err != nil {
				return fmt.Errorf("udp: error writing socket: %s", err)
			}

			return err
		}

		// Trace
		if traceIo {
			fmt.Printf("req %s %d bytes to %s\n", port, buflen, targetURL)
		}

		// Dispatch the handling of a single UDP packet to a target
		go handlePacket(targetURL, sendFunc, buf[0:buflen])

	}

}

// Handle proxying a single incoming UDP packet
func handlePacket(targetUrl string, sendFunc func(data []byte) error, data []byte) {

	// Create a new HTTP request with the hex data as the body
	req, err := http.NewRequest("POST", targetUrl, bytes.NewBufferString(hex.EncodeToString(data)))
	if err != nil {
		fmt.Printf("proxy: error creating new request for %s: %v\n", targetUrl, err)
		return
	}

	// Send the request
	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		fmt.Printf("proxy: error sending UDP packet to %s: %v\n", targetUrl, err)
		return
	}
	defer rsp.Body.Close()

	// Read the response body
	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		fmt.Printf("proxy: error reading result UDP packet body from %s: %v\n", targetUrl, err)
		return
	}

	// Process the rare case where there is a downlink reply
	if len(rspBody) != 0 {

		// Decode the hex-formatted body
		decodedData, err := hex.DecodeString(string(rspBody))
		if err != nil {
			fmt.Printf("proxy: error decoding result UDP packet body from %s: %v\n", targetUrl, err)
			return
		}

		// Send it back to the caller
		err = sendFunc(decodedData)
		if err != nil {
			fmt.Printf("proxy: error returning result %d-byte UDP packet body: %v\n", len(decodedData), err)
			return
		}

	}

}
