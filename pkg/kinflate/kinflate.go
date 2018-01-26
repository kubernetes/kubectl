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

	"github.com/spf13/cobra"

	kutil "k8s.io/kubectl/pkg/kinflate/util"
)

type kinflateOptions struct {
	manifestPath string
	namespace    string
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

	cmd.Flags().StringVarP(&o.manifestPath, "filename", "f", "", "Pass in a Kube-manifest.yaml file or a directory that contains the file.")
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
	m, err := LoadFromManifestPath(o.manifestPath)
	if err != nil {
		return err
	}
	res, err := kutil.Encode(m)
	if err != nil {
		return err
	}
	_, err = out.Write(res)
	return err
}
