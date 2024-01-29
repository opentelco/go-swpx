package core

import (
	"fmt"
	"slices"
	"sort"
	"time"

	"go.opentelco.io/go-swpx/proto/go/analysispb"
	"go.opentelco.io/go-swpx/proto/go/corepb"
	"go.opentelco.io/go-swpx/proto/go/devicepb"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

var coreActivities = Activities{}

func RunDiagnosticWorkflow(ctx workflow.Context, request *corepb.RunDiagnosticRequest) (*analysispb.Report, error) {

	logger := workflow.GetLogger(ctx)
	logger.Info("Diagnostic workflow started")
	defer logger.Info("Diagnostic workflow ended")

	var availabilityResp corepb.CheckAvailabilityResponse
	if err := workflow.ExecuteActivity(
		workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
			TaskQueue:           TemporalTaskQueue,
			StartToCloseTimeout: time.Duration(60 * time.Second),
			RetryPolicy: &temporal.RetryPolicy{
				MaximumAttempts: 2,
			},
		}),
		coreActivities.CheckAvailability,
		request.Session,
	).Get(ctx, &availabilityResp); err != nil {
		return nil, err
	}

	responses := make([]*corepb.PollResponse, request.PollTimes)

	for x := 0; x < int(request.PollTimes); x++ {
		var pollResponse corepb.PollResponse
		if err := workflow.ExecuteActivity(
			workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
				TaskQueue:           TemporalTaskQueue,
				StartToCloseTimeout: time.Duration(80 * time.Second),
				RetryPolicy: &temporal.RetryPolicy{
					MaximumAttempts: 1,
				},
			}),
			coreActivities.Poll,
			&corepb.PollRequest{
				Session:  request.Session,
				Settings: request.Settings,
				Type:     corepb.PollRequest_GET_TECHNICAL_INFO,
			},
		).Get(ctx, &pollResponse); err != nil {
			return nil, err
		}
		responses[x] = &pollResponse

		// if the poll is not the last one, sleep for 10 seconds
		if x != int(request.PollTimes)-1 {
			workflow.Sleep(ctx, time.Second*10)
		}
	}

	report := &analysispb.Report{}

	for n, resp := range responses {
		if resp.Device == nil {
			return nil, fmt.Errorf("diagnosis report was missing for array element: %d", n)
		}
	}

	// analyze the link
	r, err := analyzeLink(request.Session.Port, responses)
	if err != nil {
		return nil, err
	}
	report.Analysis = append(report.Analysis, r...)

	// analyze the errors
	r, err = analyzeErrors(request.Session.Port, responses)
	if err != nil {
		return nil, err
	}
	report.Analysis = append(report.Analysis, r...)

	// analyze the traffic
	r, err = analyzeTraffic(request.Session.Port, responses)
	if err != nil {
		return nil, err
	}
	report.Analysis = append(report.Analysis, r...)

	// analyze the transceiver
	r, err = analyzeTransceiver(request.Session.Port, responses)
	if err != nil {
		return nil, err
	}
	report.Analysis = append(report.Analysis, r...)

	// analyze the mac address table
	r, err = analyzeMacAddressTable(request.Session.Port, responses)
	if err != nil {
		return nil, err
	}
	report.Analysis = append(report.Analysis, r...)

	// analyze the dhcp table
	r, err = analyzeDhcpTable(request.Session.Port, responses)
	if err != nil {
		return nil, err
	}
	report.Analysis = append(report.Analysis, r...)
	report.Type = analysispb.Report_TYPE_DETAILED

	return report, nil
}

const (
	linkThreshold = "Up/Up"
)

