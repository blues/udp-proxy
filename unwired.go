// Copyright 2022 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

// Unwired labs exporter
package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/blues/note-go/note"
)

// Constants & statics
const unwiredStateKey = "unwired"

type unwiredState struct {
	LastModifiedMs int64 `json:"last_modified_ms,omitempty"`
}

var scanRecsAvailable *Event

// Unwired Labs exporter task
func exportUnwired() {

	// Initialize our event
	scanRecsAvailable = EventNew()

	// Read the state
	var state unwiredState
	_, err := dbGetObject(unwiredStateKey, &state)
	if err != nil {
		fmt.Printf("unwired: reading state object: %s\n", err)
	}

	// Go into a perpetual loop, exporting state
	for {

		//
		since := state.LastModifiedMs
		until := time.Now().UTC().UnixNano() / 1000000

		fmt.Printf("unwired: looking for new records from %d - %d\n", since, until)

		// Do a query to find some number of the records since last time we did an export
		var recs []RadarScan
		recs, err = dbGetChangedRecs(since, until)
		if err != nil {
			fmt.Printf("unwired: error processing records: %s\n", err)
		}

		// Process the records
		if len(recs) > 0 {

			err = exportRecs(recs)
			if err != nil {
				fmt.Printf("unwired: error exporting records: %s\n", err)
			} else {
				// Save state so we don't read the same records again
				fmt.Printf("unwired: processed %d records\n", len(recs))
				state.LastModifiedMs = until
				err = dbSetObject(unwiredStateKey, &state)
				if err != nil {
					fmt.Printf("unwired: updating state object: %s\n", err)
				}
			}

		} else {

			fmt.Printf("unwired: no more records to process\n")

			// Wait until some device says that data is available to be aggregated
			scanRecsAvailable.Wait(24 * time.Hour)

			// We are awakened here by the first event streaming in from a device.  Because events
			// stream in groups, wait a moment before aggregating.  This has a positive secondary
			// benefit in that if many devices are pounding us, we naturally take a small breather
			// so we are not constantly querying the database for changes.
			time.Sleep(15 * time.Second)

		}

	}

}

// Signal that there are events ready
func unwiredScanEventsReady() {
	scanRecsAvailable.Signal()
}

// Export the records that have changed
func exportRecs(r []RadarScan) (err error) {

	fmt.Printf("exportRecs: %d records\n", len(r))

	// Sort the records based on the scan that was performed.  By doing this, we can use the
	// begin/end/duration/etc for the first record for all of them, so we know what to aggregate.
	sort.Slice(r, func(i, j int) bool {

		// Primamry key is tile ID
		if r[i].ScanFieldZID != r[j].ScanFieldZID {
			return r[i].ScanFieldZID < r[j].ScanFieldZID
		}

		// Secondary key is source ID
		if r[i].ScanFieldSID != r[j].ScanFieldSID {
			return r[i].ScanFieldSID < r[j].ScanFieldSID
		}

		// Tertiary key is when the scan began
		return r[i].ScanFieldBegan < r[j].ScanFieldBegan

	})

	// Iterate over the records, dividing them up into aggregateable sets that were done by the
	// same source in the same tile
	i := 0
	recsRemaining := len(r)
	for recsRemaining > 0 {

		count := 0
		for j := 0; recsRemaining > 0 && r[i].ScanFieldSID == r[i+j].ScanFieldSID && r[i].ScanFieldZID == r[i+j].ScanFieldZID && r[i].ScanFieldBegan == r[i+j].ScanFieldBegan; j++ {
			count++
			recsRemaining--
		}

		err = exportScan(r[i : i+count])
		if err != nil {
			fmt.Printf("exportTileRecs: %s\n", err)
		}

		i += count

	}

	// Success
	fmt.Printf("exportRecs: done exporting %d records\n", len(r))
	return

}

