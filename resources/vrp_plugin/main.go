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
	"context"
	"fmt"
	"git.liero.se/opentelco/go-dnc/models/protobuf/metric"
	"git.liero.se/opentelco/go-swpx/proto/traffic_policy"
	"github.com/pkg/errors"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"git.liero.se/opentelco/go-dnc/client"
	"git.liero.se/opentelco/go-dnc/models/protobuf/transport"

	"git.liero.se/opentelco/go-swpx/shared/oids"

	"git.liero.se/opentelco/go-swpx/proto/networkelement"
	proto "git.liero.se/opentelco/go-swpx/proto/resource"
	"git.liero.se/opentelco/go-swpx/shared"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"
	"github.com/nats-io/nats.go"
)

var VERSION *version.Version

var logger hclog.Logger

const (
	VersionBase = "1.0-beta"
	DriverName  = "vrp-driver"
	Float64Size = 64
)

var reFindIndexinOID = regexp.MustCompile("(\\d+)$") // used to get the last number of the oid

type discoveryItem struct {
	index int
	descr string
	alias string
}

var dncChan chan string

func init() {
	var err error
	if VERSION, err = version.NewVersion(VersionBase); err != nil {
		log.Fatal(err)
	}
}

// Here is a real implementation of Driver
type VRPDriver struct {
	logger hclog.Logger
	dnc    client.Client
	conf   shared.Configuration
}

func (d *VRPDriver) GetConfiguration(ctx context.Context) (shared.Configuration, error) {
	return d.conf, nil
}

func (d *VRPDriver) SetConfiguration(ctx context.Context, conf shared.Configuration) error {
	d.conf = conf

	return nil
}

func (d *VRPDriver) Version() (string, error) {
	d.logger.Debug("message from resource-driver running at version", VERSION.String())
	dncChan <- VERSION.String()
	return fmt.Sprintf("%s@%s", DriverName, VERSION.String()), nil
}

// parse a map of description/alias and return the ID
func (d VRPDriver) parseDescriptionToIndex(port string, discoveryMap map[int]*discoveryItem) (*discoveryItem, error) {
	for _, v := range discoveryMap {
		// v.index = k
		if v.descr == port {
			d.logger.Info("parser found match", "port", v.descr)
			return v, nil
		}
	}
	return nil, fmt.Errorf("%s was not found on network element", port)
}

// Find matching OID for port
func (d *VRPDriver) MapEntityPhysical(ctx context.Context, el *proto.NetworkElement) (*proto.NetworkElementInterfaces, error) {
	conf := shared.Proto2conf(el.Conf)
	portMsg := createPortInformationMsg(el, conf)
	msg, err := d.dnc.Put(ctx, portMsg)
	if err != nil {
		d.logger.Error(err.Error())
		return nil, err
	}

	switch task := msg.Task.(type) {
	case *transport.Message_Snmpc:
		interfaces := make(map[string]*proto.NetworkElementInterface)
		for _, m := range task.Snmpc.Metrics {
			fields := strings.Split(m.Oid, ".")
			index, err := strconv.Atoi(fields[len(fields)-1])
			if err != nil {
				logger.Error("can't convert phys.port to int: ", err.Error())
				return nil, err
			}

			interfaces[m.GetStringValue()] = &proto.NetworkElementInterface{
				Alias:       m.Name,
				Index:       int64(index),
				Description: m.GetStringValue(),
			}
		}

		return &proto.NetworkElementInterfaces{Interfaces: interfaces}, nil
	}
	return nil, errors.Errorf("Unsupported message type")
}

