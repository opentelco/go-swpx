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
	"sort"
	"strconv"
	"strings"
	"time"

	"git.liero.se/opentelco/go-dnc/client"
	"git.liero.se/opentelco/go-dnc/models/pb/transport"
	"github.com/pkg/errors"

	"git.liero.se/opentelco/go-swpx/config"

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
	VersionBase            = "1.0-beta"
	DriverName             = "vrp-driver"
	Float64Size            = 64
	QueueEntryLength       = 12
	defaultDeadlineTimeout = 90 * time.Second
)

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
	conf   *config.ResourceVRP
}

func (d *VRPDriver) Version() (string, error) {
	d.logger.Debug("message from resource-driver running at version", VERSION.String())
	return fmt.Sprintf("%s@%s", DriverName, VERSION.String()), nil
}

// parse a map of description/alias and return the ID
func (d VRPDriver) parseDescriptionToIndex(port string, discoveryMap map[int]*discoveryItem) (*discoveryItem, error) {
	for _, v := range discoveryMap {
		// v.index = k
		if v.Descr == port {
			d.logger.Info("parser found match", "port", v.Descr)
			return v, nil
		}
	}
	return nil, fmt.Errorf("%s was not found on network element", port)
}

// MapInterface maps the interfaces to the snmp index to be used in the rest of the driver
func (d *VRPDriver) MapInterface(ctx context.Context, req *proto.Request) (*proto.PortIndex, error) {
	d.logger.Info("determine what index and name this interface has", "host", req.Hostname, "port", req.Port)

	var msg *transport.Message
	discoveryMap := make(map[int]*discoveryItem)
	var ports = make(map[string]*proto.PortIndexEntity)

	msg = createLogicalPortIndex(req, d.conf)
	msg, err := d.dnc.Put(ctx, msg)
	if err != nil {
		d.logger.Error("could not complete MapInterface", "error", err)
		return nil, err
	}

	task := msg.Task.GetSnmpc()
	if task == nil {
		return nil, fmt.Errorf("could not complete MapEntityPhysical: %w", errors.New("task is nil"))
	}

	d.logger.Debug("the msg returns from dnc", "status", msg.Status.String(), "completed", msg.Completed.String(), "execution_time", msg.ExecutionTime.String(), "size", len(task.Metrics))
	populateDiscoveryMap(d.logger, task, discoveryMap)

	for _, v := range discoveryMap {
		ports[v.Descr] = &proto.PortIndexEntity{
			Index:       int64(v.Index),
			Description: v.Descr,
			Alias:       v.Alias,
		}
	}

	return &proto.PortIndex{Ports: ports}, nil
}

// Find matching OID for port
func (d *VRPDriver) MapEntityPhysical(ctx context.Context, request *proto.Request) (*proto.PortIndex, error) {

	portMsg := createPhysicalPortIndex(request, d.conf)

	msg, err := d.dnc.Put(ctx, portMsg)
	if err != nil {
		d.logger.Error("got error back from the dnc", "error", err.Error())
		return nil, fmt.Errorf("could not complete MapEntityPhysical: %w", err)
	}

	task := msg.Task.GetSnmpc()
	if task == nil {
		return nil, fmt.Errorf("could not complete MapEntityPhysical: %w", errors.New("task is nil"))
	}

	physicalPorts := make(map[string]*proto.PortIndexEntity)
	for _, m := range task.Metrics {
		fields := strings.Split(m.Oid, ".")
		index, err := strconv.Atoi(fields[len(fields)-1])
		if err != nil {
			d.logger.Error("can't convert phys.port to int: ", err.Error())
			return nil, err
		}

		if m.Error != "" {
			d.logger.Error("problem with snmp collection", "error", m.Error)

		}

		physicalPorts[m.GetValue()] = &proto.PortIndexEntity{
			Alias:       m.Name,
			Index:       int64(index),
			Description: m.GetValue(),
		}

	}

	return &proto.PortIndex{Ports: physicalPorts}, nil

}

// GetInterfaceStatistics returns all transceiver information in a array of Transceivers
// each Transceiver contains the physical port index that can be mapped to the interface if needed
func (d *VRPDriver) GetAllTransceiverInformation(ctx context.Context, request *proto.Request) (*networkelement.Transceivers, error) {
	response := &networkelement.Transceivers{}

	vrpMsg := createAllVRPTransceiverMsg(request, d.conf, request.NumInterfaces)
	msg, err := d.dnc.Put(ctx, vrpMsg)
	if err != nil {
		d.logger.Error("could not complete GetAllTransceiverInformation", "error", err.Error())
		return nil, err
	}

	task := msg.Task.GetSnmpc()
	if task == nil {
		return nil, fmt.Errorf("could not complete MapEntityPhysical: %w", errors.New("task is nil"))
	}

	for i := 0; i < len(task.Metrics); i += 7 {
		index, _ := strconv.Atoi(reFindIndexinOID.FindString(task.Metrics[i].Oid))
		if transceiver := d.parseTransceiverMessage(task, i); transceiver != nil {
			transceiver.PhysicalPortIndex = int64(index)
			response.Transceivers = append(response.Transceivers, transceiver)
		}
	}

	return response, nil
}

