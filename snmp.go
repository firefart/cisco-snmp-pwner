package main

import (
	"fmt"
	"math/rand"
	"strconv"

	g "github.com/gosnmp/gosnmp"
)

// http://www.circitor.fr/Mibs/Html/C/CISCO-CONFIG-COPY-MIB.php

const ciscoConfigMib = ".1.3.6.1.4.1.9.9.96.1.1.1.1"

func randomID() string {
	low := 10000
	high := 99999
	i := low + rand.Intn(high-low)
	return strconv.Itoa(i)
}

func getDumpRequest(randomID, target, filename string) []g.SnmpPDU {
	return []g.SnmpPDU{
		{
			Name:  fmt.Sprintf("%s.2.%s", ciscoConfigMib, randomID), // Protocol
			Type:  g.Integer,
			Value: 1, // tftp(1), ftp(2), rcp(3), scp(4), sftp(5)
		},
		{
			Name:  fmt.Sprintf("%s.3.%s", ciscoConfigMib, randomID), // Source File Type
			Type:  g.Integer,
			Value: 4, // networkFile(1), iosFile(2), startupConfig(3), runningConfig(4), terminal(5), fabricStartupConfig(6)
		},
		{
			Name:  fmt.Sprintf("%s.4.%s", ciscoConfigMib, randomID), // Dest File Type
			Type:  g.Integer,
			Value: 1, // networkFile(1), iosFile(2), startupConfig(3), runningConfig(4), terminal(5), fabricStartupConfig(6)
		},
		{
			Name:  fmt.Sprintf("%s.5.%s", ciscoConfigMib, randomID), // Server Address
			Type:  g.IPAddress,
			Value: target,
		},
		{
			Name:  fmt.Sprintf("%s.6.%s", ciscoConfigMib, randomID), // Filename
			Type:  g.OctetString,
			Value: filename,
		},
		{
			Name:  fmt.Sprintf("%s.14.%s", ciscoConfigMib, randomID), // EntryRowStatus http://www.circitor.fr/Mibs/Html/C/CISCO-CONFIG-COPY-MIB.php#ccCopyEntryRowStatus
			Type:  g.Integer,
			Value: 4, // 	active(1), notInService(2), notReady(3), createAndGo(4), createAndWait(5), destroy(6)
		},
	}
}

func getMergeRequest(randomID, target, filename string) []g.SnmpPDU {
	return []g.SnmpPDU{
		{
			Name:  fmt.Sprintf("%s.2.%s", ciscoConfigMib, randomID), // Protocol
			Type:  g.Integer,
			Value: 1, // tftp(1), ftp(2), rcp(3), scp(4), sftp(5)
		},
		{
			Name:  fmt.Sprintf("%s.3.%s", ciscoConfigMib, randomID), // Source File Type
			Type:  g.Integer,
			Value: 1, // networkFile(1), iosFile(2), startupConfig(3), runningConfig(4), terminal(5), fabricStartupConfig(6)
		},
		{
			Name:  fmt.Sprintf("%s.4.%s", ciscoConfigMib, randomID), // Dest File Type
			Type:  g.Integer,
			Value: 4, // networkFile(1), iosFile(2), startupConfig(3), runningConfig(4), terminal(5), fabricStartupConfig(6)
		},
		{
			Name:  fmt.Sprintf("%s.5.%s", ciscoConfigMib, randomID), // Server Address
			Type:  g.IPAddress,
			Value: target,
		},
		{
			Name:  fmt.Sprintf("%s.6.%s", ciscoConfigMib, randomID), // Filename
			Type:  g.OctetString,
			Value: filename,
		},
		{
			Name:  fmt.Sprintf("%s.14.%s", ciscoConfigMib, randomID), // EntryRowStatus http://www.circitor.fr/Mibs/Html/C/CISCO-CONFIG-COPY-MIB.php#ccCopyEntryRowStatus
			Type:  g.Integer,
			Value: 4, // 	active(1), notInService(2), notReady(3), createAndGo(4), createAndWait(5), destroy(6)
		},
	}
}

func getStatusRequest(randomID string) []string {
	return []string{fmt.Sprintf("%s.10.%s", ciscoConfigMib, randomID)} // http://www.circitor.fr/Mibs/Html/C/CISCO-CONFIG-COPY-MIB.php#ccCopyState
}

// nolint:deadcode,unused
func getFailReasonRequest(randomID string) []string {
	return []string{fmt.Sprintf("%s.13.%s", ciscoConfigMib, randomID)} // http://www.circitor.fr/Mibs/Html/C/CISCO-CONFIG-COPY-MIB.php#ccCopyFailCause
}

func getDeleteJobRequest(randomID string) []g.SnmpPDU {
	return []g.SnmpPDU{
		{
			Name:  fmt.Sprintf("%s.14.%s", ciscoConfigMib, randomID), // http://www.circitor.fr/Mibs/Html/C/CISCO-CONFIG-COPY-MIB.php#ccCopyEntryRowStatus
			Type:  g.Integer,
			Value: 6, // 	active(1), notInService(2), notReady(3), createAndGo(4), createAndWait(5), destroy(6)
		},
	}
}
