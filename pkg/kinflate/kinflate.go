/*
Copyright 2017 The Kubernetes Authors.

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

package kinflate

import (
	"fmt"
	"io"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/common/categories"
	"k8s.io/common/resource"
)

type kinflateOptions struct {
	manifestDir string
	namespace   string
}

type groupVersionKindName struct {
	gvk schema.GroupVersionKind
	// name of the resource.
	name string
}

// NewCmdKinflate creates a new kinflate command.
func NewCmdKinflate(out, errOut io.Writer) *cobra.Command {
	var o kinflateOptions

	cmd := &cobra.Command{
		Use:   "kinflate -f [path]",
		Short: "Use a Manifest file to generate a set of api resources",
		Long:  "Use a Manifest file to generate a set of api resources",
		Example: `
		# Use the Kube-manifest.yaml file under somedir/ to generate a set of api resources.
		kubectl alpha kinflate -f somedir/`,
		Run: func(cmd *cobra.Command, args []string) {
			err := o.Validate(cmd, args)
			if err != nil {
				panic(err)
			}
			err = o.Complete(cmd, args)
			if err != nil {
				panic(err)
			}
			err = o.RunKinflate(cmd, out, errOut)
			if err != nil {
				panic(err)
			}
		},
	}

	cmd.Flags().StringVarP(&o.manifestDir, "filename", "f", "", "Pass in directory that contains the Kube-manifest.yaml file.")
	cmd.MarkFlagRequired("filename")
	cmd.Flags().StringVarP(&o.namespace, "namespace", "o", "yaml", "Output mode. Support json or yaml.")
	return cmd
}

// Validate validates kinflate command.
func (o *kinflateOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// Complete completes kinflate command.
func (o *kinflateOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

// RunKinflate runs kinflate command (do real work).
func (o *kinflateOptions) RunKinflate(cmd *cobra.Command, out, errOut io.Writer) error {
	encoder := Codecs.LegacyCodec(Registry.EnabledVersions()...)
	decoder := unstructured.UnstructuredJSONScheme

	// mapper, typer and categoryExpander are not from discovery to do it in standalone mode.
	mapper := Registry.RESTMapper()
	typer := Scheme
	categoryExpander := categories.LegacyCategoryExpander

	baseFilesOpsList, overlayFileOp, overlayPkg, err := loadBaseAndOverlayPkg(o.manifestDir)
	if err != nil {
		return err
	}

	// This func will build a visitor given filenameOptions.
	// It will visit each info and populate the map.
	populateResourceMap := func(fileOps resource.FilenameOptions, m map[groupVersionKindName]runtime.Object) error {
		r := resource.NewBuilder(mapper, categoryExpander, typer, resource.DisabledClientForMapping{}, unstructured.UnstructuredJSONScheme).
			ContinueOnError().
			NamespaceParam(o.namespace).DefaultNamespace().
			FilenameParam(true, &fileOps).
			Flatten().
			Do()
		err = r.Err()
		if err != nil {
			return err
		}
		return r.Visit(func(info *resource.Info, err error) error {
			var obj runtime.Object
			obj = info.VersionedObject
			if info.VersionedObject == nil {
				obj = info.Object
			}
			gvk := obj.GetObjectKind().GroupVersionKind()
			accessor, err := meta.Accessor(obj)
			if err != nil {
				return err
			}
			name := accessor.GetName()
			gvkn := groupVersionKindName{gvk: gvk, name: name}
			if err != nil {
				return err
			}
			if _, found := m[gvkn]; found {
				return fmt.Errorf("unexpected same groupVersionKindName: %#v", gvkn)
			}
			m[gvkn] = obj
			return nil
		})
	}

	// map from GroupVersionKind to marshaled json bytes
	overlayResouceMap := map[groupVersionKindName]runtime.Object{}
	err = populateResourceMap(overlayFileOp, overlayResouceMap)
	if err != nil {
		return err
	}

	// map from GroupVersionKind to marshaled json bytes
	baseResouceMap := map[groupVersionKindName]runtime.Object{}
	for _, baseFilesOps := range baseFilesOpsList {
		err = populateResourceMap(baseFilesOps, baseResouceMap)
		if err != nil {
			return err
		}
	}

	// Strategic merge the resources exist in both base and overlay.
	for gvkn, base := range baseResouceMap {
		// Merge overlay with base resource.
		if overlay, found := overlayResouceMap[gvkn]; found {
			versionedObj, err := Scheme.New(gvkn.gvk)
			if err != nil {
				switch {
				case runtime.IsNotRegisteredError(err):
					return fmt.Errorf("CRD and TPR are not supported now: %v", err)
				default:
					return err
				}
			}
			jsonBase, err := runtime.Encode(encoder, base)
			if err != nil {
				return err
			}
			jsonOverlay, err := runtime.Encode(encoder, overlay)
			if err != nil {
				return err
			}
			jsonMerged, err := strategicpatch.StrategicMergePatch(jsonBase, jsonOverlay, versionedObj)
			if err != nil {
				return err
			}
			merged, _, err := decoder.Decode(jsonMerged, nil, nil)
			baseResouceMap[gvkn] = merged
			delete(overlayResouceMap, gvkn)
		}
	}

	// If there are resources in overlay that are not defined in base, just add it to base.
	if len(overlayResouceMap) > 0 {
		for gvkn, jsonObj := range overlayResouceMap {
			baseResouceMap[gvkn] = jsonObj
		}
	}

	// Inject the labels, annotations and name prefix.
	// Then print the object.
	for _, obj := range baseResouceMap {
		err = updateMetadata(obj, overlayPkg)
		if err != nil {
			return err
		}
		yamlObj, err := yaml.Marshal(obj)
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "---\n%s", yamlObj)
	}
	return nil
}
