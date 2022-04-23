// Copyright 2022 Blues Inc.  All rights reserved.
// Use of this source code is governed by licenses granted by the
// copyright holder including that found in the LICENSE file.

// Database handling
package main

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"go.elastic.co/apm/module/apmsql"
	_ "go.elastic.co/apm/module/apmsql/pq"
)

// Protects db access
var dbLock sync.Mutex

// DbDesc is our own database descriptor
type DbDesc struct {
	db *sql.DB
}

var radarDb DbDesc

// Fields common across tables
const fieldDbSerial = "db_serial"
const fieldDbSerialType = "BIGSERIAL"
const fieldDbModified = "db_modified"
const fieldDbModifiedType = "TIMESTAMP WITHOUT TIME ZONE"
const fieldSID = "sid"
const fieldSIDType = "TEXT"
const fieldZID = "zid"
const fieldZIDType = "TEXT"
const fieldTime = "time"
const fieldTimeType = "INTEGER"

// Scan table params
const tableScan = "scan"
const scanFieldDbSerial = fieldDbSerial
const scanFieldDbSerialType = fieldDbSerialType
const scanFieldDbModified = fieldDbModified
const scanFieldDbModifiedType = fieldDbModifiedType
const scanFieldSID = fieldSID
const scanFieldSIDType = fieldSIDType
const scanFieldZID = fieldZID
const scanFieldZIDType = fieldZIDType
const scanFieldXID = "xid"
const scanFieldXIDType = "TEXT"
const scanFieldTime = fieldTime
const scanFieldTimeType = fieldTimeType
const scanFieldDuration = "duration"
const scanFieldDurationType = "INTEGER"
const scanFieldDistance = "distance"
const scanFieldDistanceType = "REAL"
const scanFieldBearing = "bearing"
const scanFieldBearingType = "REAL"
const scanFieldBegan = "began"
const scanFieldBeganType = "INTEGER"
const scanFieldBeganLoc = "began_loc"
const scanFieldBeganLocType = "TEXT"
const scanFieldBeganLocHDOP = "began_loc_hdop"
const scanFieldBeganLocHDOPType = "INTEGER"
const scanFieldBeganLocTime = "began_loc_time"
const scanFieldBeganLocTimeType = "INTEGER"
const scanFieldBeganMotionTime = "began_motion_time"
const scanFieldBeganMotionTimeType = "INTEGER"
const scanFieldEnded = "ended"
const scanFieldEndedType = "INTEGER"
const scanFieldEndedLoc = "ended_loc"
const scanFieldEndedLocType = "TEXT"
const scanFieldEndedLocHDOP = "ended_loc_hdop"
const scanFieldEndedLocHDOPType = "INTEGER"
const scanFieldEndedLocTime = "ended_loc_time"
const scanFieldEndedLocTimeType = "INTEGER"
const scanFieldEndedMotionTime = "ended_motion_time"
const scanFieldEndedMotionTimeType = "INTEGER"
const scanFieldDataRAT = "rat"
const scanFieldDataRATType = "TEXT"
const scanFieldDataMCC = "mcc"
const scanFieldDataMCCType = "INTEGER"
const scanFieldDataMNC = "mnc"
const scanFieldDataMNCType = "INTEGER"
const scanFieldDataTAC = "tac"
const scanFieldDataTACType = "INTEGER"
const scanFieldDataCID = "cid"
const scanFieldDataCIDType = "INTEGER"
const scanFieldDataPCI = "pci"
const scanFieldDataPCIType = "INTEGER"
const scanFieldDataBAND = "band"
const scanFieldDataBANDType = "INTEGER"
const scanFieldDataCHAN = "chan"
const scanFieldDataCHANType = "INTEGER"
const scanFieldDataFREQ = "freq"
const scanFieldDataFREQType = "INTEGER"
const scanFieldDataBSSID = "bssid"
const scanFieldDataBSSIDType = "TEXT"
const scanFieldDataPSC = "psc"
const scanFieldDataPSCType = "INTEGER"
const scanFieldDataRSSI = "rssi"
const scanFieldDataRSSIType = "INTEGER"
const scanFieldDataRSRP = "rsrp"
const scanFieldDataRSRPType = "INTEGER"
const scanFieldDataRSRQ = "rsrq"
const scanFieldDataRSRQType = "INTEGER"
const scanFieldDataRSCP = "rscp"
const scanFieldDataRSCPType = "INTEGER"
const scanFieldDataSNR = "snr"
const scanFieldDataSNRType = "INTEGER"
const scanFieldDataSSID = "ssid"
const scanFieldDataSSIDType = "TEXT"

