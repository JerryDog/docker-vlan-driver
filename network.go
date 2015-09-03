package main

import (
	"github.com/docker/libnetwork/drivers/remote/api"
)

type Network struct {
	api.CreateNetworkRequest
}

var (
	Networks map[string]*Network = make(map[string]*Network, 0)
)

func (Network) Create() error {
	return nil
}

func (Network) Delete() error {
	return nil
}
