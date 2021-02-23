/*
 * File: analysis.go
 * Project: sait
 * File Created: Sunday, 14th February 2021 2:19:48 pm
 * Author: Mathias Ehrlin (mathias.ehrlin@vx.se)
 * -----
 * Last Modified: Sunday, 14th February 2021 9:11:06 pm
 * Modified By: Mathias Ehrlin (mathias.ehrlin@vx.se>)
 * -----
 * Copyright - 2021 VX Service Delivery AB
 *
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 * -----
 */

package main

import (
	"fmt"

	"git.liero.se/opentelco/go-swpx/proto/go/analysis"
	"git.liero.se/opentelco/go-swpx/proto/go/networkelement"
)

const (
	TransceiverRXThreshold = -5.00
	TransceiverTXThreshold = -5.00
)

func doAnalysis(n *networkelement.Element, changes *int) {

	for _, i := range n.Interfaces {
		doStatusAnalysis(i, changes)

		doStatsAnalysis(i, changes)

		doTransceiverAnalysis(i, changes)
	}
}

func doStatusAnalysis(iface *networkelement.Interface, changes *int) {
	if iface.Analysis == nil {
		iface.Analysis = make([]*analysis.Analysis, 0)
	}

	if iface.AdminStatus != networkelement.InterfaceStatus_up {
		iface.Analysis = append(iface.Analysis, &analysis.Analysis{
			Level:                analysis.Analysis_FAILURE,
			Message:              "interface is not administrative up",
			Value:                iface.AdminStatus.String(),
			Threshold:            networkelement.InterfaceStatus_up.String(),
			Type:                 "admin-status",
		})
		*changes++
	}

	if iface.OperationalStatus != networkelement.InterfaceStatus_up {
		iface.Analysis = append(iface.Analysis, &analysis.Analysis{
			Level:                analysis.Analysis_FAILURE,
			Message:              "interface is not operational up",
			Value:                iface.OperationalStatus.String(),
			Threshold:            networkelement.InterfaceStatus_up.String(),
			Type:                 "operational-status",
		})
		*changes++
	}
}

func doStatsAnalysis(iface *networkelement.Interface, changes *int) {
	if iface.Stats.Analysis == nil {
		iface.Stats.Analysis = make([]*analysis.Analysis, 0)
	}

	if iface.Stats.Input.Errors > 0 {
		iface.Stats.Analysis = append(iface.Stats.Analysis, &analysis.Analysis{
			Level:                analysis.Analysis_WARNING,
			Message:              fmt.Sprintf("%d input errors", iface.Stats.Input.Errors),
			Value:                fmt.Sprintf("%d", iface.Stats.Input.Errors),
			Threshold:            "0",
			Type:                 "stats-input-errors",
		})
		*changes++
	}

	if iface.Stats.Input.CrcErrors > 0 {
		iface.Stats.Analysis = append(iface.Stats.Analysis, &analysis.Analysis{
			Level:                analysis.Analysis_WARNING,
			Message:              fmt.Sprintf("%d input crc errors", iface.Stats.Input.CrcErrors),
			Value:                fmt.Sprintf("%d", iface.Stats.Input.CrcErrors),
			Threshold:            "0",
			Type:                 "stats-input-crc-errors",
		})
	}

	if iface.Stats.Output.Errors > 0 {
		iface.Stats.Analysis = append(iface.Stats.Analysis, &analysis.Analysis{
			Level:                analysis.Analysis_WARNING,
			Message:              fmt.Sprintf("%d output errors", iface.Stats.Output.Errors),
			Value:                fmt.Sprintf("%d", iface.Stats.Output.Errors),
			Threshold:            "0",
			Type:                 "stats-output-errors",
		})
	}
}

func doTransceiverAnalysis(iface *networkelement.Interface, changes *int) {
	if iface.Transceiver == nil || iface.Transceiver.Stats == nil {
		return
	}
	for _, stats := range iface.Transceiver.Stats {
		if iface.Transceiver.Analysis == nil {
			iface.Transceiver.Analysis = make([]*analysis.Analysis, 0)
		}

		switch rx := stats.Rx; {
		// check if threshold for RX
		case rx <= TransceiverRXThreshold:
			iface.Transceiver.Analysis = append(iface.Transceiver.Analysis, &analysis.Analysis{
				Level:     analysis.Analysis_WARNING,
				Message:   "sfp rx value is below threshold",
				Value:     fmt.Sprintf("%.2f", rx),
				Threshold: fmt.Sprintf("%.2f", TransceiverRXThreshold),
				Type:      "rx-value",
			})
			*changes++
		}
		// check if threshold for TX
		switch tx := stats.Tx; {
		case tx <= TransceiverTXThreshold:
			iface.Transceiver.Analysis = append(iface.Transceiver.Analysis, &analysis.Analysis{
				Level:     analysis.Analysis_WARNING,
				Message:   "sfp TX value is below threshold",
				Value:     fmt.Sprintf("%.2f", tx),
				Threshold: fmt.Sprintf("%.2f", TransceiverTXThreshold),
				Type:      "tx-value",
			})
			*changes++
		}
	}
}