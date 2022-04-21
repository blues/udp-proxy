// Copyright 2022 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

// Ingests data sent in via notehub's route
package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/blues/note-go/note"
)

// Ingest handler
func inboundWebIngestHandler(httpRsp http.ResponseWriter, httpReq *http.Request) {

	// Get the body if supplied
	eventJSON, err := ioutil.ReadAll(httpReq.Body)
	if err != nil {
		httpRsp.WriteHeader(http.StatusNoContent)
		httpRsp.Write([]byte(fmt.Sprintf("{\"err\":\"%s\"}", err)))
		return
	}
	if len(eventJSON) == 0 {
		httpRsp.WriteHeader(http.StatusNoContent)
		return
	}

	// Unmarshal the event
	var e note.Event
	err = note.JSONUnmarshal(eventJSON, &e)
	if err != nil {
		httpRsp.WriteHeader(http.StatusBadRequest)
		httpRsp.Write([]byte(fmt.Sprintf("{\"err\":\"%s\"}", err)))
		return
	}

	// Ingest different information depending upon notefile
	if e.Body != nil {
		switch e.NotefileID {

		case ScanNotefile:
			var data RadarScan
			err = note.BodyToObject(e.Body, &data)
			if err != nil {
				err = ingestScan(e.DeviceUID, data)
			}

		case TrackNotefile:
			var data RadarTrack
			err = note.BodyToObject(e.Body, &data)
			if err != nil {
				err = ingestTrack(e.DeviceUID, data)
			}

		default:
			httpRsp.WriteHeader(http.StatusOK)
			return

		}
	}
	if err != nil {
		httpRsp.WriteHeader(http.StatusBadRequest)
		httpRsp.Write([]byte(fmt.Sprintf("{\"err\":\"%s\"}", err)))
		return
	}

	// Assuming we ingested something, also ingest the contact
	if e.When != 0 {
		if e.DeviceContact == nil {
			e.DeviceContact = &note.EventContact{}
		}
		err = ingestContact(e.When, e.DeviceUID, e.DeviceSN,
			e.DeviceContact.Name, e.DeviceContact.Affiliation, e.DeviceContact.Role, e.DeviceContact.Email)
		if err != nil {
			httpRsp.WriteHeader(http.StatusBadRequest)
			httpRsp.Write([]byte(fmt.Sprintf("{\"err\":\"%s\"}", err)))
			return
		}
	}

	// Write reply JSON
	httpRsp.WriteHeader(http.StatusOK)

}
