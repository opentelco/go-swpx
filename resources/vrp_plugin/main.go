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
	"git.liero.se/opentelco/go-dnc/client"
	"git.liero.se/opentelco/go-dnc/models/protobuf/transport"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

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
	VersionBase      = "1.0-beta"
	DriverName       = "vrp-driver"
	Float64Size      = 64
	QueueEntryLength = 12
)

var reFindIndexinOID = regexp.MustCompile("(\\d+)$") // used to get the last number of the oid

type discoveryItem struct {
	index       int
	descr       string
	alias       string
	ifType      int
	mtu         int
	physAddress string
	adminStatus int
	operStatus  int
	lastChange  *timestamppb.Timestamp
	highSpeed   int
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
		d.logger.Error("transceiver put error", err.Error())
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
				return &networkelement.Transceiver{}, nil
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
		d.populateDiscoveryMap(task, discoveryMap)

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

func (d *VRPDriver) AllPortInformation(ctx context.Context, el *proto.NetworkElement) (*networkelement.Element, error) {
	dncChan <- "ok"
	d.logger.Info("running ALL port info", "host", el.Hostname, "ip", el.Ip, "interface", el.Interface)
	conf := shared.Proto2conf(el.Conf)
	ne := &networkelement.Element{}
	ne.Hostname = el.Hostname

	sysInfoMsg := createTaskSystemInfo(el, conf)
	sysInfoMsg, err := d.dnc.Put(ctx, sysInfoMsg)
	if err != nil {
		d.logger.Error(err.Error())
		return nil, err
	}

	sysInfoTask := sysInfoMsg.Task.(*transport.Message_Snmpc)
	for _, m := range sysInfoTask.Snmpc.Metrics {
		if strings.HasPrefix(m.Oid, oids.SysPrefix) {
			getSystemInformation(m, ne)
		}
	}

	portsMsg := createAllPortsMsg(el, conf)
	portsMsg, err = d.dnc.Put(ctx, portsMsg)
	if err != nil {
		d.logger.Error(err.Error())
		return nil, err
	}

	var task *transport.Message_Snmpc
	var ok bool
	if task, ok = portsMsg.Task.(*transport.Message_Snmpc); ok {
		discoveryMap := make(map[int]*discoveryItem)
		d.populateDiscoveryMap(task, discoveryMap)

		for _, discoveryItem := range discoveryMap {
			ne.Interfaces = append(ne.Interfaces, itemToInterface(discoveryItem))
		}

		sort.Slice(ne.Interfaces, func(i, j int) bool {
			return ne.Interfaces[i].Description < ne.Interfaces[j].Description
		})
	}

	return ne, nil
}

func (d *VRPDriver) TechnicalPortInformation(ctx context.Context, el *proto.NetworkElement) (*networkelement.Element, error) {
	dncChan <- "ok"
	d.logger.Info("running technical port info", "host", el.Hostname, "ip", el.Ip, "interface", el.Interface)
	errs := make([]*networkelement.TransientError, 0)

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
			if elementInterface.MacAddressTable, err = parseMacTable(task.Telnet.Payload[0].Lookfor); err != nil {
				errs = d.logAndAppend(err, errs, task.Telnet.Payload[0].Command)
			}

			if elementInterface.DhcpTable, err = parseIPTable(task.Telnet.Payload[1].Lookfor); err != nil {
				errs = d.logAndAppend(err, errs, task.Telnet.Payload[1].Command)
			}
			elementInterface.Config = parseCurrentConfig(task.Telnet.Payload[2].Lookfor)

			if elementInterface.ConfiguredTrafficPolicy, err = parsePolicy(task.Telnet.Payload[3].Lookfor); err != nil {
				errs = d.logAndAppend(err, errs, task.Telnet.Payload[3].Command)
			}

			if err = parsePolicyStatistics(elementInterface.ConfiguredTrafficPolicy, task.Telnet.Payload[4].Lookfor); err != nil {
				errs = d.logAndAppend(err, errs, task.Telnet.Payload[4].Command)
			}

			if elementInterface.Qos, err = parseQueueStatistics(task.Telnet.Payload[5].Lookfor); err != nil {
				errs = d.logAndAppend(err, errs, task.Telnet.Payload[5].Command)
			}
		case *transport.Message_Ssh:
			if elementInterface.MacAddressTable, err = parseMacTable(task.Ssh.Payload[0].Lookfor); err != nil {
				errs = d.logAndAppend(err, errs, task.Ssh.Payload[0].Command)
			}

			if elementInterface.DhcpTable, err = parseIPTable(task.Ssh.Payload[1].Lookfor); err != nil {
				errs = d.logAndAppend(err, errs, task.Ssh.Payload[1].Command)
			}
			elementInterface.Config = parseCurrentConfig(task.Ssh.Payload[2].Lookfor)

			if elementInterface.ConfiguredTrafficPolicy, err = parsePolicy(task.Ssh.Payload[3].Lookfor); err != nil {
				errs = d.logAndAppend(err, errs, task.Ssh.Payload[3].Command)
			}

			if err = parsePolicyStatistics(elementInterface.ConfiguredTrafficPolicy, task.Ssh.Payload[4].Lookfor); err != nil {
				errs = d.logAndAppend(err, errs, task.Ssh.Payload[4].Command)
			}

			if elementInterface.Qos, err = parseQueueStatistics(task.Ssh.Payload[5].Lookfor); err != nil {
				errs = d.logAndAppend(err, errs, task.Ssh.Payload[5].Command)
			}
		}
	}
	if elementInterface.Transceiver, err = d.GetTransceiverInformation(ctx, el); err != nil {
		errs = d.logAndAppend(err, errs, "GetTransceiverInformation")
	}

	ne.Interfaces = append(ne.Interfaces, elementInterface)
	ne.TransientErrors = &networkelement.TransientErrors{Errors: errs}
	return ne, nil
}

func (d *VRPDriver) logAndAppend(err error, errs []*networkelement.TransientError, command string) []*networkelement.TransientError {
	d.logger.Error(err.Error())
	errs = append(errs, &networkelement.TransientError{
		Message: err.Error(),
		Level:   networkelement.TransientError_WARN,
		Cause:   command,
	})

	return errs
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

	sharedConf := shared.GetConfig()

	natsConf := sharedConf.NATS
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
		conf:   sharedConf,
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			shared.PluginResourceKey: &shared.ResourcePlugin{Impl: driver},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
