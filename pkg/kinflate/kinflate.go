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
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
	"k8s.io/kubectl/pkg/scheme"
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
		kinflate -f somedir/`,
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
	decoder := unstructured.UnstructuredJSONScheme

	baseFiles, overlayFiles, overlayPkg, err := loadBaseAndOverlayPkg(o.manifestDir)
	if err != nil {
		return err
	}

	// This func will build a visitor given filenameOptions.
	// It will visit each info and populate the map.
	populateResourceMap := func(files []string, m map[groupVersionKindName][]byte) error {
		for _, file := range files {
			content, err := ioutil.ReadFile(file)
			if err != nil {
				return err
			}

			// try converting to json, if there is a error, probably because the content is already json.
			jsoncontent, err := yaml.YAMLToJSON(content)
			if err != nil {
				fmt.Fprintf(errOut, "error when trying to convert yaml to json: %v\n", err)
			} else {
				content = jsoncontent
			}

			obj, gvk, err := decoder.Decode(content, nil, nil)
			if err != nil {
				return err
			}
			accessor, err := meta.Accessor(obj)
			if err != nil {
				return err
			}
			name := accessor.GetName()
			gvkn := groupVersionKindName{gvk: *gvk, name: name}
			if err != nil {
				return err
			}
			if _, found := m[gvkn]; found {
				return fmt.Errorf("unexpected same groupVersionKindName: %#v", gvkn)
			}
			m[gvkn] = content
		}
		return nil
	}

	// map from GroupVersionKind to marshaled json bytes
	overlayResouceMap := map[groupVersionKindName][]byte{}
	err = populateResourceMap(overlayFiles, overlayResouceMap)
	if err != nil {
		return err
	}

	// map from GroupVersionKind to marshaled json bytes
	baseResouceMap := map[groupVersionKindName][]byte{}
	err = populateResourceMap(baseFiles, baseResouceMap)
	if err != nil {
		return err
	}

	// Strategic merge the resources exist in both base and overlay.
	for gvkn, base := range baseResouceMap {
		// Merge overlay with base resource.
		if overlay, found := overlayResouceMap[gvkn]; found {
			versionedObj, err := scheme.Scheme.New(gvkn.gvk)
			if err != nil {
				switch {
				case runtime.IsNotRegisteredError(err):
					return fmt.Errorf("CRD and TPR are not supported now: %v", err)
				default:
					return err
				}
			}
			merged, err := strategicpatch.StrategicMergePatch(base, overlay, versionedObj)
			if err != nil {
				return err
			}
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
	for _, jsonObj := range baseResouceMap {
		yamlObj, err := updateMetadata(jsonObj, overlayPkg)
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "---\n%s", yamlObj)
	}
	return nil
}
