package main

import (
	"os/exec"

	"github.com/docker/libnetwork/drivers/remote/api"
	"github.com/vishvananda/netlink"
)

type Endpoint struct {
	*api.CreateEndpointRequest
	EndpointShortID string
}

func NewEndpoint(req *api.CreateEndpointRequest) *Endpoint {
	e := &Endpoint{req, ""}
	if len(e.EndpointID) < 5 {
		e.EndpointShortID = e.EndpointID
	} else {
		e.EndpointShortID = e.EndpointID[:5]
	}
	return e
}

func (e *Endpoint) Create(network *Network) error {
	ethLink, err := netlink.LinkByName(network.Eth)
	if err != nil {
		return err
	}

	attrs := netlink.NewLinkAttrs()
	attrs.Name = "veth" + e.EndpointShortID
	attrs.ParentIndex = ethLink.Attrs().Index

	vlan := &netlink.Vlan{
		attrs, network.VLanID,
	}

	if err := netlink.LinkAdd(vlan); err != nil {
		return err
	}

	cmd := exec.Command("ip", "netns", "add", e.EndpointShortID)
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("ip", "link", "set", "veth"+e.EndpointShortID, "netns", e.EndpointShortID)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

func (e *Endpoint) Delete() error {
	if err := exec.Command("ip", "netns", "delete", e.EndpointShortID).Run(); err != nil {
		return err
	}
	return nil
}