// Track table params
const tableTrack = "track"
const trackFieldDbSerial = fieldDbSerial
const trackFieldDbSerialType = fieldDbSerialType
const trackFieldDbModified = fieldDbModified
const trackFieldDbModifiedType = fieldDbModifiedType
const trackFieldSID = fieldSID
const trackFieldSIDType = fieldSIDType
const trackFieldZID = fieldZID
const trackFieldZIDType = fieldZIDType
const trackFieldTime = "when"
const trackFieldTimeType = "INTEGER"
const trackFieldLoc = "loc"
const trackFieldLocType = "TEXT"
const trackFieldLocTime = fieldTime
const trackFieldLocTimeType = fieldTimeType
const trackFieldLocHDOP = "hdop"
const trackFieldLocHDOPType = "INTEGER"
const trackFieldJourneyTime = "journey"
const trackFieldJourneyTimeType = "INTEGER"
const trackFieldJourneyCount = "jcount"
const trackFieldJourneyCountType = "INTEGER"
const trackFieldMotionCount = "motion"
const trackFieldMotionCountType = "INTEGER"
const trackFieldMotionTime = "motion_time"
const trackFieldMotionTimeType = "INTEGER"
const trackFieldMotionDistance = "motion_distance"
const trackFieldMotionDistanceType = "REAL"
const trackFieldMotionBearing = "motion_bearing"
const trackFieldMotionBearingType = "REAL"
const trackFieldMotionVelocity = "motion_velocity"
const trackFieldMotionVelocityType = "REAL"
const trackFieldTemperature = "temperature"
const trackFieldTemperatureType = "REAL"
const trackFieldHumidity = "humidity"
const trackFieldHumidityType = "REAL"
const trackFieldPressure = "pressure"
const trackFieldPressureType = "REAL"
const trackFieldFlagUSB = "usb"
const trackFieldFlagUSBType = "INTEGER"
const trackFieldFlagCharging = "charging"
const trackFieldFlagChargingType = "INTEGER"
const trackFieldFlagHeartbeat = "heartbeat"
const trackFieldFlagHeartbeatType = "INTEGER"

// Contact table params
const tableContact = "contact"
const contactFieldDbSerial = fieldDbSerial
const contactFieldDbSerialType = fieldDbSerialType
const contactFieldDbModified = fieldDbModified
const contactFieldDbModifiedType = fieldDbModifiedType
const contactFieldSID = fieldSID
const contactFieldSIDType = fieldSIDType
const contactFieldTime = fieldTime
const contactFieldTimeType = fieldTimeType
const contactFieldSN = "sn"
const contactFieldSNType = "INTEGER"
const contactFieldName = "name"
const contactFieldNameType = "TEXT"
const contactFieldAffiliation = "affiliation"
const contactFieldAffiliationType = "TEXT"
const contactFieldRole = "role"
const contactFieldRoleType = "TEXT"
const contactFieldEmail = "email"
const contactFieldEmailType = "TEXT"

