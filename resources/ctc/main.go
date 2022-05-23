package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"git.liero.se/opentelco/go-dnc/client"
	"git.liero.se/opentelco/go-swpx/proto/go/networkelement"
	proto "git.liero.se/opentelco/go-swpx/proto/go/resource"
	"git.liero.se/opentelco/go-swpx/resources"
	"git.liero.se/opentelco/go-swpx/shared"
	"git.liero.se/opentelco/go-swpx/shared/oids"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"git.liero.se/opentelco/go-dnc/models/pb/transport"
	"github.com/nats-io/nats.go"
)

var VERSION *version.Version

var logger hclog.Logger

const (
	VersionBase      = "1.0-beta"
	DriverName       = "ctc-driver"
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
type CTCDriver struct {
	logger hclog.Logger
	dnc    client.Client
	conf   *shared.Configuration
}

func (d *CTCDriver) GetConfiguration(ctx context.Context) (*shared.Configuration, error) {
	return d.conf, nil
}

func (d *CTCDriver) SetConfiguration(ctx context.Context, conf *shared.Configuration) error {
	d.conf = conf

	return nil
}

func (d *CTCDriver) Version() (string, error) {
	d.logger.Debug("message from resource-driver running at version", VERSION.String())
	dncChan <- VERSION.String()
	return fmt.Sprintf("%s@%s", DriverName, VERSION.String()), nil
}

// TechnicalPortInformation Gets all the technical information for a Port
// from interface name/descr a SNMP index must be found. This functions helps to solve this problem
func (d *CTCDriver) TechnicalPortInformation(ctx context.Context, el *proto.NetworkElement) (*networkelement.Element, error) {
	return nil, fmt.Errorf("[NOT IMPLEMENTED!]")
}

func (d *CTCDriver) logAndAppend(err error, errs []*networkelement.TransientError, command string) []*networkelement.TransientError {
	d.logger.Error("log and append error from dnc", "error", err.Error())
	errs = append(errs, &networkelement.TransientError{
		Message: err.Error(),
		Level:   networkelement.TransientError_WARN,
		Cause:   command,
	})

	return errs
}

// BasicPortInformation
func (d *CTCDriver) BasicPortInformation(ctx context.Context, el *proto.NetworkElement) (*networkelement.Element, error) {
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

	msgs = append(msgs, CreateCTCTelnetInterfaceTask(el, conf))

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

		}
	}
	if elementInterface.Transceiver, err = d.GetTransceiverInformation(ctx, el); err != nil {
		errs = d.logAndAppend(err, errs, "GetTransceiverInformation")
	}

	ne.Interfaces = append(ne.Interfaces, elementInterface)
	ne.TransientErrors = &networkelement.TransientErrors{Errors: errs}
	return ne, nil
}

// AllPortInformation
func (d *CTCDriver) AllPortInformation(ctx context.Context, el *proto.NetworkElement) (*networkelement.Element, error) {
	return nil, fmt.Errorf("[NOT IMPLEMENTED!]")
}

// MapInterface Map interfaces (IF-MIB) from device with the swpx model
func (d *CTCDriver) MapInterface(ctx context.Context, el *proto.NetworkElement) (*proto.NetworkElementInterfaces, error) {
	d.logger.Info("determine what index and name this interface has", "host", el.Hostname, "ip", el.Ip, "interface", el.Interface)
	var msg *transport.Message
	discoveryMap := make(map[int]*resources.DiscoveryItem)
	var interfaces = make(map[string]*proto.NetworkElementInterface)

	conf := shared.Proto2conf(el.Conf)

	msg = CreateCTCDiscoveryMsg(el, conf)
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
				Index:       v.Index,
				Description: v.Descr,
				Alias:       v.Alias,
			}
		}
	}

	return &proto.NetworkElementInterfaces{Interfaces: interfaces}, nil
}

// MapEntityPhysical Map interfcaes from Envirnment MIB to the swpx model
// Find matching OID for port
func (d *CTCDriver) MapEntityPhysical(ctx context.Context, el *proto.NetworkElement) (*proto.NetworkElementInterfaces, error) {

	return nil, status.Error(codes.Unimplemented, "MapEntityPhysical is unimplemented")

}

// GetTransceiverInformation Get SFP (transceiver) information
func (d *CTCDriver) GetTransceiverInformation(ctx context.Context, ne *proto.NetworkElement) (*networkelement.Transceiver, error) {
	return nil, status.Error(codes.Unimplemented, "GetTransceiverInformation is unimplemented")
}

// GetAllTransceiverInformation Maps transceivers to corresponding interfaces using physical port information in the wrapper
func (d *CTCDriver) GetAllTransceiverInformation(ctx context.Context, ne *proto.NetworkElementWrapper) (*networkelement.Element, error) {
	return nil, status.Error(codes.Unimplemented, "GetAllTransceiverInformation is unimplemented")
}

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
	err = enc.BindSendChan("ctc-driver", dncChan)
	if err != nil {
		log.Fatal(err)
	}

	logger.Debug("message", "version", VERSION.String())
	//dncClient, err := client.NewGRPC(DISPATCHER_ADDR)
	dncClient, err := client.NewNATS(strings.Join(natsConf.EventServers, ","))
	if err != nil {
		log.Fatal(err)
	}
	driver := &CTCDriver{
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