func (d *VRPDriver) GetTransceiverInformation(ctx context.Context, el *proto.NetworkElement) (*networkelement.Transceiver, error) {
	conf := shared.Proto2conf(el.Conf)

	vrpMsg := createVRPTransceiverMsg(el, conf)
	msg, err := d.dnc.Put(ctx, vrpMsg)
	if err != nil {
		d.logger.Error(err.Error())
		return nil, err
	}

	switch task := msg.Task.(type) {
	case *transport.Message_Snmpc:
		if len(task.Snmpc.Metrics) >= 7 {

			tempInt := task.Snmpc.Metrics[1].GetIntValue()
			voltInt := task.Snmpc.Metrics[2].GetIntValue()
			curInt := task.Snmpc.Metrics[3].GetIntValue()
			rxInt := task.Snmpc.Metrics[4].GetIntValue()
			txInt := task.Snmpc.Metrics[5].GetIntValue()
			var rx, tx, temp, volt, curr float64
			rx = float64(rxInt*-1) / 100
			tx = float64(txInt*-1) / 100
			temp = float64(tempInt)
			volt = float64(voltInt) / 1000
			curr = float64(curInt) / 1000

			// no transceiver available, return nil
			if tempInt == -255 && rxInt == -1 && txInt == -1 {
				return nil, nil
			}

			val := &networkelement.Transceiver{
				SerialNumber: strings.Trim(task.Snmpc.Metrics[0].GetStringValue(), " "),
				Stats: []*networkelement.TransceiverStatistics{
					{
						Temp:    temp,
						Voltage: volt,
						Current: curr,
						Rx:      rx,
						Tx:      tx,
					},
				},
				PartNumber: task.Snmpc.Metrics[6].GetStringValue(),
			}
			return val, nil
		}
	}
	return nil, errors.Errorf("Unsupported message type")
}

func (d *VRPDriver) MapInterface(ctx context.Context, el *proto.NetworkElement) (*proto.NetworkElementInterfaces, error) {
	d.logger.Info("got a task to determine what index and name this interface has", "host", el.Hostname, "ip", el.Ip, "interface", el.Interface)
	var msg *transport.Message
	discoveryMap := make(map[int]*discoveryItem)
	var index int
	var interfaces = make(map[string]*proto.NetworkElementInterface)

	conf := shared.Proto2conf(el.Conf)

	msg = createDiscoveryMsg(el, conf)
	msg, err := d.dnc.Put(ctx, msg)
	if err != nil {
		d.logger.Error(err.Error())
		return nil, err
	}

	switch task := msg.Task.(type) {
	case *transport.Message_Snmpc:
		d.logger.Debug("the msg returns from dnc", "status", msg.Status.String(), "completed", msg.Completed.String(), "execution_time", msg.ExecutionTime.String(), "size", len(task.Snmpc.Metrics))
		d.populateDiscoveryMap(task, index, discoveryMap)

		for _, v := range discoveryMap {
			interfaces[v.descr] = &proto.NetworkElementInterface{
				Index:       int64(v.index),
				Description: v.descr,
				Alias:       v.alias,
			}
		}
	}

	return &proto.NetworkElementInterfaces{Interfaces: interfaces}, nil
}

func (d *VRPDriver) populateDiscoveryMap(task *transport.Message_Snmpc, index int, discoveryMap map[int]*discoveryItem) {
	for _, m := range task.Snmpc.Metrics {
		index, _ = strconv.Atoi(reFindIndexinOID.FindString(m.Oid))
		switch m.GetName() {
		case "ifIndex":
			if val, ok := discoveryMap[index]; ok {
				val.index = int(m.GetIntValue())
			} else {
				discoveryMap[index] = &discoveryItem{
					index: int(m.GetIntValue()),
				}
			}
		case "ifAlias":
			if val, ok := discoveryMap[index]; ok {
				val.alias = m.GetStringValue()
			} else {
				discoveryMap[index] = &discoveryItem{
					descr: m.GetStringValue(),
				}
			}
		case "ifDescr":
			if val, ok := discoveryMap[index]; ok {
				val.descr = m.GetStringValue()
			} else {
				discoveryMap[index] = &discoveryItem{
					descr: m.GetStringValue(),
				}
			}
		}
	}
}

