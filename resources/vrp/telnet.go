/*
 * Copyright (c) 2023. Liero AB
 *
 * Permission is hereby granted, free of charge, to any person obtaining
 * a copy of this software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation
 * the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software
 * is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
 * EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES
 * OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
 * IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
 * CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
 * TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"git.liero.se/opentelco/go-swpx/proto/go/networkelementpb"
	"git.liero.se/opentelco/go-swpx/proto/go/trafficpolicypb"
)

const (
	MacRegex = "^([[:xdigit:]]{2}[:.-]?){5}[[:xdigit:]]{2}$"
	IPRegex  = "(\\b(?:\\d{1,3}\\.){3}\\d{1,3}\\b)\\s+([0-9A-Fa-f]{4}[-][a-f0-9A-F]{4}[-][a-f0-9A-F]{4})\\s+([0-9]{1,3}).*([1-9][0-9]{3}.[0-9]{2}.[0-9]{2}-[0-9]{2}:[0-9]{2})"

	// The number of lines in a traffic policy statistics header.
	StatisticsHeaderLength = 7
	// The number of lines for each traffic policy statistics metricpb.
	StatisticsMetricLength = 25
)

var (
	reDhcpTableEntry = regexp.MustCompile(`(?P<ipAddress>(^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}))\s+(?P<macAddress>([0-9aA-fF]{1,4}\-[0-9aA-fF]{1,4}\-[0-9aA-fF]{1,4}))\s+(?P<outerVlan>[0-9]{1,4})\s+/(?P<innerVlan>[0-9-]{1,4})\s+/(?P<mappedVlan>[0-9-]{1,4})\s+(?P<ifAlias>\w+[0-9]/[0-9]/[0-9]{1,2})\s+(?P<timestamp>.+)$`)
)

func parseMacTable(data string) ([]*networkelementpb.MACEntry, error) {
	if data == "" {
		return nil, fmt.Errorf("no data found")
	}

	dataRows := strings.Split(data, "\n")
	rows := make([]*networkelementpb.MACEntry, 0)

	for _, row := range dataRows {
		fields := strings.Fields(row)

		// skip rows without mac data
		if !isMacAddressRow(fields) {
			continue
		}
		vlan, err := strconv.Atoi(strings.TrimRight(fields[1], "/-"))
		if err != nil {
			logger.Error("can't convert VLAN for mac entry table: ", err.Error())
			return nil, err
		}

		rows = append(rows, &networkelementpb.MACEntry{
			HardwareAddress: fields[0],
			Vlan:            int64(vlan),
		})
	}

	return rows, nil
}

func isMacAddressRow(fields []string) bool {
	if len(fields) != 4 {
		return false
	}

	match, parseErr := regexp.Match(MacRegex, []byte(fields[0]))
	if parseErr != nil {
		return false
	}

	return match
}

func parseIPTable(data string) ([]*networkelementpb.DHCPEntry, error) {
	if data == "" {
		return nil, fmt.Errorf("no data found")
	}

	dataRows := strings.Split(data, "\n")
	rows := make([]*networkelementpb.DHCPEntry, 0)

	for _, row := range dataRows {

		matches := reDhcpTableEntry.FindStringSubmatch(row)
		if len(matches) > 1 {
			ip := matches[reDhcpTableEntry.SubexpIndex("ipAddress")]
			mac := matches[reDhcpTableEntry.SubexpIndex("macAddress")]
			vlanStr := matches[reDhcpTableEntry.SubexpIndex("outerVlan")]
			timestamp := matches[reDhcpTableEntry.SubexpIndex("timestamp")]
			vlan, _ := strconv.Atoi(vlanStr)

			rows = append(rows, &networkelementpb.DHCPEntry{
				IpAddress:       ip,
				HardwareAddress: mac,
				Vlan:            int64(vlan),
				Timestamp:       timestamp,
			})

		}

	}

	return rows, nil
}

func parseCurrentConfig(config string) string {
	configStart := strings.Index(config, "#\r\n") + 1
	configEnd := strings.LastIndex(config, "#\r\n")

	if configStart == 0 || configEnd == -1 {
		return ""
	}

	return config[configStart:configEnd]
}

func parsePolicyStatistics(policy *trafficpolicypb.ConfiguredTrafficPolicy, data string) error {
	if data == "" {
		return fmt.Errorf("no data found")
	}

	data = strings.Replace(data, "\r", "", -1) // remove line feeds
	lines := strings.Split(data, "\n")

	statistics := &trafficpolicypb.ConfiguredTrafficPolicy_Statistics{
		Classifiers: make(map[string]*trafficpolicypb.ConfiguredTrafficPolicy_Statistics_Classifier),
	}

	if !policyStatsOutputValid(lines) {
		return errors.New("output for policy statistics is malformed - skipping")
	}

	if err := parseStatisticsHeader(statistics, lines); err != nil {
		return err
	}

	parseMetrics(lines, statistics)

	policy.InboundStatistics = statistics

	return nil
}

func parseStatisticsHeader(statistics *trafficpolicypb.ConfiguredTrafficPolicy_Statistics, lines []string) error {
	statistics.TrafficPolicy = strings.Split(lines[3], ": ")[1]

	rulenumber, err := strconv.Atoi(strings.Split(lines[4], ": ")[1])
	if err != nil {
		return err
	}
	statistics.RuleNumber = int64(rulenumber)
	statistics.Status = strings.Split(lines[5], ": ")[1]
	interval, err := strconv.Atoi(strings.Split(lines[6], ": ")[1])
	if err != nil {
		return err
	}
	statistics.RuleNumber = int64(rulenumber)
	statistics.Interval = int64(interval)

	return nil
}

func parseMetrics(lines []string, statistics *trafficpolicypb.ConfiguredTrafficPolicy_Statistics) {
	var classifierName string
	for i := StatisticsHeaderLength; i < len(lines)-1; {
		if strings.HasPrefix(lines[i], "-") {
			if strings.HasPrefix(lines[i+1], " Classifier:") {
				classifierName = strings.Split(lines[i+1], "Classifier: ")[1]
				statistics.Classifiers[classifierName] = &trafficpolicypb.ConfiguredTrafficPolicy_Statistics_Classifier{
					Classifier: classifierName,
					Behavior:   strings.Split(lines[i+2], "Behavior: ")[1],
					Board:      strings.Split(lines[i+3], "Board : ")[1],
					Metrics:    make(map[string]*trafficpolicypb.ConfiguredTrafficPolicy_Statistics_Classifier_Metric),
				}
				i += 3
			}
			i++
		}

		var metricName string
		for !strings.HasPrefix(lines[i], "-") && i < len(lines)-1 {
			fields := strings.Fields(lines[i])

			if len(fields) == 4 {
				metricName = fields[0] //passed, dropped etc
				metric := &trafficpolicypb.ConfiguredTrafficPolicy_Statistics_Classifier_Metric{
					Values: make(map[string]float64),
				}
				statistics.Classifiers[classifierName].Metrics[metricName] = metric
			}
			metricKey := fields[len(fields)-2]
			metricValue, _ := strconv.ParseFloat(strings.Replace(fields[len(fields)-1], ",", "", -1), Float64Size)

			statistics.Classifiers[classifierName].Metrics[metricName].Values[metricKey] = metricValue

			i++
		}
	}
}

func parsePolicy(data string) (*trafficpolicypb.ConfiguredTrafficPolicy, error) {
	if data == "" {
		return nil, fmt.Errorf("no data found")
	}

	policy := &trafficpolicypb.ConfiguredTrafficPolicy{}

	data = strings.Replace(data, "\r", "", -1) // remove line feeds
	lines := strings.Split(data, "\n")

	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) <= 11 {
			logger.Warn("malformed output in policy command")
			continue
		}
		if strings.Contains(line, "inbound") {
			policy.Inbound = fields[1]
		}

		if strings.Contains(line, "outbound") {
			policy.Outbound = fields[1]
		}

		if strings.Contains(line, "qos") {
			queue, _ := strconv.Atoi(fields[2])
			cir, _ := strconv.ParseFloat(fields[5], Float64Size)
			pir, _ := strconv.ParseFloat(fields[7], Float64Size)
			cbs, _ := strconv.ParseFloat(fields[9], Float64Size)
			pbs, _ := strconv.ParseFloat(fields[11], Float64Size)

			policy.Qos = &trafficpolicypb.ConfiguredTrafficPolicy_QOS{
				Queue: int64(queue),
				Shaping: &trafficpolicypb.ConfiguredTrafficPolicy_QOS_Shaping{
					Cir: cir,
					Pir: pir,
					Cbs: cbs,
					Pbs: pbs,
				},
			}
		}
	}

	return policy, nil
}

func parseQueueStatistics(data string) (*trafficpolicypb.QOS, error) {
	if data == "" {
		return nil, fmt.Errorf("no data found")
	}

	data = strings.Replace(data, ",", "", -1)
	lines := strings.Split(data, "\n")
	qos := &trafficpolicypb.QOS{
		QueueStatistics: make([]*trafficpolicypb.QOS_QueueStatistics, len(lines)/QueueEntryLength),
	}

	for i := 2; i < len(lines)-1; i += QueueEntryLength {
		if len(lines) <= i+10 {
			return nil, fmt.Errorf("malformed output in queue statistics command")
		}

		id, err := parseQOSLineInt(lines[i])
		cir, err := parseQOSLineFloat(lines[i+1])
		pir, err := parseQOSLineFloat(lines[i+2])
		passedPackets, err := parseQOSLineInt(lines[i+3])
		passedRatePps, err := parseQOSLineFloat(lines[i+4])
		passedBytes, err := parseQOSLineInt(lines[i+5])
		passedRateBps, err := parseQOSLineFloat(lines[i+6])
		droppedPackets, err := parseQOSLineInt(lines[i+7])
		droppedRatePps, err := parseQOSLineFloat(lines[i+8])
		droppedBytes, err := parseQOSLineInt(lines[i+9])
		droppedRateBps, err := parseQOSLineFloat(lines[i+10])

		if err != nil {
			return nil, fmt.Errorf("error parsing qos output: %v", err)
		}

		qos.QueueStatistics[i/QueueEntryLength] = &trafficpolicypb.QOS_QueueStatistics{
			QueueId:        id,
			Cir:            cir,
			Pir:            pir,
			PassedPackets:  passedPackets,
			PassedRatePps:  passedRatePps,
			PassedBytes:    passedBytes,
			PassedRateBps:  passedRateBps,
			DroppedPackets: droppedPackets,
			DroppedRatePps: droppedRatePps,
			DroppedBytes:   droppedBytes,
			DroppedRateBps: droppedRateBps,
		}
	}

	return qos, nil
}

func parseQOSLineInt(line string) (int64, error) {
	fields := strings.Fields(line)

	val, err := strconv.ParseInt(fields[len(fields)-1], 10, Float64Size)
	if err != nil {
		return 0, err
	}

	return val, nil
}

func parseQOSLineFloat(line string) (float64, error) {
	fields := strings.Fields(line)

	val, err := strconv.ParseFloat(fields[len(fields)-1], Float64Size)
	if err != nil {
		return 0, err
	}

	return val, nil

}

func policyStatsOutputValid(lines []string) bool {
	return (len(lines)-StatisticsHeaderLength-2)%StatisticsMetricLength == 0
}
