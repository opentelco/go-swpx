package core

import (
	"context"
	"fmt"

	"go.opentelco.io/go-swpx/proto/go/corepb"
	"go.opentelco.io/go-swpx/proto/go/resourcepb"
	"go.opentelco.io/go-swpx/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Core) doPollRequest(ctx context.Context, request *corepb.PollRequest) (*corepb.PollResponse, error) {

	response := &corepb.PollResponse{}

	selectedProviders, err := c.selectProviders(ctx, request.Settings)
	if err != nil {
		return nil, err
	}
	if len(selectedProviders) == 0 {
		response.Error = &corepb.Error{Message: "the provider is missing/does not exist", Code: ErrCodeInvalidProvider}
		return nil, NewError(response.Error.Message, ErrorCode(response.Error.Code))
	}

	var plugin shared.Resource
	plugin, err = c.resolveResourcePlugin(ctx, request.Session, request.Settings, selectedProviders)
	if err != nil {
		response.Error = &corepb.Error{
			Message: "selected driver is missing/does not exist",
			Code:    ErrCodeInvalidResource,
		}
		return nil, NewError(response.Error.Message, ErrorCode(response.Error.Code))
	}

	request.Session, err = c.resolveSession(ctx, request.Session, selectedProviders)
	if err != nil {
		return nil, fmt.Errorf("could not resolve resource session request: %w", err)
	}

	c.logger.Debug("request processed",
		"hostname", request.Session.Hostname,
	)

	// implementation of different requests that SWP-X can handle right now
	switch request.Type {
	case corepb.PollRequest_GET_BASIC_INFO:

		if request.Session.Port != "" {
			response, err = c.handleGetBasicInformationPort(ctx, request, plugin)
			if err != nil {
				return nil, err
			}

		} else {
			err := c.handleGetPasicInformationElement(request, response, plugin)
			if err != nil {
				return nil, err
			}

		}

	case corepb.PollRequest_GET_TECHNICAL_INFO:

		if request.Session.Port != "" {
			response, err = c.handleGetTechnicalInformationPort(ctx, request, plugin)
			if err != nil {
				return nil, err
			}

		} else {
			err := c.getTechnicalInformationElement(request, response, plugin)
			if err != nil {
				return nil, err
			}
		}

	}

	// process the response with the selected providers (post-process for polling)
	if err := providerPollPostProcess(ctx, selectedProviders, response); err != nil {
		return nil, err
	}

	return response, nil

}

// handleGetTechnicalInformationElement gets full information of an Element
func (c *Core) getTechnicalInformationElement(msg *corepb.PollRequest, resp *corepb.PollResponse, plugin shared.Resource) error {
	c.logger.Warn("getTechnicalInformationElement is not implemented")
	return nil
}

