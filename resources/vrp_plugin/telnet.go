/*
 * Copyright (c) 2020. Liero AB
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
	"git.liero.se/opentelco/go-swpx/proto/networkelement"
	"git.liero.se/opentelco/go-swpx/proto/traffic_policy"
	"regexp"
	"strconv"
	"strings"
)

const (
	MacRegex = "^([[:xdigit:]]{2}[:.-]?){5}[[:xdigit:]]{2}$"
	IPRegex  = "(\\b(?:\\d{1,3}\\.){3}\\d{1,3}\\b)\\s+([0-9A-Fa-f]{4}[-][a-f0-9A-F]{4}[-][a-f0-9A-F]{4})\\s+([0-9]{1,3}).*([1-9][0-9]{3}.[0-9]{2}.[0-9]{2}-[0-9]{2}:[0-9]{2})"

	// The number of lines in a traffic policy statistics header.
	StatisticsHeaderLength = 7
	// The number of lines for each traffic policy statistics metric.
	StatisticsMetricLength = 25
)

func parseMacTable(data string) ([]*networkelement.MACEntry, error) {
	dataRows := strings.Split(data, "\n")
	rows := make([]*networkelement.MACEntry, 0)

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

		rows = append(rows, &networkelement.MACEntry{
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

func parseIPTable(data string) ([]*networkelement.DHCPEntry, error) {
	dataRows := strings.Split(data, "\n")
	rows := make([]*networkelement.DHCPEntry, 0)

	for _, row := range dataRows {
		fields := strings.Fields(row)

		// skip rows without IP data
		if !isIPAddressRow(fields) {
			continue
		}

		vlan, err := getVLAN(fields)
		if err != nil || vlan == 0 {
			return nil, fmt.Errorf("can't convert vlan string to int: %v", err)
		}

		rows = append(rows, &networkelement.DHCPEntry{
			IpAddress:       fields[0],
			HardwareAddress: fields[1],
			Vlan:            int64(vlan),
			Timestamp:       fields[6],
		})
	}

	return rows, nil
}

// looks for populated VLAN field
func getVLAN(fields []string) (int, error) {
	var err error

	for i := 2; i < 5; i++ {
		if fields[i] == "/--" {
			continue
		}
		val, err := strconv.Atoi(fields[i])

		if err == nil && val != 0 {
			return val, nil
		}
	}

	return 0, err
}

func isIPAddressRow(fields []string) bool {
	if len(fields) != 7 {
		return false
	}

	match, parseErr := regexp.Match(IPRegex, []byte(strings.Join(fields, " ")))
	if parseErr != nil {
		return false
	}

	return match
}

func parseCurrentConfig(config string) string {
	configStart := strings.Index(config, "#\r\n") + 1
	configEnd := strings.LastIndex(config, "#\r\n")

	if configStart == 0 || configEnd == -1 {
		return ""
	}

	return config[configStart:configEnd]
}

func parsePolicyStatistics(policy *traffic_policy.ConfiguredTrafficPolicy, data string) error {
	data = strings.Replace(data, "\r", "", -1) // remove line feeds
	lines := strings.Split(data, "\n")

	statistics := &traffic_policy.ConfiguredTrafficPolicy_Statistics{
		Classifiers: make(map[string]*traffic_policy.ConfiguredTrafficPolicy_Statistics_Classifier),
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

func parseStatisticsHeader(statistics *traffic_policy.ConfiguredTrafficPolicy_Statistics, lines []string) error {
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

func parseMetrics(lines []string, statistics *traffic_policy.ConfiguredTrafficPolicy_Statistics) {
	var classifierName string
	for i := StatisticsHeaderLength; i < len(lines)-1; {
		if strings.HasPrefix(lines[i], "-") {
			if strings.HasPrefix(lines[i+1], " Classifier:") {
				classifierName = strings.Split(lines[i+1], "Classifier: ")[1]
				statistics.Classifiers[classifierName] = &traffic_policy.ConfiguredTrafficPolicy_Statistics_Classifier{
					Classifier: classifierName,
					Behavior:   strings.Split(lines[i+2], "Behavior: ")[1],
					Board:      strings.Split(lines[i+3], "Board : ")[1],
					Metrics:    make(map[string]*traffic_policy.ConfiguredTrafficPolicy_Statistics_Classifier_Metric),
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
				metric := &traffic_policy.ConfiguredTrafficPolicy_Statistics_Classifier_Metric{
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

func parsePolicy(data string) (*traffic_policy.ConfiguredTrafficPolicy, error) {
	policy := &traffic_policy.ConfiguredTrafficPolicy{}

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

			policy.Qos = &traffic_policy.ConfiguredTrafficPolicy_QOS{
				Queue: int64(queue),
				Shaping: &traffic_policy.ConfiguredTrafficPolicy_QOS_Shaping{
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

func parseQueueStatistics(data string) (*traffic_policy.QOS, error) {
	data = strings.Replace(data, ",", "", -1)
	lines := strings.Split(data, "\n")
	qos := &traffic_policy.QOS{
		QueueStatistics: make([]*traffic_policy.QOS_QueueStatistics, len(lines)/QueueEntryLength),
	}

	for i := 2; i < len(lines)-1; i += QueueEntryLength {
		if len(lines) <= i+10 {
			logger.Warn("malformed output in queue statistics command")
			continue
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
			return nil, err
		}

		qos.QueueStatistics[i/QueueEntryLength] = &traffic_policy.QOS_QueueStatistics{
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
