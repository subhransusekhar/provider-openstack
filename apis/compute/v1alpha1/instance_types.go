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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// Known Instance statuses.
const (
	StatusActive = "ACTIVE"
	StatusBuild = "BUILD"
	StatusDeleted = "DELETED"
	StatusError = "ERROR"
	StatusHardReboot = "HARD_REBOOT"
	StatusMigrating = "MIGRATING"
	StatusPassword = "PASSWORD"
	StatusPaused = "PAUSED"
	StatusReboot = "REBOOT"
	StatusRebuild = "REBUILD"
	StatusRescue = "RESCUE"
	StatusResize = "RESIZE"
	StatusRevertResize = "REVERT_RESIZE"
	StatusShelved = "SHELVED"
	StatusShelvedOffloaded = "SHELVED_OFFLOADED"
	StatusShutoff = "SHUTOFF"
	StatusSoftDeleted = "SOFT_DELETED"
	StatusSuspended = "SUSPENDED"
	StatusUnknown = "UNKNOWN"
	StatusVerifyResize = "VERIFY_RESIZE"
)

// InstanceParameters define the desired state of a OpenStack Instance.
// Most fields map directly to a Instance:
type InstanceParameters struct {
	// Region: The unique slug identifier for the region that you wish to
	// deploy in.
	// +immutable
	Region string `json:"region"`

	// Flavor: The unique slug identifier for the Flavor that you wish to select
	// for this Instance.
	// +immutable
	Flavor string `json:"flavor"`

	// Image: The image ID of a public or private image, or the unique slug
	// identifier for a public image. This image will be the base image for
	// your Instance.
	// +immutable
	Image string `json:"image"`

	// SSHKeys: An array containing the IDs or fingerprints of the SSH keys
	// that you wish to embed in the Instance's root account upon creation.
	// +optional
	// +immutable
	SSHKeys []string `json:"ssh_keys"`

	// Backups: A boolean indicating whether automated backups should be enabled
	// for the Instance. Automated backups can only be enabled when the Instance is
	// created.
	// +optional
	// +immutable
	Backups *bool `json:"backups"`

	// FloatingNetworking:  A boolean indicating whether to have floating IP on the instance.
	// +optional
	// +immutable
	FloatingNetworking *bool `json:"floating_networking"`

	// Monitoring: A boolean indicating whether to install the nodeexporter
	// agent for monitoring.
	// +optional
	// +immutable
	Monitoring *bool `json:"monitoring"`

	// Volumes: A flat array including the unique string identifier for each block
	// storage volume to be attached to the Instance. At the moment a volume can only
	// be attached to a single Instance.
	// +optional
	// +immutable
	Volumes []string `json:"volumes,omitempty"`

	// Tags: A flat array of tag names as strings to apply to the Instance after it
	// is created. Tag names can either be existing or new tags.
	// +optional
	// +immutable
	Tags []string `json:"tags"`

	// NETUUID: A string specifying the UUID of the VPC to which the Instance
	// will be assigned.
	// +optional
	// +immutable
	NETUUID *string `json:"network_uuid,omitempty"`

}

// A InstanceObservation reflects the observed state of a Instance on OpenStack.
type InstanceObservation struct {
	// CreationTimestamp in RFC3339 text format.
	CreationTimestamp string `json:"creationTimestamp,omitempty"`

	// ID for the resource. This identifier is defined by the server.
	ID int `json:"id,omitempty"`

	// A Status string indicating the state of the Instance instance.
	//
	// Possible values:
	// "ACTIVE"
	// "BUILD"
	// "DELETED"
	// "ERROR"
	// "HARD_REBOOT"
	// "MIGRATING"
	// "PASSWORD"
	// "PAUSED"
	// "REBOOT"
	// "REBUILD"
	// "RESCUE"
	// "RESIZE"
	// "REVERT_RESIZE"
	// "SHELVED"
	// "SHELVED_OFFLOADED"
	// "SHUTOFF"
	// "SOFT_DELETED"
	// "SUSPENDED"
	// "UNKNOWN"
	// "VERIFY_RESIZE"
	Status string `json:"status,omitempty"`
}

// A InstanceSpec defines the desired state of a Instance.
type InstanceSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       InstanceParameters `json:"forProvider"`
}

// A InstanceStatus represents the observed state of a Instance.
type InstanceStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          InstanceObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Instance is a managed resource that represents a OpenStack Instance.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,do}
type Instance struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InstanceSpec   `json:"spec"`
	Status InstanceStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// InstanceList contains a list of Instance.
type InstanceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Instance `json:"items"`
}
