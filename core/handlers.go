package core

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	pb_core "git.liero.se/opentelco/go-swpx/proto/go/core"
	"git.liero.se/opentelco/go-swpx/proto/go/provider"
	"git.liero.se/opentelco/go-swpx/proto/go/resource"
	"git.liero.se/opentelco/go-swpx/shared"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (c *Core) RequestHandler(ctx context.Context, request *Request, response *pb_core.Response) error {

	var err error

	// TODO what to do if this is empty? Should fallback on default? change to pointer so we can check if == nil ?
	var providerConf *shared.Configuration
	defaultConf := shared.GetConfig()

	var selectedProviders []shared.Provider
	// check if a providers are selected in the request
	if len(request.Settings.ProviderPlugin) > 0 {
		c.logger.Info("request has selected providers", "providers", strings.Join(request.Settings.ProviderPlugin, ","))

		for _, provider := range request.Settings.ProviderPlugin {
			var err error
			selectedProvider := c.providers[provider]
			if selectedProvider == nil {
				response.Error = &pb_core.Error{Message: "the provider is missing/does not exist", Code: ErrInvalidProvider}
				return NewError(response.Error.Message, ErrorCode(response.Error.Code))
			}

			// pre-process the request with provider func
			request.Request, err = selectedProvider.PreHandler(ctx, request.Request)
			if err != nil {
				return err
			}

			// add the provider to a slice for usage in the end
			selectedProviders = append(selectedProviders, selectedProvider)
			providerConf = defaultConf
		}
	} else {

		c.logger.Info("request has selected provider and default provider is set in config", "default_provider", c.config.DefaultProvider)
		if provider, ok := c.providers[c.config.DefaultProvider]; ok {

			request.Request, err = provider.PreHandler(ctx, request.Request)
			if err != nil {
				return err
			}

			selectedProviders = append(selectedProviders, provider)
		}

		providerConf = defaultConf
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
		return nil
	}

	err = plugin.SetConfiguration(ctx, providerConf)
	if err != nil {
		return nil
	}

	// implementation of different requests that SWP-X can handle right now
	switch request.Type {
	case pb_core.Request_GET_BASIC_INFO:

		if request.Port != "" {
			err := c.handleGetBasicInformationPort(request, response, plugin, providerConf)
			if err != nil {
				return err
			}

		} else {
			err := c.handleGetPasicInformationElement(request, response, plugin, providerConf)
			if err != nil {
				return err
			}

		}

	case pb_core.Request_GET_TECHNICAL_INFO:

		if request.Port != "" {
			err := c.handleGetTechnicalInformationPort(request, response, plugin, providerConf)
			if err != nil {
				return err
			}
		} else {
			err := c.handleGetTechnicalInformationElement(request, response, plugin, providerConf)
			if err != nil {
				return err
			}
		}

	}

	// PostProcess the response with the selected Providers
	for _, selectedProvider := range selectedProviders {
		if selectedProvider != nil {
			nr, err := selectedProvider.PostHandler(ctx, response)
			if err != nil {
				return nil
			}
			response.NetworkElement = nr.NetworkElement
		}
	}

	return nil

}

func CreateRequestConfig(msg *Request, conf *shared.Configuration) *provider.ConfigRequest {

	if msg.Request.Settings.Timeout != "" {
		dur, _ := time.ParseDuration(msg.Request.Settings.Timeout)
		if dur.Seconds() != 0 {
			return &provider.ConfigRequest{
				Deadline: timestamppb.New(time.Now().Add(dur)),
			}
		}

	}

	return &provider.ConfigRequest{
		Deadline: timestamppb.New(time.Now().Add(conf.DefaultRequestTimeout)),
	}

}

// handleGetTechnicalInformationElement gets full information of an Element
func (c *Core) handleGetTechnicalInformationElement(msg *Request, resp *pb_core.Response, plugin shared.Resource, conf *shared.Configuration) error {
	protoConf := shared.Conf2proto(conf)

	protoConf.Request = CreateRequestConfig(msg, conf) // set deadline
	req := &resource.NetworkElement{
		Interface: "",
		Hostname:  msg.Hostname,
		Conf:      protoConf,
	}

	physPortResponse, err := plugin.MapEntityPhysical(msg.ctx, req)
	if err != nil {
		c.logger.Error("error fetching physical entities:", err.Error())
		return err
	}

	allPortInformation, err := plugin.AllPortInformation(msg.ctx, req)
	if err != nil {
		c.logger.Error("error fetching port information for all interfaces:", err.Error())
		return err
	}

	var matchingInterfaces int32 = 0
	for _, iface := range allPortInformation.Interfaces {
		if _, ok := physPortResponse.Interfaces[iface.Description]; ok {
			matchingInterfaces++
		}
	}
	allPortInformation, err = plugin.GetAllTransceiverInformation(msg.ctx, &resource.NetworkElementWrapper{
		Element:        req,
		NumInterfaces:  matchingInterfaces,
		FullElement:    allPortInformation,
		PhysInterfaces: physPortResponse,
	})
	if err != nil {
		c.logger.Error("error fetching transceiver information: ", err)
	}

	resp.NetworkElement = allPortInformation

	return nil
}

