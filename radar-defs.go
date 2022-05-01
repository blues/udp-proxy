// Copyright 2022 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

// Definitions for data flowing in from radar devices
package main

// Scan formats
const ScanRatGSM = "gsm"
const ScanRatCDMA = "cdma"
const ScanRatUMTS = "umts"
const ScanRatWCDMA = "wcdma"
const ScanRatLTE = "lte"
const ScanRatEMTC = "emtc"
const ScanRatNBIOT = "nbiot"
const ScanRatNR = "nr"
const ScanRatWIFI = "wifi"

// Body of the notes to be enqueued to the service for scanning
const ScanNotefile = "scan.qo"

type RadarScan struct {
	ScanFieldSID             string  `json:"sid,omitempty"`
	ScanFieldTID             string  `json:"tid,omitempty"`
	ScanFieldXID             string  `json:"xid,omitempty"`
	ScanFieldTime            int64   `json:"time,omitempty"`
	ScanFieldDuration        int64   `json:"duration,omitempty"`
	ScanFieldDistance        float64 `json:"distance,omitempty"`
	ScanFieldBearing         float64 `json:"bearing,omitempty"`
	ScanFieldBegan           int64   `json:"began,omitempty"`
	ScanFieldBeganLoc        string  `json:"began_loc,omitempty"`
	ScanFieldBeganLocHDOP    int64   `json:"began_loc_hdop,omitempty"`
	ScanFieldBeganLocTime    int64   `json:"began_loc_time,omitempty"`
	ScanFieldBeganMotionTime int64   `json:"began_motion_time,omitempty"`
	ScanFieldEnded           int64   `json:"ended,omitempty"`
	ScanFieldEndedLoc        string  `json:"ended_loc,omitempty"`
	ScanFieldEndedLocHDOP    int64   `json:"ended_loc_hdop,omitempty"`
	ScanFieldEndedLocTime    int64   `json:"ended_loc_time,omitempty"`
	ScanFieldEndedMotionTime int64   `json:"ended_motion_time,omitempty"`
	ScanFieldDataRAT         string  `json:"rat,omitempty"`
	ScanFieldDataMCC         int64   `json:"mcc,omitempty"`
	ScanFieldDataMNC         int64   `json:"mnc,omitempty"`
	ScanFieldDataTAC         int64   `json:"tac,omitempty"`
	ScanFieldDataCID         int64   `json:"cid,omitempty"`
	ScanFieldDataPCI         int64   `json:"pci,omitempty"`
	ScanFieldDataBAND        int64   `json:"band,omitempty"`
	ScanFieldDataCHAN        int64   `json:"chan,omitempty"`
	ScanFieldDataFREQ        int64   `json:"freq,omitempty"`
	ScanFieldDataBSSID       string  `json:"bssid,omitempty"`
	ScanFieldDataPSC         int64   `json:"psc,omitempty"`
	ScanFieldDataRSSI        int64   `json:"rssi,omitempty"`
	ScanFieldDataRSRP        int64   `json:"rsrp,omitempty"`
	ScanFieldDataRSRQ        int64   `json:"rsrq,omitempty"`
	ScanFieldDataRSCP        int64   `json:"rscp,omitempty"`
	ScanFieldDataSNR         int64   `json:"snr,omitempty"`
	ScanFieldDataSSID        string  `json:"ssid,omitempty"`
}

// For standard tracking, the data format of a single point
const TrackTypeNormal = ""
const TrackTypeHeartbeat = "heartbeat"
const TrackTypeUSBChange = "usb"
const TrackTypeNoSat = "no-sat"

// Body of the notes to be enqueued to the service for tracking
const TrackNotefile = "track.qo"

type RadarTrack struct {
	TrackFieldTime           int64   `json:"when,omitempty"`
	TrackFieldLoc            string  `json:"loc,omitempty"`
	TrackFieldLocTime        int64   `json:"time,omitempty"`
	TrackFieldLocHDOP        int64   `json:"hdop,omitempty"`
	TrackFieldJourneyTime    int64   `json:"journey,omitempty"`
	TrackFieldJourneyCount   int64   `json:"jcount,omitempty"`
	TrackFieldMotionCount    int64   `json:"motion,omitempty"`
	TrackFieldMotionTime     int64   `json:"motion_time,omitempty"`
	TrackFieldMotionDistance float64 `json:"motion_distance,omitempty"`
	TrackFieldMotionBearing  float64 `json:"motion_bearing,omitempty"`
	TrackFieldMotionVelocity float64 `json:"motion_velocity,omitempty"`
	TrackFieldTemperature    float64 `json:"temperature,omitempty"`
	TrackFieldHumidity       float64 `json:"humidity,omitempty"`
	TrackFieldPressure       float64 `json:"pressure,omitempty"`
	TrackFieldFlagUSB        bool    `json:"usb,omitempty"`
	TrackFieldFlagCharging   bool    `json:"charging,omitempty"`
	TrackFieldFlagHeartbeat  bool    `json:"heartbeat,omitempty"`
}
