// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package openstack

import (
	"context"
	"fmt"
	"github.com/drone/autoscaler"
	"github.com/h2non/gock"
	"os"
	"testing"
)

func TestCreate(t *testing.T) {

	defer gock.Off()
	oldEnv := swapEnv(testEnv)
	defer swapEnv(oldEnv)

	setupEndpoints(t)

	v, err := New(
		WithSSHKey(os.Getenv("DRONE_OPENSTACK_SSH_KEY")),
		WithFloatingIpPool(os.Getenv("DRONE_OPENSTACK_IP_POOL")),
		WithFlavor(os.Getenv("DRONE_OPENSTACK_FLAVOR")),
		WithImage(os.Getenv("DRONE_OPENSTACK_IMAGE")),
		WithFlavor(os.Getenv("DRONE_OPENSTACK_FLAVOR")),
		WithMetadata(map[string]string{"name": "agent", "owner": "drone-ci"}),
		WithRegion(os.Getenv("OS_REGION_NAME")),
		)
	if err != nil {
		t.Error(err)
		return
	}
	p := v.(*provider)
	p.init.Do(func() {})

	instance, err := p.Create(context.TODO(), autoscaler.InstanceCreateOpts{Name: "agent-097i6IDf"})
	if err != nil {
		t.Error(err)
	}

	if !gock.IsDone() {
		t.Errorf("Expected http requests not detected")
	}
	t.Run("Instance Attributes", testInstance(instance))
}

func testInstance(instance *autoscaler.Instance) func(t *testing.T) {
	return func(t *testing.T) {
		if instance == nil {
			t.Errorf("Expect non-nil instance even if error")
		}
		if got, want := instance.ID, "951172e9-ac8b-4df8-a649-2f52107f7b5f"; got != want {
			t.Errorf("Want instance ID %v, got %v", want, got)
		}
		if got, want := instance.Image, "ubuntu-16.04-server-latest"; got != want {
			t.Errorf("Want Image %v, got %v", want, got)
		}
		if got, want := instance.Name, "agent-097i6IDf"; got != want {
			t.Errorf("Want instance Name %v, got %v", want, got)
		}
		if got, want := instance.Region, "my-region"; got != want {
			t.Errorf("Want instance Region %v, got %v", want, got)
		}
		if got, want := instance.Provider, autoscaler.ProviderOpenStack; got != want {
			t.Errorf("Want instance Provider %v, got %v", want, got)
		}
	}
}
func setupEndpoints(t *testing.T) {
	authResp1 := helperLoadTestData(t, "authresp1.json")

	gock.New(testEnv["OS_AUTH_URL"]).
		Get("/").
		Reply(300).
		SetHeader("Content-Type", "application/json").
		BodyString(authResp1)

	authResp2 := helperLoadTestData(t, "authresp2.json")
	gock.New(testEnv["OS_AUTH_URL"]).
		Post("/v3/auth/tokens").
		Reply(201).
		SetHeader("Content-Type", "application/json").
		SetHeader("X-Subject-Token", testToken).
		BodyString(authResp2)

	// Sorry, not unmarshalling Json to get this :)
	computeUrl := "https://my.openstack.cloud:8774"
	fipResp := helperLoadTestData(t, "fipresp1.json")
	gock.New(computeUrl).
		MatchHeader("X-Auth-Token", testToken).
		Post("/v2.1/os-floating-ips").
		JSON(map[string]string{"pool": os.Getenv("DRONE_OPENSTACK_IP_POOL")}).
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(fipResp)

	imageListResp := helperLoadTestData(t, "imagelistresp1.json")
	fmt.Println("Here")
	gock.New(computeUrl).
		MatchHeader("X-Auth-Token", testToken).
		Get("/v2.1/images/detail").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(imageListResp)

	flavorListResp := helperLoadTestData(t, "flavorlistresp1.json")
	gock.New(computeUrl).
		MatchHeader("X-Auth-Token", testToken).
		Get("/v2.1/flavors/detail").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(flavorListResp)

	serverCreateResp := helperLoadTestData(t, "servercreateresp1.json")
	gock.New(computeUrl).
		MatchHeader("X-Auth-Token", testToken).
		Post("/v2.1/servers").
		Reply(202).
		SetHeader("Content-Type", "application/json").
		BodyString(serverCreateResp)

	serverStatusResp := helperLoadTestData(t, "serverstatusresp1.json")
	gock.New(computeUrl).
		MatchHeader("X-Auth-Token", testToken).
		Get("/v2.1/servers/951172e9-ac8b-4df8-a649-2f52107f7b5f").
		Reply(200).
		SetHeader("Content-Type", "application/json").
		BodyString(serverStatusResp)
}
