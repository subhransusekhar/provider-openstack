/*
Copyright 2020 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package compute

import (
	"strconv"

	"github.com/gophercloud/gophercloud"
  "github.com/gophercloud/gophercloud/openstack"
  "github.com/gophercloud/gophercloud/openstack/utils"
  "github.com/gophercloud/gophercloud/openstack/compute/v2/servers"

	"github.com/subhransusekhar/provider-openstack/apis/compute/v1alpha1"
	so "github.com/subhransusekhar/provider-openstack/pkg/clients"
)

func GenerateInstance(name string, in v1alpha1.InstanceParameters, create *goos.InstanceCreateRequest) {
	create.Name = name
	create.Region = in.Region
	create.Flavor = in.Flavor
	create.Image = in.Image
	create.SSHKeys = in.SSHKeys
	create.Backups = do.BoolValue(in.Backups)
	create.FloatingNetworking = do.BoolValue(in.FloatingNetworking)
	create.Monitoring = do.BoolValue(in.Monitoring)
	create.Volumes = in.Volumes
	create.Tags = in.Tags
	create.NETUUID = do.StringValue(in.NETUUID)
}

// LateInitializeSpec updates any unset (i.e. nil) optional fields of the
// supplied InstanceParameters that are set (i.e. non-zero) on the supplied
// Instance.
func LateInitializeSpec(p *v1alpha1.InstanceParameters, observed goos.Instance) {
	p.Volumes = os.LateInitializeStringSlice(p.Volumes, observed.VolumeIDs)
	p.Tags = os.LateInitializeStringSlice(p.Tags, observed.Tags)
	p.NETUUID = os.LateInitializeString(p.NETUUID, observed.NETUUID)
}
