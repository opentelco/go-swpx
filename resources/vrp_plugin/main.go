package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/opentelco/go-swpx/shared/oids"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"
	"github.com/nats-io/go-nats"
	"github.com/opentelco/go-dnc/client"
	"github.com/opentelco/go-dnc/models/protobuf/transport"
	"github.com/opentelco/go-swpx/proto/networkelement"
	proto "github.com/opentelco/go-swpx/proto/resource"

	"github.com/opentelco/go-swpx/shared"
)

var VERSION *version.Version

const (
	VERSION_BASE string = "1.0-beta"
	DRIVER_NAME  string = "vrp-driver"
)

var reFindIndexinOID = regexp.MustCompile("(\\d+)$") // used to get the last number of the oid

type discoveryItem struct {
	index int
	descr string
	alias string
}

var (
	dncChan       chan string
	EVENT_SERVERS []string = []string{"nats://localhost:14222", "nats://localhost:24222", "nats://localhost:34222"}
)

func init() {
	var err error
	if VERSION, err = version.NewVersion(VERSION_BASE); err != nil {
		log.Fatal(err)
	}
}

// Here is a real implementation of Driver
type VRPDriver struct {
	logger hclog.Logger
	dnc    client.Client
}

func (d *VRPDriver) Version() (string, error) {
	d.logger.Debug("message from resource-driver running at version", VERSION.String())
	dncChan <- VERSION.String()
	return fmt.Sprintf("%s@%s", DRIVER_NAME, VERSION.String()), nil
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

func (d *VRPDriver) MapInterface(ctx context.Context, el *proto.NetworkElement) (*proto.NetworkElementInterface, error) {
	d.logger.Info("got a task to determine what index and name this interface has", "host", el.Hostname, "ip", el.Ip, "interface", el.Interface)
	var msg *transport.Message
	discoveryMap := make(map[int]*discoveryItem)
	var index int

	msg = createDiscoveryMsg(el)
	msg, err := d.dnc.Put(msg)
	if err != nil {
		d.logger.Error(err.Error())
		return nil, err
	}

	switch task := msg.Task.(type) {
	case *transport.Message_Snmpc:
		d.logger.Debug("the msg returns from dnc", "status", msg.Status.String(), "completed", msg.Completed.String(), "execution_time", msg.ExecutionTime.String(), "size", len(task.Snmpc.Metrics))
		for _, m := range task.Snmpc.Metrics {
			index, _ = strconv.Atoi(reFindIndexinOID.FindString(m.Oid))
			// d.logger.Debug("metric data", "oid", m.Oid, "index", index, "base", m.BaseOid)
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

	item, err := d.parseDescriptionToIndex(el.Interface, discoveryMap)
	if err != nil {
		d.logger.Error("failed to parse port name", err.Error())
		return nil, err
	}

	return &proto.NetworkElementInterface{Index: int64(item.index), Description: item.descr}, nil
}

// GIMME DATA!!! InterfaceMetrics
func (d *VRPDriver) TechnicalPortInformation(ctx context.Context, el *proto.NetworkElement) (*networkelement.Element, error) {
	dncChan <- "ok"
	d.logger.Info("running technical port info", "host", el.Hostname, "ip", el.Ip, "interface", el.Interface)

	var msgs []*transport.Message
	if el.InterfaceIndex != 0 {
		msgs = append(msgs, createSinglePortMsg(el.InterfaceIndex, el))
		msgs = append(msgs, createTaskGetPortStats(el.InterfaceIndex, el))
	} else {
		msgs = append(msgs, createMsg())
	}

	msgs = append(msgs, createTaskSystemInfo(el))

	ne := &networkelement.Element{}
	ne.Hostname = el.Hostname
	elif := &networkelement.Interface{
		Stats: &networkelement.InterfaceStatistics{
			Input:  &networkelement.InterfaceStatisticsInput{},
			Output: &networkelement.InterfaceStatisticsOutput{},
		},
	}
	var err error

	for _, msg := range msgs {
		d.logger.Debug("sending msg")
		msg, err = d.dnc.Put(msg)
		if err != nil {
			d.logger.Error(err.Error())
			return nil, err
		}

		switch task := msg.Task.(type) {
		case *transport.Message_Snmpc:
			d.logger.Debug("the msg returns from dnc", "status", msg.Status.String(), "completed", msg.Completed.String(), "execution_time", msg.ExecutionTime.String(), "size", len(task.Snmpc.Metrics))

			elif.Index = el.InterfaceIndex

			for _, m := range task.Snmpc.Metrics {
				d.logger.Debug(m.GetStringValue())
				switch {
				// sys
				case m.Oid == oids.SysContact:
					ne.Contact = m.GetStringValue()
				case m.Oid == oids.SysDescr:
					ne.Version = m.GetStringValue()
				case m.Oid == oids.SysLocation:
					ne.Location = m.GetStringValue()
				case m.Oid == oids.SysName:
					ne.Sysname = m.GetStringValue()
				// case m.Oid == oids.SysORLastChange:
				// case m.Oid == oids.SysObjectID:
				case m.Oid == oids.SysUpTime:
					ne.Uptime = m.GetStringValue()
				
				for _, metric := range metrics {
					setValue(ne, Target, Value)
				}c

				// Output
				case strings.HasPrefix(m.Oid, oids.IfOutOctets):
					elif.Stats.Output.Bytes = m.GetIntValue()

				case strings.HasPrefix(m.Oid, oids.IfOutBroadcastPkts):
					elif.Stats.Output.Broadcast = m.GetIntValue()

				case strings.HasPrefix(m.Oid, oids.IfOutMulticastPkts):
					elif.Stats.Output.Multicast = m.GetIntValue()

				case strings.HasPrefix(m.Oid, oids.IfOutUcastPkts):
					elif.Stats.Output.Unicast = m.GetIntValue()

					// Input
				case strings.HasPrefix(m.Oid, oids.IfInOctets):
					elif.Stats.Input.Bytes = m.GetIntValue()

				case strings.HasPrefix(m.Oid, oids.IfInBroadcastPkts):
					elif.Stats.Input.Broadcast = m.GetIntValue()

				case strings.HasPrefix(m.Oid, oids.IfInMulticastPkts):
					elif.Stats.Input.Multicast = m.GetIntValue()

				case strings.HasPrefix(m.Oid, oids.IfInUcastPkts):
					elif.Stats.Input.Unicast = m.GetIntValue()

					// Rest
				case strings.HasPrefix(m.Oid, oids.IfInErrors):
					elif.Stats.Input.Errors = m.GetIntValue()

				case strings.HasPrefix(m.Oid, oids.HuaIfEtherStatInCRCPkts):
					elif.Stats.Input.CrcErrors = m.GetIntValue()

				case strings.HasPrefix(m.Oid, oids.HuaIfEtherStatInPausePkts):
					elif.Stats.Input.Pauses = m.GetIntValue()

				case strings.HasPrefix(m.Oid, oids.HuaIfEthIfStatReset):
					elif.Stats.Resets = m.GetIntValue()

				case strings.HasPrefix(m.Oid, oids.HuaIfEtherStatOutPausePkts):
					elif.Stats.Output.Pauses = m.GetIntValue()

				case strings.HasPrefix(m.Oid, oids.IfOutErrors):
					elif.Stats.Output.Errors = m.GetIntValue()

				case strings.HasPrefix(m.Oid, oids.IfAlias):
					elif.Alias = m.GetStringValue()

				case strings.HasPrefix(m.Oid, oids.IfDescr):
					elif.Description = m.GetStringValue()

				case strings.HasPrefix(m.Oid, oids.IfType):
					elif.Type = networkelement.InterfaceType(m.GetIntValue())

				case strings.HasPrefix(m.Oid, oids.IfMtu):
					elif.Mtu = m.GetIntValue()

				case strings.HasPrefix(m.Oid, oids.IfLastChange):
					elif.LastChanged = m.GetTimestampValue()

				case strings.HasPrefix(m.Oid, oids.IfPhysAddress):
					elif.Hwaddress = m.GetStringValue()

				case strings.HasPrefix(m.Oid, oids.IfOperStatus):
					elif.AdminStatus = networkelement.InterfaceStatus(m.GetIntValue())

				case strings.HasPrefix(m.Oid, oids.IfAdminStatus):
					elif.OperationalStatus = networkelement.InterfaceStatus(m.GetIntValue())

				case strings.HasPrefix(m.Oid, oids.IfHighSpeed):
					elif.Speed = m.GetIntValue()

				case strings.HasPrefix(m.Oid, oids.IfConnectorPresent):
					elif.ConnectorPresent = m.GetBoolValue()
				}

			}
		}
	}

	ne.Interfaces = append(ne.Interfaces, elif)

	return ne, nil
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:       fmt.Sprintf("%s@%s", DRIVER_NAME, VERSION.String()),
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	nc, _ := nats.Connect(strings.Join(EVENT_SERVERS, ","))
	dncChan = make(chan string, 0)
	enc, _ := nats.NewEncodedConn(nc, "json")
	enc.BindSendChan("vrp-driver", dncChan)

	logger.Debug("message", "message from resource-driver", "version", VERSION.String())
	driver := &VRPDriver{
		logger: logger,
		dnc:    client.New(":50051"),
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			shared.PLUGIN_RESOURCE_KEY: &shared.ResourcePlugin{Impl: driver},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