// handleGetTechnicalInformationPort gets information related to the selected interface
func (c *Core) handleGetTechnicalInformationPort(ctx context.Context, msg *corepb.PollRequest, plugin shared.Resource) (*corepb.PollResponse, error) {
	var resp corepb.PollResponse
	req := resourcepb.Request{
		NetworkRegion: msg.Session.NetworkRegion,
		Hostname:      msg.Session.Hostname,
		Port:          msg.Session.Port,
		Timeout:       msg.Settings.Timeout,
	}

	var (
		mapInterfaceResponse *resourcepb.PortIndex
		cachedInterface      *CachedInterface
		err                  error
	)

	if c.cacheEnabled && !msg.Settings.RecreateIndex {
		c.logger.Info("cache is enabled, pop index from cache")
		cachedInterface, err = c.interfaceCache.Pop(context.TODO(), req.Hostname, req.Port)

		if cachedInterface != nil {
			c.logger.Info("cached interface indexs",
				"physicalPort", cachedInterface.Port,
				"physicalIndex", cachedInterface.PhysicalEntityIndex,
				"interfaceIndex", cachedInterface.InterfaceIndex,
			)

			resp.PhysicalPort = cachedInterface.Port
			req.PhysicalPortIndex = cachedInterface.PhysicalEntityIndex
			req.LogicalPortIndex = cachedInterface.InterfaceIndex
		}
	}

	// did not find cached item or cached is disabled
	if cachedInterface == nil || !c.cacheEnabled {
		c.logger.Info("run mapEntity to get physical entity index on device")

		physPortResponse, err := plugin.MapEntityPhysical(ctx, &req)
		if err != nil {
			if status.Code(err) == codes.Unimplemented {
				c.logger.Error("warn MapEntityPhysical is not implemented, skipping", "err", err.Error())
			} else {
				c.logger.Error("error running MapEntityPhysical", "err", err.Error())
				resp.Error = &corepb.Error{
					Message: fmt.Sprintf("could not run MapEntityPhyiscal: %s", err.Error()),
					Code:    ErrCodeInvalidPort,
				}
				return nil, err
			}

		} else {

			if val, ok := physPortResponse.Ports[req.Port]; ok {
				resp.PhysicalPort = val.Description
				req.PhysicalPortIndex = val.Index

				c.logger.Debug("found physInterface",
					"port", req.Port,
					"resp.physicalPort", val.Description,
					"req.physicalIndex", val.Index,
				)

			}
		}

		if mapInterfaceResponse, err = plugin.MapInterface(ctx, &req); err != nil {
			c.logger.Error("error running map interface", "err", err.Error())
			resp.Error = &corepb.Error{
				Message: fmt.Sprintf("could not run MapInterface: %s", err.Error()),
				Code:    ErrCodeInvalidPort,
			}
			return nil, err
		}

		if val, ok := mapInterfaceResponse.Ports[req.Port]; ok {
			req.LogicalPortIndex = val.Index
			c.logger.Info("found ifMIB interface index", "index", val.Index)
		}

		// save in cache upon success (if enabled)
		if c.cacheEnabled {
			if err = c.interfaceCache.Upsert(ctx, req.Hostname, mapInterfaceResponse, physPortResponse); err != nil {
				return nil, err
			}
		}

	} else if err != nil {

		resp.Error = &corepb.Error{
			Message: fmt.Sprintf("could handle request: %s", err.Error()),
			Code:    ErrCodeUnknownError,
		}

		c.logger.Error("error fetching from cache:", err.Error())
		return nil, err
	}

	//if the return is 0 something went wrong
	if req.LogicalPortIndex == 0 {
		c.logger.Error("error running map interface", "err", "index is zero")
		resp.Error = &corepb.Error{
			Message: "interface index returned zero",
			Code:    ErrCodeInvalidPort,
		}
		return nil, err
	}

	c.logger.Info("found index for selected interface", "index", req.LogicalPortIndex)

	ne, err := plugin.TechnicalPortInformation(ctx, &req)
	if err != nil {
		c.logger.Error(err.Error())
		return nil, err
	}
	resp.Device = ne

	return &resp, nil
}

