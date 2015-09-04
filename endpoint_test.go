package main

import (
	"os/exec"
	"testing"

	"github.com/docker/libnetwork/drivers/remote/api"
	"github.com/vishvananda/netlink"
)

func TestEndpointCreate(t *testing.T) {
	network := NewNetwork(&api.CreateNetworkRequest{})
	network.Eth = "eth0"
	network.VLanID = 101

	if _, err := netlink.LinkByName(network.Eth); err != nil {
		if err := exec.Command("ip", "link", "add", "eth0", "type", "veth", "peer", "name", "eth1").Run(); err != nil {
			t.Error(err)
		}
	}

	endpoint := NewEndpoint(&api.CreateEndpointRequest{})
	endpoint.EndpointID = "1234567890123456890"
	endpoint.EndpointShortID = endpoint.EndpointID[:5]
	err := endpoint.Create(network)
	if err != nil {
		t.Error(err)
	}

	// Tear down
	/*
			if err := exec.Command("ip", "link", "delete", "veth"+endpoint.EndpointShortID).Run(); err != nil {
				t.Error(err)
			}
		if err := exec.Command("ip", "netns", "delete", endpoint.EndpointShortID).Run(); err != nil {
			t.Error(err)
		}
	*/
}