// has the link changed between the polls (gone from up to down, or down to up)
// has all polls been up
// has all polls been down
func analyzeLink(port string, data []*corepb.PollResponse) ([]*analysispb.Analysis, error) {
	var anyalysis []*analysispb.Analysis
	linkOK := []devicepb.Port_Status{devicepb.Port_up, devicepb.Port_up}
	linkDown := []devicepb.Port_Status{devicepb.Port_up, devicepb.Port_down}
	linkShut := []devicepb.Port_Status{devicepb.Port_down, devicepb.Port_down}

	var statues [][]devicepb.Port_Status

	for _, d := range data {
		p, err := getPortFromElement(port, d.Device)
		if err != nil {
			return nil, err
		}
		if p == nil {
			return nil, fmt.Errorf("no port found for %s", port)
		}

		statues = append(statues, []devicepb.Port_Status{p.AdminStatus, p.OperationalStatus})
	}

	/*
	 check if the link is up through out the whole diagnosis
	*/
	var linkIsDown bool
	for _, s := range statues {

		if slices.Compare(s, linkOK) != 0 {
			linkIsDown = true

		}
	}
	if !linkIsDown {
		anyalysis = append(anyalysis, &analysispb.Analysis{
			Result:    analysispb.Analysis_RESULT_OK,
			Type:      analysispb.Analysis_TYPE_LINK,
			Note:      "Link is up on the port throughout the whole diagnosis",
			Value:     []string{"Up/Up"},
			Threshold: linkThreshold,
		})
		return anyalysis, nil
	}

	/*
		check if link is down through out the whole diagnosis
	*/
	var linkIsUp bool
	for _, s := range statues {
		if slices.Compare(s, linkDown) != 0 {
			linkIsUp = true

		}
	}
	if !linkIsUp {
		anyalysis = append(anyalysis, &analysispb.Analysis{
			Result:    analysispb.Analysis_RESULT_ERROR,
			Type:      analysispb.Analysis_TYPE_LINK,
			Note:      "Link is down on the port through out the whole diagnosis",
			Value:     []string{"Up/Down"},
			Threshold: linkThreshold,
		})
		return anyalysis, nil
	}

	/*
		check if the link has been shut throughout the whole diagnosis
	*/
	var linkIsShut bool
	for _, s := range statues {
		if slices.Compare(s, linkShut) != 0 {
			linkIsShut = true

		}
	}
	if !linkIsShut {
		anyalysis = append(anyalysis, &analysispb.Analysis{
			Result:    analysispb.Analysis_RESULT_ERROR,
			Type:      analysispb.Analysis_TYPE_LINK,
			Note:      "Link has been shut throughout the whole diagnosis, please check the port configuration",
			Value:     []string{"Down/Down"},
			Threshold: linkThreshold,
		})
		return anyalysis, nil
	}
	/*
		link has changed between the polls
	*/

	anyalysis = append(anyalysis, &analysispb.Analysis{
		Result:    analysispb.Analysis_RESULT_WARNING,
		Type:      analysispb.Analysis_TYPE_LINK,
		Note:      "Link has been changing state under the diagnosis",
		Value:     []string{"Up/Down", "Up/Up"},
		Threshold: linkThreshold,
	})
	return anyalysis, nil
}

// check if the errors has increased during the diagnosis period, if so return error
// check if the errors are above 0, if so return warning
// check if the errors are 0, if so return ok
func analyzeErrors(port string, data []*corepb.PollResponse) ([]*analysispb.Analysis, error) {
	var anyalysis []*analysispb.Analysis

	var (
		inputErrors  = make([]int64, len(data))
		outputErrors = make([]int64, len(data))
	)

	for n, d := range data {
		p, err := getPortFromElement(port, d.Device)
		if err != nil {
			return nil, err
		}

		if p.Stats == nil {
			continue
		}

		// if the port has no stats, set the errors to 0
		if p.Stats.Input == nil {
			inputErrors[n] = 0
		} else {
			inputErrors[n] = p.Stats.Input.Errors
		}

		// if the port has no stats, set the errors to 0
		if p.Stats.Output == nil {
			outputErrors[n] = 0
		} else {
			outputErrors[n] = p.Stats.Output.Errors
		}
	}

	//
	if biggest(inputErrors) == 0 && biggest(outputErrors) == 0 {
		anyalysis = append(anyalysis, &analysispb.Analysis{
			Result:    analysispb.Analysis_RESULT_OK,
			Type:      analysispb.Analysis_TYPE_LINK_ERROR,
			Note:      "no errors on the port",
			Value:     []string{"0"},
			Threshold: "0",
		})
		return anyalysis, nil
	}

	if hasIncreasingErrors(inputErrors) {
		anyalysis = append(anyalysis, &analysispb.Analysis{
			Result:    analysispb.Analysis_RESULT_ERROR,
			Type:      analysispb.Analysis_TYPE_LINK_ERROR,
			Note:      "input errors has increased on the port during the diagnosis",
			Value:     []string{fmt.Sprintf("%d", biggest(inputErrors))},
			Threshold: "0",
		})
	}
	if hasIncreasingErrors(outputErrors) {
		anyalysis = append(anyalysis, &analysispb.Analysis{
			Result:    analysispb.Analysis_RESULT_ERROR,
			Type:      analysispb.Analysis_TYPE_LINK_ERROR,
			Note:      "output errors has increased on the port during the diagnosis",
			Value:     []string{fmt.Sprintf("%d", biggest(outputErrors))},
			Threshold: "0",
		})
	}

	// if errors are increasing, no need to check if they are above 0
	if len(anyalysis) > 0 {
		return anyalysis, nil
	}

	// check if we have had any errors (ever) could be old errors.
	if biggest(inputErrors) > 0 {
		anyalysis = append(anyalysis, &analysispb.Analysis{
			Result:    analysispb.Analysis_RESULT_WARNING,
			Type:      analysispb.Analysis_TYPE_LINK_ERROR,
			Note:      "input errors on the port",
			Value:     []string{fmt.Sprintf("%d", biggest(inputErrors))},
			Threshold: "0",
		})
	}

	if biggest(outputErrors) > 0 {
		anyalysis = append(anyalysis, &analysispb.Analysis{
			Result:    analysispb.Analysis_RESULT_WARNING,
			Type:      analysispb.Analysis_TYPE_LINK_ERROR,
			Note:      "output errors on the port",
			Value:     []string{fmt.Sprintf("%d", biggest(outputErrors))},
			Threshold: "0",
		})
	}

	return anyalysis, nil
}

