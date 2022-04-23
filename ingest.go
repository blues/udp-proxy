// Copyright 2022 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

// Ingests data generated by a radar device
package main

// Ingest a contact
func ingestContact(deviceUID string, when int64, deviceSN string, contactName string, contactAffiliation string, contactRole string, contactEmail string) (err error) {
	err = dbAddContact(deviceUID, when, deviceSN, contactName, contactAffiliation, contactRole, contactEmail)
	return
}

// Ingest a scan entry
func ingestScan(deviceUID string, scan RadarScan) (err error) {
	err = dbAddScan(deviceUID, scan)
	return
}

// Ingest a track entry
func ingestTrack(deviceUID string, track RadarTrack) (err error) {
	err = dbAddTrack(deviceUID, track)
	return
}
