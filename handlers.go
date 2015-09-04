package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/docker/libnetwork/drivers/remote/api"
	"github.com/labstack/echo"
)

func _json(c *echo.Context, req interface{}) error {
	decoder := json.NewDecoder(c.Request().Body)
	return decoder.Decode(req)
}

func _error(c *echo.Context, err error) error {
	fmt.Println("_error:", err)
	return c.JSON(http.StatusInternalServerError, api.Response{Err: err.Error()})
}

func PluginActivate(c *echo.Context) error {
	resp := map[string][]string{
		"Implements": []string{"NetworkDriver"},
	}
	return c.JSON(http.StatusOK, resp)
}

func NetworkDriverCreateNetwork(c *echo.Context) error {
	var req api.CreateNetworkRequest
	err := _json(c, &req)
	if err != nil {
		return _error(c, err)
	}

	network := NewNetwork(&req)
	if err := network.Create(); err != nil {
		return _error(c, err)
	}

	//TODO(anarcher):
	Networks[network.NetworkID] = network
	fmt.Println(network)
	fmt.Println(Networks, len(Networks))

	return c.JSON(http.StatusOK, api.CreateNetworkResponse{})
}

func NetworkDriverDeleteNetwork(c *echo.Context) error {
	var req api.DeleteEndpointRequest
	err := _json(c, &req)
	if err != nil {
		return _error(c, err)
	}

	//Networks[req.NetworkID].Delete()
	delete(Networks, req.NetworkID)

	return c.JSON(http.StatusOK, api.DeleteEndpointResponse{})
}

func NetworkDriverCreateEndpoint(c *echo.Context) error {
	var req api.CreateEndpointRequest
	if err := _json(c, &req); err != nil {
		return _error(c, err)
	}

	endpointID := req.EndpointID
	networkID := req.NetworkID

	network, ok := Networks[networkID]
	if !ok {
		return _error(c, fmt.Errorf("network not exist"))
	}

	endpoint := NewEndpoint(&req)
	if err := endpoint.Create(network); err != nil {
		return _error(c, err)
	}
	network.Endpoints[endpointID] = endpoint

	//TODO
	resp := api.CreateEndpointResponse{}
	/*
			resp.Interfaces = []*api.EndpointInterface{
				&api.EndpointInterface{
					ID:         0,
					Address:    "10.1.1.1/24",
					MacAddress: "4a:6a:c1:57:af:98",
				},
			}
		fmt.Println("resp:", resp)
	*/

	return c.JSON(http.StatusOK, resp)
}

func NetworkDriverEndpointOperInfo(c *echo.Context) error {
	fmt.Println("NetworkDriverEndpointOperInfo.....")

	var req api.EndpointInfoRequest
	if err := _json(c, &req); err != nil {
		return _error(c, err)
	}

	return c.JSON(http.StatusOK, api.EndpointInfoResponse{})
}

func NetworkDriverDeleteEndpoint(c *echo.Context) error {
	var req api.DeleteEndpointRequest
	if err := _json(c, &req); err != nil {
		return _error(c, err)
	}

	network, ok := Networks[req.NetworkID]
	if !ok {
		return _error(c, fmt.Errorf("network id not found"))
	}
	endpoint := network.Endpoints[req.EndpointID]
	if !ok {
		return _error(c, fmt.Errorf("endpoint id not found"))
	}

	if err := endpoint.Delete(); err != nil {
		return _error(c, err)
	}

	delete(network.Endpoints, req.EndpointID)
	return c.JSON(http.StatusOK, api.DeleteEndpointResponse{})
}

func NetworkDriverJoin(c *echo.Context) error {
	var req api.JoinRequest
	if err := _json(c, &req); err != nil {
		return _error(c, err)
	}
	fmt.Println("Join:", req)

	network, ok := Networks[req.NetworkID]
	if !ok {
		return _error(c, fmt.Errorf("network id not found"))
	}
	endpoint := network.Endpoints[req.EndpointID]
	if !ok {
		return _error(c, fmt.Errorf("endpoint id not found"))
	}

	err := endpoint.Join(network, req.SandboxKey)
	if err != nil {
		return _error(c, err)
	}

	resp := api.JoinResponse{}
	//TODO:
	/*
		resp.InterfaceNames = []*api.InterfaceName{
			&api.InterfaceName{
				SrcName:   "veth1",
				DstPrefix: "eth",
			},
		}
		resp.Gateway = "10.1.1.1"
	*/

	return c.JSON(http.StatusOK, resp)
}

func NetworkDriverLeave(c *echo.Context) error {
	var req api.LeaveRequest
	if err := _json(c, &req); err != nil {
		return _error(c, err)
	}

	network, ok := Networks[req.NetworkID]
	if !ok {
		return _error(c, fmt.Errorf("network id not found"))
	}
	endpoint := network.Endpoints[req.EndpointID]
	if !ok {
		return _error(c, fmt.Errorf("endpoint id not found"))
	}

	err := endpoint.Leave()
	if err != nil {
		return _error(c, err)
	}

	fmt.Println("Leave:", req)
	return c.JSON(http.StatusOK, api.LeaveResponse{})
}
