// Copyright 2022 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

// Unwired labs exporter
package main

import (
	"fmt"
	"time"
)

const unwiredStateKey = "unwired"

type unwiredState struct {
	LastModifiedMs int64 `json:"last_modified_ms,omitempty"`
}

// Unwired Labs exporter task
func exportUnwired() {

	// Read the state
	var state unwiredState
	_, err := dbGetObject(unwiredStateKey, &state)
	if err != nil {
		fmt.Printf("unwired: reading state object: %s\n", err)
	}

	// Go into a perpetual loop, exporting state
	for {

		// Do a query to find some number of the records since last time we did an export
		var recs int
		recs, err = dbEnumNewScanRecs(state.LastModifiedMs, 100, unwiredExportScanRec, &state)
		if err != nil {
			fmt.Printf("unwired: error processing records: %s\n", err)
		}
		if recs > 0 {

			// Success
			fmt.Printf("unwired: processed %d records\n", recs)

			// If any recs were added, update state
			err = dbSetObject(unwiredStateKey, &state)
			if err != nil {
				fmt.Printf("unwired: updating state object: %s\n", err)
			} else {

				// Sleep a little just to be sociable, and loop to process more
				time.Sleep(10 * time.Second)
				continue
			}

		}

		// Wait for a more substantial amount of time before trying again
		time.Sleep(120 * time.Second)

	}

}

// Export a single record
func unwiredExportScanRec(state *unwiredState, deviceUID string, recordModifiedMs int64, r RadarScan) (err error) {

	fmt.Printf("EXPORT SCAN: %s %d\n", deviceUID, recordModifiedMs)

	// Begin to formulate an item
	var item ulItem
	item.TimestampMs = recordModifiedMs

	// Update the modified MS under the assumption that these are enumerated in ASC sequence
	state.LastModifiedMs = recordModifiedMs

	// Success
	return

}