func (d *VRPDriver) GetTransceiverInformation(ctx context.Context, req *proto.Request) (*networkelement.Transceiver, error) {

	vrpMsg := createVRPTransceiverMsg(req, d.conf)
	msg, err := d.dnc.Put(ctx, vrpMsg)
	if err != nil {
		d.logger.Error("could not complete GetTransceiverInformation", "error", err)
		return nil, err
	}

	task := msg.Task.GetSnmpc()
	if task == nil {
		return nil, fmt.Errorf("could not complete MapEntityPhysical: %w", errors.New("task is nil"))
	}

	if len(task.Metrics) >= 7 {
		t := d.parseTransceiverMessage(task, 0)
		t.PhysicalPortIndex = int64(req.PhysicalPortIndex)
		return t, nil
	}

	return nil, errors.Errorf("Unsupported message type")
}

func (d *VRPDriver) AllPortInformation(ctx context.Context, req *proto.Request) (*networkelement.Element, error) {
	d.logger.Info("running ALL port info", "host", req.Hostname, "port", req.Port)
	ne := &networkelement.Element{}
	ne.Hostname = req.Hostname

	sysInfoMsg := createTaskSystemInfo(req, d.conf)
	sysInfoMsg, err := d.dnc.Put(ctx, sysInfoMsg)
	if err != nil {
		d.logger.Error("could not complete AllPortInformation", "error", err)
		return nil, err
	}

	sysInfoTask := sysInfoMsg.Task.GetSnmpc()
	for _, m := range sysInfoTask.Metrics {
		if strings.HasPrefix(m.Oid, oids.SysPrefix) {
			getSystemInformation(m, ne)
		}
	}

	portsMsg := createAllPortsMsg(req, d.conf)
	portsMsg, err = d.dnc.Put(ctx, portsMsg)
	if err != nil {
		d.logger.Error(err.Error())
		return nil, err
	}

	task := portsMsg.Task.GetSnmpc()
	discoveryMap := make(map[int]*discoveryItem)
	populateDiscoveryMap(d.logger, task, discoveryMap)

	for _, discoveryItem := range discoveryMap {
		ne.Interfaces = append(ne.Interfaces, itemToInterface(discoveryItem))
	}

	sort.Slice(ne.Interfaces, func(i, j int) bool {
		return ne.Interfaces[i].Description < ne.Interfaces[j].Description
	})

	return ne, nil
}

func (d *VRPDriver) TechnicalPortInformation(ctx context.Context, req *proto.Request) (*networkelement.Element, error) {
	d.logger.Info("running technical port info", "host", req.Hostname, "port", req.Port)
	errs := make([]*networkelement.TransientError, 0)

	var msgs []*transport.Message
	if req.LogicalPortIndex != 0 {
		msgs = append(msgs, createSinglePortMsg(req.LogicalPortIndex, req, d.conf))
		msgs = append(msgs, createTaskGetPortStats(req.LogicalPortIndex, req, d.conf))
	} else {
		msgs = append(msgs, createMsg(req, d.conf))
	}

	msgs = append(msgs, createTaskSystemInfo(req, d.conf))
	msgs = append(msgs, createBasicSSHInterfaceTask(req, d.conf))

	ne := &networkelement.Element{}
	ne.Hostname = req.Hostname

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

		switch task := reply.Task.Task.(type) {
		case *transport.Task_Snmpc:
			d.logger.Debug("the reply returns from dnc",
				"status", reply.Status.String(),
				"completed", reply.Completed.String(),
				"execution_time", reply.ExecutionTime.AsDuration().String(),
				"size", len(task.Snmpc.Metrics))

			elementInterface.Index = req.LogicalPortIndex

			for _, m := range task.Snmpc.Metrics {
				switch {
				case strings.HasPrefix(m.Oid, oids.SysPrefix):
					getSystemInformation(m, ne)

				case strings.HasPrefix(m.Oid, oids.HuaPrefix):
					d.getHuaweiInterfaceStats(m, elementInterface)

				case strings.HasPrefix(m.Oid, oids.IfEntryPrefix):
					d.getIfEntryInformation(m, elementInterface)

				case strings.HasPrefix(m.Oid, oids.IfXEntryPrefix):
					getIfXEntryInformation(m, elementInterface)
				}
			}
		case *transport.Task_Terminal:

			if reply.Error != "" {
				logger.Error("error back from dnc", "errors", reply.Error, "command", task.Terminal.Payload[0].Command)
				errs = d.logAndAppend(fmt.Errorf(reply.Error), errs, task.Terminal.Payload[0].Command)
				continue
			}

			if elementInterface.MacAddressTable, err = parseMacTable(task.Terminal.Payload[0].Data); err != nil {
				errs = d.logAndAppend(err, errs, task.Terminal.Payload[0].Command)
			}

			if elementInterface.DhcpTable, err = parseIPTable(task.Terminal.Payload[1].Data); err != nil {
				errs = d.logAndAppend(err, errs, task.Terminal.Payload[1].Command)
			}

			elementInterface.Config = parseCurrentConfig(task.Terminal.Payload[2].Data)

			if elementInterface.ConfiguredTrafficPolicy, err = parsePolicy(task.Terminal.Payload[3].Data); err != nil {
				errs = d.logAndAppend(err, errs, task.Terminal.Payload[3].Command)
			}

			if err = parsePolicyStatistics(elementInterface.ConfiguredTrafficPolicy, task.Terminal.Payload[4].Data); err != nil {
				errs = d.logAndAppend(err, errs, task.Terminal.Payload[4].Command)
			}

			if elementInterface.Qos, err = parseQueueStatistics(task.Terminal.Payload[5].Data); err != nil {
				errs = d.logAndAppend(err, errs, task.Terminal.Payload[5].Command)
			}

		}
	}
	if elementInterface.Transceiver, err = d.GetTransceiverInformation(ctx, req); err != nil {
		errs = d.logAndAppend(err, errs, "GetTransceiverInformation")
	}

	ne.Interfaces = append(ne.Interfaces, elementInterface)
	ne.TransientErrors = &networkelement.TransientErrors{Errors: errs}
	return ne, nil
}

