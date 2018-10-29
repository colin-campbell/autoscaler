package openstack

import (
	"errors"
	nova "github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	neutron "github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

func getFloatingIP(p *provider) (interface{}, error) {
	if p.networkType == "neutron" {
		return getNeutronIP(p)
	} else {
		return getNovaIP(p)
	}
}

func releaseFloatingIP(p *provider, ip interface{}) error {
	if fIP, ok := p.floatingIps[ip]; ok  {
		switch v := fIP.(type) {
		case neutron.FloatingIP:
			return deleteNeutronFloatingIP(v)
		case nova.FloatingIP:
			return deleteNovaFloatingIP(v)
		}
	}
	return nil
}

func associateFloatingIP(p *provider, serverId string, ip string) error {
	if p.networkType == "neutron" {
		opts := ports.ListOpts{
			DeviceID: serverId,
		}
		allPages, err := ports.List(p.networkClient, opts).AllPages()
		if err != nil {
			return err
		}
		allPorts, err := ports.ExtractPorts(allPages)
		neutron.Update(p.networkClient, )
	} else {

	}
}



func deleteNeutronFloatingIP(ip neutron.FloatingIP) error {
	return nil
}

func deleteNovaFloatingIP(ip nova.FloatingIP) error {
	return nil
}

func getNeutronIP(p *provider) (string, error) {
	if p.usePreAllocated == true {
		allPages, err := neutron.List(p.networkClient, neutron.ListOpts{
			FloatingNetworkID: p.networkId,
		}).AllPages()
		if err != nil {
			return "", err
		}
		allFloatingIPs, err := neutron.ExtractFloatingIPs(allPages)
		if err != nil {
			return "", err
		}
		for _, ip := range allFloatingIPs {
			if ip.FixedIP == "" {
				return ip.FloatingIP, nil
			}
		}
		return "", errors.New("no preallocated floating ips available in pool")
	} else {
		ip, err := neutron.Create(p.networkClient, neutron.CreateOpts{
			FloatingNetworkID: p.networkId,
		}).Extract()
		if err != nil {
			return "", err
		}
		p.floatingIps[ip.FloatingIP] = ip
		return ip.FloatingIP, nil
	}
}

func getNovaIP(p *provider) (string, error) {

	if p.usePreAllocated == true {
		allPages, err := nova.List(p.computeClient).AllPages()
		if err != nil {
			return "", err
		}
		allFloatingIPs, err := nova.ExtractFloatingIPs(allPages)
		for _, ip := range allFloatingIPs {
			if ip.InstanceID == "" {
				return ip.IP, nil
			}
		}
		return "", errors.New("no preallocated floating ips available in pool")

	} else {
		ip, err := nova.Create(p.computeClient, nova.CreateOpts{
			Pool: p.pool,
		}).Extract()
		if err != nil {
			return "", err
		}
		p.floatingIps[ip.IP] = ip
		return ip.IP, nil
	}
}
