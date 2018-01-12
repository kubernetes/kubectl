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
	"os"
	"sort"

	"github.com/spf13/cobra"

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
				fmt.Fprintf(errOut, "error: %v\n", err)
				os.Exit(1)
			}
			err = o.Complete(cmd, args)
			if err != nil {
				fmt.Fprintf(errOut, "error: %v\n", err)
				os.Exit(1)
			}
			err = o.RunKinflate(cmd, out, errOut)
			if err != nil {
				fmt.Fprintf(errOut, "error: %v\n", err)
				os.Exit(1)
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
	baseResources, overlayResource, overlayPkg, err := loadBaseAndOverlayPkg(o.manifestDir)
	if err != nil {
		return err
	}

	gvknToNewNameObject := map[groupVersionKindName]newNameObject{}

	// map from GroupVersionKind to marshaled json bytes
	overlayResouceMap := map[groupVersionKindName][]byte{}
	err = populateResourceMap(overlayResource.resources, overlayResouceMap, errOut)
	if err != nil {
		return err
	}

	err = populateMapOfConfigMapAndSecret(overlayResource, gvknToNewNameObject)
	if err != nil {
		return err
	}

	// map from GroupVersionKind to marshaled json bytes
	baseResouceMap := map[groupVersionKindName][]byte{}
	for _, baseResource := range baseResources {
		err = populateResourceMap(baseResource.resources, baseResouceMap, errOut)
		if err != nil {
			return err
		}
		err = populateMapOfConfigMapAndSecret(baseResource, gvknToNewNameObject)
		if err != nil {
			return err
		}
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

	cmAndSecretGVKN := []groupVersionKindName{}
	for gvkn := range gvknToNewNameObject {
		cmAndSecretGVKN = append(cmAndSecretGVKN, gvkn)
	}
	sort.Sort(ByGVKN(cmAndSecretGVKN))
	for _, gvkn := range cmAndSecretGVKN {
		nameAndobj := gvknToNewNameObject[gvkn]
		yamlObj, err := updateObjectMetadata(nameAndobj.obj, overlayPkg)
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "---\n%s", yamlObj)
	}

	// Inject the labels, annotations and name prefix.
	// Then print the object.
	resourceGVKN := []groupVersionKindName{}
	for gvkn := range baseResouceMap {
		resourceGVKN = append(resourceGVKN, gvkn)
	}
	sort.Sort(ByGVKN(resourceGVKN))
	for _, gvkn := range resourceGVKN {
		yamlObj, err := updateMetadata(baseResouceMap[gvkn], overlayPkg, gvknToNewNameObject)
		if err != nil {
			return err
		}
		fmt.Fprintf(out, "---\n%s", yamlObj)
	}
	return nil
}