func (d *VRPDriver) logAndAppend(err error, errs []*networkelement.TransientError, command string) []*networkelement.TransientError {
	d.logger.Error("log and append error from dnc", "error", err.Error(), "command", command)
	errs = append(errs, &networkelement.TransientError{
		Message: err.Error(),
		Level:   networkelement.TransientError_WARN,
		Cause:   command,
	})

	return errs
}

func (d *VRPDriver) BasicPortInformation(ctx context.Context, req *proto.Request) (*networkelement.Element, error) {
	d.logger.Info("running basic port info", "host", req.Hostname, "port", req.Port)
	errs := make([]*networkelement.TransientError, 0)

	var msgs []*transport.Message
	if req.LogicalPortIndex != 0 {
		msgs = append(msgs, createSinglePortMsgShort(req.LogicalPortIndex, req, d.conf))
		// msgs = append(msgs, createTaskGetPortStats(req.LogicalPortIndex, req,d.conf))
	} else {
		msgs = append(msgs, createMsg(req, d.conf))
	}

	t := createBasicSSHInterfaceTask(req, d.conf)

	msgs = append(msgs, t)

	ne := &networkelement.Element{}
	ne.Hostname = req.Hostname

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
			return nil, fmt.Errorf("could not complete BasicTechnicalPortInformation: %w", err)
		}

		switch task := reply.Task.Task.(type) {
		case *transport.Task_Snmpc:
			d.logger.Debug("the reply returns from dnc",
				"status", reply.Status.String(),
				"completed", reply.Completed.String(),
				"execution_time", reply.ExecutionTime.AsDuration().String(),
				"size", len(task.Snmpc.Metrics))

			elementInterface.Index = req.LogicalPortIndex

			for _, m := range task.Snmpc.Metrics {
				switch {
				case strings.HasPrefix(m.Oid, oids.SysPrefix):
					getSystemInformation(m, ne)

				case strings.HasPrefix(m.Oid, oids.HuaPrefix):
					d.getHuaweiInterfaceStats(m, elementInterface)

				case strings.HasPrefix(m.Oid, oids.IfEntryPrefix):
					d.getIfEntryInformation(m, elementInterface)

				case strings.HasPrefix(m.Oid, oids.IfXEntryPrefix):
					getIfXEntryInformation(m, elementInterface)
				}
			}
		case *transport.Task_Terminal:
			if reply.Error != "" {
				logger.Error("error back from dnc", "errors", reply.Error, "command", task.Terminal.Payload[0].Command)
				errs = d.logAndAppend(fmt.Errorf(reply.Error), errs, task.Terminal.Payload[0].Command)
				continue
			}

			if elementInterface.MacAddressTable, err = parseMacTable(task.Terminal.Payload[0].Data); err != nil {
				errs = d.logAndAppend(err, errs, task.Terminal.Payload[0].Command)
			}

		}
	}
	if elementInterface.Transceiver, err = d.GetTransceiverInformation(ctx, req); err != nil {
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

	var appConf config.ResourceVRP
	configPath := os.Getenv("VRP_CONFIG_FILE")
	if configPath == "" {
		configPath = "vrp.hcl"
	}
	err := config.LoadConfig(configPath, &appConf)
	if err != nil {
		logger.Error("failed to loadd.config", "error", err)
		os.Exit(1)
	}

	logger.Info("connecting to DNC", "address", appConf.DNC.Addr)
	dncClient, err := client.NewGRPC(appConf.DNC.Addr)
	if err != nil {
		logger.Error("failed to create DNC Client", "error", err)
		os.Exit(1)
	}
	driver := &VRPDriver{
		logger: logger,
		dnc:    dncClient,
		conf:   &appConf,
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			shared.PluginResourceKey: &shared.ResourcePlugin{Impl: driver},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

func validateEOLTimeout(req *proto.Request, defaultDuration time.Duration) time.Duration {
	dur, _ := time.ParseDuration(req.Timeout)

	if dur.Seconds() == 0 {
		dur = defaultDuration

	}

	return dur

}
