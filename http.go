// Copyright Blues Inc.	 All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

// Common support for all HTTP topic handlers
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"time"
)

// HTTPInboundHandler kicks off inbound messages coming from all sources, then serve HTTP
func HTTPInboundHandler(port string) {

	// Topics
	http.HandleFunc("/github", inboundWebGithubHandler)
	http.HandleFunc("/ping", inboundWebPingHandler)
	http.HandleFunc("/", inboundWebRootHandler)

	// HTTP
	fmt.Printf("Now handling inbound HTTP on %s\n", port)
	go http.ListenAndServe(port, nil)

}

// Ping handler
func inboundWebPingHandler(httpRsp http.ResponseWriter, httpReq *http.Request) {

	// Get the body if supplied
	reqJSON, err := ioutil.ReadAll(httpReq.Body)
	if err != nil {
		reqJSON = []byte("{}")
	}
	_ = reqJSON

	// Write reply JSON
	rspJSON := []byte(time.Now().UTC().Format("2006-01-02T15:04:05Z"))
	httpRsp.Write(rspJSON)

}

// Root handler
func inboundWebRootHandler(httpRsp http.ResponseWriter, httpReq *http.Request) {

	// Process the request URI, looking for things that will indicate "dev"
	method := httpReq.Method
	if method == "" {
		method = "GET"
	}

	// Get the target
	parsedURL, _ := url.Parse(httpReq.RequestURI)
	target := path.Base(parsedURL.Path)

	// Exit if just the favicon
	if target == "favicon.ico" {
		return
	}

	// Done
	httpRsp.Write([]byte("I'm watching you."))

}
