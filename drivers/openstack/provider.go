// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package openstack

import (
	"github.com/pkg/errors"
	"sync"
	"text/template"

	"github.com/drone/autoscaler"
	"github.com/drone/autoscaler/drivers/internal/userdata"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
)

// provider implements an OpenStack provider
type provider struct {
	init sync.Once

	key             string
	region          string
	zone            []string
	image           string
	flavor          string
	subnet          string
	pool            string
	networkType     string
	networkId       string
	usePreAllocated bool
	userdata        *template.Template
	groups          []string
	metadata        map[string]string

	floatingIps   map[string]interface{}
	networkClient *gophercloud.ServiceClient
	computeClient *gophercloud.ServiceClient
}

// New returns a new OpenStack provider.
func New(opts ...Option) (autoscaler.Provider, error) {
	var err error
	var authOpts gophercloud.AuthOptions
	var endpointOpts gophercloud.EndpointOpts
	var authClient *gophercloud.ProviderClient
	p := new(provider)

	for _, opt := range opts {
		opt(p)
	}

	if p.region == "" {
		endpointOpts.Region = p.region
	}

	authOpts, err = openstack.AuthOptionsFromEnv()
	if err != nil {
		return nil, err
	}

	authClient, err = openstack.AuthenticatedClient(authOpts)
	if err != nil {
		return nil, err
	}

	if p.userdata == nil {
		p.userdata = userdata.T
	}

	if p.computeClient == nil {
		p.computeClient, err = openstack.NewComputeV2(authClient, endpointOpts)
		if err != nil {
			return nil, err
		}
	}

	if p.networkType == "" {
		p.networkType = "neutron"
	} else if p.networkType != "nova" && p.networkType != "neutron" {
		return nil,
			errors.New("unsupported network type (\"neutron\" or \"nova\". Default: \"neutron\")")
	}

	if p.networkType == "neutron" && p.networkId == "" {
		return nil, errors.New("external network id must be set for neutron networking")
	}

	if p.networkClient == nil && p.networkType == "neutron" {
		p.networkClient, err = openstack.NewNetworkV2(authClient, endpointOpts)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}
