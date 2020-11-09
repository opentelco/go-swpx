package core

import (
	pb_core "git.liero.se/opentelco/go-swpx/proto/go/core"
	"git.liero.se/opentelco/go-swpx/proto/go/resource"
	"git.liero.se/opentelco/go-swpx/shared"
	"context"
	"fmt"
)

// TODO this just runs some functions.. not a real implementation
func providerFunc(provider shared.Provider, msg *Request) {
	name, err := provider.Name()
	if err != nil {
		logger.Debug("getting provider name failed", "error", err)
	}
	ver, err := provider.Version()
	if err != nil {
		logger.Debug("getting provider name failed", "error", err)
	}
	
	logger.Debug("data from provider plugin", "provider", name, "version", ver)
}

func handle(ctx context.Context, msg *Request, resp *pb_core.Response, f func(msg *Request, resp *pb_core.Response) error) error {
	c := make(chan error, 1)
	go func() { c <- f(msg, resp) }()
	select {
	case <-ctx.Done():
		logger.Error("got a timeout, letrs go")
		return ctx.Err()
	case err := <-c:
		if err != nil {
			logger.Error("err: ", err.Error())
		}
		return err
	}
}

func handleMsg(msg *Request, resp *pb_core.Response) error {
	logger.Debug("worker has payload")
	logger.Info("selected provider", "provider", msg.ProviderPlugin)
	
	// TODO what to do if this is empty? Should fallback on default? change to pointer so we can check if == nil ?
	var providerConf *shared.Configuration
	defaultConf := shared.GetConfig()
	
	// check if a provider is selected in the request
	if msg.ProviderPlugin != "" {
		provider := providers[msg.ProviderPlugin]
		if provider == nil {
			resp.Error = &pb_core.Error{Message: "the provider is missing/does not exist", Code: ErrInvalidProvider}
			return NewError(resp.Error.Message, ErrorCode(resp.Error.Code))
		}
		// run some provider funcs
		providerFunc(provider, msg)
		providerConf = defaultConf
		
	} else {
		// no provider selected, walk all providers
		for pname, provider := range providers {
			logger.Debug("parsing provider", "provider", pname)
			providerFunc(provider, msg)
		}
	}
	
	// select resource-plugin to send the requests to
	plugin := resources[msg.ResourcePlugin]
	if plugin == nil {
		logger.Error("selected driver is not a installed resource-driver-plugin", "selected-driver", msg.ResourcePlugin)
		resp.Error = &pb_core.Error{
			Message: "selected driver is missing/does not exist",
			Code:    ErrInvalidResource,
		}
		return nil
	}
	plugin.SetConfiguration(msg.Context, providerConf)
	
	// implementation of different messages that SWP-X can handle right now
	// TODO is this the best way to to this.. ?
	switch msg.Type {
	case pb_core.Request_GET_TECHNICAL_INFO:
		return handleGetTechnicalInformationElement(msg, resp, plugin, providerConf)
	case pb_core.Request_GET_TECHNICAL_INFO_PORT:
		return handleGetTechnicalInformationPort(msg, resp, plugin, providerConf)
	}
	
	return nil
}

// handleGetTechnicalInformationElement gets full information of an Element
func handleGetTechnicalInformationElement(msg *Request, resp *pb_core.Response, plugin shared.Resource, conf *shared.Configuration) error {
	protoConf := shared.Conf2proto(conf)
	
	req := &resource.NetworkElement{
		Interface: "",
		Hostname:  msg.Hostname,
		Conf:      protoConf,
	}
	
	physPortResponse, err := plugin.MapEntityPhysical(msg.Context, req)
	if err != nil {
		logger.Error("error fetching physical entities:", err.Error())
		return err
	}
	
	allPortInformation, err := plugin.AllPortInformation(msg.Context, req)
	if err != nil {
		logger.Error("error fetching port information for all interfaces:", err.Error())
		return err
	}
	
	var matchingInterfaces int32 = 0
	for _, iface := range allPortInformation.Interfaces {
		if _, ok := physPortResponse.Interfaces[iface.Description]; ok {
			matchingInterfaces++
		}
	}
	allPortInformation, err = plugin.GetAllTransceiverInformation(msg.Context, &resource.NetworkElementWrapper{
		Element:        req,
		NumInterfaces:  matchingInterfaces,
		FullElement:    allPortInformation,
		PhysInterfaces: physPortResponse,
	})
	if err != nil {
		logger.Error("error fetching transceiver information: ", err)
	}
	
	resp.NetworkElement = allPortInformation
	
	return nil
}

// handleGetTechnicalInformationPort gets information related to the selected interface
func handleGetTechnicalInformationPort(msg *Request, resp *pb_core.Response, plugin shared.Resource, conf *shared.Configuration) error {
	protConf := shared.Conf2proto(conf)
	req := &resource.NetworkElement{
		Hostname:  msg.Hostname,
		Interface: msg.Port,
		Conf:      protConf,
	}
	
	mapInterfaceResponse := &resource.NetworkElementInterfaces{}
	var cachedInterface *CachedInterface
	var err error
	
	if useCache && !msg.RecreateIndex{
		logger.Debug("cache enabled, pop object from cache")
		cachedInterface, err = InterfaceCache.PopInterface(req.Hostname, req.Interface)
		if cachedInterface != nil {
			resp.PhysicalPort = cachedInterface.Port
			req.PhysicalIndex = cachedInterface.PhysicalEntityIndex
			req.InterfaceIndex = cachedInterface.InterfaceIndex
		}
	}
	
	// did not find cached item or cached is disabled
	if cachedInterface == nil || !useCache {
		var physPortResponse *resource.NetworkElementInterfaces
		logger.Debug("run mapEntity")
		if physPortResponse, err = plugin.MapEntityPhysical(msg.Context, req); err != nil {
			logger.Error("error running getphysport", "err", err.Error())
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
		
		if mapInterfaceResponse, err = plugin.MapInterface(msg.Context, req); err != nil {
			logger.Error("error running map interface", "err", err.Error())
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
		if useCache {
			if err = InterfaceCache.SetInterface(req, mapInterfaceResponse, physPortResponse); err != nil {
				return err
			}
		}
		
	} else if err != nil {
		logger.Error("error fetching from cache:", err.Error())
		return err
	}
	
	fmt.Println(req)
	//if the return is 0 something went wrong
	if req.InterfaceIndex == 0 {
		logger.Error("error running map interface", "err", "index is zero")
		resp.Error = &pb_core.Error{
			Message: "interface index returned zero",
			Code:    ErrInvalidPort,
		}
		return err
	}
	
	logger.Info("found index for selected interface", "index", req.InterfaceIndex)
	
	ti, err := plugin.TechnicalPortInformation(msg.Context, req)
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	logger.Info("calling technical info ok ", "result", ti)
	resp.NetworkElement = ti
	
	return nil
}