// Initialize the db subsystem and make sure the tables are created
func dbInit() (err error) {
	var exists bool

	// Open the db
	db, err := dbContext()
	if err != nil {
		return
	}

	// Lock
	dbLock.Lock()
	defer dbLock.Unlock()

	// Initialize the scan table
	exists, err = uTableExists(db, tableScan)
	if err != nil {
		return
	}
	if !exists {

		// Create the scan table
		query := fmt.Sprintf("CREATE TABLE \"%s\" ( \n", tableScan)
		query += fmt.Sprintf("%s %s NOT NULL UNIQUE, \n", scanFieldDbSerial, scanFieldDbSerialType)

		query += fmt.Sprintf("%s %s, \n", scanFieldSID, scanFieldSIDType)
		query += fmt.Sprintf("%s %s, \n", scanFieldZID, scanFieldZIDType)
		query += fmt.Sprintf("%s %s, \n", scanFieldXID, scanFieldXIDType)
		query += fmt.Sprintf("%s %s, \n", scanFieldTime, scanFieldTimeType)
		query += fmt.Sprintf("%s %s, \n", scanFieldDuration, scanFieldDurationType)
		query += fmt.Sprintf("%s %s, \n", scanFieldDistance, scanFieldDistanceType)
		query += fmt.Sprintf("%s %s, \n", scanFieldBearing, scanFieldBearingType)
		query += fmt.Sprintf("%s %s, \n", scanFieldBegan, scanFieldBeganType)
		query += fmt.Sprintf("%s %s, \n", scanFieldBeganLoc, scanFieldBeganLocType)
		query += fmt.Sprintf("%s %s, \n", scanFieldBeganLocHDOP, scanFieldBeganLocHDOPType)
		query += fmt.Sprintf("%s %s, \n", scanFieldBeganLocTime, scanFieldBeganLocTimeType)
		query += fmt.Sprintf("%s %s, \n", scanFieldBeganMotionTime, scanFieldBeganMotionTimeType)
		query += fmt.Sprintf("%s %s, \n", scanFieldEnded, scanFieldEndedType)
		query += fmt.Sprintf("%s %s, \n", scanFieldEndedLoc, scanFieldEndedLocType)
		query += fmt.Sprintf("%s %s, \n", scanFieldEndedLocHDOP, scanFieldEndedLocHDOPType)
		query += fmt.Sprintf("%s %s, \n", scanFieldEndedLocTime, scanFieldEndedLocTimeType)
		query += fmt.Sprintf("%s %s, \n", scanFieldEndedMotionTime, scanFieldEndedMotionTimeType)
		query += fmt.Sprintf("%s %s, \n", scanFieldDataRAT, scanFieldDataRATType)
		query += fmt.Sprintf("%s %s, \n", scanFieldDataMCC, scanFieldDataMCCType)
		query += fmt.Sprintf("%s %s, \n", scanFieldDataMNC, scanFieldDataMNCType)
		query += fmt.Sprintf("%s %s, \n", scanFieldDataTAC, scanFieldDataTACType)
		query += fmt.Sprintf("%s %s, \n", scanFieldDataCID, scanFieldDataCIDType)
		query += fmt.Sprintf("%s %s, \n", scanFieldDataPCI, scanFieldDataPCIType)
		query += fmt.Sprintf("%s %s, \n", scanFieldDataBAND, scanFieldDataBANDType)
		query += fmt.Sprintf("%s %s, \n", scanFieldDataCHAN, scanFieldDataCHANType)
		query += fmt.Sprintf("%s %s, \n", scanFieldDataFREQ, scanFieldDataFREQType)
		query += fmt.Sprintf("%s %s, \n", scanFieldDataBSSID, scanFieldDataBSSIDType)
		query += fmt.Sprintf("%s %s, \n", scanFieldDataPSC, scanFieldDataPSCType)
		query += fmt.Sprintf("%s %s, \n", scanFieldDataRSSI, scanFieldDataRSSIType)
		query += fmt.Sprintf("%s %s, \n", scanFieldDataRSRP, scanFieldDataRSRPType)
		query += fmt.Sprintf("%s %s, \n", scanFieldDataRSRQ, scanFieldDataRSRQType)
		query += fmt.Sprintf("%s %s, \n", scanFieldDataRSCP, scanFieldDataRSCPType)
		query += fmt.Sprintf("%s %s, \n", scanFieldDataSNR, scanFieldDataSNRType)
		query += fmt.Sprintf("%s %s, \n", scanFieldDataSSID, scanFieldDataSSIDType)

		query += fmt.Sprintf("%s %s NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC') \n",
			scanFieldDbModified, scanFieldDbModifiedType)
		query += "); \n"
		_, err = db.db.Exec(query)
		if err != nil {
			return fmt.Errorf("%s table creation error: %s", tableScan, err)
		}

		// Create the scan table indexes
		query = fmt.Sprintf("CREATE INDEX ia_%s ON %s ( %s ASC );", scanFieldSID, tableScan, scanFieldSID)
		_, err = db.db.Exec(query)
		if err != nil {
			return fmt.Errorf("%s %s index creation error: %s", tableScan, scanFieldSID, err)
		}
		query = fmt.Sprintf("CREATE INDEX ia_%s ON %s ( %s ASC );", scanFieldZID, tableScan, scanFieldZID)
		_, err = db.db.Exec(query)
		if err != nil {
			return fmt.Errorf("%s %s index creation error: %s", tableScan, scanFieldZID, err)
		}
		query = fmt.Sprintf("CREATE INDEX ia_%s ON %s ( %s ASC );", scanFieldXID, tableScan, scanFieldXID)
		_, err = db.db.Exec(query)
		if err != nil {
			return fmt.Errorf("%s %s index creation error: %s", tableScan, scanFieldXID, err)
		}
		query = fmt.Sprintf("CREATE INDEX ia_%s ON %s ( %s ASC );", scanFieldTime, tableScan, scanFieldTime)
		_, err = db.db.Exec(query)
		if err != nil {
			return fmt.Errorf("%s %s index creation error: %s", tableScan, scanFieldTime, err)
		}

	}

	// Initialize the track table
	exists, err = uTableExists(db, tableTrack)
	if err != nil {
		return
	}
	if !exists {

		query := fmt.Sprintf("CREATE TABLE \"%s\" ( \n", tableTrack)
		query += fmt.Sprintf("%s %s NOT NULL UNIQUE, \n", trackFieldDbSerial, trackFieldDbSerialType)

		query += fmt.Sprintf("%s %s, \n", trackFieldSID, trackFieldSIDType)
		query += fmt.Sprintf("%s %s, \n", trackFieldZID, trackFieldZIDType)
		query += fmt.Sprintf("%s %s, \n", trackFieldTime, trackFieldTimeType)
		query += fmt.Sprintf("%s %s, \n", trackFieldLoc, trackFieldLocType)
		query += fmt.Sprintf("%s %s, \n", trackFieldLocTime, trackFieldLocTimeType)
		query += fmt.Sprintf("%s %s, \n", trackFieldLocHDOP, trackFieldLocHDOPType)
		query += fmt.Sprintf("%s %s, \n", trackFieldJourneyTime, trackFieldJourneyTimeType)
		query += fmt.Sprintf("%s %s, \n", trackFieldJourneyCount, trackFieldJourneyCountType)
		query += fmt.Sprintf("%s %s, \n", trackFieldMotionCount, trackFieldMotionCountType)
		query += fmt.Sprintf("%s %s, \n", trackFieldMotionTime, trackFieldMotionTimeType)
		query += fmt.Sprintf("%s %s, \n", trackFieldMotionDistance, trackFieldMotionDistanceType)
		query += fmt.Sprintf("%s %s, \n", trackFieldMotionBearing, trackFieldMotionBearingType)
		query += fmt.Sprintf("%s %s, \n", trackFieldMotionVelocity, trackFieldMotionVelocityType)
		query += fmt.Sprintf("%s %s, \n", trackFieldTemperature, trackFieldTemperatureType)
		query += fmt.Sprintf("%s %s, \n", trackFieldHumidity, trackFieldHumidityType)
		query += fmt.Sprintf("%s %s, \n", trackFieldPressure, trackFieldPressureType)
		query += fmt.Sprintf("%s %s, \n", trackFieldFlagUSB, trackFieldFlagUSBType)
		query += fmt.Sprintf("%s %s, \n", trackFieldFlagCharging, trackFieldFlagChargingType)
		query += fmt.Sprintf("%s %s, \n", trackFieldFlagHeartbeat, trackFieldFlagHeartbeatType)

		query += fmt.Sprintf("%s %s NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC') \n",
			trackFieldDbModified, trackFieldDbModifiedType)
		query += "); \n"
		_, err = db.db.Exec(query)
		if err != nil {
			return fmt.Errorf("%s table creation error: %s", tableTrack, err)
		}

		// Create the track table indexes
		query = fmt.Sprintf("CREATE INDEX ia_%s ON %s ( %s ASC );", trackFieldSID, tableTrack, trackFieldSID)
		_, err = db.db.Exec(query)
		if err != nil {
			return fmt.Errorf("%s %s index creation error: %s", tableTrack, trackFieldSID, err)
		}
		query = fmt.Sprintf("CREATE INDEX ia_%s ON %s ( %s ASC );", trackFieldZID, tableTrack, trackFieldZID)
		_, err = db.db.Exec(query)
		if err != nil {
			return fmt.Errorf("%s %s index creation error: %s", tableTrack, trackFieldZID, err)
		}
		query = fmt.Sprintf("CREATE INDEX ia_%s ON %s ( %s ASC );", trackFieldTime, tableTrack, trackFieldTime)
		_, err = db.db.Exec(query)
		if err != nil {
			return fmt.Errorf("%s %s index creation error: %s", tableTrack, trackFieldTime, err)
		}
	}

	// Initialize the contacts table
	exists, err = uTableExists(db, tableContact)
	if err != nil {
		return
	}
	if !exists {

		query := fmt.Sprintf("CREATE TABLE \"%s\" ( \n", tableContact)
		query += fmt.Sprintf("%s %s NOT NULL UNIQUE, \n", contactFieldDbSerial, contactFieldDbSerialType)
		query += fmt.Sprintf("%s %s NOT NULL DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'), \n",
			contactFieldDbModified, contactFieldDbModifiedType)

		query += fmt.Sprintf("%s %s, \n", contactFieldTime, contactFieldTimeType)

		query += fmt.Sprintf("%s %s, \n", contactFieldSID, contactFieldSIDType)
		query += fmt.Sprintf("%s %s, \n", contactFieldSN, contactFieldSNType)
		query += fmt.Sprintf("%s %s, \n", contactFieldName, contactFieldNameType)
		query += fmt.Sprintf("%s %s, \n", contactFieldAffiliation, contactFieldAffiliationType)
		query += fmt.Sprintf("%s %s, \n", contactFieldRole, contactFieldRoleType)
		query += fmt.Sprintf("%s %s, \n", contactFieldEmail, contactFieldEmailType)

		query += fmt.Sprintf("PRIMARY KEY (%s, %s, %s, %s, %s, %s) \n",
			contactFieldSID, contactFieldSN, contactFieldName,
			contactFieldAffiliation, contactFieldRole, contactFieldEmail)

		query += "); \n"
		_, err = db.db.Exec(query)
		if err != nil {
			return fmt.Errorf("%s table creation error: %s", tableTrack, err)
		}

		// Create the track table indexes
		query = fmt.Sprintf("CREATE INDEX ia_%s ON %s ( %s ASC );", contactFieldSID, tableContact, contactFieldSID)
		_, err = db.db.Exec(query)
		if err != nil {
			return fmt.Errorf("%s %s index creation error: %s", tableContact, trackFieldSID, err)
		}
		query = fmt.Sprintf("CREATE INDEX ia_%s ON %s ( %s ASC );", contactFieldTime, tableContact, contactFieldTime)
		_, err = db.db.Exec(query)
		if err != nil {
			return fmt.Errorf("%s %s index creation error: %s", tableContact, trackFieldTime, err)
		}

	}

	// Done
	return

}

