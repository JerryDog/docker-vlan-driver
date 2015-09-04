package main

import (
	"github.com/docker/libnetwork/drivers/remote/api"
)

var (
	Networks map[string]*Network = make(map[string]*Network, 0)
)

type Network struct {
	*api.CreateNetworkRequest
	Name      string
	Eth       string
	VLanID    int
	Endpoints map[string]*Endpoint
}

func NewNetwork(req *api.CreateNetworkRequest) *Network {
	n := &Network{
		req, "", "", 100,
		make(map[string]*Endpoint, 0),
	}
	return n
}

func (n *Network) Create() error {
	//TODO:
	return nil
}

func (Network) Delete() error {
	//TODO:
	return nil
}
