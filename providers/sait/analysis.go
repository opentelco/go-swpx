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

	for ix, i := range n.Interfaces {
		if i.Transceiver == nil || i.Transceiver.Stats == nil {
			continue
		}
		for _, stats := range i.Transceiver.Stats {
			if n.Interfaces[ix].Transceiver.Analysis == nil {
				n.Interfaces[ix].Transceiver.Analysis = make([]*analysis.Analysis, 0)
			}

			switch rx := stats.Rx; {
			// check if threshold for RX
			case rx <= TransceiverRXThreshold:
				n.Interfaces[ix].Transceiver.Analysis = append(n.Interfaces[ix].Transceiver.Analysis, &analysis.Analysis{
					Level:     analysis.Analysis_WARNING,
					Message:   "sfp rx value is below threshold",
					Value:     fmt.Sprintf("%.2f", rx),
					Threshold: fmt.Sprintf("%.2f", TransceiverRXThreshold),
					Type:      "rx-value",
				})
				*changes++
				// check if threshold for TX

			}
			switch tx := stats.Tx; {
			case tx <= TransceiverTXThreshold:
				n.Interfaces[ix].Transceiver.Analysis = append(n.Interfaces[ix].Transceiver.Analysis, &analysis.Analysis{
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

}
