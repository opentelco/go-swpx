/*
 * File: analysispb.go
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

	"git.liero.se/opentelco/go-swpx/proto/go/analysispb"
	"git.liero.se/opentelco/go-swpx/proto/go/networkelementpb"
)

const (
	TransceiverRXThreshold = -5.00
	TransceiverTXThreshold = -5.00
)

func doAnalysis(n *networkelementpb.Element, changes *int) {

	for _, i := range n.Interfaces {
		doStatusAnalysis(i, changes)

		doStatsAnalysis(i, changes)

		doTransceiverAnalysis(i, changes)
	}
}

func doStatusAnalysis(iface *networkelementpb.Port, changes *int) {
	if iface.Analysis == nil {
		iface.Analysis = make([]*analysispb.Analysis, 0)
	}

	if iface.AdminStatus != networkelementpb.Port_up {
		iface.Analysis = append(iface.Analysis, &analysispb.Analysis{
			Note: "interface is not administrative up",
			// Value: iface.AdminStatus.String(),
			// Code:  analysispb.Analysis_CODE_ERR_LINK_DOWN,
		})
		*changes++
	}

	if iface.OperationalStatus != networkelementpb.Port_up {
		iface.Analysis = append(iface.Analysis, &analysispb.Analysis{
			Note: "interface is not operational up",
			// Level: analysispb.Analysis_LEVEL_ERROR,
			// Value: iface.OperationalStatus.String(),
			// Code:  analysispb.Analysis_CODE_ERR_LINK_DOWN,
		})
		*changes++
	}
}

func doStatsAnalysis(iface *networkelementpb.Port, changes *int) {
	if iface.Stats.Analysis == nil {
		iface.Stats.Analysis = make([]*analysispb.Analysis, 0)
	}

	if iface.Stats.Input.Errors > 0 {
		iface.Stats.Analysis = append(iface.Stats.Analysis, &analysispb.Analysis{
			Note: fmt.Sprintf("%d input errors", iface.Stats.Input.Errors),
			// Value: fmt.Sprintf("%d", iface.Stats.Input.Errors),
			// Level: analysispb.Analysis_LEVEL_WARNING,
		})
		*changes++
	}

	if iface.Stats.Input.CrcErrors > 0 {
		iface.Stats.Analysis = append(iface.Stats.Analysis, &analysispb.Analysis{
			Note: fmt.Sprintf("%d input crc errors", iface.Stats.Input.CrcErrors),
			// Level: analysispb.Analysis_LEVEL_WARNING,
			// Value: fmt.Sprintf("%d", iface.Stats.Input.CrcErrors),
		})
	}

	if iface.Stats.Output.Errors > 0 {
		iface.Stats.Analysis = append(iface.Stats.Analysis, &analysispb.Analysis{
			Note: fmt.Sprintf("%d output errors", iface.Stats.Output.Errors),
			// Value: fmt.Sprintf("%d", iface.Stats.Output.Errors),
			// Level: analysispb.Analysis_LEVEL_WARNING,
		})
	}
}

func doTransceiverAnalysis(iface *networkelementpb.Port, changes *int) {
	if iface.Transceiver == nil || iface.Transceiver.Stats == nil {
		return
	}
	for _, stats := range iface.Transceiver.Stats {
		if iface.Transceiver.Analysis == nil {
			iface.Transceiver.Analysis = make([]*analysispb.Analysis, 0)
		}

		switch rx := stats.Rx; {
		// check if threshold for RX
		case rx <= TransceiverRXThreshold:
			iface.Transceiver.Analysis = append(iface.Transceiver.Analysis, &analysispb.Analysis{
				Note: "sfp rx value is below threshold",
				// Value: fmt.Sprintf("%.2f", rx),
				// Level: analysispb.Analysis_LEVEL_WARNING,
			})
			*changes++
		}
		// check if threshold for TX
		switch tx := stats.Tx; {
		case tx <= TransceiverTXThreshold:
			iface.Transceiver.Analysis = append(iface.Transceiver.Analysis, &analysispb.Analysis{
				Note: "sfp TX value is below threshold",
				// Level: analysispb.Analysis_LEVEL_WARNING,
				// Value: fmt.Sprintf("%.2f", tx),
			})
			*changes++
		}
	}
}