// handleGetBasicInformationPort gets information related to the selected interface
func (c *Core) handleGetBasicInformationPort(ctx context.Context, msg *corepb.PollRequest, plugin shared.Resource) (*corepb.PollResponse, error) {

	var resp corepb.PollResponse
	req := resourcepb.Request{
		NetworkRegion: msg.Session.NetworkRegion,
		Hostname:      msg.Session.Hostname,
		Port:          msg.Session.Port,
		Timeout:       msg.Settings.Timeout,
	}

	var (
		mapInterfaceResponse *resourcepb.PortIndex
		cachedInterface      *CachedInterface
		err                  error
	)

	if c.cacheEnabled && !msg.Settings.RecreateIndex {
		c.logger.Info("cache is enabled, pop index from cache")
		cachedInterface, err = c.interfaceCache.Pop(context.TODO(), req.Hostname, req.Port)

		if cachedInterface != nil {
			c.logger.Info("cached interface indexs",
				"physicalPort", cachedInterface.Port,
				"physicalIndex", cachedInterface.PhysicalEntityIndex,
				"interfaceIndex", cachedInterface.InterfaceIndex,
			)

			resp.PhysicalPort = cachedInterface.Port

			req.PhysicalPortIndex = cachedInterface.PhysicalEntityIndex
			req.LogicalPortIndex = cachedInterface.InterfaceIndex
		}
	}

	// did not find cached item or cached is disabled
	if cachedInterface == nil || !c.cacheEnabled {
		c.logger.Info("run mapEntity to get physical entity index on device")

		physPortResponse, err := plugin.MapEntityPhysical(ctx, &req)
		if err != nil {
			if status.Code(err) == codes.Unimplemented {
				c.logger.Error("warn MapEntityPhysical is not implemented, skipping", "err", err.Error())
			} else {
				c.logger.Error("error running MapEntityPhysical", "err", err.Error())
				resp.Error = &corepb.Error{
					Message: fmt.Sprintf("could not run MapEntityPhyiscal: %s", err.Error()),
					Code:    ErrCodeInvalidPort,
				}
				return nil, err
			}

		} else {

			if val, ok := physPortResponse.Ports[req.Port]; ok {
				resp.PhysicalPort = val.Description
				req.PhysicalPortIndex = val.Index

				c.logger.Debug("found physInterface",
					"port", req.Port,
					"resp.physicalPort", val.Description,
					"req.physicalIndex", val.Index,
				)

			}
		}

		if mapInterfaceResponse, err = plugin.MapInterface(ctx, &req); err != nil {
			c.logger.Error("error running map interface", "err", err.Error())
			resp.Error = &corepb.Error{
				Message: fmt.Sprintf("could not run MapInterface: %s", err.Error()),
				Code:    ErrCodeInvalidPort,
			}
			return nil, err
		}

		if val, ok := mapInterfaceResponse.Ports[req.Port]; ok {
			req.LogicalPortIndex = val.Index
			c.logger.Info("found ifMIB interface index", "index", val.Index)
		}

		// save in cache upon success (if enabled)
		if c.cacheEnabled {
			if err = c.interfaceCache.Upsert(ctx, req.Hostname, mapInterfaceResponse, physPortResponse); err != nil {
				return nil, err
			}
		}

	} else if err != nil {

		resp.Error = &corepb.Error{
			Message: fmt.Sprintf("could handle request: %s", err.Error()),
			Code:    ErrCodeUnknownError,
		}

		c.logger.Error("error fetching from cache:", err.Error())
		return nil, err
	}

	//if the return is 0 something went wrong
	if req.LogicalPortIndex == 0 {
		c.logger.Error("error running map interface", "err", "index is zero")
		resp.Error = &corepb.Error{
			Message: "interface index returned zero",
			Code:    ErrCodeInvalidPort,
		}
		return nil, err
	}

	c.logger.Info("found index for selected interface", "index", req.LogicalPortIndex)

	ne, err := plugin.BasicPortInformation(ctx, &req)
	if err != nil {
		c.logger.Error(err.Error())
		return nil, err
	}
	resp.Device = ne

	return &resp, nil
}

// handleGetTechnicalInformationElement gets full information of an Element
func (c *Core) handleGetPasicInformationElement(msg *corepb.PollRequest, resp *corepb.PollResponse, plugin shared.Resource) error {
	// resourcepbConf := shared.Conf2resourcepb(conf)
	// resourcepbConf.Request = c.createRequestConfig(msg, conf) // set deadline
	// req := &resource.Device{
	// 	Interface: "",
	// 	Hostname:  msg.Hostname,
	// 	Conf:      resourcepbConf,
	// }

	// physPortResponse, err := plugin.MapEntityPhysical(msg.ctx, req)
	// if err != nil {
	// 	c.logger.Error("error fetching physical entities:", err.Error())
	// 	return err
	// }

	// allPortInformation, err := plugin.AllPortInformation(msg.ctx, req)
	// if err != nil {
	// 	c.logger.Error("error fetching port information for all interfaces:", err.Error())
	// 	return err
	// }

	// var matchingInterfaces int32 = 0
	// for _, iface := range allPortInformation.Interfaces {
	// 	if _, ok := physPortResponse.Interfaces[iface.Description]; ok {
	// 		matchingInterfaces++
	// 	}
	// }
	// allPortInformation, err = plugin.GetAllTransceiverInformation(msg.ctx, &resource.DeviceWrapper{
	// 	Element:        req,
	// 	NumInterfaces:  matchingInterfaces,
	// 	FullElement:    allPortInformation,
	// 	PhysInterfaces: physPortResponse,
	// })
	// if err != nil {
	// 	c.logger.Error("error fetching transceiver information: ", err)
	// }

	// resp.Device = allPortInformation

	return nil
}
