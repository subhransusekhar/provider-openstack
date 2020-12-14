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
	"context"

	"github.com/gophercloud/gophercloud"
  "github.com/gophercloud/gophercloud/openstack"
  "github.com/gophercloud/gophercloud/openstack/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	"github.com/subhransusekhar/provider-openstack/apis/compute/v1alpha1"
	os "github.com/subhransusekhar/provider-openstack/pkg/clients"
	oscompute "github.com/subhransusekhar/provider-openstack/pkg/clients/compute"
)

const (
	// Error strings.
	errNotInstance = "managed resource is not a Instance resource"
	errGetInstance = "cannot get instance"

	errInstanceCreateFailed = "creation of Instance resource has failed"
	errInstanceDeleteFailed = "deletion of Instance resource has failed"
	errInstanceUpdate       = "cannot update managed Instance resource"
)

// SetupInstance adds a controller that reconciles Instance managed
// resources.
func SetupInstance(mgr ctrl.Manager, l logging.Logger) error {
	name := managed.ControllerName(v1alpha1.InstanceGroupKind)

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&v1alpha1.Instance{}).
		Complete(managed.NewReconciler(mgr,
			resource.ManagedKind(v1alpha1.InstanceGroupVersionKind),
			managed.WithExternalConnecter(&instanceConnector{kube: mgr.GetClient()}),
			managed.WithReferenceResolver(managed.NewAPISimpleReferenceResolver(mgr.GetClient())),
			managed.WithConnectionPublishers(),
			managed.WithLogger(l.WithValues("controller", name)),
			managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name)))))
}

type instanceConnector struct {
	kube client.Client
}

func (c *instanceConnector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	endpoint, project, user, token, err := os.GetAuthInfo(ctx, c.kube, mg)
	if err != nil {
		return nil, err
	}
	opts := gophercloud.AuthOptions {
  IdentityEndpoint: endpoint,
  Username: user,
  Password: token,
	}
	client := openstack.AuthenticatedClient(opts)
	return &instanceExternal{Client: client, kube: c.kube}, nil
}

type instanceExternal struct {
	kube client.Client
	*goos.Client
}

func (c *instanceExternal) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) {
	cr, ok := mg.(*v1alpha1.Instance)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotInstance)
	}
	observed, _, err := c.Instances.Get(ctx, cr.Status.AtProvider.ID)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errGetInstance)
	}

	currentSpec := cr.Spec.ForProvider.DeepCopy()
	oscompute.LateInitializeSpec(&cr.Spec.ForProvider, *observed)
	if !cmp.Equal(currentSpec, &cr.Spec.ForProvider) {
		if err := c.kube.Update(ctx, cr); err != nil {
			return managed.ExternalObservation{}, errors.Wrap(err, errInstanceUpdate)
		}
	}

	cr.Status.AtProvider = v1alpha1.InstanceObservation{
		CreationTimestamp: observed.Created,
		ID:                observed.ID,
		Status:            observed.Status,
	}

	switch cr.Status.AtProvider.Status {
	case v1alpha1.StatusBuild:
		cr.SetConditions(xpv1.Creating())
	case v1alpha1.StatusActive:
		cr.SetConditions(xpv1.Available())
	}

	// Instances are always "up to date" because they can't be updated. ¯\_(ツ)_/¯
	return managed.ExternalObservation{
		ResourceExists:   true,
		ResourceUpToDate: true,
	}, nil
}

func (c *instanceExternal) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.Instance)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotInstance)
	}

	cr.Status.SetConditions(xpv1.Creating())

	create := &goos.InstanceCreateRequest{}
	oscompute.GenerateInstance(meta.GetExternalName(cr), cr.Spec.ForProvider, create)

	instance, _, err := c.Instances.Create(ctx, create)
	if err != nil {
		cr.Status.AtProvider.ID = instance.ID
		cr.Status.AtProvider.CreationTimestamp = instance.Created
		cr.Status.AtProvider.Status = instance.Status
	}
	return managed.ExternalCreation{}, errors.Wrap(err, errInstanceCreateFailed)
}

func (c *instanceExternal) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	// Instances cannot be updated.
	return managed.ExternalUpdate{}, nil
}

func (c *instanceExternal) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.Instance)
	if !ok {
		return errors.New(errNotInstance)
	}

	cr.Status.SetConditions(xpv1.Deleting())
	_, err := c.Instances.Delete(ctx, cr.Status.AtProvider.ID)
	return errors.Wrap(err, errInstanceDeleteFailed)
}
