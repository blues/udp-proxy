// Copyright 2018 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

// Unwired Submit API docs:
// https://unwiredlabs.com/data-apiv2

package main

type ulItem struct {
	TimestampMs int64      `json:"timestamp,omitempty"`
	Position    ulPosition `json:"position,omitempty"`
	Cells       []ulCell   `json:"cells,omitempty"`
	WiFi        []ulWiFi   `json:"wifi,omitempty"`
}

type ulPosition struct {
	Latitude               float64 `json:"latitude,omitempty"`
	Longitude              float64 `json:"longitude,omitempty"`
	AccuracyMeters         float64 `json:"accuracy,omitempty"`         // optional
	AltitudeMeters         float64 `json:"altitude,omitempty"`         // optional
	AltitudeAccuracyMeters float64 `json:"altitudeAccuracy,omitempty"` // optional
	SpeedMetersPerSec      float64 `json:"speed,omitempty"`            // optional
	HeadingDeg             float64 `json:"heading,omitempty"`          // optional
}

type ulWiFi struct {
	BSSID     string `json:"bssid,omitempty"`              // xx:xx:xx:xx:xx:xx
	SSID      string `json:"ssid,omitempty"`               // optional
	Channel   int    `json:"channel,omitempty"`            // optional
	Frequency int    `json:"frequency,omitempty"`          // optional
	Signal    int    `json:"signal,omitempty"`             // optional
	SNR       int    `json:"signalToNoiseRatio,omitempty"` // optional
}

type ulCell struct {
	Radio   string `json:"radio,omitempty"`
	MCC     int    `json:"mcc,omitempty"`
	MNC     int    `json:"mnc,omitempty"`
	LAC     int    `json:"lac,omitempty"`
	CID     int    `json:"cid,omitempty"`
	Serving int    `json:"serving,omitempty"`
	Signal  int    `json:"signal,omitempty"`
}

const ulRadioGSM = "gsm"     // GSM, EDGE, GPRS, 2G
const ulRadioLTE = "lte"     // LTE, 4G
const ulRadioCDMA = "cdma"   // 1xRTT, CDMA, eHRPD, EVDO_0, EVDO_A, EVDO_B, IS95A, IS95B
const ulRadioUMTS = "umts"   // UMTS, HSPA, HSDPA, HSPA+, HSUPA, WCDMA, 3G
const ulRadioNBIOT = "nbiot" // NB-IoT
const ulRadioNR = "nr"       // NR, 5G
