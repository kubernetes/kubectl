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

package app

import (
	"k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
	"k8s.io/kubectl/pkg/kinflate/resource"
	"k8s.io/kubectl/pkg/loader"
)

type Application interface {
	Resources() []resource.Resource
}

// Private implementation of the Application interface
type applicationImpl struct {
	manifest v1alpha1.Manifest
	loader   loader.Loader
}

// NewApp parses the manifest at the path using the loader.
func New(loader loader.Loader) (Application, error) {
	// load the manifest using the loader
	return &applicationImpl{loader: loader}, nil
}

// Resources computes and returns the resources from the manifest.
func (a *applicationImpl) Resources() []resource.Resource {
	// This is where all the fun happens.
	// From the manifest create the resources.
	return nil
}
