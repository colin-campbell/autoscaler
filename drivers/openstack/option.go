// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package openstack

import (
	"github.com/drone/autoscaler/drivers/internal/userdata"
	"github.com/gophercloud/gophercloud"
	"io/ioutil"
)

type Option func(*provider)

// WithImage returns an option to set the instance image.
func WithImage(image string) Option {
	return func(p *provider) {
		p.image = image
	}
}

// WithRegion returns an option to set the OpenStack target region.
func WithRegion(region string) Option {
	return func(p *provider) {
		p.region = region
	}
}

// WithAvailabilityZone returns an option to set the OpenStack availability zone.
func WithAvailabilityZone(zone string) Option {
	return func(p *provider) {
		p.zone = zone
	}
}

// WithFlavor returns an option to set the instance flavor.
func WithFlavor(flavor string) Option {
	return func(p *provider) {
		p.flavor = flavor
	}
}

// WithNetworkType returns an option to set the OpenStack networking type.
func WithNetworkType(networkType string) Option {
	return func(p *provider) {
		p.networkType = networkType
	}
}

// WithNetworkId returns an option to set the OpenStack network id.
func WithNetworkId(id string) Option {
	return func(p *provider) {
		p.networkId = id
	}
}

// WithUsePreallocated returns an option to use preallocated ip addresses.
func WithUsePreAllocated(usePreAllocated bool) Option {
	return func(p *provider) {
		p.usePreAllocated = usePreAllocated
	}
}

// WithSecurityGroup returns an option to set the instance security groups.
func WithSecurityGroup(group ...string) Option {
	return func(p *provider) {
		p.groups = group
	}
}

// WithNetworkClient returns an option to set the
// GopherCloud network ServiceClient.
func WithNetworkClient(networkClient *gophercloud.ServiceClient) Option {
	return func(p *provider) {
		p.networkClient = networkClient
	}
}

// WithComputeClient returns an option to set the
// GopherCloud compute ServiceClient.
func WithComputeClient(computeClient *gophercloud.ServiceClient) Option {
	return func(p *provider) {
		p.computeClient = computeClient
	}
}

// WithSSHKey returns an option to set the ssh key.
func WithSSHKey(key string) Option {
	return func(p *provider) {
		p.key = key
	}
}

// WithSubnet returns an option to set the subnet id.
func WithSubnet(id string) Option {
	return func(p *provider) {
		p.subnet = id
	}
}

func WithFloatingIpPool(pool string) Option {
	return func(p *provider) {
		p.pool = pool
	}
}

// WithMetadata returns an option to set the instance metadata.
func WithMetadata(metadata map[string]string) Option {
	return func(p *provider) {
		p.metadata = metadata
	}
}

// WithUserData returns an option to set the cloud-init
// template from text.
func WithUserData(text string) Option {
	return func(p *provider) {
		if text != "" {
			p.userdata = userdata.Parse(text)
		}
	}
}

// WithUserDataFile returns an option to set the cloud-init
// template from file.
func WithUserDataFile(filepath string) Option {
	return func(p *provider) {
		if filepath != "" {
			b, err := ioutil.ReadFile(filepath)
			if err != nil {
				panic(err)
			}
			p.userdata = userdata.Parse(string(b))
		}
	}
}
