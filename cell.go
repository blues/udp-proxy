// Copyright 2022 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package main

import (
	"fmt"

	olc "github.com/google/open-location-code/go"
	"github.com/uber/h3-go"
)

// Consistent with what others use when doing mapping
const cellResolution = 8

// Convert a lat/lon into a cell ID
func cellFromLatLon(latDegrees float64, lonDegrees float64) (cid string) {
	if latDegrees == 0 && lonDegrees == 0 {
		return "?"
	}
	geo := h3.GeoCoord{
		Latitude:  latDegrees,
		Longitude: lonDegrees,
	}
	resolution := cellResolution
	return fmt.Sprintf("%016X", h3.FromGeo(geo, resolution))
}

// Convert an OLC into a cell ID
func cellFromOLC(loc string) (cid string) {
	area, err := olc.Decode(loc)
	if err != nil {
		return
	}
	lat, lon := area.Center()
	return cellFromLatLon(lat, lon)
}
