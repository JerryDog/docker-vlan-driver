package main

import (
	"github.com/docker/libnetwork/drivers/remote/api"
	"github.com/fsouza/go-dockerclient"

	"strconv"
	"strings"
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
	dockerEp := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(dockerEp)
	if err != nil {
		return err
	}
	nwInfo, err := client.NetworkInfo(n.NetworkID)
	if err != nil {
		return err
	}

	n.Name = nwInfo.Name
	if name := strings.Split(n.Name, "-"); len(name) > 2 {
		n.Eth = name[0]
		n.VLanID, err = strconv.Atoi(name[1])
		if err != nil {
			return err
		}
	} else {
		n.Eth = "eth0"
		n.VLanID = 100
	}

	return nil
}

func (Network) Delete() error {
	//TODO:
	return nil
}