const (
	thresholdInputPkt  = 200
	thresholdOutputPkt = 200
)

// check if the traffic has increased with atleast X bits and return OK if it has
// if no data has been received or sent on the port, return warning
func analyzeTraffic(port string, data []*corepb.PollResponse) ([]*analysispb.Analysis, error) {
	var anyalysis []*analysispb.Analysis

	var (
		input  = make([]int64, len(data))
		output = make([]int64, len(data))
	)

	for n, d := range data {
		p, err := getPortFromElement(port, d.Device)
		if err != nil {
			return nil, err
		}

		if p.Stats == nil {
			continue
		}

		// if the port has no stats, set the data to 0
		if p.Stats.Input == nil {
			input[n] = 0
		} else {
			input[n] = p.Stats.Input.Bits
		}

		// if the port has no stats, set the data to 0
		if p.Stats.Output == nil {
			output[n] = 0
		} else {
			output[n] = p.Stats.Output.Bits
		}
	}

	if hasIncreasing(input, thresholdInputPkt) {
		anyalysis = append(anyalysis, &analysispb.Analysis{
			Result:    analysispb.Analysis_RESULT_OK,
			Type:      analysispb.Analysis_TYPE_TRAFFIC,
			Note:      fmt.Sprintf("input traffic has increased with more than %d bits", thresholdInputPkt),
			Value:     []string{fmt.Sprintf("%d", biggest(input))},
			Threshold: fmt.Sprintf("%d", thresholdInputPkt),
		})
	}

	if hasIncreasing(output, thresholdOutputPkt) {
		anyalysis = append(anyalysis, &analysispb.Analysis{
			Result:    analysispb.Analysis_RESULT_OK,
			Type:      analysispb.Analysis_TYPE_TRAFFIC,
			Note:      fmt.Sprintf("output traffic has increased with more than %d bits", thresholdOutputPkt),
			Value:     []string{fmt.Sprintf("%d", biggest(input))},
			Threshold: fmt.Sprintf("%d", thresholdOutputPkt),
		})
	}

	if len(anyalysis) == 0 {
		inputTraffic := biggest(input) - smallest(input)
		anyalysis = append(anyalysis, &analysispb.Analysis{
			Result:    analysispb.Analysis_RESULT_WARNING,
			Type:      analysispb.Analysis_TYPE_TRAFFIC,
			Note:      "no (or not enough) input traffic on the port",
			Value:     []string{fmt.Sprintf("%d bits", inputTraffic)},
			Threshold: fmt.Sprintf("%d", thresholdInputPkt),
		})
		outputTraffic := biggest(output) - smallest(output)
		anyalysis = append(anyalysis, &analysispb.Analysis{
			Result:    analysispb.Analysis_RESULT_WARNING,
			Type:      analysispb.Analysis_TYPE_TRAFFIC,
			Note:      "no (or not enough) output traffic on the port",
			Value:     []string{fmt.Sprintf("%d bits", outputTraffic)},
			Threshold: fmt.Sprintf("%d", thresholdOutputPkt),
		})
	}

	return anyalysis, nil
}

