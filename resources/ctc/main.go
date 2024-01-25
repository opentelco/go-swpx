package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/hashicorp/go-version"
	"go.opentelco.io/go-dnc/client"
	"go.opentelco.io/go-swpx/config"
	"go.opentelco.io/go-swpx/proto/go/devicepb"
	"go.opentelco.io/go-swpx/proto/go/resourcepb"
	"go.opentelco.io/go-swpx/shared"
	"go.opentelco.io/go-swpx/shared/oids"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.opentelco.io/go-dnc/models/pb/transportpb"
)

var VERSION *version.Version

var logger hclog.Logger

const (
	defaultDeadlineTimeout = 90 * time.Second
	VersionBase            = "1.0-beta"
	DriverName             = "ctc-driver"
	Float64Size            = 64
	QueueEntryLength       = 12
)

func init() {
	var err error
	if VERSION, err = version.NewVersion(VersionBase); err != nil {
		log.Fatal(err)
	}
}

// Here is a real implementation of Driver
type driver struct {
	logger hclog.Logger
	dnc    client.Client
	conf   *config.ResourceCTC
}

func (d *driver) Version() (string, error) {
	d.logger.Debug("message from resource-driver running", "version", VERSION.String())
	return fmt.Sprintf("%s@%s", DriverName, VERSION.String()), nil
}

// TechnicalPortInformation Gets all the technical information for a Port
// from interface name/descr a SNMP index must be found. This functions helps to solve this problem
func (d *driver) TechnicalPortInformation(ctx context.Context, req *resourcepb.Request) (*devicepb.Device, error) {
	return nil, fmt.Errorf("[NOT IMPLEMENTED!]")
}

func (d *driver) logAndAppend(err error, errs []*devicepb.TransientError, command string) []*devicepb.TransientError {
	d.logger.Error("log and append error from dnc", "error", err.Error())
	errs = append(errs, &devicepb.TransientError{
		Message: err.Error(),
		Level:   devicepb.TransientError_WARN,
		Cause:   command,
	})

	return errs
}

// BasicPortInformation
func (d *driver) BasicPortInformation(ctx context.Context, req *resourcepb.Request) (*devicepb.Device, error) {
	d.logger.Info("running basic port info", "host", req.Hostname, "port", req.Port, "index", req.LogicalPortIndex)
	errs := make([]*devicepb.TransientError, 0)

	var msgs []*transportpb.Message
	if req.LogicalPortIndex != 0 {
		msgs = append(msgs, createSinglePortMsgShort(req.LogicalPortIndex, req, d.conf))
		// msgs = append(msgs, createTaskGetPortStats(el.InterfaceIndex, el, conf))
	} else {
		msgs = append(msgs, createMsg(req, d.conf))
	}

	msgs = append(msgs, createCTCSSHInterfaceTask(req, d.conf))

	ne := &devicepb.Device{}
	ne.Hostname = req.Hostname

	// Create the model
	port := &devicepb.Port{
		Stats: &devicepb.Port_Statistics{
			Input:  &devicepb.Port_Statistics_Metrics{},
			Output: &devicepb.Port_Statistics_Metrics{},
		},
	}
	for _, msg := range msgs {
		reply, err := d.dnc.Put(ctx, msg)
		if err != nil {
			d.logger.Error("could not complete BasicTechnicalPortInformation", "error", err.Error())
			return nil, err
		}

		switch task := reply.Task.Task.(type) {
		case *transportpb.Task_Snmpc:
			d.logger.Debug("the reply returns from dnc",
				"status", reply.Status.String(),
				"completed", reply.Completed.String(),
				"execution_time", reply.ExecutionTime.AsDuration().String(),
				"size", len(task.Snmpc.Metrics))

			port.Index = req.LogicalPortIndex

			for _, m := range task.Snmpc.Metrics {
				switch {
				case strings.HasPrefix(m.Oid, oids.SysPrefix):
					getSystemInformation(m, ne)

				case strings.HasPrefix(m.Oid, oids.IfEntryPrefix):
					d.getIfEntryInformation(m, port)

				case strings.HasPrefix(m.Oid, oids.IfXEntryPrefix):
					getIfXEntryInformation(m, port)
				}
			}

		case *transportpb.Task_Terminal:
			if reply.Error != "" {
				logger.Error("error back from dnc", "errors", reply.Error, "command", task.Terminal.Payload[0].Command)
				errs = d.logAndAppend(fmt.Errorf(reply.Error), errs, task.Terminal.Payload[0].Command)
				continue
			}

			if port.MacAddressTable, err = parseMacTable(task.Terminal.Payload[0].Data); err != nil {
				errs = d.logAndAppend(err, errs, task.Terminal.Payload[0].Command)
			}

		}
	}
	// todo: add support for transceiver information
	// transceiver information is not implemented CTC
	// if port.Transceiver, err = d.GetTransceiverInformation(ctx, el); err != nil {
	// 	errs = d.logAndAppend(err, errs, "GetTransceiverInformation")
	// }

	ne.Ports = append(ne.Ports, port)
	ne.TransientErrors = &devicepb.TransientErrors{Errors: errs}
	return ne, nil
}

