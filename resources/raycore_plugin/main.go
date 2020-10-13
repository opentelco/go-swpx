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
	"git.liero.se/opentelco/go-swpx/proto/networkelement"
	proto "git.liero.se/opentelco/go-swpx/proto/resource"
	"git.liero.se/opentelco/go-swpx/resources"
	"git.liero.se/opentelco/go-swpx/shared"
	"git.liero.se/opentelco/go-swpx/shared/oids"
	"github.com/pkg/errors"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"
)

var VERSION *version.Version

const (
	VERSION_BASE string = "1.0-beta"
	DRIVER_NAME  string = "raycore-driver"
)

func init() {
	var err error
	if VERSION, err = version.NewVersion(VERSION_BASE); err != nil {
		log.Fatal(err)
	}
}

type RaycoreDriver struct {
	logger hclog.Logger
	dnc    client.Client
	conf   shared.Configuration
}

func (d *RaycoreDriver) Version() (string, error) {
	d.logger.Debug("message from resource-driver running at version:", VERSION)
	return fmt.Sprintf("%s@%s", DRIVER_NAME, VERSION.String()), nil
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   shared.MagicCookieKey,
	MagicCookieValue: shared.MagicCookieValue,
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:       fmt.Sprintf("%s@%s", DRIVER_NAME, VERSION.String()),
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})
	logger.Debug("message", "message from resource-driver", "version", VERSION.String())

	sharedConf := shared.GetConfig()

	natsConf := sharedConf.NATS
	dncClient, err := client.NewNATS(strings.Join(natsConf.EventServers, ","))
	if err != nil {
		log.Fatal(err)
	}

	resource := &RaycoreDriver{
		logger: logger,
		dnc:    dncClient,
		conf:   sharedConf,
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"resource": &shared.ResourcePlugin{Impl: resource},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

func (d *RaycoreDriver) MapEntityPhysical(ctx context.Context, el *proto.NetworkElement) (*proto.NetworkElementInterfaces, error) {
	conf := shared.Proto2conf(el.Conf)
	portMsg := resources.CreatePortInformationMsg(el, conf)
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
				d.logger.Error("can't convert phys.port to int: ", err.Error())
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

func (d *RaycoreDriver) AllPortInformation(ctx context.Context, el *proto.NetworkElement) (*networkelement.Element, error) {
	d.logger.Info("running ALL port info", "host", el.Hostname, "ip", el.Ip, "interface", el.Interface)
	conf := shared.Proto2conf(el.Conf)
	ne := &networkelement.Element{}
	ne.Hostname = el.Hostname

	sysInfoMsg := resources.CreateTaskSystemInfo(el, conf)
	sysInfoMsg, err := d.dnc.Put(ctx, sysInfoMsg)
	if err != nil {
		d.logger.Error(err.Error())
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
		resources.PopulateDiscoveryMap(task, discoveryMap)

		for _, discoveryItem := range discoveryMap {
			ne.Interfaces = append(ne.Interfaces, resources.ItemToInterface(discoveryItem))
		}

		sort.Slice(ne.Interfaces, func(i, j int) bool {
			return ne.Interfaces[i].Description < ne.Interfaces[j].Description
		})
	}

	return ne, nil
}

// Gets all the technical information for a Port
func (d *RaycoreDriver) TechnicalPortInformation(context.Context, *proto.NetworkElement) (*networkelement.Element, error) {
	return nil, nil
}

func (d *RaycoreDriver) MapInterface(context.Context, *proto.NetworkElement) (*proto.NetworkElementInterfaces, error) {
	return nil, nil
}

func (d *RaycoreDriver) GetTransceiverInformation(ctx context.Context, el *proto.NetworkElement) (*networkelement.Transceiver, error) {
	return nil, nil
}

func (d *RaycoreDriver) SetConfiguration(ctx context.Context, conf shared.Configuration) error {
	d.conf = conf

	return nil
}
func (d *RaycoreDriver) GetConfiguration(ctx context.Context) (shared.Configuration, error) {
	return d.conf, nil
}

func (d *RaycoreDriver) GetAllTransceiverInformation(ctx context.Context, el *proto.NetworkElementWrapper) (*proto.Transceivers, error) {
	conf := shared.Proto2conf(el.Element.Conf)

	vrpMsg := resources.CreateRaycoreTelnetTransceiverTask(el.Element, conf)
	msg, err := d.dnc.Put(ctx, vrpMsg)
	if err != nil {
		d.logger.Error("transceiver put error", err.Error())
		return nil, err
	}

	// TODO parse telnet output

	switch task := msg.Task.(type) {
	case *transport.Message_Telnet:
		if len(task.Telnet.Payload) > 0 {

			return &proto.Transceivers{Transceivers: make(map[int32]*networkelement.Transceiver)}, nil
		}
		// TODO append the transceiver to the last interface with type ETHERNETCSMACD once it is found

	}
	return nil, errors.Errorf("Unsupported message type")
}