const (
	transceiverRXThreshold = -20.0
	transceiverTXThreshold = -20.0
	transceiverFluctuation = 5.0
)

// Analyzez the transceiver diagnostics for the port
// check if the RX and TX are within the threshold (comparing an average of all polls)
func analyzeTransceiver(port string, data []*corepb.PollResponse) ([]*analysispb.Analysis, error) {
	var anyalysis []*analysispb.Analysis

	var (
		rx = make([]float64, len(data))
		tx = make([]float64, len(data))
	)

	for n, d := range data {
		p, err := getPortFromElement(port, d.Device)
		if err != nil {
			return nil, err
		}

		if p.Transceiver == nil {
			continue
		}

		rx[n] = p.Transceiver.Stats.Rx
		tx[n] = p.Transceiver.Stats.Tx

	}

	switch avgRx := average(rx); {
	case avgRx == 0:
	case avgRx < transceiverRXThreshold:

		anyalysis = append(anyalysis, &analysispb.Analysis{
			Result:    analysispb.Analysis_RESULT_WARNING,
			Type:      analysispb.Analysis_TYPE_TRANSCEIVER_DIAGNOSTICS,
			Note:      fmt.Sprintf("the average RX (%3.f) is below threshold", avgRx),
			Value:     float64ToString(rx),
			Threshold: fmt.Sprintf("%3.f", transceiverTXThreshold),
		})
	default: // no errors or warnings
		anyalysis = append(anyalysis, &analysispb.Analysis{
			Result:    analysispb.Analysis_RESULT_OK,
			Type:      analysispb.Analysis_TYPE_TRANSCEIVER_DIAGNOSTICS,
			Note:      fmt.Sprintf("the average RX (%3.f) is within threshold", avgRx),
			Value:     float64ToString(rx),
			Threshold: fmt.Sprintf("%3.f", transceiverTXThreshold),
		})

	}

	switch avgTx := average(tx); {
	case avgTx == 0:
	case avgTx < transceiverTXThreshold:
		anyalysis = append(anyalysis, &analysispb.Analysis{
			Result:    analysispb.Analysis_RESULT_WARNING,
			Type:      analysispb.Analysis_TYPE_TRANSCEIVER_DIAGNOSTICS,
			Note:      fmt.Sprintf("the average TX (%3.f) is below threshold", avgTx),
			Value:     float64ToString(tx),
			Threshold: fmt.Sprintf("%3.f", transceiverTXThreshold),
		})
	default:
		anyalysis = append(anyalysis, &analysispb.Analysis{
			Result:    analysispb.Analysis_RESULT_OK,
			Type:      analysispb.Analysis_TYPE_TRANSCEIVER_DIAGNOSTICS,
			Note:      fmt.Sprintf("the average TX (%3.f) is within threshold", avgTx),
			Value:     float64ToString(tx),
			Threshold: fmt.Sprintf("%3.f", transceiverTXThreshold),
		})
	}

	if len(anyalysis) == 0 {
		anyalysis = append(anyalysis, &analysispb.Analysis{
			Result: analysispb.Analysis_RESULT_OK,
			Type:   analysispb.Analysis_TYPE_TRANSCEIVER_DIAGNOSTICS,
			Note:   "transceiver diagnostic is ok",
		})
	}

	return anyalysis, nil
}

// if no mac address, return error
func analyzeMacAddressTable(port string, data []*corepb.PollResponse) ([]*analysispb.Analysis, error) {
	var anyalysis []*analysispb.Analysis

	var (
		macAddressTable = make([]*devicepb.MACEntry, 0)
	)

	for _, d := range data {
		p, err := getPortFromElement(port, d.Device)
		if err != nil {
			return nil, err
		}

		if p.MacAddressTable == nil {
			continue
		}

		for _, entry := range p.MacAddressTable {
			macAddressTable = appendMacEntry(macAddressTable, entry)
		}

	}

	if len(macAddressTable) == 0 {
		anyalysis = append(anyalysis, &analysispb.Analysis{
			Result: analysispb.Analysis_RESULT_ERROR,
			Type:   analysispb.Analysis_TYPE_MAC_ADDRESS,
			Note:   "no mac address found on port",
		})
	} else {
		anyalysis = append(anyalysis, &analysispb.Analysis{
			Result: analysispb.Analysis_RESULT_OK,
			Type:   analysispb.Analysis_TYPE_MAC_ADDRESS,
			Note:   fmt.Sprintf("found %d mac address on port", len(macAddressTable)),
		})
	}

	return anyalysis, nil
}