// AllPortInformation
func (d *driver) AllPortInformation(ctx context.Context, req *resourcepb.Request) (*devicepb.Device, error) {
	return nil, fmt.Errorf("[NOT IMPLEMENTED!]")
}

// MapInterface Map interfaces (IF-MIB) from device with the swpx model
func (d *driver) MapInterface(ctx context.Context, req *resourcepb.Request) (*resourcepb.PortIndex, error) {
	d.logger.Info("determine what index and name this interface has", "host", req.Hostname, "port", req.Port)
	var msg *transportpb.Message
	discoveryMap := make(map[int]*discoveryItem)
	var interfaces = make(map[string]*resourcepb.PortIndexEntity)

	msg = createCTCDiscoveryMsg(req, d.conf)
	msg, err := d.dnc.Put(ctx, msg)
	if err != nil {
		d.logger.Error("could not complete Mapport", "error", err)
		return nil, err
	}

	switch task := msg.Task.Task.(type) {
	case *transportpb.Task_Snmpc:
		d.logger.Debug("the msg returns from dnc", "status", msg.Status.String(), "completed", msg.Completed.String(), "execution_time", msg.ExecutionTime.AsDuration().String(), "size", len(task.Snmpc.Metrics))

		populateDiscoveryMap(d.logger, task.Snmpc, discoveryMap)

		for _, v := range discoveryMap {
			interfaces[v.Descr] = &resourcepb.PortIndexEntity{
				Index:       v.Index,
				Description: v.Descr,
				Alias:       v.Alias,
			}
		}
	}

	return &resourcepb.PortIndex{Ports: interfaces}, nil
}

// MapEntityPhysical Map interfcaes from Envirnment MIB to the swpx model
// Find matching OID for port
func (d *driver) MapEntityPhysical(ctx context.Context, req *resourcepb.Request) (*resourcepb.PortIndex, error) {
	return nil, status.Error(codes.Unimplemented, "MapEntityPhysical is unimplemented")
}

// GetTransceiverInformation Get SFP (transceiver) information
func (d *driver) GetTransceiverInformation(ctx context.Context, req *resourcepb.Request) (*devicepb.Transceiver, error) {
	return nil, status.Error(codes.Unimplemented, "GetTransceiverInformation is unimplemented")
}

// GetAllTransceiverInformation Maps transceivers to corresponding interfaces using physical port information in the wrapper
func (d *driver) GetAllTransceiverInformation(ctx context.Context, req *resourcepb.Request) (*devicepb.Transceivers, error) {
	return nil, status.Error(codes.Unimplemented, "GetAllTransceiverInformation is unimplemented")
}

func (d *driver) Discover(ctx context.Context, req *resourcepb.Request) (*devicepb.Device, error) {
	return &devicepb.Device{}, status.Error(codes.Unimplemented, "discover not implemented")
}

func main() {
	logger = hclog.New(&hclog.LoggerOptions{
		Name:       fmt.Sprintf("%s@%s", DriverName, VERSION.String()),
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	var appConf config.ResourceCTC
	configPath := os.Getenv("CTC_CONFIG_FILE")
	if configPath == "" {
		configPath = "ctc.hcl"
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

	driver := &driver{
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

func (d *driver) GetRunningConfig(ctx context.Context, req *resourcepb.GetRunningConfigParameters) (*resourcepb.GetRunningConfigResponse, error) {
	return nil, nil
}

func (d *driver) GetDeviceInformation(ctx context.Context, req *resourcepb.Request) (*devicepb.Device, error) {
	return nil, status.Error(codes.Unimplemented, "discover not implemented")
}
