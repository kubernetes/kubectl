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

package transformers

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/kubectl/pkg/kinflate/types"
	"k8s.io/kubectl/pkg/scheme"
)

// OverlayTransformer contains a map of overlay objects
type OverlayTransformer struct {
	overlay types.KObject
}

var _ Transformer = &OverlayTransformer{}

// NewNamePrefixTransformer construct a NamePrefixTransformer.
func NewOverlayTransformer(overlay types.KObject) (Transformer, error) {
	if len(overlay) == 0 {
		return nil, nil
	}
	return &OverlayTransformer{overlay}, nil
}

// Transform prepends the name prefix.
func (o *OverlayTransformer) Transform(baseResourceMap types.KObject) error {
	// Strategic merge the resources exist in both base and overlay.
	for gvkn, base := range baseResourceMap {
		// Merge overlay with base resource.
		if overlay, found := o.overlay[gvkn]; found {
			versionedObj, err := scheme.Scheme.New(gvkn.GVK)
			if err != nil {
				switch {
				case runtime.IsNotRegisteredError(err):
					return fmt.Errorf("CRD and TPR are not supported now: %v", err)
				default:
					return err
				}
			}
			// TODO: Change this to use the new Merge package.
			// Store the name of the base object, because this name may have been munged.
			// Apply this name to the StrategicMergePatched object.
			baseName := base.GetName()
			merged, err := strategicpatch.StrategicMergeMapPatch(
				base.UnstructuredContent(),
				overlay.UnstructuredContent(),
				versionedObj)
			if err != nil {
				return err
			}
			base.SetName(baseName)
			baseResourceMap[gvkn].Object = merged
			delete(o.overlay, gvkn)
		}
	}
	// If there are resources in overlay that are not defined in base, just add it to base.
	if len(o.overlay) > 0 {
		for gvkn, jsonObj := range o.overlay {
			baseResourceMap[gvkn] = jsonObj
		}
	}
	return nil
}
