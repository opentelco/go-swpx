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
	"log"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"git.liero.se/opentelco/go-dnc/client"
	"git.liero.se/opentelco/go-dnc/models/pb/transport"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"

	"git.liero.se/opentelco/go-swpx/resources"

	"git.liero.se/opentelco/go-swpx/shared/oids"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"

	"git.liero.se/opentelco/go-swpx/proto/go/networkelement"
	proto "git.liero.se/opentelco/go-swpx/proto/go/resource"
	"git.liero.se/opentelco/go-swpx/shared"
)

var VERSION *version.Version

var logger hclog.Logger

const (
	VersionBase      = "1.0-beta"
	DriverName       = "vrp-driver"
	Float64Size      = 64
	QueueEntryLength = 12
)

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
	conf   *shared.Configuration
}

func (d *VRPDriver) GetConfiguration(ctx context.Context) (*shared.Configuration, error) {
	return d.conf, nil
}

func (d *VRPDriver) SetConfiguration(ctx context.Context, conf *shared.Configuration) error {
	d.conf = conf

	return nil
}

func (d *VRPDriver) Version() (string, error) {
	d.logger.Debug("message from resource-driver running at version", VERSION.String())
	dncChan <- VERSION.String()
	return fmt.Sprintf("%s@%s", DriverName, VERSION.String()), nil
}

// parse a map of description/alias and return the ID
func (d VRPDriver) parseDescriptionToIndex(port string, discoveryMap map[int]*resources.DiscoveryItem) (*resources.DiscoveryItem, error) {
	for _, v := range discoveryMap {
		// v.index = k
		if v.Descr == port {
			d.logger.Info("parser found match", "port", v.Descr)
			return v, nil
		}
	}
	return nil, fmt.Errorf("%s was not found on network element", port)
}

// Find matching OID for port
func (d *VRPDriver) MapEntityPhysical(ctx context.Context, el *proto.NetworkElement) (*proto.NetworkElementInterfaces, error) {
	conf := shared.Proto2conf(el.Conf)
	portMsg := resources.CreatePortInformationMsg(el, conf)
	msg, err := d.dnc.Put(ctx, portMsg)
	if err != nil {
		d.logger.Error("got error back from the dnc", "error", err.Error())
		return nil, fmt.Errorf("could not complete MapEntityPhysical: %w", err)
	}

	switch task := msg.Task.(type) {
	case *transport.Message_Snmpc:
		interfaces := make(map[string]*proto.NetworkElementInterface)
		for _, m := range task.Snmpc.Metrics {
			fields := strings.Split(m.Oid, ".")
			index, err := strconv.Atoi(fields[len(fields)-1])
			if err != nil {
				d.logger.Error("can't convert phys.port to int: ", err.Error())
				return nil, err
			}

			if m.Error != "" {
				d.logger.Error("problem with snmp collection", "error", m.Error)

			}
			interfaces[m.GetStringValue()] = &proto.NetworkElementInterface{
				Alias:       m.Name,
				Index:       int64(index),
				Description: m.GetStringValue(),
			}
		}

		return &proto.NetworkElementInterfaces{Interfaces: interfaces}, nil
	case nil:
		return nil, errors.Errorf("msg task back from DNC is nil")
	default:
		return nil, errors.Errorf("unsupported message task back from dnc: %s", reflect.TypeOf(msg.Task).String())

	}

}

func (d *VRPDriver) GetAllTransceiverInformation(ctx context.Context, wrapper *proto.NetworkElementWrapper) (*networkelement.Element, error) {
	el := wrapper.Element
	conf := shared.Proto2conf(el.Conf)
	result := make(map[int32]*networkelement.Transceiver)

	vrpMsg := resources.CreateAllVRPTransceiverMsg(el, conf, wrapper.NumInterfaces)
	msg, err := d.dnc.Put(ctx, vrpMsg)
	if err != nil {
		d.logger.Error("could not complete GetAllTransceiverInformation", "error", err.Error())
		return wrapper.FullElement, err
	}

	switch task := msg.Task.(type) {
	case *transport.Message_Snmpc:
		for i := 0; i < len(task.Snmpc.Metrics); i += 7 {
			index, _ := strconv.Atoi(resources.ReFindIndexinOID.FindString(task.Snmpc.Metrics[i].Oid))
			if transceiver := d.parseTransceiverMessage(task, i); transceiver != nil {
				result[int32(index)] = transceiver
			}
		}
	}

	// match transceiver to interface using phys. indexes
	for _, iface := range wrapper.FullElement.Interfaces {
		if matchingPhysInterface, ok := wrapper.PhysInterfaces.Interfaces[iface.Description]; ok {
			iface.Transceiver = result[int32(matchingPhysInterface.Index)]
		}
	}

	return wrapper.FullElement, nil
}

