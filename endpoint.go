package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
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
	NetNS           string
}

func NewEndpoint(req *api.CreateEndpointRequest) *Endpoint {
	e := &Endpoint{req, "", "", "", "", "", ""}
	if len(e.EndpointID) < 5 {
		e.EndpointShortID = e.EndpointID
	} else {
		e.EndpointShortID = e.EndpointID[:5]
	}
	return e
}

func (e *Endpoint) Create(network *Network) error {

	e.updateNetworkInfo(network)

	return nil
}

func (e *Endpoint) makeIface(network *Network, netns string) error {
	if vethLink, _ := netlink.LinkByName("veth" + e.EndpointShortID); vethLink != nil {
		log.Println("veth"+e.EndpointShortID, "already exist")
	} else {
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
			log.Println("Err: ", err, "LinkAdd")
			return err
		}
	}

	if err := exec.Command("ip", "link", "set", "veth"+e.EndpointShortID, "netns", netns).Run(); err != nil {
		log.Println("Err: ", err, "LinkSet NS")
		return err
	}
	return nil
}

func (e *Endpoint) activate(network *Network, netns string) error {

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

	_name := strings.Split(e.Name, "-")
	e.IPStr = strings.Replace(_name[0], "_", ".", -1)
	e.SubnetStr = _name[1]
	ipstr := e.IPStr + "/" + e.SubnetStr

	fmt.Println("ip:", e.IPStr, "subnet:", e.SubnetStr)

	cmd := []string{
		"netns", "exec", netns,
		"ip", "addr", "add", ipstr, "dev", "veth" + e.EndpointShortID}

	if err := exec.Command("ip", cmd...).Run(); err != nil {
		return err
	}

	cmd2 := []string{
		"netns", "exec", netns,
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

func (e *Endpoint) Join(network *Network, sandboxKey string) error {
	e.SandboxKey = sandboxKey
	if e.SandboxKey == "" {
		return fmt.Errorf("SandboxKey is empty!")
	}

	e.NetNS = filepath.Base(e.SandboxKey)
	log.Println("netns:", e.NetNS)

	//ln -sf to /var/run/netns
	if err := os.MkdirAll("/var/run/netns", 0777); err != nil {
		return err
	}
	if err := exec.Command("ln", "-sf", e.SandboxKey, "/var/run/netns").Run(); err != nil {
		return err
	}

	if err := e.makeIface(network, e.NetNS); err != nil {
		return err
	}
	if err := e.activate(network, e.NetNS); err != nil {
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
