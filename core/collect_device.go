package core

import (
	"context"
	"fmt"

	"go.opentelco.io/go-swpx/proto/go/corepb"
	"go.opentelco.io/go-swpx/shared"
)

func (c *Core) CollectBasicDeviceInformation(ctx context.Context, request *corepb.CollectBasicDeviceInformationRequest) (*corepb.DeviceInformationResponse, error) {
	c.logger.Debug("get basic device information",
		"hostname", request.Session.Hostname,
		"region", request.Session.NetworkRegion,
		"recreateIndex", request.Settings.RecreateIndex,
		"cacheTTL", request.Settings.CacheTtl,
		"timeout", request.Settings.Timeout,
	)

	resourceReq := assemblyResourceRequest(request.Session, request.Settings)

	/*
		Pre Processing with providers
	*/
	selectedProviders, err := c.selectProviders(ctx, request.Settings)
	if err != nil {
		return nil, err
	}

	request.Session, err = c.resolveSession(ctx, request.Session, selectedProviders)
	if err != nil {
		return nil, fmt.Errorf("could not resolve resource session request: %w", err)
	}

	// get resource plugin (this is done with the help of providers)
	var plugin shared.Resource
	plugin, err = c.resolveResourcePlugin(ctx, request.Session, request.Settings, selectedProviders)
	if err != nil {
		return nil, err
	}

	// device := &devicepb.Device{}
	// // map interfaces snmp-id (IF MIB) to ports
	// ports, err := plugin.MapInterface(ctx, resourceReq)
	// if err != nil {
	// 	return nil, err
	// }

	// physPorts, err := plugin.MapEntityPhysical(ctx, resourceReq)
	// if err != nil {
	// 	return nil, err
	// }

	// /*
	// 	Map Ports
	// 	IF-MIB is mapped with ENTITY-MIB
	// */
	// for _, p := range ports.Ports {
	// 	physPort, ok := physPorts.Ports[p.GetDescription()]
	// 	var physPortIndex int64
	// 	if ok {
	// 		physPortIndex = physPort.GetIndex()
	// 	}

	// 	device.Ports = append(device.Ports, &devicepb.Port{
	// 		Index:         p.Index,
	// 		IndexPhysical: physPortIndex,
	// 		Alias:         p.Alias,
	// 		Description:   p.Description,
	// 	})

	// }

	// // sort ports by index
	// sort.Slice(device.Ports, func(I, j int) bool {
	// 	return device.Ports[I].Index < device.Ports[j].Index
	// })

	/*
		Collect Basic Port Information
	*/
	allPortResult, err := plugin.AllPortInformation(ctx, resourceReq)
	if err != nil {
		return nil, err
	}
	// for _, p := range allPortResult.Ports {
	// 	for _, port := range device.Ports {
	// 		if p.Alias == port.Alias {
	// 		}
	// 	}
	// }
	// map physical entity snmp-id (ENTITY MIB) to ports

	// collect basic device information

	// collect port information

	// collect sfps

	// post process data (*devicepb.Device) with providers

	return &corepb.DeviceInformationResponse{
		Device: allPortResult,
	}, nil
}

func (c *Core) CollectDeviceInformation(ctx context.Context, request *corepb.CollectDeviceInformationRequest) (*corepb.DeviceInformationResponse, error) {

	c.logger.Debug("get extended device information",
		"hostname", request.Session.Hostname,
		"region", request.Session.NetworkRegion,
		"recreateIndex", request.Settings.RecreateIndex,
		"cacheTTL", request.Settings.CacheTtl,
		"timeout", request.Settings.Timeout,
	)

	// how to handle the cache now?
	// cacheTTLduration, _ := time.ParseDuration(request.Settings.CacheTtl)
	// if !request.Settings.RecreateIndex && cacheTTLduration != 0 {
	// 	cr, err := c.pollResponseCache.Pop(ctx, request.Session.Hostname, request.Session.Port, request.Session.AccessId, request.Type)
	// 	if err != nil {
	// 		c.logger.Warn("could not pop from cache", "error", err)
	// 	}
	// 	// if a cached response exists
	// 	if cr != nil {
	// 		c.logger.Debug("found cached item", "age", time.Since(cr.Timestamp))
	// 		if time.Since(cr.Timestamp) < cacheTTLduration {
	// 			c.logger.Info("found response in cache")
	// 			return cr.Response, nil
	// 		}

	// 	}
	// }

	// timeoutDur, _ := time.ParseDuration(request.Settings.Timeout)
	// if timeoutDur == 0 {
	// 	c.logger.Debug("using default timeout, since none was specified", "timeout", c.config.Request.DefaultRequestTimeout.AsDuration())
	// 	timeoutDur = c.config.Request.DefaultRequestTimeout.AsDuration()
	// }

	// ctxTimeout, cancel := context.WithTimeout(ctx, timeoutDur)
	// defer cancel()
	return &corepb.DeviceInformationResponse{}, nil

}