// handleGetTechnicalInformationPort gets information related to the selected interface
func (c *Core) handleGetTechnicalInformationPort(msg *Request, resp *pb_core.Response, plugin shared.Resource, conf *shared.Configuration) error {
	protoConf.Request = CreateRequestConfig(msg, conf) // set deadline
	req := &resource.NetworkElement{
		Hostname:  msg.Hostname,
		Interface: msg.Port,
		Conf:      protConf,
	}

	var mapInterfaceResponse *resource.NetworkElementInterfaces
	var cachedInterface *CachedInterface
	var err error

	if c.cacheEnabled && !msg.Settings.RecreateIndex {
		c.logger.Info("cache is enabled, pop index from cache")
		cachedInterface, err = c.interfaceCache.Pop(context.TODO(), req.Hostname, req.Interface)
		if cachedInterface != nil {
			resp.PhysicalPort = cachedInterface.Port
			req.PhysicalIndex = cachedInterface.PhysicalEntityIndex
			req.InterfaceIndex = cachedInterface.InterfaceIndex
		}
	}

	// did not find cached item or cached is disabled
	if cachedInterface == nil || !c.cacheEnabled {
		var physPortResponse *resource.NetworkElementInterfaces
		c.logger.Info("run mapEntity to get physical entity index on device")
		if physPortResponse, err = plugin.MapEntityPhysical(msg.ctx, req); err != nil {
			c.logger.Error("error running MapEntityPhysical", "err", err.Error())
			resp.Error = &pb_core.Error{
				Message: err.Error(),
				Code:    ErrInvalidPort,
			}
			return err
		}

		if val, ok := physPortResponse.Interfaces[req.Interface]; ok {
			resp.PhysicalPort = val.Description
			req.PhysicalIndex = val.Index
		}

		if mapInterfaceResponse, err = plugin.MapInterface(msg.ctx, req); err != nil {
			c.logger.Error("error running map interface", "err", err.Error())
			resp.Error = &pb_core.Error{
				Message: err.Error(),
				Code:    ErrInvalidPort,
			}
			return err
		}
		if val, ok := mapInterfaceResponse.Interfaces[req.Interface]; ok {
			req.InterfaceIndex = val.Index
		}

		// save in cache upon success (if enabled)
		if c.cacheEnabled {
			if err = c.interfaceCache.Upsert(context.TODO(), req, mapInterfaceResponse, physPortResponse); err != nil {
				return err
			}
		}

	} else if err != nil {
		c.logger.Error("error fetching from cache:", err.Error())
		return err
	}

	//if the return is 0 something went wrong
	if req.InterfaceIndex == 0 {
		c.logger.Error("error running map interface", "err", "index is zero")
		resp.Error = &pb_core.Error{
			Message: "interface index returned zero",
			Code:    ErrInvalidPort,
		}
		return err
	}

	c.logger.Info("found index for selected interface", "index", req.InterfaceIndex)

	ti, err := plugin.TechnicalPortInformation(msg.ctx, req)
	if err != nil {
		c.logger.Error(err.Error())
		return err
	}
	resp.NetworkElement = ti

	return nil
}

