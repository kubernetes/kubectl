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
	"fmt"

	"github.com/ghodss/yaml"

	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
	"k8s.io/kubectl/pkg/kinflate/constants"
	interror "k8s.io/kubectl/pkg/kinflate/internal/error"
	"k8s.io/kubectl/pkg/kinflate/resource"
	"k8s.io/kubectl/pkg/kinflate/transformers"
	"k8s.io/kubectl/pkg/kinflate/types"
	"k8s.io/kubectl/pkg/loader"
)

type Application interface {
	// Resources computes and returns the resources for the app.
	Resources() (types.ResourceCollection, error)
	// RawResources computes and returns the raw resources from the manifest.
	// It contains resources from 1) untransformed resources from current manifest 2) transformed resources from sub packages
	RawResources() (types.ResourceCollection, error)
}

var _ Application = &applicationImpl{}

// Private implementation of the Application interface
type applicationImpl struct {
	manifest *manifest.Manifest
	loader   loader.Loader
}

// NewApp parses the manifest at the path using the loader.
func New(loader loader.Loader) (Application, error) {
	// load the manifest using the loader
	manifestBytes, err := loader.Load(constants.KubeManifestFileName)
	if err != nil {
		return nil, err
	}

	var m manifest.Manifest
	err = yaml.Unmarshal(manifestBytes, &m)
	if err != nil {
		return nil, err
	}
	return &applicationImpl{manifest: &m, loader: loader}, nil
}

// Resources computes and returns the resources from the manifest.
func (a *applicationImpl) Resources() (types.ResourceCollection, error) {
	errs := &interror.ManifestErrors{}
	raw, err := a.RawResources()
	if err != nil {
		errs.Append(err)
	}

	resources := []*resource.Resource{}
	cms, err := resource.NewFromConfigMaps(a.loader, a.manifest.Configmaps)
	if err != nil {
		errs.Append(err)
	} else {
		resources = append(resources, cms...)
	}
	secrets, err := resource.NewFromSecretGenerators(a.loader.Root(), a.manifest.SecretGenerators)
	if err != nil {
		errs.Append(err)
	} else {
		resources = append(resources, secrets...)
	}

	ps, err := resource.NewFromPaths(a.loader, a.manifest.Patches)
	if err != nil {
		errs.Append(err)
	}
	// TODO: remove this func after migrating to ResourceCollection
	patches, err := resourceCollectionFromResources(ps)
	if err != nil {
		errs.Append(err)
	}

	if len(errs.Get()) > 0 {
		return nil, errs
	}

	allResources, err := resourceCollectionFromResources(resources)
	if err != nil {
		return nil, err
	}
	err = types.Merge(allResources, raw)
	if err != nil {
		return nil, err
	}

	t, err := a.getTransformer(patches)
	if err != nil {
		return nil, err
	}
	err = t.Transform(types.KObject(allResources))
	if err != nil {
		return nil, err
	}
	return allResources, nil
}

// RawResources computes and returns the raw resources from the manifest.
func (a *applicationImpl) RawResources() (types.ResourceCollection, error) {
	subAppResources, errs := a.subAppResources()
	allRes, err := resource.NewFromPaths(a.loader, a.manifest.Resources)
	if err != nil {
		errs.Append(err)
	}
	// TODO: remove this func after migrating to ResourceCollection
	allResources, err := resourceCollectionFromResources(allRes)
	if err != nil {
		errs.Append(err)
	}
	if len(errs.Get()) > 0 {
		return nil, errs
	}

	err = types.Merge(allResources, subAppResources)
	return allResources, err
}

func (a *applicationImpl) subAppResources() (types.ResourceCollection, *interror.ManifestErrors) {
	allResources := types.ResourceCollection{}
	errs := &interror.ManifestErrors{}
	for _, pkgPath := range a.manifest.Packages {
		subloader, err := a.loader.New(pkgPath)
		if err != nil {
			errs.Append(err)
			continue
		}
		subapp, err := New(subloader)
		if err != nil {
			errs.Append(err)
			continue
		}
		// Gather all transformed resources from subpackages.
		subAppResources, err := subapp.Resources()
		if err != nil {
			errs.Append(err)
			continue
		}
		types.Merge(allResources, subAppResources)
	}
	return allResources, errs
}

// getTransformer generates the following transformers:
// 1) append hash for configmaps ans secrets
// 2) apply overlay
// 3) name prefix
// 4) apply labels
// 5) apply annotations
// 6) update name reference
func (a *applicationImpl) getTransformer(patches types.ResourceCollection) (transformers.Transformer, error) {
	ts := []transformers.Transformer{}

	nht := transformers.NewNameHashTransformer()
	ts = append(ts, nht)

	ot, err := transformers.NewOverlayTransformer(types.KObject(patches))
	if err != nil {
		return nil, err
	}
	ts = append(ts, ot)

	npt, err := transformers.NewDefaultingNamePrefixTransformer(string(a.manifest.NamePrefix))
	if err != nil {
		return nil, err
	}
	ts = append(ts, npt)

	lt, err := transformers.NewDefaultingLabelsMapTransformer(a.manifest.ObjectLabels)
	if err != nil {
		return nil, err
	}
	ts = append(ts, lt)

	at, err := transformers.NewDefaultingAnnotationsMapTransformer(a.manifest.ObjectAnnotations)
	if err != nil {
		return nil, err
	}
	ts = append(ts, at)

	nrt, err := transformers.NewDefaultingNameReferenceTransformer()
	if err != nil {
		return nil, err
	}
	ts = append(ts, nrt)
	return transformers.NewMultiTransformer(ts), nil
}

func resourceCollectionFromResources(res []*resource.Resource) (types.ResourceCollection, error) {
	out := types.ResourceCollection{}
	for _, res := range res {
		gvkn := res.GVKN()
		if _, found := out[gvkn]; found {
			return nil, fmt.Errorf("duplicated %#v is not allowed", gvkn)
		}
		out[gvkn] = res.Data
	}
	return out, nil
}
