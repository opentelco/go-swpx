package core

import (
	"context"
	"fmt"

	pb_core "git.liero.se/opentelco/go-swpx/proto/go/core"
	proto "git.liero.se/opentelco/go-swpx/proto/go/resource"
	"git.liero.se/opentelco/go-swpx/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *Core) doRequest(ctx context.Context, request *pb_core.Request) (*pb_core.Response, error) {

	response := &pb_core.Response{}

	selectedProviders, err := c.selectProviders(ctx, request)
	if err != nil {
		return nil, err
	}
	if len(selectedProviders) == 0 {
		response.Error = &pb_core.Error{Message: "the provider is missing/does not exist", Code: ErrInvalidProvider}
		return nil, NewError(response.Error.Message, ErrorCode(response.Error.Code))
	}

	// select resource-plugin to send the requests to
	c.logger.Info("selected resource plugin", "plugin", request.Settings.ResourcePlugin)
	plugin := c.resources[request.Settings.ResourcePlugin]
	if plugin == nil {
		c.logger.Error("selected driver is not a installed resource-driver-plugin", "selected-driver", request.Settings.ResourcePlugin)
		response.Error = &pb_core.Error{
			Message: "selected driver is missing/does not exist",
			Code:    ErrInvalidResource,
		}
		return nil, NewError(response.Error.Message, ErrorCode(response.Error.Code))
	}

	// implementation of different requests that SWP-X can handle right now
	switch request.Type {
	case pb_core.Request_GET_BASIC_INFO:

		if request.Port != "" {
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

	case pb_core.Request_GET_TECHNICAL_INFO:

		if request.Port != "" {
			err := c.handleGetTechnicalInformationPort(request, response, plugin)
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

	// process the response with the selected providers (post-process)
	if err := providerPostProcess(ctx, selectedProviders, response); err != nil {
		return nil, err
	}

	return response, nil

}

// handleGetTechnicalInformationElement gets full information of an Element
func (c *Core) getTechnicalInformationElement(msg *pb_core.Request, resp *pb_core.Response, plugin shared.Resource) error {

	// req := &resource.Request{
	// 	Port:     "",
	// 	Hostname: msg.Hostname,
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
	// allPortInformation, err = plugin.GetAllTransceiverInformation(msg.ctx, &resource.NetworkElementWrapper{
	// 	Element:        req,
	// 	NumInterfaces:  matchingInterfaces,
	// 	FullElement:    allPortInformation,
	// 	PhysInterfaces: physPortResponse,
	// })
	// if err != nil {
	// 	c.logger.Error("error fetching transceiver information: ", err)
	// }

	// resp.NetworkElement = allPortInformation

	return nil
}

// handleGetTechnicalInformationPort gets information related to the selected interface
func (c *Core) handleGetTechnicalInformationPort(msg *pb_core.Request, resp *pb_core.Response, plugin shared.Resource) error {
	// protoConf := shared.Conf2proto(conf)
	// protoConf.Request = c.createRequestConfig(msg, conf) // set deadline
	// req := &resource.NetworkElement{
	// 	Hostname:  msg.Hostname,
	// 	Interface: msg.Port,
	// 	Conf:      protoConf,
	// }

	// var mapInterfaceResponse *resource.NetworkElementInterfaces
	// var cachedInterface *CachedInterface
	// var err error

	// if c.cacheEnabled && !msg.Settings.RecreateIndex {
	// 	c.logger.Info("cache is enabled, pop index from cache")
	// 	cachedInterface, err = c.interfaceCache.Pop(context.TODO(), req.Hostname, req.Interface)
	// 	if cachedInterface != nil {
	// 		resp.PhysicalPort = cachedInterface.Port
	// 		req.PhysicalIndex = cachedInterface.PhysicalEntityIndex
	// 		req.InterfaceIndex = cachedInterface.InterfaceIndex
	// 	}
	// }

	// // did not find cached item or cached is disabled
	// if cachedInterface == nil || !c.cacheEnabled {
	// 	var physPortResponse *resource.NetworkElementInterfaces
	// 	c.logger.Info("run mapEntity to get physical entity index on device")
	// 	if physPortResponse, err = plugin.MapEntityPhysical(msg.ctx, req); err != nil {
	// 		c.logger.Error("error running MapEntityPhysical", "err", err.Error())
	// 		resp.Error = &pb_core.Error{
	// 			Message: err.Error(),
	// 			Code:    ErrInvalidPort,
	// 		}
	// 		return err
	// 	}

	// 	if val, ok := physPortResponse.Interfaces[req.Interface]; ok {
	// 		resp.PhysicalPort = val.Description
	// 		req.PhysicalIndex = val.Index
	// 	}

	// 	if mapInterfaceResponse, err = plugin.MapInterface(msg.ctx, req); err != nil {
	// 		c.logger.Error("error running map interface", "err", err.Error())
	// 		resp.Error = &pb_core.Error{
	// 			Message: err.Error(),
	// 			Code:    ErrInvalidPort,
	// 		}
	// 		return err
	// 	}
	// 	if val, ok := mapInterfaceResponse.Interfaces[req.Interface]; ok {
	// 		req.InterfaceIndex = val.Index
	// 	}

	// 	// save in cache upon success (if enabled)
	// 	if c.cacheEnabled {
	// 		if err = c.interfaceCache.Upsert(context.TODO(), req, mapInterfaceResponse, physPortResponse); err != nil {
	// 			return err
	// 		}
	// 	}

	// } else if err != nil {
	// 	c.logger.Error("error fetching from cache:", err.Error())
	// 	return err
	// }

	// //if the return is 0 something went wrong
	// if req.InterfaceIndex == 0 {
	// 	c.logger.Error("error running map interface", "err", "index is zero")
	// 	resp.Error = &pb_core.Error{
	// 		Message: "interface index returned zero",
	// 		Code:    ErrInvalidPort,
	// 	}
	// 	return err
	// }

	// c.logger.Info("found index for selected interface", "index", req.InterfaceIndex)

	// ti, err := plugin.TechnicalPortInformation(msg.ctx, req)
	// if err != nil {
	// 	c.logger.Error(err.Error())
	// 	return err
	// }
	// resp.NetworkElement = ti

	return nil
}

// handleGetBasicInformationPort gets information related to the selected interface
func (c *Core) handleGetBasicInformationPort(ctx context.Context, msg *pb_core.Request, plugin shared.Resource) (*pb_core.Response, error) {

	var resp pb_core.Response
	req := proto.Request{
		Hostname: msg.Hostname,
		Port:     msg.Port,
	}

	var (
		mapInterfaceResponse *proto.PortIndex
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
				resp.Error = &pb_core.Error{
					Message: fmt.Sprintf("could not run MapEntityPhyiscal: %s", err.Error()),
					Code:    ErrInvalidPort,
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
			resp.Error = &pb_core.Error{
				Message: fmt.Sprintf("could not run MapInterface: %s", err.Error()),
				Code:    ErrInvalidPort,
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

		resp.Error = &pb_core.Error{
			Message: fmt.Sprintf("could handle request: %s", err.Error()),
			Code:    ErrUnknownError,
		}

		c.logger.Error("error fetching from cache:", err.Error())
		return nil, err
	}

	//if the return is 0 something went wrong
	if req.LogicalPortIndex == 0 {
		c.logger.Error("error running map interface", "err", "index is zero")
		resp.Error = &pb_core.Error{
			Message: "interface index returned zero",
			Code:    ErrInvalidPort,
		}
		return nil, err
	}

	c.logger.Info("found index for selected interface", "index", req.LogicalPortIndex)

	ne, err := plugin.BasicPortInformation(ctx, &req)
	if err != nil {
		c.logger.Error(err.Error())
		return nil, err
	}
	resp.NetworkElement = ne

	return &resp, nil
}

// handleGetTechnicalInformationElement gets full information of an Element
func (c *Core) handleGetPasicInformationElement(msg *pb_core.Request, resp *pb_core.Response, plugin shared.Resource) error {
	// protoConf := shared.Conf2proto(conf)
	// protoConf.Request = c.createRequestConfig(msg, conf) // set deadline
	// req := &resource.NetworkElement{
	// 	Interface: "",
	// 	Hostname:  msg.Hostname,
	// 	Conf:      protoConf,
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
	// allPortInformation, err = plugin.GetAllTransceiverInformation(msg.ctx, &resource.NetworkElementWrapper{
	// 	Element:        req,
	// 	NumInterfaces:  matchingInterfaces,
	// 	FullElement:    allPortInformation,
	// 	PhysInterfaces: physPortResponse,
	// })
	// if err != nil {
	// 	c.logger.Error("error fetching transceiver information: ", err)
	// }

	// resp.NetworkElement = allPortInformation

	return nil
}