// Acquire the context of the database
func dbContext() (db *DbDesc, err error) {

	// Exit if it's already open
	db = &radarDb
	if db.db != nil {
		return
	}

	// Lock, and check again
	dbLock.Lock()
	if db.db != nil {
		dbLock.Unlock()
		return
	}

	// Connect to the database
	// Construct the metabase connection string
	var conn strings.Builder
	conn.WriteString(fmt.Sprintf("host=%s ", Config.PostgresHost))
	conn.WriteString(fmt.Sprintf("port=%d ", Config.PostgresPort))
	conn.WriteString(fmt.Sprintf("user=%s ", Config.PostgresUsername))
	conn.WriteString(fmt.Sprintf("password=%s ", Config.PostgresPassword))
	conn.WriteString(fmt.Sprintf("dbname=%s ", Config.PostgresDatabase))
	conn.WriteString("sslmode=disable")

	// Open the database
	db.db, err = apmsql.Open("postgres", conn.String())
	if err != nil {
		db.db = nil
		dbLock.Unlock()
		return
	}
	dbLock.Unlock()

	// Make sure the connection is alive
	err = db.Ping()
	if err != nil {
		return
	}

	// Done
	return

}

// dbPing will make sure the database connection is alive. After trying
// to connect to the database, if this function gets a "connection refused", "no
// such host", or "the database system is starting up" error it will retry 29
// times over 29 seconds before giving up.
func (db *DbDesc) Ping() (err error) {

	dbLock.Lock()
	defer dbLock.Unlock()

	maxTries := 30
	for i := 0; i < maxTries; i++ {
		if i != 0 {
			time.Sleep(1 * time.Second)
		}
		err = db.db.Ping()
		if err == nil {
			break
		}
		if !strings.Contains(err.Error(), "connection refused") &&
			!strings.Contains(err.Error(), "no such host") &&
			!strings.Contains(err.Error(), "the database system is starting up") {
			break
		}
	}
	return
}

