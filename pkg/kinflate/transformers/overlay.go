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
	"encoding/json"
	"fmt"

	jsonpatch "github.com/evanphx/json-patch"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/kubectl/pkg/kinflate/resource"
	"k8s.io/kubectl/pkg/scheme"
)

// overlayTransformer contains a map of overlay objects
type overlayTransformer struct {
	overlay []*resource.Resource
}

var _ Transformer = &overlayTransformer{}

// NewOverlayTransformer constructs a overlayTransformer.
func NewOverlayTransformer(overlay []*resource.Resource) (Transformer, error) {
	if len(overlay) == 0 {
		return NewNoOpTransformer(), nil
	}
	return &overlayTransformer{overlay}, nil
}

// Transform apply the overlay on top of the base resources.
func (o *overlayTransformer) Transform(baseResourceMap resource.ResourceCollection) error {
	// Strategic merge the resources exist in both base and overlay.
	for _, overlay := range o.overlay {
		// Merge overlay with base resource.
		gvkn := overlay.GVKN()
		base, found := baseResourceMap[gvkn]
		if !found {
			return fmt.Errorf("failed to find an object with %#v to apply the patch", gvkn.GVK)
		}
		merged := map[string]interface{}{}
		versionedObj, err := scheme.Scheme.New(gvkn.GVK)
		baseName := base.Data.GetName()
		switch {
		case runtime.IsNotRegisteredError(err):
			// Use JSON merge patch to handle types w/o schema
			baseBytes, err := json.Marshal(base.Data)
			if err != nil {
				return err
			}
			patchBytes, err := json.Marshal(overlay.Data)
			if err != nil {
				return err
			}
			mergedBytes, err := jsonpatch.MergePatch(baseBytes, patchBytes)
			if err != nil {
				return err
			}
			err = json.Unmarshal(mergedBytes, &merged)
			if err != nil {
				return err
			}
		case err != nil:
			return err
		default:
			// Use Strategic Merge Patch to handle types w/ schema
			// TODO: Change this to use the new Merge package.
			// Store the name of the base object, because this name may have been munged.
			// Apply this name to the StrategicMergePatched object.
			merged, err = strategicpatch.StrategicMergeMapPatch(
				base.Data.UnstructuredContent(),
				overlay.Data.UnstructuredContent(),
				versionedObj)
			if err != nil {
				return err
			}
		}
		base.Data.SetName(baseName)
		baseResourceMap[gvkn].Data.Object = merged
	}
	return nil
}
