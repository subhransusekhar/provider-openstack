# provider-openstack

## Overview

`provider-openstack` is the Crossplane infrastructure provider for the
[OpenStack](https://openstack.org/). The provider that is built from the source
code in this repository can be installed into a Crossplane control plane and
adds the following new functionality:

* Custom Resource Definitions (CRDs) that model OpenStack infrastructure and services
  (e.g. Instances, etc.)
* Controllers to provision these resources in OpenStack based on the users
  desired state captured in CRDs they create
* Implementations of Crossplane's portable resource abstractions, enabling OpenStack
  resources to fulfill a user's general need for cloud services

## Getting Started and Documentation

For getting started guides, installation, deployment, and administration, see
our [Documentation](https://crossplane.io/docs/latest).

## Developing

Run against a Kubernetes cluster:

```console
make run
```

Install `latest` into Kubernetes cluster where Crossplane is installed:

```console
make install
```

Install local build into [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/)
cluster where Crossplane is installed:

```console
make install-local
```

Build, push, and install:

```console
make all
```

Build image:

```console
make image
```

Push image:

```console
make push
```

Build binary:

```console
make build
```