// Export records from within a single source within a single tile
func exportScan(r []RadarScan) (err error) {

	// Defensive, because we reference [0]
	if len(r) == 0 {
		return
	}

	fmt.Printf("exportScan: exporting %d-record scan done by %s in %s\n", len(r), r[0].ScanFieldSID, r[0].ScanFieldZID)

	// Begin to formulate an item by using a position at the midpoint of the line traveled during the scan
	var item ulItem
	item.Token = "<token>"
	timestampMs := (r[0].ScanFieldBegan + (r[0].ScanFieldDuration / 2)) * 1000

	// Add GPS array
	var pos ulPosition
	pos.Source = ulPositionSourceGPS
	pos.Timestamp = timestampMs
	if r[0].ScanFieldEndedLoc == "" {
		pos.Latitude, pos.Longitude = gpsFromOLC(r[0].ScanFieldBeganLoc)
	} else {
		pos.Latitude, pos.Longitude = gpsMidpointFromOLC(r[0].ScanFieldBeganLoc, r[0].ScanFieldEndedLoc)
		distanceMeters := olcDistanceMeters(r[0].ScanFieldBeganLoc, r[0].ScanFieldEndedLoc)
		if r[0].ScanFieldDuration != 0 && distanceMeters != 0 {
			pos.AccuracyMeters = distanceMeters / 2
			pos.SpeedMetersPerSec = distanceMeters / float64(r[0].ScanFieldDuration)
			pos.HeadingDeg = olcInitialBearing(r[0].ScanFieldBeganLoc, r[0].ScanFieldEndedLoc)
		}
	}
	item.GPS = append(item.GPS, pos)

	// Append the records from the various tiles
	for _, rec := range r {
		var c ulCell
		var w ulWiFi

		switch rec.ScanFieldDataRAT {
		case ScanRatGSM:
			c.Radio = ulRadioGSM
			c.Timestamp = timestampMs
			c.MCC = int(rec.ScanFieldDataMCC)
			c.MNC = int(rec.ScanFieldDataMNC)
			c.LAC = int(rec.ScanFieldDataTAC)
			c.CID = int(rec.ScanFieldDataCID)
			c.Signal = int(rec.ScanFieldDataRSSI)
		case ScanRatCDMA:
			c.Radio = ulRadioCDMA
			c.Timestamp = timestampMs
			c.MCC = int(rec.ScanFieldDataMCC)
			c.MNC = int(rec.ScanFieldDataMNC)
			c.LAC = int(rec.ScanFieldDataTAC)
			c.CID = int(rec.ScanFieldDataCID)
			c.Signal = int(rec.ScanFieldDataRSSI)
		case ScanRatUMTS:
			c.Radio = ulRadioUMTS
			c.Timestamp = timestampMs
			c.MCC = int(rec.ScanFieldDataMCC)
			c.MNC = int(rec.ScanFieldDataMNC)
			c.LAC = int(rec.ScanFieldDataTAC)
			c.CID = int(rec.ScanFieldDataCID)
			c.Signal = int(rec.ScanFieldDataRSSI)
		case ScanRatWCDMA:
			c.Radio = ulRadioCDMA
			c.Timestamp = timestampMs
			c.MCC = int(rec.ScanFieldDataMCC)
			c.MNC = int(rec.ScanFieldDataMNC)
			c.LAC = int(rec.ScanFieldDataTAC)
			c.CID = int(rec.ScanFieldDataCID)
			c.Signal = int(rec.ScanFieldDataRSSI)
		case ScanRatLTE:
			c.Radio = ulRadioLTE
			c.Timestamp = timestampMs
			c.MCC = int(rec.ScanFieldDataMCC)
			c.MNC = int(rec.ScanFieldDataMNC)
			c.LAC = int(rec.ScanFieldDataTAC)
			c.CID = int(rec.ScanFieldDataCID)
			c.PCI = int(rec.ScanFieldDataPCI)
			c.Band = int(rec.ScanFieldDataBAND)
			c.Channel = int(rec.ScanFieldDataCHAN)
			c.Signal = int(rec.ScanFieldDataRSSI)
		case ScanRatEMTC:
			c.Radio = ulRadioLTE
			c.Timestamp = timestampMs
			c.MCC = int(rec.ScanFieldDataMCC)
			c.MNC = int(rec.ScanFieldDataMNC)
			c.LAC = int(rec.ScanFieldDataTAC)
			c.CID = int(rec.ScanFieldDataCID)
			c.PCI = int(rec.ScanFieldDataPCI)
			c.Band = int(rec.ScanFieldDataBAND)
			c.Channel = int(rec.ScanFieldDataCHAN)
			c.Signal = int(rec.ScanFieldDataRSSI)
		case ScanRatNBIOT:
			c.Radio = ulRadioNBIOT
			c.Timestamp = timestampMs
			c.MCC = int(rec.ScanFieldDataMCC)
			c.MNC = int(rec.ScanFieldDataMNC)
			c.LAC = int(rec.ScanFieldDataTAC)
			c.CID = int(rec.ScanFieldDataCID)
			c.PCI = int(rec.ScanFieldDataPCI)
			c.Band = int(rec.ScanFieldDataBAND)
			c.Channel = int(rec.ScanFieldDataCHAN)
			c.Signal = int(rec.ScanFieldDataRSSI)
		case ScanRatNR:
			c.Radio = ulRadioNR
			c.Timestamp = timestampMs
			c.MCC = int(rec.ScanFieldDataMCC)
			c.MNC = int(rec.ScanFieldDataMNC)
			c.LAC = int(rec.ScanFieldDataTAC)
			c.CID = int(rec.ScanFieldDataCID)
			c.PCI = int(rec.ScanFieldDataPCI)
			c.Band = int(rec.ScanFieldDataBAND)
			c.Channel = int(rec.ScanFieldDataCHAN)
			c.Signal = int(rec.ScanFieldDataRSSI)
		case ScanRatWIFI:
			w.Timestamp = timestampMs
			w.BSSID = rec.ScanFieldDataBSSID
			w.SSID = rec.ScanFieldDataSSID
			w.Channel = int(rec.ScanFieldDataCHAN)
			w.Frequency = int(rec.ScanFieldDataFREQ)
			w.Signal = int(rec.ScanFieldDataRSSI)
			w.SNR = int(rec.ScanFieldDataSNR)
		}

		if c.Radio != "" {
			item.Cells = append(item.Cells, c)
		}
		if w.BSSID != "" {
			item.WiFi = append(item.WiFi, w)
		}

	}

	// Marshal
	var ulJSON []byte
	ulJSON, err = note.JSONMarshalIndent(item, "", "    ")

	// Trace
	fmt.Printf("%s\n", string(ulJSON))

	// Done
	return

}
