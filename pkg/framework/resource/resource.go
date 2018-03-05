/*
Copyright 2018 The Kubernetes Authors.

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

package resource

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	openapi "k8s.io/kube-openapi/pkg/util/proto"
)

// Resource is an API Resource.
type Resource struct {
	Resource        v1.APIResource
	SubResources    []*SubResource
	apiGroupVersion schema.GroupVersion
	schema          openapi.Schema
}

// SubResource is an API subresource.
type SubResource struct {
	Resource        v1.APIResource
	parent          *Resource
	apiGroupVersion schema.GroupVersion
	schema          openapi.Schema
}

// EndpointGroupVersionKind creates a GroupVersionKind object based on
// the GroupVersion and the Kind of the URL endpoint.
func (r *Resource) EndpointGroupVersionKind() schema.GroupVersionKind {
	return r.apiGroupVersion.WithKind(r.Resource.Kind)
}

// EndpointGroupVersionKind creates a GroupVersionKind object based on
// the GroupVersion and the Kind of the URL endpoint.
func (sr *SubResource) EndpointGroupVersionKind() schema.GroupVersionKind {
	return sr.apiGroupVersion.WithKind(sr.Resource.Kind)
}

// ResourceGroupVersionKind returns a GVK based on the request object.
func (r *Resource) ResourceGroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{Group: r.Resource.Group, Version: r.Resource.Version, Kind: r.Resource.Kind}
}

// ResourceGroupVersionKind returns a GVK based on the request object.
func (sr *SubResource) RequestGroupVersionKind() schema.GroupVersionKind {
	return schema.GroupVersionKind{Group: sr.Resource.Group, Version: sr.Resource.Version, Kind: sr.Resource.Kind}
}