// uTableExists sees if a table exists
func uTableExists(db *DbDesc, tableName string) (exists bool, err error) {

	dbLock.Lock()
	defer dbLock.Unlock()

	var row string
	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM pg_tables WHERE tablename = '%s')", tableName)
	err = db.db.QueryRow(query).Scan(&row)
	if err != nil {
		return
	}
	if row != "true" && row != "t" {
		return
	}
	exists = true
	return
}

// TableExists sees if a table exists
func (db *DbDesc) TableExists(tableName string) (exists bool, err error) {
	dbLock.Lock()
	defer dbLock.Unlock()
	return uTableExists(db, tableName)
}

// Drop drops the table
func (db *DbDesc) Drop(tableName string) (err error) {

	dbLock.Lock()
	defer dbLock.Unlock()

	_, err = db.db.Exec(fmt.Sprintf("drop table \"%s\"", tableName))
	if err != nil {
		return
	}

	return
}

// Add or update a contact entry in the DB
func dbAddContact(deviceUID string, when int64, deviceSN string, contactName string, contactAffiliation string, contactRole string, contactEmail string) (err error) {

	// Generate the query that will replace or update the contact
	query := fmt.Sprintf("REPLACE INTO %s (", tableContact)
	query += fmt.Sprintf("%s, ", contactFieldSID)
	query += fmt.Sprintf("%s, ", contactFieldSN)
	query += fmt.Sprintf("%s, ", contactFieldName)
	query += fmt.Sprintf("%s, ", contactFieldAffiliation)
	query += fmt.Sprintf("%s, ", contactFieldRole)
	query += fmt.Sprintf("%s, ", contactFieldEmail)
	query += fmt.Sprintf("%s) VALUES (", contactFieldTime)
	query += fmt.Sprintf("'%s', ", deviceUID)
	query += fmt.Sprintf("'%s', ", deviceSN)
	query += fmt.Sprintf("'%s', ", contactName)
	query += fmt.Sprintf("'%s', ", contactAffiliation)
	query += fmt.Sprintf("'%s', ", contactRole)
	query += fmt.Sprintf("'%s', ", contactEmail)
	query += fmt.Sprintf("%d);", when)

	// Get DB context and lock
	db, err := dbContext()
	if err != nil {
		return
	}
	dbLock.Lock()
	defer dbLock.Unlock()

	// Add or replace the contact
	_, err = db.db.Exec(query)
	if err != nil {
		return fmt.Errorf("dbAddContact: %s", err)
	}

	// Done
	return

}

