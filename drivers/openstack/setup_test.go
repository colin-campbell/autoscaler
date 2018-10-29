// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package openstack

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// swapEnv replaces environment variables with new values and returns the old values
func swapEnv(newEnv map[string]string) map[string]string {
	oldEnv := make(map[string]string)
	for k, v := range newEnv {
		oldEnv[k] = os.Getenv(k)
		_ = os.Setenv(k,v)
	}
	return oldEnv
}

// helperLoadTestData loads test data from files in folder "testdata".
// "testdata" is ignored by `go build`.
func helperLoadTestData(t *testing.T, name string) string {
	path := filepath.Join("testdata", name)
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(bytes)
}

// testEnv contains our OpenStack connection details for the tests.
var testEnv = map[string]string{
	"OS_AUTH_URL": "http://ops.my.cloud",
	"OS_ENDPOINT_TYPE": "publicURL",
	//"OS_IDENTITY_API_VERSION": "2",
	"OS_PASSWORD": "admin",
	"OS_DOMAIN_ID": "default",
	"OS_REGION_NAME": "RegionOne",
	"OS_TENANT_NAME": "demo",
	"OS_USERNAME": "admin",
	"DRONE_OPENSTACK_IP_POOL": "my-ip-pool",
	"DRONE_OPENSTACK_IP_USE_PREALLOCATED": "False",
	"DRONE_OPENSTACK_IP_CREATE_NEW": "False",
	"DRONE_OPENSTACK_SSH_KEY": "drone-ci-key",
	"DRONE_OPENSTACK_SECURITY_GROUP": "drone-agent",
	"DRONE_OPENSTACK_FLAVOR": "v1-standard-2",
	"DRONE_OPENSTACK_IMAGE": "ubuntu-16.04-server-latest",
	"DRONE_OPENSTACK_METADATA": "name:agent,owner:drone-ci",
}

// testToken contains our test OS authentication token.
var testToken = "gAAAAABbx5Nf3c0AXk2JDNdW6t534lqHtk0xKvpcjuSgBOTSYHF_y-Q2nwLSD5l8AfDAKtxAVsOla9gLZfy3uGWfkgr2rCiTe7cnUe3mH6fAz9nsNz1LcR6TyCCIMv__cNoTYXLMkFr2e7G2-f_FgjpoCP_WgZeK4InBydHKGA0UBvTZwYHXMaM"