// if no ip lease, return error
func analyzeDhcpTable(port string, data []*corepb.PollResponse) ([]*analysispb.Analysis, error) {
	var anyalysis []*analysispb.Analysis

	var (
		dhcpEntryTable = make([]*devicepb.DHCPEntry, 0)
	)

	for _, d := range data {
		p, err := getPortFromElement(port, d.Device)
		if err != nil {
			return nil, err
		}

		if p.DhcpTable == nil {
			continue
		}

		for _, entry := range p.DhcpTable {
			dhcpEntryTable = appendDHCPEntry(dhcpEntryTable, entry)
		}

	}

	if len(dhcpEntryTable) == 0 {
		anyalysis = append(anyalysis, &analysispb.Analysis{
			Result: analysispb.Analysis_RESULT_ERROR,
			Type:   analysispb.Analysis_TYPE_DHCP_LEASE,
			Note:   "no mac address with IP Lease found on port",
		})
	} else {
		anyalysis = append(anyalysis, &analysispb.Analysis{
			Result: analysispb.Analysis_RESULT_OK,
			Type:   analysispb.Analysis_TYPE_DHCP_LEASE,
			Note:   fmt.Sprintf("found %d mac address with IP lease on port", len(dhcpEntryTable)),
		})
	}
	return anyalysis, nil
}

// from the array of interfaces get the affected for the diagnosis
func getPortFromElement(port string, data *devicepb.Device) (*devicepb.Port, error) {
	for _, i := range data.Ports {
		if i.Description == port {
			return i, nil
		}

	}
	return nil, fmt.Errorf("no port found for %s", port)
}

// calculate the average of a slice of float64
func average(slice []float64) float64 {
	total := 0.0
	if len(slice) == 0 {
		return 0
	}

	for _, value := range slice {
		total += value
	}
	return total / float64(len(slice))
}

// slice of float64 to slice of "-3.222" string
func float64ToString(slice []float64) []string {
	var strSlice []string
	for _, value := range slice {
		strSlice = append(strSlice, fmt.Sprintf("%.3f", value))
	}
	return strSlice
}

// check if last input errors is bigger than the first one
func hasIncreasingErrors(slice []int64) bool {
	if len(slice) < 2 {
		return false
	}
	return slice[len(slice)-1] > slice[0]
}

// check if the last value is bigger than the first one and that it has increased with atleast X
func hasIncreasing(slice []int64, fluctuation int64) bool {
	if len(slice) < 2 {
		return false
	}
	return slice[len(slice)-1] > slice[0] && slice[len(slice)-1]-slice[0] >= fluctuation
}

// return biggest int from a slice of ints
func biggest(slice []int64) int64 {
	sort.Slice(slice, func(i, j int) bool { return slice[i] < slice[j] })
	return slice[len(slice)-1]
}

// return smallest int from a slice of ints
func smallest(slice []int64) int64 {
	sort.Slice(slice, func(i, j int) bool { return slice[i] < slice[j] })
	return slice[0]
}

// append *devicepb.MACEntry to the slice if it does not exist ( by mac address )
func appendMacEntry(slice []*devicepb.MACEntry, entry *devicepb.MACEntry) []*devicepb.MACEntry {
	for _, s := range slice {
		if s.HardwareAddress == entry.HardwareAddress {
			return slice
		}
	}
	return append(slice, entry)
}

// append *devicepb.DHCPEntry to the slice if it does not exist ( by mac address )
func appendDHCPEntry(slice []*devicepb.DHCPEntry, entry *devicepb.DHCPEntry) []*devicepb.DHCPEntry {
	for _, s := range slice {
		if s.HardwareAddress == entry.HardwareAddress {
			return slice
		}
	}

	return append(slice, entry)
}