// Add a scan entry to the db
func dbAddScan(deviceUID string, scan RadarScan) (err error) {

	// Generate the query that will replace or update the contact
	query := fmt.Sprintf("INSERT INTO %s (", tableScan)
	query += fmt.Sprintf("%s, ", scanFieldSID)
	query += fmt.Sprintf("%s, ", scanFieldZID)
	query += fmt.Sprintf("%s, ", scanFieldXID)
	query += fmt.Sprintf("%s, ", scanFieldTime)
	query += fmt.Sprintf("%s, ", scanFieldDuration)
	query += fmt.Sprintf("%s, ", scanFieldDistance)
	query += fmt.Sprintf("%s, ", scanFieldBearing)
	query += fmt.Sprintf("%s, ", scanFieldBegan)
	query += fmt.Sprintf("%s, ", scanFieldBeganLoc)
	query += fmt.Sprintf("%s, ", scanFieldBeganLocHDOP)
	query += fmt.Sprintf("%s, ", scanFieldBeganLocTime)
	query += fmt.Sprintf("%s, ", scanFieldBeganMotionTime)
	query += fmt.Sprintf("%s, ", scanFieldEnded)
	query += fmt.Sprintf("%s, ", scanFieldEndedLoc)
	query += fmt.Sprintf("%s, ", scanFieldEndedLocHDOP)
	query += fmt.Sprintf("%s, ", scanFieldEndedLocTime)
	query += fmt.Sprintf("%s, ", scanFieldEndedMotionTime)
	query += fmt.Sprintf("%s, ", scanFieldDataRAT)
	query += fmt.Sprintf("%s, ", scanFieldDataMCC)
	query += fmt.Sprintf("%s, ", scanFieldDataMNC)
	query += fmt.Sprintf("%s, ", scanFieldDataTAC)
	query += fmt.Sprintf("%s, ", scanFieldDataCID)
	query += fmt.Sprintf("%s, ", scanFieldDataPCI)
	query += fmt.Sprintf("%s, ", scanFieldDataBAND)
	query += fmt.Sprintf("%s, ", scanFieldDataCHAN)
	query += fmt.Sprintf("%s, ", scanFieldDataFREQ)
	query += fmt.Sprintf("%s, ", scanFieldDataBSSID)
	query += fmt.Sprintf("%s, ", scanFieldDataPSC)
	query += fmt.Sprintf("%s, ", scanFieldDataRSSI)
	query += fmt.Sprintf("%s, ", scanFieldDataRSRP)
	query += fmt.Sprintf("%s, ", scanFieldDataRSRQ)
	query += fmt.Sprintf("%s, ", scanFieldDataRSCP)
	query += fmt.Sprintf("%s, ", scanFieldDataSNR)
	query += fmt.Sprintf("%s) VALUES (", scanFieldDataSSID)
	query += fmt.Sprintf("'%s', ", deviceUID)
	query += fmt.Sprintf("'%s', ", scan.ScanFieldZID)
	query += fmt.Sprintf("'%s', ", scan.ScanFieldXID)
	query += fmt.Sprintf("%d, ", scan.ScanFieldTime)
	query += fmt.Sprintf("%d, ", scan.ScanFieldDuration)
	query += fmt.Sprintf("%f, ", scan.ScanFieldDistance)
	query += fmt.Sprintf("%f, ", scan.ScanFieldBearing)
	query += fmt.Sprintf("%d, ", scan.ScanFieldBegan)
	query += fmt.Sprintf("'%s', ", scan.ScanFieldBeganLoc)
	query += fmt.Sprintf("%d, ", scan.ScanFieldBeganLocHDOP)
	query += fmt.Sprintf("%d, ", scan.ScanFieldBeganLocTime)
	query += fmt.Sprintf("%d, ", scan.ScanFieldBeganMotionTime)
	query += fmt.Sprintf("%d, ", scan.ScanFieldEnded)
	query += fmt.Sprintf("'%s', ", scan.ScanFieldEndedLoc)
	query += fmt.Sprintf("%d, ", scan.ScanFieldEndedLocHDOP)
	query += fmt.Sprintf("%d, ", scan.ScanFieldEndedLocTime)
	query += fmt.Sprintf("%d, ", scan.ScanFieldEndedMotionTime)
	query += fmt.Sprintf("'%s', ", scan.ScanFieldDataRAT)
	query += fmt.Sprintf("%d, ", scan.ScanFieldDataMCC)
	query += fmt.Sprintf("%d, ", scan.ScanFieldDataMNC)
	query += fmt.Sprintf("%d, ", scan.ScanFieldDataTAC)
	query += fmt.Sprintf("%d, ", scan.ScanFieldDataCID)
	query += fmt.Sprintf("%d, ", scan.ScanFieldDataPCI)
	query += fmt.Sprintf("%d, ", scan.ScanFieldDataBAND)
	query += fmt.Sprintf("%d, ", scan.ScanFieldDataCHAN)
	query += fmt.Sprintf("%d, ", scan.ScanFieldDataFREQ)
	query += fmt.Sprintf("'%s', ", scan.ScanFieldDataBSSID)
	query += fmt.Sprintf("%d, ", scan.ScanFieldDataPSC)
	query += fmt.Sprintf("%d, ", scan.ScanFieldDataRSSI)
	query += fmt.Sprintf("%d, ", scan.ScanFieldDataRSRP)
	query += fmt.Sprintf("%d, ", scan.ScanFieldDataRSRQ)
	query += fmt.Sprintf("%d, ", scan.ScanFieldDataRSCP)
	query += fmt.Sprintf("%d, ", scan.ScanFieldDataSNR)
	query += fmt.Sprintf("'%s');", scan.ScanFieldDataSSID)

	// Get DB context and lock
	db, err := dbContext()
	if err != nil {
		return
	}
	dbLock.Lock()
	defer dbLock.Unlock()

	// Add the record
	_, err = db.db.Exec(query)
	if err != nil {
		return fmt.Errorf("dbAddScan: %s", err)
	}
	return
}

