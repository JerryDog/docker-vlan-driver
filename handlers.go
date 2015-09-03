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
	return c.JSON(http.StatusInternalServerError, api.Response{Err: err.Error()})
}

func PluginActivate(c *echo.Context) error {
	resp := map[string][]string{
		"Implements": []string{"NetworkDriver"},
	}
	return c.JSON(http.StatusOK, resp)
}

func NetworkDriverCreateNetwork(c *echo.Context) error {
	req := api.CreateNetworkRequest{}
	err := _json(c, &req)
	if err != nil {
		return _error(c, err)
	}

	network := &Network{req}
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
	req := api.DeleteEndpointRequest{}
	err := _json(c, &req)
	if err != nil {
		return _error(c, err)
	}

	//Networks[req.NetworkID].Delete()
	delete(Networks, req.NetworkID)

	return c.JSON(http.StatusOK, api.DeleteEndpointResponse{})
}

func NetworkDriverCreateEndpoint(c *echo.Context) error {
	req := api.CreateEndpointRequest{}
	if err := _json(c, &req); err != nil {
		return _error(c, err)
	}
	fmt.Println("CreateEndpointRequest:", req)
	return c.JSON(http.StatusOK, api.CreateEndpointResponse{})
}

func NetworkDriverEndpointOperInfo(c *echo.Context) error {
	req := api.EndpointInfoRequest{}
	if err := _json(c, &req); err != nil {
		return _error(c, err)
	}
	fmt.Println("EndpointInfoResponse:", req)
	return c.JSON(http.StatusOK, api.EndpointInfoResponse{})
}

func NetworkDriverDeleteEndpoint(c *echo.Context) error {
	req := api.DeleteEndpointRequest{}
	if err := _json(c, &req); err != nil {
		return _error(c, err)
	}
	fmt.Println("DeleteEndpoint:", req)
	return c.JSON(http.StatusOK, api.DeleteEndpointResponse{})
}

func NetworkDriverJoin(c *echo.Context) error {
	req := api.JoinRequest{}
	if err := _json(c, &req); err != nil {
		return _error(c, err)
	}
	fmt.Println("Join:", req)
	return c.JSON(http.StatusOK, api.JoinResponse{})
}

func NetworkDriverLeave(c *echo.Context) error {
	req := api.LeaveRequest{}
	if err := _json(c, &req); err != nil {
		return _error(c, err)
	}
	fmt.Println("Leave:", req)
	return c.JSON(http.StatusOK, api.LeaveResponse{})
}