// handleGetBasicInformationPort gets information related to the selected interface
func (c *Core) handleGetBasicInformationPort(msg *Request, resp *pb_core.Response, plugin shared.Resource, conf *shared.Configuration) error {
	protConf := shared.Conf2proto(conf)
	protoConf.Request = CreateRequestConfig(msg, conf) // set deadline
	req := &resource.NetworkElement{
		Hostname:  msg.Hostname,
		Interface: msg.Port,
		Conf:      protConf,
	}

	var mapInterfaceResponse *resource.NetworkElementInterfaces
	var cachedInterface *CachedInterface
	var err error

	if c.cacheEnabled && !msg.Settings.RecreateIndex {
		c.logger.Info("cache is enabled, pop index from cache")
		cachedInterface, err = c.interfaceCache.Pop(context.TODO(), req.Hostname, req.Interface)
		if cachedInterface != nil {
			c.logger.Info("cached interface indexs",
				"physicalPort", resp.PhysicalPort,
				"physicalIndex", req.PhysicalIndex,
				"interfaceIndex", req.InterfaceIndex,
			)
			resp.PhysicalPort = cachedInterface.Port
			req.PhysicalIndex = cachedInterface.PhysicalEntityIndex
			req.InterfaceIndex = cachedInterface.InterfaceIndex
		}
	}

	// did not find cached item or cached is disabled
	if cachedInterface == nil || !c.cacheEnabled {
		c.logger.Info("run mapEntity to get physical entity index on device")

		physPortResponse, err := plugin.MapEntityPhysical(msg.ctx, req)
		if err != nil {
			if status.Code(err) == codes.Unimplemented {
				c.logger.Error("warn MapEntityPhysical is not implemented, skipping", "err", err.Error())
			} else {
				c.logger.Error("error running MapEntityPhysical", "err", err.Error())
				resp.Error = &pb_core.Error{
					Message: err.Error(),
					Code:    ErrInvalidPort,
				}
				return err
			}

		} else {
			js, _ := json.MarshalIndent(physPortResponse, "", "  ")
			c.logger.Error("map", "data", string(js))

			if val, ok := physPortResponse.Interfaces[req.Interface]; ok {
				resp.PhysicalPort = val.Description
				req.PhysicalIndex = val.Index

				c.logger.Debug("found physInterface",
					"port", req.Interface,
					"resp.physicalPort", val.Description,
					"req.physicalIndex", req.PhysicalIndex,
				)

			}
		}

		if mapInterfaceResponse, err = plugin.MapInterface(msg.ctx, req); err != nil {
			c.logger.Error("error running map interface", "err", err.Error())
			resp.Error = &pb_core.Error{
				Message: err.Error(),
				Code:    ErrInvalidPort,
			}
			return err
		}

		if val, ok := mapInterfaceResponse.Interfaces[req.Interface]; ok {
			req.InterfaceIndex = val.Index
			c.logger.Info("found ifMIB interface index", "index", val.Index)
		}

		// save in cache upon success (if enabled)
		if c.cacheEnabled {
			if err = c.interfaceCache.Upsert(context.TODO(), req, mapInterfaceResponse, physPortResponse); err != nil {
				return err
			}
		}

	} else if err != nil {
		c.logger.Error("error fetching from cache:", err.Error())
		return err
	}

	//if the return is 0 something went wrong
	if req.InterfaceIndex == 0 {
		c.logger.Error("error running map interface", "err", "index is zero")
		resp.Error = &pb_core.Error{
			Message: "interface index returned zero",
			Code:    ErrInvalidPort,
		}
		return err
	}

	c.logger.Info("found index for selected interface", "index", req.InterfaceIndex)

	ti, err := plugin.BasicPortInformation(msg.ctx, req)
	if err != nil {
		c.logger.Error(err.Error())
		return err
	}
	resp.NetworkElement = ti

	return nil
}

// handleGetTechnicalInformationElement gets full information of an Element
func (c *Core) handleGetPasicInformationElement(msg *Request, resp *pb_core.Response, plugin shared.Resource, conf *shared.Configuration) error {
	protoConf := shared.Conf2proto(conf)
	protoConf.Request = CreateRequestConfig(msg, conf) // set deadline
	req := &resource.NetworkElement{
		Interface: "",
		Hostname:  msg.Hostname,
		Conf:      protoConf,
	}

	physPortResponse, err := plugin.MapEntityPhysical(msg.ctx, req)
	if err != nil {
		c.logger.Error("error fetching physical entities:", err.Error())
		return err
	}

	allPortInformation, err := plugin.AllPortInformation(msg.ctx, req)
	if err != nil {
		c.logger.Error("error fetching port information for all interfaces:", err.Error())
		return err
	}

	var matchingInterfaces int32 = 0
	for _, iface := range allPortInformation.Interfaces {
		if _, ok := physPortResponse.Interfaces[iface.Description]; ok {
			matchingInterfaces++
		}
	}
	allPortInformation, err = plugin.GetAllTransceiverInformation(msg.ctx, &resource.NetworkElementWrapper{
		Element:        req,
		NumInterfaces:  matchingInterfaces,
		FullElement:    allPortInformation,
		PhysInterfaces: physPortResponse,
	})
	if err != nil {
		c.logger.Error("error fetching transceiver information: ", err)
	}

	resp.NetworkElement = allPortInformation

	return nil
}