// Add a track entry to the DB
func dbAddTrack(deviceUID string, track RadarTrack) (err error) {

	// Generate the query that will replace or update the contact
	query := fmt.Sprintf("INSERT INTO %s (", tableTrack)
	query += fmt.Sprintf("%s, ", scanFieldSID)
	query += fmt.Sprintf("%s, ", trackFieldTime)
	query += fmt.Sprintf("%s, ", trackFieldLoc)
	query += fmt.Sprintf("%s, ", trackFieldLocTime)
	query += fmt.Sprintf("%s, ", trackFieldLocHDOP)
	query += fmt.Sprintf("%s, ", trackFieldJourneyTime)
	query += fmt.Sprintf("%s, ", trackFieldJourneyCount)
	query += fmt.Sprintf("%s, ", trackFieldMotionCount)
	query += fmt.Sprintf("%s, ", trackFieldMotionTime)
	query += fmt.Sprintf("%s, ", trackFieldMotionDistance)
	query += fmt.Sprintf("%s, ", trackFieldMotionBearing)
	query += fmt.Sprintf("%s, ", trackFieldMotionVelocity)
	query += fmt.Sprintf("%s, ", trackFieldTemperature)
	query += fmt.Sprintf("%s, ", trackFieldHumidity)
	query += fmt.Sprintf("%s, ", trackFieldPressure)
	query += fmt.Sprintf("%s, ", trackFieldFlagUSB)
	query += fmt.Sprintf("%s, ", trackFieldFlagCharging)
	query += fmt.Sprintf("%s) VALUES (", trackFieldFlagHeartbeat)
	query += fmt.Sprintf("'%s', ", deviceUID)
	query += fmt.Sprintf("%d, ", track.TrackFieldTime)
	query += fmt.Sprintf("'%s', ", track.TrackFieldLoc)
	query += fmt.Sprintf("%d, ", track.TrackFieldLocTime)
	query += fmt.Sprintf("%d, ", track.TrackFieldLocHDOP)
	query += fmt.Sprintf("%d, ", track.TrackFieldJourneyTime)
	query += fmt.Sprintf("%d, ", track.TrackFieldJourneyCount)
	query += fmt.Sprintf("%d, ", track.TrackFieldMotionCount)
	query += fmt.Sprintf("%d, ", track.TrackFieldMotionTime)
	query += fmt.Sprintf("%f, ", track.TrackFieldMotionDistance)
	query += fmt.Sprintf("%f, ", track.TrackFieldMotionBearing)
	query += fmt.Sprintf("%f, ", track.TrackFieldMotionVelocity)
	query += fmt.Sprintf("%f, ", track.TrackFieldTemperature)
	query += fmt.Sprintf("%f, ", track.TrackFieldHumidity)
	query += fmt.Sprintf("%f, ", track.TrackFieldPressure)
	if track.TrackFieldFlagUSB {
		query += "1, "
	} else {
		query += "0, "
	}
	if track.TrackFieldFlagCharging {
		query += "1, "
	} else {
		query += "0, "
	}
	if track.TrackFieldFlagHeartbeat {
		query += "1 "
	} else {
		query += "0 "
	}
	query += ");"

	// Get DB context and lock
	db, err := dbContext()
	if err != nil {
		return
	}
	dbLock.Lock()
	defer dbLock.Unlock()

	// Add the record
	_, err = db.db.Exec(query)
	if err != nil {
		return fmt.Errorf("dbAddTrack: %s", err)
	}
	return

}
