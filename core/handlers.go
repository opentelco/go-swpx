package core

import (
	"context"
	"encoding/json"
	"fmt"

	pb_core "git.liero.se/opentelco/go-swpx/proto/go/core"
	"git.liero.se/opentelco/go-swpx/proto/go/resource"
	"git.liero.se/opentelco/go-swpx/shared"
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
		logger.Error("timeout reached or context cancelled")
		return ctx.Err()
	case err := <-c:
		if err != nil {
			logger.Error("err: ", err.Error())
		}
		return err
	}
}

func handleMsg(msg *Request, resp *pb_core.Response) error {
	var err error
	logger.Info("selected provider", "provider", msg.ProviderPlugin)

	// TODO what to do if this is empty? Should fallback on default? change to pointer so we can check if == nil ?
	var providerConf *shared.Configuration
	defaultConf := shared.GetConfig()

	var selectedProvider shared.Provider
	// check if a provider is selected in the request
	if msg.ProviderPlugin != "" {
		var err error
		selectedProvider = providers[msg.ProviderPlugin]
		if selectedProvider == nil {
			resp.Error = &pb_core.Error{Message: "the provider is missing/does not exist", Code: ErrInvalidProvider}
			return NewError(resp.Error.Message, ErrorCode(resp.Error.Code))
		}
		// Pre-process the request with provider func
		msg.Request, err = selectedProvider.PreHandler(context.Background(), msg.Request)
		if err != nil {
			return err
		}
		providerConf = defaultConf

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

	err = plugin.SetConfiguration(msg.Context, providerConf)
	if err != nil {
		return nil
	}

	// implementation of different messages that SWP-X can handle right now
	// TODO is this the best way to to this.. ?
	switch msg.Type {
	case pb_core.Request_GET_TECHNICAL_INFO:
		err := handleGetTechnicalInformationElement(msg, resp, plugin, providerConf)
		if err != nil {
			return err
		}
	case pb_core.Request_GET_TECHNICAL_INFO_PORT:
		err := handleGetTechnicalInformationPort(msg, resp, plugin, providerConf)
		if err != nil {
			return err
		}
	}

	if selectedProvider != nil {
		// Post-Process the provider

		nr, err := selectedProvider.PostHandler(context.Background(), resp)
		if err != nil {
			return nil
		}
		// cant keep the pointer over the wire
		resp.NetworkElement = nr.NetworkElement
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

	if useCache && !msg.RecreateIndex {
		logger.Debug("cache enabled, pop object from cache")
		cachedInterface, err = CacheInterface.Pop(context.TODO(), req.Hostname, req.Interface)
		if cachedInterface != nil {
			resp.PhysicalPort = cachedInterface.Port
			req.PhysicalIndex = cachedInterface.PhysicalEntityIndex
			req.InterfaceIndex = cachedInterface.InterfaceIndex
		}
	}

	js, _ := json.MarshalIndent(req, "", "  ")
	fmt.Println(string(js))
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
			if err = CacheInterface.Upsert(context.TODO(), req, mapInterfaceResponse, physPortResponse); err != nil {
				return err
			}
		}

	} else if err != nil {
		logger.Error("error fetching from cache:", err.Error())
		return err
	}

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
