package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/docker/libnetwork/drivers/remote/api"
	"github.com/fsouza/go-dockerclient"
	"github.com/vishvananda/netlink"
)

type Endpoint struct {
	*api.CreateEndpointRequest
	Name            string
	EndpointShortID string
	IPStr           string
	SubnetStr       string
	SandboxKey      string
}

func NewEndpoint(req *api.CreateEndpointRequest) *Endpoint {
	e := &Endpoint{req, "", "", "", "", ""}
	if len(e.EndpointID) < 5 {
		e.EndpointShortID = e.EndpointID
	} else {
		e.EndpointShortID = e.EndpointID[:5]
	}
	return e
}

func (e *Endpoint) Create(network *Network) error {

	e.updateNetworkInfo(network)

	ethLink, err := netlink.LinkByName(network.Eth)
	if err != nil {
		fmt.Println("[Err] LinkByName:", err)
		return err
	}

	attrs := netlink.NewLinkAttrs()
	attrs.Name = "veth" + e.EndpointShortID
	attrs.ParentIndex = ethLink.Attrs().Index

	vlan := &netlink.Vlan{
		attrs, network.VLanID,
	}

	if err := netlink.LinkAdd(vlan); err != nil {
		fmt.Println("[Err] LinkAdd:", err)
		return err
	}

	cmd := exec.Command("ip", "netns", "add", e.EndpointShortID)
	if err := cmd.Run(); err != nil {
		fmt.Println("cmd", cmd, "err", err)
		return err
	}

	cmd = exec.Command("ip", "link", "set", "veth"+e.EndpointShortID, "netns", e.EndpointShortID)
	if err := cmd.Run(); err != nil {
		fmt.Println("cmd", cmd, "err", err)
		return err
	}
	return nil
}

func (e *Endpoint) Activate(network *Network) error {

	networkInfo, err := e.getNetworkInfoFromDocker(network.NetworkID)
	if err != nil {
		return err
	}
	for _, eInfo := range networkInfo.Endpoints {
		if eInfo.ID == e.EndpointID {
			e.Name = eInfo.Name
			break
		}
	}

	fmt.Println("Activate: name:", e.Name, "endpointID", e.EndpointID)

	_name := strings.Split(e.Name, "/")
	e.IPStr = strings.Replace(_name[0], "_", ".", -1)
	e.SubnetStr = _name[1]
	ipstr := e.IPStr + "/" + e.SubnetStr

	fmt.Println("ip:", e.IPStr, "subnet:", e.SubnetStr)

	cmd := []string{
		"netns", "exec", e.EndpointShortID,
		"ip", "addr", "add", ipstr, "dev", "veth" + e.EndpointShortID}

	if err := exec.Command("ip", cmd...).Run(); err != nil {
		return err
	}

	cmd2 := []string{
		"netns", "exec", e.EndpointShortID,
		"ip", "link", "set", "veth" + e.EndpointShortID, "up"}

	if err := exec.Command("ip", cmd2...).Run(); err != nil {
		return err
	}

	return nil
}

func (Endpoint) getNetworkInfoFromDocker(networkID string) (*docker.Network, error) {
	//TODO:
	dockerEp := "unix:///var/run/docker.sock"
	client, err := docker.NewClient(dockerEp)
	if err != nil {
		return nil, err
	}
	nwInfo, err := client.NetworkInfo(networkID)
	if err != nil {
		return nil, err
	}

	return nwInfo, nil
}

func (e *Endpoint) updateNetworkInfo(network *Network) error {
	//TODO:
	nwInfo, err := e.getNetworkInfoFromDocker(network.NetworkID)
	if err != nil {
		return err
	}
	network.Name = nwInfo.Name

	if name := strings.Split(network.Name, "-"); len(name) > 2 {
		network.Eth = name[0]
		network.VLanID, err = strconv.Atoi(name[1])
		if err != nil {
			return err
		}
	} else {
		network.Eth = "eth0"
		network.VLanID = 100
	}

	return nil
}

func (e *Endpoint) Delete() error {
	if err := exec.Command("ip", "netns", "delete", e.EndpointShortID).Run(); err != nil {
		return err
	}
	return nil
}

func (e *Endpoint) Join() error {
	cmd := []string{
		"netns", "exec", e.EndpointShortID,
		"ip", "link", "add", "veth0", "peer", "name", "veth1",
	}
	if err := exec.Command("ip", cmd...).Run(); err != nil {
		return err
	}

	return nil
}

func (e *Endpoint) Leave() error {
	cmd := []string{
		"netns", "exec", e.EndpointShortID,
		"ip", "link", "delete", "veth0",
	}
	if err := exec.Command("ip", cmd...).Run(); err != nil {
		return err
	}

	return nil
}
