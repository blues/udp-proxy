// Copyright 2022 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

package main

import (
	"math"

	olc "github.com/google/open-location-code/go"
)

// OLC precision definitions
const olcPrecision2Meters = 2226000.0
const olcPrecision4Meters = 111321.0
const olcPrecision6Meters = 5566.0
const olcPrecision8Meters = 278.0
const olcPrecision10Meters = 13.9
const olcPrecision11Meters = 3.5   // 2.8m x 3.5m
const olcPrecision12Meters = 0.87  // 56cm x 87cm
const olcPrecision13Meters = 0.22  // 11cm x 22cm
const olcPrecision14Meters = 0.05  // 2cm x 5cm
const olcPrecision15Meters = 0.014 // 4mm x 14mm
const olcPrecisionDefault = 13

// Constants
const R = 6371
const degreesToRadians = (3.1415926536 / 180)
const radiansToDegrees = (180 / 3.1415926536)

func gpsFromOLC(loc string) (lat float64, lon float64) {
	area, err := olc.Decode(loc)
	if err == nil {
		lat, lon = area.Center()
	}
	return
}

func olcFromGPS(lat float64, lon float64, precision int) (loc string) {
	loc = olc.Encode(lat, lon, precision)
	return
}

func olcDistanceMeters(loc1 string, loc2 string) (distanceMeters float64) {
	lat1, lon1 := gpsFromOLC(loc1)
	lat2, lon2 := gpsFromOLC(loc1)
	distanceMeters = gpsDistanceMeters(lat1, lon1, lat2, lon2)
	return
}

func gpsDistanceMeters(lat1 float64, lon1 float64, lat2 float64, lon2 float64) (distanceMeters float64) {
	var dx, dy, dz float64
	lon1 = lon1 - lon2
	lon1 = lon1 * degreesToRadians
	lat1 = lat1 * degreesToRadians
	lat2 = lat2 * degreesToRadians
	dz = math.Sin(lat1) - math.Sin(lat2)
	dx = math.Cos(lon1)*math.Cos(lat1) - math.Cos(lat2)
	dy = math.Sin(lon1) * math.Cos(lat1)
	distanceMeters = 1000 * (math.Asin(math.Sqrt(math.Abs(dx*dx+dy*dy+dz*dz))/2) * 2 * R)
	return
}

func olcInitialBearing(loc1 string, loc2 string) (bearingDegrees float64) {
	lat1, lon1 := gpsFromOLC(loc1)
	lat2, lon2 := gpsFromOLC(loc1)
	bearingDegrees = gpsInitialBearing(lat1, lon1, lat2, lon2)
	return
}

func gpsInitialBearing(lat1 float64, lon1 float64, lat2 float64, lon2 float64) (bearingDegrees float64) {
	lat1 = lat1 * degreesToRadians
	lat2 = lat2 * degreesToRadians
	deltaLambda := (lon2 * degreesToRadians) - (lon1 * degreesToRadians)
	y := math.Sin(deltaLambda) * math.Cos(lat2)
	x := (math.Cos(lat1) * math.Sin(lat2)) - (math.Sin(lat1) * math.Cos(lat2) * math.Cos(deltaLambda))
	theta := math.Atan2(y, x)
	bearingDegrees = math.Mod((theta*radiansToDegrees)+360.0, 360.0)
	return
}

func gpsMidpointFromOLC(loc1 string, loc2 string) (lat3 float64, lon3 float64) {
	lat1, lon1 := gpsFromOLC(loc1)
	lat2, lon2 := gpsFromOLC(loc1)
	lat3, lon3 = gpsMidpoint(lat1, lon1, lat2, lon2)
	return
}

func gpsMidpoint(lat1 float64, lon1 float64, lat2 float64, lon2 float64) (lat3 float64, lon3 float64) {
	if lat1 == lat2 && lon1 == lon2 {
		lat3 = lat1
		lon3 = lon1
	} else {
		dLon := (lon2 - lon1) * degreesToRadians
		lat1 = lat1 * degreesToRadians
		lat2 = lat2 * degreesToRadians
		lon1 = lon1 * degreesToRadians
		Bx := math.Cos(lat2) * math.Cos(dLon)
		By := math.Cos(lat2) * math.Sin(dLon)
		lat3r := math.Atan2(math.Sin(lat1)+math.Sin(lat2), math.Sqrt((math.Cos(lat1)+Bx)*(math.Cos(lat1)+Bx)+By*By))
		lon3r := lon1 + math.Atan2(By, math.Cos(lat1)+Bx)
		lat3 = lat3r * radiansToDegrees
		lon3 = lon3r * radiansToDegrees
	}
	return

}

// For go vet
var _ = olcPrecision2Meters
var _ = olcPrecision4Meters
var _ = olcPrecision6Meters
var _ = olcPrecision8Meters
var _ = olcPrecision10Meters
var _ = olcPrecision11Meters
var _ = olcPrecision12Meters
var _ = olcPrecision13Meters
var _ = olcPrecision14Meters
var _ = olcPrecision15Meters
var _ = olcPrecisionDefault
var _ = gpsFromOLC
var _ = olcFromGPS
var _ = olcDistanceMeters
var _ = olcInitialBearing
var _ = gpsInitialBearing
var _ = gpsMidpointFromOLC
var _ = gpsMidpoint
