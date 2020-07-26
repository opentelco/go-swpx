package main

import (
	"context"
	"fmt"
	"git.liero.se/opentelco/go-dnc/models/protobuf/metric"
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
	VERSION_BASE    string = "1.0-beta"
	DRIVER_NAME     string = "vrp-driver"
	DISPATCHER_ADDR string = "127.0.0.1:50051"
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

// Find matching OID for port
func (d *VRPDriver) MapEntityPhysical(ctx context.Context, el *proto.NetworkElement) (*proto.NetworkElementInterfaces, error) {
	conf := shared.Proto2conf(*el.Conf)
	portMsg := createPortInformationMsg(el, conf)
	msg, err := d.dnc.Put(&portMsg)
	if err != nil {
		d.logger.Error(err.Error())
		return nil, err
	}

	switch task := msg.Task.(type) {
	case *transport.Message_Snmpc:
		data := make([]*proto.NetworkElementInterface, len(task.Snmpc.Metrics))

		for i, m := range task.Snmpc.Metrics {
			fields := strings.Split(m.Oid, ".")
			index, err := strconv.Atoi(fields[len(fields)-1])
			if err != nil {
				logger.Error("can't convert phys.port to int: ", err.Error())
				return nil, err
			}

			data[i] = &proto.NetworkElementInterface{
				Alias:       m.Name,
				Index:       int64(index),
				Description: m.GetStringValue(),
			}
		}

		return &proto.NetworkElementInterfaces{Interfaces: data}, nil
	}
	return nil, errors.Errorf("Unsupported message type")
}

func (d *VRPDriver) GetTransceiverInformation(ctx context.Context, el *proto.NetworkElement) (*networkelement.Transceiver, error) {
	conf := shared.Proto2conf(*el.Conf)

	vrpMsg := createVRPTransceiverMsg(el, conf)
	msg, err := d.dnc.Put(&vrpMsg)
	if err != nil {
		d.logger.Error(err.Error())
		return nil, err
	}

	switch task := msg.Task.(type) {
	case *transport.Message_Snmpc:
		if len(task.Snmpc.Metrics) >= 7 {
			val := &networkelement.Transceiver{
				SerialNumber: task.Snmpc.Metrics[0].GetStringValue(),
				Stats: []*networkelement.TransceiverStatistics{
					{
						Temp:    task.Snmpc.Metrics[1].GetIntValue(),
						Voltage: task.Snmpc.Metrics[2].GetIntValue(),
						Current: task.Snmpc.Metrics[3].GetIntValue(),
						Rx:      task.Snmpc.Metrics[4].GetIntValue(),
						Tx:      task.Snmpc.Metrics[5].GetIntValue(),
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
	var interfaces = make([]*proto.NetworkElementInterface, 1)

	conf := shared.Proto2conf(*el.Conf)

	msg = createDiscoveryMsg(el, conf)
	msg, err := d.dnc.Put(msg)
	if err != nil {
		d.logger.Error(err.Error())
		return nil, err
	}

	switch task := msg.Task.(type) {
	case *transport.Message_Snmpc:
		d.logger.Debug("the msg returns from dnc", "status", msg.Status.String(), "completed", msg.Completed.String(), "execution_time", msg.ExecutionTime.String(), "size", len(task.Snmpc.Metrics))
		d.populateDiscoveryMap(task, index, discoveryMap)

		for _, v := range discoveryMap {
			interfaces = append(interfaces, &proto.NetworkElementInterface{
				Index:       int64(v.index),
				Description: v.descr,
				Alias:       v.alias,
			})
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

// GIMME DATA!!! InterfaceMetrics
func (d *VRPDriver) TechnicalPortInformation(ctx context.Context, el *proto.NetworkElement) (*networkelement.Element, error) {
	dncChan <- "ok"
	d.logger.Info("running technical port info", "host", el.Hostname, "ip", el.Ip, "interface", el.Interface)

	conf := shared.Proto2conf(*el.Conf)

	var msgs []*transport.Message
	if el.InterfaceIndex != 0 {
		msgs = append(msgs, createSinglePortMsg(el.InterfaceIndex, el, conf))
		msgs = append(msgs, createTaskGetPortStats(el.InterfaceIndex, el, conf))
	} else {
		msgs = append(msgs, createMsg(conf))
	}

	msgs = append(msgs, createTaskSystemInfo(el, conf))

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
				case strings.HasPrefix(m.Oid, oids.SysPrefix):
					getSystemInformation(m, ne)

				case strings.HasPrefix(m.Oid, oids.HuaPrefix):
					getHuaweiInformation(m, elif)

				case strings.HasPrefix(m.Oid, oids.IfEntryPrefix):
					getIfEntryInformation(m, elif)

				case strings.HasPrefix(m.Oid, oids.IfXEntryPrefix):
					getIfXEntryInformation(m, elif)
				}
			}
		}
	}

	ne.Interfaces = append(ne.Interfaces, elif)

	return ne, nil
}

func getIfXEntryInformation(m *metric.Metric, elif *networkelement.Interface) {

	switch {
	case strings.HasPrefix(m.Oid, oids.IfOutUcastPkts):
		elif.Stats.Output.Unicast = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfInBroadcastPkts):
		elif.Stats.Input.Broadcast = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfInMulticastPkts):
		elif.Stats.Input.Multicast = m.GetIntValue()
	case strings.HasPrefix(m.Oid, oids.IfOutBroadcastPkts):
		elif.Stats.Output.Broadcast = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfOutMulticastPkts):
		elif.Stats.Output.Multicast = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfAlias):
		elif.Alias = m.GetStringValue()

	case strings.HasPrefix(m.Oid, oids.IfHighSpeed):
		elif.Speed = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfConnectorPresent):
		elif.ConnectorPresent = m.GetBoolValue()
	}

}

func getIfEntryInformation(m *metric.Metric, elif *networkelement.Interface) {
	switch {
	case strings.HasPrefix(m.Oid, oids.IfOutOctets):
		elif.Stats.Output.Bytes = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfInOctets):
		elif.Stats.Input.Bytes = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfInUcastPkts):
		elif.Stats.Input.Unicast = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfInErrors):
		elif.Stats.Input.Errors = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.IfOutErrors):
		elif.Stats.Output.Errors = m.GetIntValue()

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

	}
}

func getHuaweiInformation(m *metric.Metric, elif *networkelement.Interface) {
	switch {
	case strings.HasPrefix(m.Oid, oids.HuaIfEtherStatInCRCPkts):
		elif.Stats.Input.CrcErrors = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.HuaIfEtherStatInPausePkts):
		elif.Stats.Input.Pauses = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.HuaIfEthIfStatReset):
		elif.Stats.Resets = m.GetIntValue()

	case strings.HasPrefix(m.Oid, oids.HuaIfEtherStatOutPausePkts):
		elif.Stats.Output.Pauses = m.GetIntValue()
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
	dncClient, err := client.New(DISPATCHER_ADDR)
	if err != nil {
		log.Fatal(err)
	}
	driver := &VRPDriver{
		logger: logger,
		dnc:    dncClient,
		conf: shared.Configuration{
			SNMP:   shared.ConfigSNMP{
				Community:          "semipublic",
				Version:            2,
				Timeout:            time.Second*20,
				Retries:            2,
				DynamicRepetitions: true,
			},
			Telnet: shared.ConfigTelnet{},
		},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			shared.PLUGIN_RESOURCE_KEY: &shared.ResourcePlugin{Impl: driver},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
