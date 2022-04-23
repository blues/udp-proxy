// Copyright 2022 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

// Ingests data generated by a radar device
package main

// Contact cache which lowers the frequency of update of the sql database.  This is important
// because otherwise there would be one contact update per ingest of a scan or track item
var contactCache map[string]string

// Ingest a contact
func ingestContact(deviceUID string, when int64, deviceSN string, cName string, cAffiliation string, cRole string, cEmail string) (err error) {

	// See if it has changed since what's in the cache
	if contactCache == nil {
		contactCache = map[string]string{}
	}
	prevValue, _ := contactCache[deviceUID]
	thisValue := deviceSN + ":" + cName + ":" + cAffiliation + ":" + cRole + ":" + cEmail
	if prevValue == thisValue {
		return
	}
	contactCache[deviceUID] = thisValue

	// Add the contact
	err = dbAddContact(deviceUID, when, deviceSN, cName, cAffiliation, cRole, cEmail)
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
