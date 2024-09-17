// Copyright 2024 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
)

// Display the proxy server activity on console
const traceIo = true

// Lookup the proxy for a given server
func httpProxyLookupHandler(w http.ResponseWriter, r *http.Request) {
	if traceIo {
		fmt.Print(getNowTimestamp(), " ", r.RemoteAddr, " ", r.Method, " ", r.URL.Path, " ")
	}

	if r.Method != "GET" || r.URL.Path == "/favicon.ico" {
		if traceIo {
			fmt.Println(http.StatusNotImplemented)
		}
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	headers, present := proxyData[strings.TrimPrefix(r.URL.Path, "/")]

	if !present {
		if traceIo {
			fmt.Println(http.StatusNotFound)
		}
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for _, v := range headers {
		w.Header().Set(v.Key, v.Value)
	}

	if traceIo {
		fmt.Println(http.StatusOK, headers[0].Value, headers[1].Value) // ipv4, port (headers should really be changed to a real structure)
	}

	w.WriteHeader(http.StatusOK)

}

// Register a UDP handler for each target
func udpProxyHandlers() {

	for target, headers := range proxyData {
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
		loggedExit(43, "can't resolve UDP port", port, "for target", target, ":", err)
	}

	sock, err := net.ListenUDP("udp", addr)
	if err != nil {
		loggedExit(44, "can't listen on UDP port", port, "for target", target, ":", err)
	}
	defer sock.Close()

	targetURL := "https://" + target

	for {

		// Receive a packet
		buf := make([]byte, 65535)
		buflen, addr, err := sock.ReadFrom(buf)
		if err != nil {
			fmt.Printf("%s error reading from UDP port %s for target %s: %v", getNowTimestamp(), port, target, err)
			continue
		}

		// Instantiate a send func as a closure that sends a UDP message back to caller
		sendFunc := func(data []byte) error {

			ip, err := net.ResolveUDPAddr("udp", addr.String())
			if err != nil {
				return fmt.Errorf("udp: resolving return IP: %s", err)
			}

			if traceIo {
				fmt.Printf("%s %s rsp(%s) %d bytes to %s\n", getNowTimestamp(), addr.String(), port, len(data), addr.String())
			}
			_, err = sock.WriteTo(data, ip)
			if err != nil {
				return fmt.Errorf("udp: error writing socket: %s", err)
			}

			return err
		}

		// Trace
		if traceIo {
			fmt.Printf("%s %s req(%s) %d bytes to %s\n", getNowTimestamp(), addr.String(), port, buflen, targetURL)
		}

		// Dispatch the handling of a single UDP packet to a target
		go handlePacket(targetURL, addr, sendFunc, buf[0:buflen])

	}

}

// Handle proxying a single incoming UDP packet
func handlePacket(targetUrl string, addr net.Addr, sendFunc func(data []byte) error, data []byte) {

	// Create a new HTTP request with the hex data as the body
	req, err := http.NewRequest("POST", targetUrl, bytes.NewBufferString(hex.EncodeToString(data)))
	if err != nil {
		fmt.Printf("%s proxy: error creating new request for %s: %v\n", getNowTimestamp(), targetUrl, err)
		return
	}

	// Send the request
	client := &http.Client{}
	rsp, err := client.Do(req)
	if err != nil {
		fmt.Printf("%s proxy: error sending UDP packet to %s: %v\n", getNowTimestamp(), targetUrl, err)
		return
	}
	defer rsp.Body.Close()

	if traceIo {
		fmt.Printf("%s %s fwd %d bytes to %s: status %s\n", getNowTimestamp(), addr.String(), len(data), targetUrl, rsp.Status)
	}

	// Read the response body
	rspBody, err := io.ReadAll(rsp.Body)
	if err != nil {
		fmt.Printf("%s proxy: error reading result UDP packet body from %s: %v\n", getNowTimestamp(), targetUrl, err)
		return
	}

	// Process the rare case where there is a downlink reply
	if len(rspBody) != 0 {

		// Decode the hex-formatted body
		decodedData, err := hex.DecodeString(string(rspBody))
		if err != nil {
			fmt.Printf("%s proxy: error decoding result UDP packet body from %s: %v\n", getNowTimestamp(), targetUrl, err)
			return
		}

		// Send it back to the caller
		err = sendFunc(decodedData)
		if err != nil {
			fmt.Printf("%s proxy: error returning result %d-byte UDP packet body: %v\n", getNowTimestamp(), len(decodedData), err)
			return
		}

	}

}