func (d *VRPDriver) GetTransceiverInformation(ctx context.Context, el *proto.NetworkElement) (*networkelement.Transceiver, error) {
	conf := shared.Proto2conf(el.Conf)

	vrpMsg := createVRPTransceiverMsg(el, conf)
	msg, err := d.dnc.Put(ctx, vrpMsg)
	if err != nil {
		d.logger.Error("could not complete GetTransceiverInformation", "error", err)
		return nil, err
	}

	switch task := msg.Task.(type) {
	case *transport.Message_Snmpc:
		if len(task.Snmpc.Metrics) >= 7 {
			return d.parseTransceiverMessage(task, 0), nil
		}
	}
	return nil, errors.Errorf("Unsupported message type")
}

func (d *VRPDriver) MapInterface(ctx context.Context, el *proto.NetworkElement) (*proto.NetworkElementInterfaces, error) {
	d.logger.Info("determine what index and name this interface has", "host", el.Hostname, "ip", el.Ip, "interface", el.Interface)
	var msg *transport.Message
	discoveryMap := make(map[int]*resources.DiscoveryItem)
	var interfaces = make(map[string]*proto.NetworkElementInterface)

	conf := shared.Proto2conf(el.Conf)

	msg = resources.CreateDiscoveryMsg(el, conf)
	msg, err := d.dnc.Put(ctx, msg)
	if err != nil {
		d.logger.Error("could not complete MapInterface", "error", err)
		return nil, err
	}

	switch task := msg.Task.(type) {
	case *transport.Message_Snmpc:
		d.logger.Debug("the msg returns from dnc", "status", msg.Status.String(), "completed", msg.Completed.String(), "execution_time", msg.ExecutionTime.String(), "size", len(task.Snmpc.Metrics))
		resources.PopulateDiscoveryMap(d.logger, task, discoveryMap)

		for _, v := range discoveryMap {
			interfaces[v.Descr] = &proto.NetworkElementInterface{
				Index:       int64(v.Index),
				Description: v.Descr,
				Alias:       v.Alias,
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

	sysInfoMsg := resources.CreateTaskSystemInfo(el, conf)
	sysInfoMsg, err := d.dnc.Put(ctx, sysInfoMsg)
	if err != nil {
		d.logger.Error("could not complete AllPortInformation", "error", err)
		return nil, err
	}

	sysInfoTask := sysInfoMsg.Task.(*transport.Message_Snmpc)
	for _, m := range sysInfoTask.Snmpc.Metrics {
		if strings.HasPrefix(m.Oid, oids.SysPrefix) {
			resources.GetSystemInformation(m, ne)
		}
	}

	portsMsg := resources.CreateAllPortsMsg(el, conf)
	portsMsg, err = d.dnc.Put(ctx, portsMsg)
	if err != nil {
		d.logger.Error(err.Error())
		return nil, err
	}

	if task, ok := portsMsg.Task.(*transport.Message_Snmpc); ok {
		discoveryMap := make(map[int]*resources.DiscoveryItem)
		resources.PopulateDiscoveryMap(d.logger, task, discoveryMap)

		for _, discoveryItem := range discoveryMap {
			ne.Interfaces = append(ne.Interfaces, resources.ItemToInterface(discoveryItem))
		}

		sort.Slice(ne.Interfaces, func(i, j int) bool {
			return ne.Interfaces[i].Description < ne.Interfaces[j].Description
		})
	}

	return ne, nil
}

func (d *VRPDriver) TechnicalPortInformation(ctx context.Context, el *proto.NetworkElement) (*networkelement.Element, error) {
	d.logger.Info("running technical port info", "host", el.Hostname, "ip", el.Ip, "interface", el.Interface)
	errs := make([]*networkelement.TransientError, 0)

	conf := shared.Proto2conf(el.Conf)

	var msgs []*transport.Message
	if el.InterfaceIndex != 0 {
		msgs = append(msgs, resources.CreateSinglePortMsg(el.InterfaceIndex, el, conf))
		msgs = append(msgs, createTaskGetPortStats(el.InterfaceIndex, el, conf))
	} else {
		msgs = append(msgs, resources.CreateMsg(conf))
	}

	msgs = append(msgs, resources.CreateTaskSystemInfo(el, conf))
	msgs = append(msgs, resources.CreateTelnetInterfaceTask(el, conf))

	ne := &networkelement.Element{}
	ne.Hostname = el.Hostname

	// Create the model
	elementInterface := &networkelement.Interface{
		Stats: &networkelement.InterfaceStatistics{
			Input:  &networkelement.InterfaceStatisticsInput{},
			Output: &networkelement.InterfaceStatisticsOutput{},
		},
	}
	var err error

	for _, msg := range msgs {
		reply, err := d.dnc.Put(ctx, msg)
		if err != nil {
			d.logger.Error("could not complete TechnicalPortInformation", "error", err.Error())
			return nil, err
		}

		switch task := reply.Task.(type) {
		case *transport.Message_Snmpc:
			d.logger.Debug("the reply returns from dnc",
				"status", reply.Status.String(),
				"completed", reply.Completed.String(),
				"execution_time", reply.ExecutionTime.String(),
				"size", len(task.Snmpc.Metrics))

			elementInterface.Index = el.InterfaceIndex

			for _, m := range task.Snmpc.Metrics {
				switch {
				case strings.HasPrefix(m.Oid, oids.SysPrefix):
					resources.GetSystemInformation(m, ne)

				case strings.HasPrefix(m.Oid, oids.HuaPrefix):
					d.getHuaweiInterfaceStats(m, elementInterface)

				case strings.HasPrefix(m.Oid, oids.IfEntryPrefix):
					d.getIfEntryInformation(m, elementInterface)

				case strings.HasPrefix(m.Oid, oids.IfXEntryPrefix):
					resources.GetIfXEntryInformation(m, elementInterface)
				}
			}
		case *transport.Message_Telnet:
			if reply.Error != "" {
				logger.Error("error back from dnc", "errors", reply.Error, "command", task.Telnet.Payload[0].Command)
				errs = d.logAndAppend(fmt.Errorf(reply.Error), errs, task.Telnet.Payload[0].Command)
				continue
			}

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
			if reply.Error != "" {
				errs = d.logAndAppend(fmt.Errorf(reply.Error), errs, task.Ssh.Payload[0].Command)
				continue
			}

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
	d.logger.Error("log and append error from dnc", "error", err.Error())
	errs = append(errs, &networkelement.TransientError{
		Message: err.Error(),
		Level:   networkelement.TransientError_WARN,
		Cause:   command,
	})

	return errs
}

func (d *VRPDriver) BasicPortInformation(ctx context.Context, el *proto.NetworkElement) (*networkelement.Element, error) {
	d.logger.Info("running basic port info", "host", el.Hostname, "ip", el.Ip, "interface", el.Interface)
	errs := make([]*networkelement.TransientError, 0)

	conf := shared.Proto2conf(el.Conf)

	var msgs []*transport.Message
	if el.InterfaceIndex != 0 {
		msgs = append(msgs, resources.CreateSinglePortMsg(el.InterfaceIndex, el, conf))
		msgs = append(msgs, createTaskGetPortStats(el.InterfaceIndex, el, conf))
	} else {
		msgs = append(msgs, resources.CreateMsg(conf))
	}

	msgs = append(msgs, resources.CreateBasicTelnetInterfaceTask(el, conf))

	ne := &networkelement.Element{}
	ne.Hostname = el.Hostname

	// Create the model
	elementInterface := &networkelement.Interface{
		Stats: &networkelement.InterfaceStatistics{
			Input:  &networkelement.InterfaceStatisticsInput{},
			Output: &networkelement.InterfaceStatisticsOutput{},
		},
	}
	var err error

	for _, msg := range msgs {
		reply, err := d.dnc.Put(ctx, msg)
		if err != nil {
			d.logger.Error("could not complete BasicTechnicalPortInformation", "error", err.Error())
			return nil, err
		}

		switch task := reply.Task.(type) {
		case *transport.Message_Snmpc:
			d.logger.Debug("the reply returns from dnc",
				"status", reply.Status.String(),
				"completed", reply.Completed.String(),
				"execution_time", reply.ExecutionTime.String(),
				"size", len(task.Snmpc.Metrics))

			elementInterface.Index = el.InterfaceIndex

			for _, m := range task.Snmpc.Metrics {
				switch {
				case strings.HasPrefix(m.Oid, oids.SysPrefix):
					resources.GetSystemInformation(m, ne)

				case strings.HasPrefix(m.Oid, oids.HuaPrefix):
					d.getHuaweiInterfaceStats(m, elementInterface)

				case strings.HasPrefix(m.Oid, oids.IfEntryPrefix):
					d.getIfEntryInformation(m, elementInterface)

				case strings.HasPrefix(m.Oid, oids.IfXEntryPrefix):
					resources.GetIfXEntryInformation(m, elementInterface)
				}
			}
		case *transport.Message_Telnet:
			if reply.Error != "" {
				logger.Error("error back from dnc", "errors", reply.Error, "command", task.Telnet.Payload[0].Command)
				errs = d.logAndAppend(fmt.Errorf(reply.Error), errs, task.Telnet.Payload[0].Command)
				continue
			}

			if elementInterface.MacAddressTable, err = parseMacTable(task.Telnet.Payload[0].Lookfor); err != nil {
				errs = d.logAndAppend(err, errs, task.Telnet.Payload[0].Command)
			}

			if elementInterface.DhcpTable, err = parseIPTable(task.Telnet.Payload[1].Lookfor); err != nil {
				errs = d.logAndAppend(err, errs, task.Telnet.Payload[1].Command)
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
	enc, err := nats.NewEncodedConn(nc, "json")
	if err != nil {
		logger.Error("failed to create dnc connection", "error", err)
		os.Exit(1)
	}
	err = enc.BindSendChan("vrp-driver", dncChan)
	if err != nil {
		log.Fatal(err)
	}

	logger.Debug("message", "version", VERSION.String())
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