func (d *VRPDriver) TechnicalPortInformation(ctx context.Context, el *proto.NetworkElement) (*networkelement.Element, error) {
	dncChan <- "ok"
	d.logger.Info("running technical port info", "host", el.Hostname, "ip", el.Ip, "interface", el.Interface)

	conf := shared.Proto2conf(el.Conf)

	var msgs []*transport.Message
	if el.InterfaceIndex != 0 {
		msgs = append(msgs, createSinglePortMsg(el.InterfaceIndex, el, conf))
		msgs = append(msgs, createTaskGetPortStats(el.InterfaceIndex, el, conf))
	} else {
		msgs = append(msgs, createMsg(conf))
	}

	msgs = append(msgs, createTaskSystemInfo(el, conf))
	msgs = append(msgs, createTelnetInterfaceTask(el, conf))

	ne := &networkelement.Element{}
	ne.Hostname = el.Hostname
	elementInterface := &networkelement.Interface{
		Stats: &networkelement.InterfaceStatistics{
			Input:  &networkelement.InterfaceStatisticsInput{},
			Output: &networkelement.InterfaceStatisticsOutput{},
		},
	}
	var err error

	for _, msg := range msgs {
		d.logger.Debug("sending msg")
		msg, err = d.dnc.Put(ctx, msg)
		if err != nil {
			d.logger.Error(err.Error())
			return nil, err
		}

		switch task := msg.Task.(type) {
		case *transport.Message_Snmpc:
			d.logger.Debug("the msg returns from dnc", "status", msg.Status.String(), "completed", msg.Completed.String(), "execution_time", msg.ExecutionTime.String(), "size", len(task.Snmpc.Metrics))

			elementInterface.Index = el.InterfaceIndex

			for _, m := range task.Snmpc.Metrics {
				d.logger.Debug(m.GetStringValue())
				switch {
				case strings.HasPrefix(m.Oid, oids.SysPrefix):
					getSystemInformation(m, ne)

				case strings.HasPrefix(m.Oid, oids.HuaPrefix):
					getHuaweiInformation(m, elementInterface)

				case strings.HasPrefix(m.Oid, oids.IfEntryPrefix):
					getIfEntryInformation(m, elementInterface)

				case strings.HasPrefix(m.Oid, oids.IfXEntryPrefix):
					getIfXEntryInformation(m, elementInterface)
				}
			}
		case *transport.Message_Telnet:
			if elementInterface.MacAddressTable, err = ParseMacTable(task.Telnet.Payload[0].Lookfor); err != nil {
				logger.Error("can't parse MAC address table: ", err.Error())
				return nil, err
			}

			if elementInterface.DhcpTable, err = ParseIPTable(task.Telnet.Payload[1].Lookfor); err != nil {
				logger.Error("can't parse DHCP table: ", err.Error())
				return nil, err
			}
			elementInterface.Config = parseCurrentConfig(task.Telnet.Payload[2].Lookfor)

			if elementInterface.ConfiguredTrafficPolicy, err = parsePolicy(task.Telnet.Payload[3].Lookfor); err != nil {
				logger.Error("can't parse policy: ", err.Error())
				return nil, err
			}

			if err = parsePolicyStatistics(elementInterface.ConfiguredTrafficPolicy, task.Telnet.Payload[4].Lookfor); err != nil {
				logger.Error("can't parse policy statistics: ", err.Error())
				return nil, err
			}
		}
	}

	ne.Interfaces = append(ne.Interfaces, elementInterface)

	return ne, nil
}

func parsePolicyStatistics(policy *traffic_policy.ConfiguredTrafficPolicy, data string) error {
	lines := strings.Split(data, "\r\n")

	statistics := &traffic_policy.ConfiguredTrafficPolicy_Statistics{
		Classifiers: make(map[string]*traffic_policy.ConfiguredTrafficPolicy_Statistics_Classifier),
	}

	if err := parseStatisticsHeader(statistics, lines); err != nil {
		return err
	}

	parseMetrics(lines, statistics)

	policy.InboundStatistics = statistics

	return nil
}

func parseMetrics(lines []string, statistics *traffic_policy.ConfiguredTrafficPolicy_Statistics) {
	var classifierName string
	for i := 7; i < len(lines)-1; {
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

func parsePolicy(data string) (*traffic_policy.ConfiguredTrafficPolicy, error) {
	policy := &traffic_policy.ConfiguredTrafficPolicy{}

	lines := strings.Split(data, "\r\n")

	for _, line := range lines {
		fields := strings.Fields(line)
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

func parseCurrentConfig(config string) string {
	configStart := strings.Index(config, "#\r\n") + 1
	configEnd := strings.LastIndex(config, "#\r\n")

	return config[configStart:configEnd]
}

func getIfXEntryInformation(m *metric.Metric, elementInterface *networkelement.Interface) {

	switch {
	case strings.HasPrefix(m.Oid, oids.IfOutUcastPkts):
		elementInterface.Stats.Output.Unicast = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfInBroadcastPkts):
		elementInterface.Stats.Input.Broadcast = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfInMulticastPkts):
		elementInterface.Stats.Input.Multicast = m.GetIntValue()
	case strings.HasPrefix(m.Oid, oids.IfOutBroadcastPkts):
		elementInterface.Stats.Output.Broadcast = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfOutMulticastPkts):
		elementInterface.Stats.Output.Multicast = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfAlias):
		elementInterface.Alias = m.GetStringValue()

	case strings.HasPrefix(m.Oid, oids.IfHighSpeed):
		elementInterface.Speed = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfConnectorPresent):
		elementInterface.ConnectorPresent = m.GetBoolValue()
	}

}

