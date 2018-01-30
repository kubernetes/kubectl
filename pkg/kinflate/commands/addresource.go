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

package commands

import (
	"io"

	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/kinflate/util/fs"
)

type addResourceOptions struct {
}

// NewCmdAddResource adds the name of a file containing a resource to the manifest.
func NewCmdAddResource(out, errOut io.Writer, fs fs.FileSystem) *cobra.Command {
	var o addResourceOptions

	cmd := &cobra.Command{
		Use:   "addresource",
		Short: "Add the name of a file containing a resource to the manifest.",
		Long:  "Add the name of a file containing a resource to the manifest.",
		Example: `
		addresource {filepath}`,
		RunE: func(cmd *cobra.Command, args []string) error {
			err := o.Validate(cmd, args)
			if err != nil {
				return err
			}
			err = o.Complete(cmd, args)
			if err != nil {
				return err
			}
			return o.RunAddResource(cmd, out, errOut, fs)
		},
	}
	return cmd
}

// Validate validates addResource command.
func (o *addResourceOptions) Validate(cmd *cobra.Command, args []string) error {
	return nil
}

// Complete completes addResource command.
func (o *addResourceOptions) Complete(cmd *cobra.Command, args []string) error {
	return nil
}

// RunAddResource runs addResource command (do real work).
func (o *addResourceOptions) RunAddResource(cmd *cobra.Command, out, errOut io.Writer, fs fs.FileSystem) error {
	// error if unable to read the file argument
	// error if unable to find kubemanist in current directory
	// error if unable to parse Kube-manifest from the current directory
	// error (or just INFO) and exit if the resource is already in the manifest
	// add the resource
	// write the new manifest, error if trouble writing it.
	return nil
}