func getIfEntryInformation(m *metric.Metric, elementInterface *networkelement.Interface) {
	switch {
	case strings.HasPrefix(m.Oid, oids.IfOutOctets):
		elementInterface.Stats.Output.Bytes = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfInOctets):
		elementInterface.Stats.Input.Bytes = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfInUcastPkts):
		elementInterface.Stats.Input.Unicast = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfInErrors):
		elementInterface.Stats.Input.Errors = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfOutErrors):
		elementInterface.Stats.Output.Errors = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfDescr):
		elementInterface.Description = m.GetStringValue()

	case strings.HasPrefix(m.Oid, oids.IfType):
		elementInterface.Type = networkelement.InterfaceType(m.GetIntValue())

	case strings.HasPrefix(m.Oid, oids.IfMtu):
		elementInterface.Mtu = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfLastChange):
		elementInterface.LastChanged = m.GetTimestampValue()

	case strings.HasPrefix(m.Oid, oids.IfPhysAddress):
		elementInterface.Hwaddress = m.GetStringValue()

	case strings.HasPrefix(m.Oid, oids.IfOperStatus):
		elementInterface.AdminStatus = networkelement.InterfaceStatus(m.GetIntValue())

	case strings.HasPrefix(m.Oid, oids.IfAdminStatus):
		elementInterface.OperationalStatus = networkelement.InterfaceStatus(m.GetIntValue())

	}
}

func getHuaweiInformation(m *metric.Metric, elementInterface *networkelement.Interface) {
	switch {
	case strings.HasPrefix(m.Oid, oids.HuaIfEtherStatInCRCPkts):
		elementInterface.Stats.Input.CrcErrors = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.HuaIfEtherStatInPausePkts):
		elementInterface.Stats.Input.Pauses = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.HuaIfEthIfStatReset):
		elementInterface.Stats.Resets = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.HuaIfEtherStatOutPausePkts):
		elementInterface.Stats.Output.Pauses = m.GetIntValue()
	}
}

func getSystemInformation(m *metric.Metric, ne *networkelement.Element) {
	switch m.Oid {
	case oids.SysContact:
		ne.Contact = m.GetStringValue()
	case oids.SysDescr:
		ne.Version = m.GetStringValue()
	case oids.SysLocation:
		ne.Location = m.GetStringValue()
	case oids.SysName:
		ne.Sysname = m.GetStringValue()
	// case oids.SysORLastChange:
	// case oids.SysObjectID:
	case oids.SysUpTime:
		ne.Uptime = m.GetStringValue()
	}
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.

func main() {
	logger = hclog.New(&hclog.LoggerOptions{
		Name:       fmt.Sprintf("%s@%s", DriverName, VERSION.String()),
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})
	natsConf := shared.GetConfig().NATS
	nc, _ := nats.Connect(strings.Join(natsConf.EventServers, ","))
	dncChan = make(chan string)
	enc, _ := nats.NewEncodedConn(nc, "json")
	enc.BindSendChan("vrp-driver", dncChan)

	logger.Debug("message", "message from resource-driver", "version", VERSION.String())
	//dncClient, err := client.NewGRPC(DISPATCHER_ADDR)
	dncClient, err := client.NewNATS(strings.Join(natsConf.EventServers, ","))
	if err != nil {
		log.Fatal(err)
	}
	driver := &VRPDriver{
		logger: logger,
		dnc:    dncClient,
		conf: shared.Configuration{
			SNMP: shared.ConfigSNMP{
				Community:          "semipublic",
				Version:            2,
				Timeout:            time.Second * 20,
				Retries:            2,
				DynamicRepetitions: true,
			},
			Telnet: shared.ConfigTelnet{},
		},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			shared.PluginResourceKey: &shared.ResourcePlugin{Impl: driver},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
